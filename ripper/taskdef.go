package ripper

import (
	"go-cli/task"
)

func CreateTasks() (task.TaskSequence, error){
	taskTasks := task.NewTask("tasks","show all available tasks and their dependencies", task.TasksOverviewHandler )
	taskClean := task.NewTask("clean","cleans specified output folders", cleanHandler)

	taskScanAudio := task.NewTask("scanAudio","scan folder and direct sub-folders for audio input", nil)
	taskScanVideo := task.NewTask("scanVideo","scan folder and direct sub-folders for video input", nil)
	taskScan      := task.NewTask("scan","scan folder and direct sub-folders for audio and video input", nil).WithDependencies(taskScanAudio, taskScanVideo)

	taskResolveAudio := task.NewTask("resolveAudio","resolve & download audio meta-info from FreeDB", nil )
	taskResolveVideo := task.NewTask("resolveVideo","resolve & download video meta-info from IMDB", nil )
	taskResolve      := task.NewTask("resolve","resolve & download audio and video meta-info from various sources", nil).WithDependencies(taskScan, taskResolveAudio, taskResolveVideo)

	taskRipAudio := task.NewTask("ripAudio","digitalize (\"rip\") audio", nil)
	taskRipVideo := task.NewTask("ripVideo","digitalize (\"rip\") video", nil)
	taskRip      := task.NewTask("rip","digitalize (\"rip\") audio and video", nil).WithDependencies(taskResolve, taskRipAudio, taskRipVideo)

	taskTagAudio := task.NewTask("tagAudio","apply meta-info from local file to audio", nil)
	taskTagVideo := task.NewTask("tagVideo","apply meta-info from local file to video", nil)
	taskTag      := task.NewTask("tag","apply meta-info from local file to audio and video", nil).WithDependencies(taskRip, taskTagAudio, taskTagVideo)

	taskAudio  := task.NewTask("audio","process all audio files in folder and direct sub-folders", nil).WithDependencies(taskScanAudio, taskResolveAudio, taskRipAudio, taskTagAudio)
	taskVideo  := task.NewTask("video","process all video files in folder and direct sub-folders", nil).WithDependencies(taskScanVideo, taskResolveVideo, taskRipVideo, taskTagVideo)
	taskAll    := task.NewTask("all","process all audio and video files in folder and direct sub-folders", nil).WithDependencies(taskTag)

	return task.LoadTasks(
		taskTasks,
		taskClean,
		taskScanAudio, taskScanVideo, taskScan,
		taskResolveAudio, taskResolveVideo, taskResolve,
		taskRipAudio, taskRipVideo, taskRip,
		taskTagAudio, taskTagVideo, taskTag,
		taskAudio, taskVideo, taskAll)
}
