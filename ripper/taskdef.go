package ripper

import (
	"go-cli/task"
)

func CreateTasks() (tasks.TaskSequence, error){
	taskTasks := tasks.NewTask("tasks","show all available tasks and their dependencies", tasks.TasksOverviewHandler )
	taskClean := tasks.NewTask("clean","cleans all intermediate processing artifacts except contents in results folder", cleanHandler)

	taskScanAudio := tasks.NewTask("scanAudio","scan folder and direct sub-folders for audio input", nil)
	taskScanVideo := tasks.NewTask("scanVideo","scan folder and direct sub-folders for video input", nil)
	taskScan      := tasks.NewTask("scan","scan folder and direct sub-folders for audio and video input", nil).WithDependencies(taskScanAudio, taskScanVideo)

	taskResolveAudio := tasks.NewTask("resolveAudio","resolve & download audio meta-info from FreeDB", nil )
	taskResolveVideo := tasks.NewTask("resolveVideo","resolve & download video meta-info from IMDB", nil )
	taskResolve      := tasks.NewTask("resolve","resolve & download audio and video meta-info from various sources", nil).WithDependencies(taskScan, taskResolveAudio, taskResolveVideo)

	taskRipAudio := tasks.NewTask("ripAudio","digitalize (\"rip\") audio", nil)
	taskRipVideo := tasks.NewTask("ripVideo","digitalize (\"rip\") video", nil)
	taskRip      := tasks.NewTask("rip","digitalize (\"rip\") audio and video", nil).WithDependencies(taskResolve, taskRipAudio, taskRipVideo)

	taskTagAudio := tasks.NewTask("tagAudio","apply meta-info from local file to audio", nil)
	taskTagVideo := tasks.NewTask("tagVideo","apply meta-info from local file to video", nil)
	taskTag      := tasks.NewTask("tag","apply meta-info from local file to audio and video", nil).WithDependencies(taskRip, taskTagAudio, taskTagVideo)

	taskAudio  := tasks.NewTask("audio","process all audio files in folder and direct sub-folders", nil).WithDependencies(taskScanAudio, taskResolveAudio, taskRipAudio, taskTagAudio)
	taskVideo  := tasks.NewTask("video","process all video files in folder and direct sub-folders", nil).WithDependencies(taskScanVideo, taskResolveVideo, taskRipVideo, taskTagVideo)
	taskAll    := tasks.NewTask("all","process all audio and video files in folder and direct sub-folders", nil).WithDependencies(taskTag)

	return tasks.LoadTasks(
		taskTasks,
		taskClean,
		taskScanAudio, taskScanVideo, taskScan,
		taskResolveAudio, taskResolveVideo, taskResolve,
		taskRipAudio, taskRipVideo, taskRip,
		taskTagAudio, taskTagVideo, taskTag,
		taskAudio, taskVideo, taskAll)
}

