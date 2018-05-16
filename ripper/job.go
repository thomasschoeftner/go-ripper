package ripper

import (
	"go-cli/pipeline"
	"fmt"
	"go-cli/task"
	"path/filepath"
)

const (
	JobField_Path = "path"
)

//func GetTempPathFor(job task.Job, conf *AppConf) string {
//	return getWorkPathFor(job, conf.TempDirectoryName)
//}
//
//func GetOutputPathFor(job task.Job, conf *AppConf) string {
//	return getWorkPathFor(job, conf.OutputDirectoryName)
//}

func GetWorkPathFor(job task.Job, subdir string) string {
	folder, _ := filepath.Split(job[JobField_Path])
	return filepath.Join(folder, subdir)
}

func ProcessPath(path string) pipeline.Command {
	return pipeline.Process(map[string]string {JobField_Path: path}, fmt.Sprintf("process multi-media sources at %s", path))
}