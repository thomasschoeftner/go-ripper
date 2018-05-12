package main

import (
	"os"
	"github.com/google/logger"
	"io/ioutil"
	"go-cli/config"
	"go-cli/cli"
	"go-cli/task"
	"go-ripper/ripper"
	"go-cli/pipeline"
	"flag"
	"fmt"
	"go-cli/require"
)

const (
	exit_success = 0
	exit_failure = 1

	cliFlagOmdbTokens          = "omdbtoken"
	cliFlagConfigFile          = "config"
	cliFlagWithoutDependencies = "without-dependencies"

	ApplicationName            = "go-ripper"
)


var omdbTokenFlags  = cli.FromFlag(cliFlagOmdbTokens, "the access token für connecting to OMDB - can also be specified as ENV variable").OrEnvironmentVar(cliFlagOmdbTokens).GetArray().WithDefault()
var configFlag = cli.FromFlag(cliFlagConfigFile,"the config file location").OrEnvironmentVar(ApplicationName + "-" + cliFlagConfigFile).GetString().WithDefault(ApplicationName + ".conf")
var runWithoutDependencies = cli.FromFlag(cliFlagWithoutDependencies, "run specified tasks without their dependencies").GetBoolean().WithDefault(false)

func main() {
	os.Exit(launch())
}

func launch() int {
	logger.Init(ApplicationName, true, false, ioutil.Discard)

	// create task Tree
	allTasks, err := ripper.CreateTasks()
	require.NotFailed(err)
	taskMap, errs := task.ValidateTasks(allTasks)
	require.NoneFailed(errs)

	// read command line params (flags & args)
	syntax := "[flags] task[...] target[...]"
	cli.Setup(&syntax, allTasks)
	taskNames, targets, err := cli.ParseCommandLineArguments(taskMap.TaskNamesDefined())
	require.TrueOrDieAfter(err == nil, "", flag.Usage)
	require.TrueOrDieAfter(len(taskNames) > 0, "no task(s) specified", flag.Usage)
	if len(targets) == 0 {
		// add default target "."
		targets = append(targets, ".")
	}
	require.True(len(targets) != 0, "no target(s) specified")


	// read config
	configFile := *configFlag
	omdbTokens := *omdbTokenFlags
	conf := ripper.AppConf{}
	err = config.FromFile(&conf, configFile, map[string]string {})
	require.NotFailed(err)
	conf.Omdb.OmdbTokens = omdbTokens

	// calculate tasks to be invoked
	invokedTasks, err := taskMap.GetTasksForNames(taskNames...)
	require.NotFailed(err)

	// materialize pipelines
	pipe, err := pipeline.Materialize(invokedTasks).WithConfig(conf.Processing, conf, allTasks)
	require.NotFailed(err)

	go func() {
		// feed work to pipeline asynchronously
		for _, target := range targets {
			pipe.Commands <- ripper.ProcessPath(target)
		}
	}()

	pipeClosed := false
	for !pipeClosed {
		event := <-pipe.Events
		if isClosed, statistics := event.IsClosed(); isClosed {
			pipeClosed = true
			fmt.Printf("statistis: %+v", *statistics)
		} else if isCanceled, reason := event.IsCanceled(); isCanceled {
			logger.Infof("processing canceled due to reason: %s", reason)
		} else if isError, err, job := event.IsError(); isError {
			logger.Errorf("job %v failed with err %s", job, err)
		} else if isDone, job := event.IsDone(); isDone {
			logger.Infof("job %v is completed", job)
		} else {
			logger.Fatalf(fmt.Sprintf("unknown event received: %+v", event))
		}
	}

	//todo check required flags & target per task!!!

	return exit_success
}
