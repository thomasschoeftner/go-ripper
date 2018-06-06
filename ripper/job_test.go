package ripper

import (
	"testing"
	"path/filepath"
	"go-cli/task"
	"strings"
	"go-cli/test"
)

func TestGetWorkPathForFile(t *testing.T) {
	var drive string
	if filepath.Separator == '/' {
		drive = "/"
	} else {
		drive = "c:/"
	}
	pathPart := "testdata/a/b/c/x.y"

	workDir := ".workdir"
	targetPath := drive + pathPart
	job := task.Job{JobField_Path : targetPath}

	expectedWorkPath := filepath.Join(workDir, strings.Replace(drive, ":", "", 1), filepath.Dir(pathPart))
	workPath, err := GetWorkPathFor(workDir, job)
	test.CheckError(t, err)
	if expectedWorkPath != workPath {
		t.Errorf("expected work path\n  %s\n but got\n  %s", expectedWorkPath, workPath)
	}
}
