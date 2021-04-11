package ripper

import (
	"testing"
	"path/filepath"
	"github.com/thomasschoeftner/go-cli/task"
	"strings"
	"github.com/thomasschoeftner/go-cli/test"
	"github.com/thomasschoeftner/go-ripper/files"
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

func TestGetProcessingArtifactPath(t *testing.T) {
	var drive string
	if filepath.Separator == '/' {
		drive = "/"
	} else {
		drive = "c:/"
	}

	targetDir := drive + "my/private/videos"
	targetFile := "cut7.mov"

	workDir := "/my/work/dir"
	requiredExtension := "xyz"

	fName, _ := files.SplitExtension(targetFile)
	expectedPath := filepath.Join(workDir, strings.Replace(targetDir, ":", "", 1), files.WithExtension(fName, requiredExtension))

	artifactPath, err := GetProcessingArtifactPathFor(workDir, targetDir, targetFile, requiredExtension)
	assert := test.AssertOn(t)
	assert.NotError(err)
	assert.StringsEqual(expectedPath, artifactPath)
}