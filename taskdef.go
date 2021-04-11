package main

import (
	"github.com/thomasschoeftner/go-cli/task"
	"github.com/thomasschoeftner/go-ripper/clean"
	"github.com/thomasschoeftner/go-ripper/scan"
	"errors"
	"github.com/thomasschoeftner/go-ripper/metainfo/video"
	"github.com/thomasschoeftner/go-ripper/tag"
	"github.com/thomasschoeftner/go-ripper/rip"
)

const TaskName_Tasks = "tasks"

func NotImplementedYetHandler(ctx task.Context) task.HandlerFunc {
	return func (job task.Job) ([]task.Job, error) {
		return nil, errors.New("not implemented yet")
	}
}

func CreateTasks() task.TaskSequence {
	taskTasks := task.NewTask(TaskName_Tasks,"show all available tasks and their dependencies", task.TasksOverviewHandler )

	//taskScanAudio := task.NewTask("scanAudio","scan folder and direct sub-folders for audio input", NotImplementedYetHandler)
	taskScanVideo := task.NewTask("scanVideo","scan folder and direct sub-folders for video input", scan.ScanVideo)
	taskScan      := task.NewTask("scan","scan folder and direct sub-folders for audio and video input", nil).WithDependencies(/*taskScanAudio,*/ taskScanVideo)

	//taskResolveAudio := task.NewTask("resolveAudio","resolve & download audio meta-info from FreeDB", NotImplementedYetHandler )
	taskResolveVideo := task.NewTask("resolveVideo","resolve & download video meta-info from IMDB", video.ResolveVideo)
	taskResolve      := task.NewTask("resolve","resolve & download audio and video meta-info from various sources", nil).WithDependencies(taskScan, /*taskResolveAudio, */ taskResolveVideo)

	//taskRipAudio := task.NewTask("ripAudio","digitalize (\"rip\") audio", NotImplementedYetHandler)
	taskRipVideo := task.NewTask("ripVideo","digitalize (\"rip\") video", rip.RipVideo)
	taskRip      := task.NewTask("rip","digitalize (\"rip\") audio and video", nil).WithDependencies(taskResolve, /*taskRipAudio, */ taskRipVideo)

	//taskTagAudio := task.NewTask("tagAudio","apply meta-info from local file to audio", NotImplementedYetHandler)
	taskTagVideo := task.NewTask("tagVideo","apply meta-info from local file to video", tag.TagVideo)
	taskTag      := task.NewTask("tag","apply meta-info from local file to audio and video", nil).WithDependencies(taskRip, /* taskTagAudio, */ taskTagVideo)

	//taskAudio  := task.NewTask("audio","process all audio files in folder and direct sub-folders", nil).WithDependencies(taskScanAudio, taskResolveAudio, taskRipAudio, taskTagAudio, taskRemoveOriginalAudio )
	taskVideo  := task.NewTask("video","process all video files in folder and direct sub-folders", nil).WithDependencies(taskScanVideo, taskResolveVideo, taskRipVideo, taskTagVideo)

	taskClean := task.NewTask("clean","cleans all processing artifacts related to a specific input file from work folder", clean.CleanHandler)
	taskRemoveOriginal := task.NewTask("removeOriginal", "deletes original input file", NotImplementedYetHandler)


	return task.LoadTasks(
		taskTasks,
		/* taskScanAudio, */ taskScanVideo, taskScan,
		/* taskResolveAudio, */ taskResolveVideo, taskResolve,
		/* taskRipAudio, */ taskRipVideo, taskRip,
		/* taskTagAudio, */ taskTagVideo, taskTag,
		taskClean, taskRemoveOriginal,
		/* taskAudio, */ taskVideo)
}
