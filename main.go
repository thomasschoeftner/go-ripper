package main

import (
	"os"
	"github.com/google/logger"
	"io/ioutil"
	"go-cli/config"
	"go-cli/cli"
	"go-cli/commons"
	"go-cli/task"
	"go-ripper/ripper"
	"go-cli/pipeline"
	"flag"
	"fmt"
)

const (
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
	allTasks, error := ripper.CreateTasks()
	commons.Check(error)
	taskMap, errs := task.ValidateTasks(allTasks)
	commons.CheckMultiple(errs)

	// read command line params (flags & args)
	syntax := "[flags] task[...] target[...]"
	cli.Setup(&syntax, allTasks)
	taskNames, targets, error := cli.ParseCommandLineArguments(taskMap.TaskNamesDefined())
	if error !=  nil {
		flag.Usage()
	}
	commons.Check(error)

	// read config
	configFile := *configFlag
	omdbTokens := *omdbTokenFlags
	conf := ripper.AppConf{}
	error = config.FromFile(&conf, configFile, map[string]string {})
	commons.Check(error)
	conf.Omdb.OmdbTokens = omdbTokens

	// run "tasks" by default if no other task is specified
	if len(taskNames) == 0 {
		taskNames = append(taskNames, "tasks")
	}

	// calculate tasks to be invoked
	invokedTasks, error := taskMap.GetTasksForNames(taskNames...)
	commons.Check(error)


	//{   // TODO remove
	//	for idx, t := range invokedTasks {
	//		logger.Infof("%d --- %s", idx, t.Name)
	//		results := t.Handler(task.Context{allTasks, conf, commons.Printf}, task.Process(task.Param{"folder", "franz"}))
	//		for _, r := range results {
	//			if r.Error != nil {
	//				logger.Error(r.Error)
	//			}
	//		}
	//	}
	//	logger.Infof("targets: %s", targets)
	//}

	// materialize pipelines
	pipe, error := pipeline.Materialize(invokedTasks).WithConfig(conf.Processing, conf, allTasks)
	commons.Check(error)

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
			logger.Errorf("job %v failed with error %s", job, err)
		} else if isDone, job := event.IsDone(); isDone {
			logger.Infof("job %v is completed", job)
		} else {
			logger.Fatalf(fmt.Sprintf("unknown event received: %+v", event))
		}
	}

	//todo check required flags & target per task!!!

	return 0
}

//type CheckResult struct {
//	Success bool
//	Event pipeline.Event
//
//}
//type Check func(event pipeline.Event) CheckResult
//
//func (c CheckResult) OrElse(check Check) CheckResult {
//	if !c.Success {
//		return check(c.Event)
//	}
//}
//
//func handleClosed(event pipeline.Event, pipelineClosed *bool) CheckResult {
//
//}
