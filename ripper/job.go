package ripper

import (
	"go-cli/pipeline"
	"fmt"
	"go-cli/task"
	"os"
)

const (
	JobField_Location = "location"
)

func getPathFor(job task.Job, subDir string) string {
	folderPath := job[JobField_Location]
	return fmt.Sprintf("%s%c%s", folderPath, os.PathSeparator, subDir)
}

func GetTempPathFor(job task.Job, conf *AppConf) string {
	return getPathFor(job, conf.TempDirectoryName)
}

func GetOutputPathFor(job task.Job, conf *AppConf) string {
	return getPathFor(job, conf.OutputDirectoryName)
}


func ProcessPath(path string) pipeline.Command {
	return pipeline.Process(map[string]string {JobField_Location : path}, fmt.Sprintf("process multi-media sources at %s", path))
}