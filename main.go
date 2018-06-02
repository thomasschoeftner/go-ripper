package main

import (
	"os"
	"github.com/google/logger"
	"io/ioutil"
	"go-cli/config"
	"go-cli/task"
	"go-ripper/ripper"
	"go-cli/pipeline"
	"flag"
	"fmt"
	"go-cli/require"
	"errors"
	"go-cli/cli"
	"go-ripper/omdb"
	"strings"
)

const (
	exit_success = 0
	exit_failure = 1

	cliFlagVerbose             = "verbose"
	cliFlagOmdbTokens          = "omdbtoken"
	cliFlagConfigFile          = "config"
	cliFlagWithoutDependencies = "without-dependencies"
	cliFlagLazy                = "lazy"
	ApplicationName            = "go-ripper"
)



var verbose = cli.FromFlag(cliFlagVerbose, "full log output in console").GetBoolean().WithDefault(false)
var lazy = cli.FromFlag(cliFlagLazy, "execute task only if no output from previous execution is available").GetBoolean().WithDefault(true)
var omdbTokenFlags  = cli.FromFlag(cliFlagOmdbTokens, "the access token für connecting to OMDB - can also be specified as ENV variable").OrEnvironmentVar(cliFlagOmdbTokens).GetArray().WithDefault()
var configFlag = cli.FromFlag(cliFlagConfigFile,"the config file location").OrEnvironmentVar(ApplicationName + "-" + cliFlagConfigFile).GetString().WithDefault(ApplicationName + ".conf")
var runWithoutDependencies = cli.FromFlag(cliFlagWithoutDependencies, "run specified tasks without their dependencies").GetBoolean().WithDefault(false)

func main() {
	os.Exit(launch())
}

func launch() int {
	syntax := "[flags] task[...] target[...]"
	cli.Setup(&syntax, "use task \"tasks\" to display all available tasks:", fmt.Sprintf("  %s tasks\n", os.Args[0]))

	// read config
	conf := getConfig()
	vmiqf, err := omdb.NewOmdbVideoQueryFactory(conf.Resolve.Video.Omdb, *omdbTokenFlags)
	require.NotFailed(err)

	allTasks  := CreateTasks(vmiqf)

	logger.Init(ApplicationName, *verbose, false, ioutil.Discard)

	// create task Tree
	taskMap, errs := task.ValidateTasks(allTasks)
	require.NoneFailed(errs)

	// read command line params (flags & args)
	taskNames, targets := getCliTasksAndTargets(allTasks, taskMap)

	// calculate tasks to be invoked
	invokedTasks, err := taskMap.GetTasksForNames(taskNames...)
	require.NotFailed(err)

	// materialize processing pipeline
	//todo check required flags per task!!!
	tasksToRun := invokedTasks
	if !(*runWithoutDependencies) {
		tasksToRun = invokedTasks.Flatten()
	}
	pipe, err := pipeline.Materialize(tasksToRun).WithConfig(conf.Processing, conf, allTasks, *lazy)
	require.NotFailed(err)

	// ASYNCHRONOUSLY send a processing command for each target to pipeline
	go fillPipelineAndClose(pipe, targets)

	err = handleProcessingEvents(pipe)
	require.NotFailed(err)

	return exit_success
}

func fillPipelineAndClose(pipe *pipeline.Pipeline, targets []string) {
	// feed processing command for each target to pipeline
	for _, target := range targets {
		pipe.Commands <- ripper.ProcessPath(target)
	}
	pipe.Commands <- pipeline.Stop()

	// close pipeline
	close(pipe.Commands)
}

func handleProcessingEvents(pipe *pipeline.Pipeline) error {
	pipeClosed := false
	for !pipeClosed {
		event, notClosed := <-pipe.Events
		if !notClosed {
			return errors.New("event channel was closed prematurely without sending \"closed\" events earlier")
		}
		if isClosed, statistics := event.IsClosed(); isClosed {
			pipeClosed = true
			fmt.Printf("statistis: %+v", *statistics) //TODO improve
		} else if isCanceled, reason := event.IsCanceled(); isCanceled {
			logger.Infof("processing canceled due to reason: %s", reason)
		} else if isError, err, job := event.IsError(); isError {
			logger.Errorf("job %v failed with %s", job, err)
		} else if isDone, job := event.IsDone(); isDone {
			logger.Infof("processing job %v is completed", job)
		} else {
			return errors.New(fmt.Sprintf("unknown event received: %+v", event))
		}
	}
	return nil
}

func getCliTasksAndTargets(allTasks task.TaskSequence, taskMap task.TaskMap) ([]string, []string) {
	taskNames, targets, err := cli.ParseCommandLineArguments(taskMap.TaskNamesDefined())
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	require.TrueOrDieAfter(err == nil, errStr, flag.Usage)
	require.TrueOrDieAfter(len(taskNames) > 0, "no task(s) specified", flag.Usage)
	if len(targets) == 0 {
		// add default target "."
		targets = append(targets, ".")
	}
	require.True(len(targets) != 0, "no target(s) specified")
	return taskNames, targets
}

func getConfig() *ripper.AppConf {
	configFile := *configFlag
	conf := ripper.AppConf{}
	require.NotFailed(config.FromFile(&conf, configFile, map[string]string {}))

	// hardcode ignore-prefix on temp and output dirs to avoid configuration issues
	conf.TempDirectoryName = appendIgnorePrefix(conf.TempDirectoryName, conf)
	conf.OutputDirectoryName = appendIgnorePrefix(conf.OutputDirectoryName, conf)

	return &conf
}

func appendIgnorePrefix(s string, conf ripper.AppConf) string {
	if strings.HasPrefix(s, conf.IgnoreFolderPrefix) {
		return s
	} else {
		return fmt.Sprintf("%s%s", conf.IgnoreFolderPrefix, s)
	}
}
