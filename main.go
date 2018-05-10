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
)

const (
	omdbTokenFlag  = "omdbtoken"
	configFileFlag = "config"
	ApplicationName = "go-ripper"
)


var omdbTokenFlags  = cli.FromFlag(omdbTokenFlag, "the access token für connecting to OMDB - can also be specified as ENV variable").OrEnvironmentVar(omdbTokenFlag).GetArray().WithDefault()
var configFlag = cli.FromFlag(configFileFlag,"the config file location").OrEnvironmentVar(ApplicationName + "-" + configFileFlag).GetString().WithDefault(ApplicationName + ".conf")

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

	// run "tasks" by default if no other task is specified
	if len(taskNames) == 0 {
		taskNames = append(taskNames, "tasks")
	}

	// calculate tasks to be invoked
	invokedTasks, error := taskMap.GetTasksForNames(taskNames...)
	commons.Check(error)


	// TODO remove
	for idx, t := range invokedTasks {
		logger.Infof("%d --- %s", idx, t.Name)
		results := t.Handler(task.Context{allTasks, conf, commons.Printf}, task.Process(task.Param{"folder", "franz"}))
		for _, r := range results {
			if r.Error != nil {
				logger.Error(r.Error)
			}
		}
	}
	logger.Infof("omdbtokens: %v", omdbTokens)
	logger.Infof("targets: %s", targets)

	// materialize pipelines
	pipeline , error := pipeline.Materialize(invokedTasks).WithConfig(conf.Processing, conf)
	commons.Check(error)
	if pipeline != nil {
		//TODO
	}
	//todo check required flags & target per task!!!

	return 0
}
