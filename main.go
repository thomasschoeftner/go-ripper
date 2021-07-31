package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/logger"
	"github.com/thomasschoeftner/go-cli/cli"
	"github.com/thomasschoeftner/go-cli/pipeline"
	"github.com/thomasschoeftner/go-cli/require"
	"github.com/thomasschoeftner/go-cli/task"
	"github.com/thomasschoeftner/go-ripper/files"
	"github.com/thomasschoeftner/go-ripper/metainfo/video"
	"github.com/thomasschoeftner/go-ripper/omdb"
	"github.com/thomasschoeftner/go-ripper/ripper"
)

const (
	cliFlagVerbose    = "verbose"
	cliFlagLazy       = "lazy"
	cliFlagConfigFile = "config"
	ApplicationName   = "go-ripper"
)

var isVerbose = cli.FromFlag(cliFlagVerbose, "full log output in console").GetBoolean().WithDefault(false)
var isLazy = cli.FromFlag(cliFlagLazy, "avoid re-execution of task, if output from previous execution is available - defaults to true").GetBoolean().WithDefault(true)
var configFile = cli.FromFlag(cliFlagConfigFile, "the config file location").OrEnvironmentVar(ApplicationName + "-" + cliFlagConfigFile).GetString().WithDefault("/" + ApplicationName + "/config/" + ApplicationName + ".conf")

func main() {
	os.Exit(launch())
}

func launch() int {

	syntax := "[flags] task[...] target[...]"
	cli.Setup(&syntax, "use task \"tasks\" to display all available tasks:", fmt.Sprintf("  %s tasks\n", os.Args[0]))
	logger.Init(ApplicationName, *isVerbose, false, ioutil.Discard)

	// read config
	conf := ripper.GetConfig(*configFile)
	require.NotFailed(files.CreateFolderStructure(conf.OutputDirectory))

	switch conf.Resolve.Video.Resolver {
	case omdb.CONF_OMDB_RESOLVER:
		video.NewVideoMetaInfoSource = omdb.NewOmdbVideoMetaInfoSource
	default:
		logger.Fatalf("unknown video resolver configured: %s", conf.Resolve.Video.Resolver)
	}

	// create task Tree
	allTasks := CreateTasks()
	taskMap, err := task.ValidateTasks(allTasks)
	require.NotFailed(err)

	// read command line params (flags & args)
	taskNames, targets := getCliTasksAndTargets(taskMap)
	//TODO add validation not to use workDir or repoDir as target folder

	// calculate tasks to be invoked
	tasksToRun, err := taskMap.CompileTasksForNamesCompact(taskNames...)
	require.NotFailed(err)

	// materialize processing pipeline
	pipe, err := pipeline.Materialize(tasksToRun).WithConfig(conf.Processing, conf, allTasks, *isLazy)
	require.NotFailed(err)

	// ASYNCHRONOUSLY send a processing command for each target to pipeline
	go fillPipelineAndClose(pipe, targets)

	err = handleProcessingEvents(pipe)
	require.NotFailed(err)

	return 0
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
			fmt.Printf("statistics: %+v\n", *statistics) //TODO improve
		} else if isCanceled, reason := event.IsCanceled(); isCanceled {
			logger.Infof("processing canceled due to reason: %s\n", reason)
		} else if isError, err, job := event.IsError(); isError {
			logger.Errorf("job %v failed with %s\n", job, err)
		} else if isDone, job := event.IsDone(); isDone {
			logger.Infof("processing job %v is completed\n", job)
		} else {
			return fmt.Errorf("unknown event received: %+v\n", event)
		}
	}
	return nil
}

func getCliTasksAndTargets(taskMap task.TaskMap) ([]string, []string) {
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

	absoluteTargetPaths := make([]string, 0, len(targets))
	for _, t := range targets {
		abs, err := filepath.Abs(t)
		require.NotFailed(err)
		absoluteTargetPaths = append(absoluteTargetPaths, abs)
	}
	return taskNames, absoluteTargetPaths
}
