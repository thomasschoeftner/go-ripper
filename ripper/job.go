package ripper

import (
	"github.com/thomasschoeftner/go-cli/pipeline"
	"fmt"
	"github.com/thomasschoeftner/go-cli/task"
	"path/filepath"
	"strings"
	"github.com/thomasschoeftner/go-ripper/files"
)

const (
	JobField_Path = "path" //location of target file
)

func GetTargetFileFromJob(job task.Job) string {
	return job[JobField_Path]
}

func GetWorkPathForJob(workDir string, job task.Job) (string, error) {
	folder, _ := filepath.Split(GetTargetFileFromJob(job))
	return GetWorkPathForTargetFolder(workDir, folder)
}

func GetProcessingArtifactPathFor(workDir string, targetDir string, targetFile string, expectedExtension string) (string, error) {
	dir, err := GetWorkPathForTargetFolder(workDir, targetDir)
	if err != nil {
		return "", err
	}

	fname, _ := files.SplitExtension(targetFile)
	preprocessedArtifact := filepath.Join(dir, files.WithExtension(fname, expectedExtension))
	return preprocessedArtifact, nil
}



func GetWorkPathForTargetFolder(workDir, targetFolder string) (string, error) {
	targetPath, err := filepath.Abs(targetFolder)
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