package ripper

import (
	"go-cli/pipeline"
	"fmt"
	"go-cli/task"
	"path/filepath"
	"strings"
)

const (
	JobField_Path = "path" //location of target file
)

func GetTargetFileFromJob(job task.Job) string {
	return job[JobField_Path]
}

func GetWorkPathForJob(workDir string, job task.Job) (string, error) {
	folder, _ := filepath.Split(GetTargetFileFromJob(job))
	return GetWorkPathForTargetFileFolder(workDir, folder)
}

func GetWorkPathForTargetFileFolder(workDir, folder string) (string, error) {
	targetPath, err := filepath.Abs(folder)
	if err !=  nil {
		return "", err
	}

	drive := fmt.Sprintf("%s%c", filepath.VolumeName(targetPath), filepath.Separator)
	relativeToDrive, err := filepath.Rel(drive, targetPath)
	if err != nil {
		return "", err
	}

	//driveletter will be empty string in linux
	driveletter := strings.Replace(filepath.VolumeName(targetPath), ":", "", 1)
	return filepath.Join(workDir, driveletter, relativeToDrive), nil
}

func ProcessPath(path string) pipeline.Command {
	return pipeline.Process(map[string]string {JobField_Path: path}, fmt.Sprintf("process multi-media sources at %s", path))
}