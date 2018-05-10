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
	"fmt"
)

const (
	omdbTokenFlag  = "omdbtoken"
	configFileFlag = "config"
	ApplicationName = "go-ripper"
)


var omdbFlags = cli.Array(cli.FromFlag(omdbTokenFlag, "the access token für connecting to OMDB - can also be specified as ENV variable").OrEnvironmentVar(omdbTokenFlag)).WithDefault("invalid", "test", "ids!")
var configFlag = cli.String(cli.FromFlag(configFileFlag, "the config file location").OrEnvironmentVar(ApplicationName + "-" + configFileFlag)).WithDefault(ApplicationName + ".conf")

func main() {
	os.Exit(launch())
}


func launch() int {
	logger.Init(ApplicationName, true, false, ioutil.Discard)

	println("before" , *omdbFlags, *configFlag)
	// read command line params (flags & args)
	flags, taskNames := cli.GetFlagsAndTasks()
	println("after " , *omdbFlags, *configFlag)
	error, _:= cli.CheckRequiredFlags(omdbTokenFlag)
	commons.Check(error)
	cli.PrintFlagsAndArgs(logger.Infof)
	fmt.Printf("omdboken is: \"%s\"", *omdbFlags)

	// read config
	configFile := flags[configFileFlag]
	omdbToken := flags[omdbTokenFlag]
	conf := ripper.AppConf{}
	commons.Check(config.FromFile(&conf, configFile, map[string]string {"omdbtoken" : omdbToken}))

	// create task Tree
	allTasks, error := ripper.CreateTasks()
	commons.Check(error)
	taskMap, errs := task.ValidateTasks(allTasks)
	commons.CheckMultiple(errs)

	// run "tasks" by default if no other task is specified
	if len(taskNames) == 0 {
		taskNames = append(taskNames, "tasks")
	}

	// calculate tasks to be invoked
	invokedTasks, error := taskMap.GetTasksForNames(taskNames...)
	commons.Check(error)

	// materialize pipelines
	pipeline , error := pipeline.Materialize(invokedTasks).WithConfig(conf.Processing, conf)
	commons.Check(error)
	if pipeline != nil {
		//TODO
	}

	// TODO remove following
	for idx, t := range invokedTasks {
		logger.Infof("%d --- %s", idx, t.Name)
		results := t.Handler(task.Context{allTasks, conf, commons.Printf}, task.Process(task.Param{"folder", "franz"}))
		for _, r := range results {
			if r.Error != nil {
				logger.Error(r.Error)
			}
		}
		//TODO
	}

	return 0
}
