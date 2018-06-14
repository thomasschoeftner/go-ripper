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
	filePart := "x.y"
	pathPart := "testdata/a/b/c/" + filePart

	workDir := ".workdir"
	targetPath := drive + pathPart
	job := task.Job{JobField_Path : targetPath}

	expectedWorkPath := filepath.Join(workDir, strings.Replace(drive, ":", "", 1), filepath.Dir(pathPart))
	assert := test.AssertOn(t)
	workPath := assert.StringNotError(GetWorkPathForJob(workDir, job))
	assert.StringsEqual(expectedWorkPath, workPath)
}
