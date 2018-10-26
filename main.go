package main

import (
	"os"
	"github.com/google/logger"
	"go-cli/config"
	"go-cli/task"
	"go-ripper/ripper"
	"go-cli/pipeline"
	"flag"
	"fmt"
	"go-cli/require"
	"errors"
	"go-cli/cli"
	"path/filepath"
	"go-ripper/metainfo/video"
	"go-ripper/omdb"
	"io/ioutil"
	"go-ripper/tag"
	"strings"
)

const (
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

func main() {
	os.Exit(launch())
}

func launch() int {
	syntax := "[flags] task[...] target[...]"
	cli.Setup(&syntax, "use task \"tasks\" to display all available tasks:", fmt.Sprintf("  %s tasks\n", os.Args[0]))
	logger.Init(ApplicationName, *verbose, false, ioutil.Discard)

	// read config
	conf := getConfig()
	conf.Resolve.Video.Omdb.OmdbTokens = *omdbTokenFlags

	switch conf.Resolve.Video.Resolver {
	case omdb.CONF_OMDB_RESOLVER:
		video.NewVideoMetaInfoSource = omdb.NewOmdbVideoMetaInfoSource
	default:
		logger.Fatalf("unknown video resolver configured: %s", conf.Resolve.Video.Resolver)
	}

	switch conf.Tag.Video.Tagger {
	case tag.CONF_ATOMICPARSLEY_TAGGER:
		tag.NewVideoTagger = tag.NewAtomicParsleyVideoTagger
	default:
		logger.Fatalf("unknown video tagger configured: %s", conf.Tag.Video.Tagger)
	}

	// create task Tree
	allTasks  := CreateTasks()
	taskMap, err := task.ValidateTasks(allTasks)
	require.NotFailed(err)

	// read command line params (flags & args)
	taskNames, targets := getCliTasksAndTargets(taskMap)

	// calculate tasks to be invoked
	tasksToRun, err := taskMap.CompileTasksForNamesCompact(taskNames...)
	require.NotFailed(err)

	// materialize processing pipeline
	pipe, err := pipeline.Materialize(tasksToRun).WithConfig(conf.Processing, conf, allTasks, *lazy)
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

func getConfig() *ripper.AppConf {
	configFile := *configFlag
	conf := ripper.AppConf{}
	require.NotFailed(config.FromFile(&conf, configFile, map[string]string {}))
	require.NotFailed(validateConfig(&conf))
	return &conf
}

func validateConfig(c *ripper.AppConf) error {
	if c == nil {
		return fmt.Errorf("config is nil - no config available")
	}

	validatePath:= func(path string, fieldName string) error {
		if 0 == len(path) {
			return fmt.Errorf("[config error] \"%s\" is empty", fieldName)
		}
		if strings.ContainsRune(path, ' ') {
			return fmt.Errorf("[config error] \"%s\" must not contain spaces", fieldName)
		}
		return nil
	}

	c.WorkDirectory = strings.Trim(c.WorkDirectory, " ")
	if err := validatePath(c.WorkDirectory, "workDirectory"); err != nil {
		return err
	}

	//validate metainforepo
	c.MetaInfoRepo = strings.Trim(c.MetaInfoRepo, " ")
	if err := validatePath(c.MetaInfoRepo, "metaInfoRepo"); err != nil {
		return err
	}

	return nil
}