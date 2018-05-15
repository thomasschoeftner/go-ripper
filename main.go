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
)

const (
	exit_success = 0
	exit_failure = 1

	cliFlagVerbose             = "verbose"
	cliFlagOmdbTokens          = "omdbtoken"
	cliFlagConfigFile          = "config"
	cliFlagWithoutDependencies = "without-dependencies"

	ApplicationName            = "go-ripper"
)



var verbose = cli.FromFlag(cliFlagVerbose, "full log output in console").GetBoolean().WithDefault(false)
var omdbTokenFlags  = cli.FromFlag(cliFlagOmdbTokens, "the access token für connecting to OMDB - can also be specified as ENV variable").OrEnvironmentVar(cliFlagOmdbTokens).GetArray().WithDefault()
var configFlag = cli.FromFlag(cliFlagConfigFile,"the config file location").OrEnvironmentVar(ApplicationName + "-" + cliFlagConfigFile).GetString().WithDefault(ApplicationName + ".conf")
var runWithoutDependencies = cli.FromFlag(cliFlagWithoutDependencies, "run specified tasks without their dependencies").GetBoolean().WithDefault(false)

func main() {
	os.Exit(launch())
}

func launch() int {
	syntax := "[flags] task[...] target[...]"
	allTasks  := CreateTasks()
	cli.Setup(&syntax, allTasks)

	logger.Init(ApplicationName, *verbose, false, ioutil.Discard)

	// create task Tree
	taskMap, errs := task.ValidateTasks(allTasks)
	require.NoneFailed(errs)

	// read command line params (flags & args)
	taskNames, targets := getCliTasksAndTargets(allTasks, taskMap)

	// read config
	conf := getConfig()

	// calculate tasks to be invoked
	invokedTasks, err := taskMap.GetTasksForNames(taskNames...)
	require.NotFailed(err)

	// materialize processing pipeline
	//todo check required flags & target per task!!!
	pipe, err := pipeline.Materialize(invokedTasks.Flatten()).WithConfig(conf.Processing, conf, allTasks)
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
			logger.Errorf("job %v failed with err %s", job, err)
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
	omdbTokens := *omdbTokenFlags
	conf := ripper.AppConf{}
	require.NotFailed(config.FromFile(&conf, configFile, map[string]string {}))
	conf.Omdb.OmdbTokens = omdbTokens
	return &conf
}