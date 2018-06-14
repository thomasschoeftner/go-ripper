package clean

import (
	"testing"
	"go-cli/task"
	"go-ripper/ripper"
	"os"
	"io/ioutil"
	"path/filepath"
	"go-cli/test"
	"go-ripper/files"
	"go-cli/commons"
)

func TestClean(t *testing.T) {
	assert := test.AssertOn(t)
	workDir := test.MkTempFolder(t)
	defer test.RmTempFolder(t, workDir)

	targetFolder := "a/b"
	targetFile := "to.remove"
	targetPath := filepath.Join(targetFolder, targetFile)
	job := task.Job{ripper.JobField_Path : targetPath}

	//create target files
	workPath := assert.StringNotError(ripper.GetWorkPathFor(workDir, job))
	assert.NotError(files.CreateFolderStructure(workPath))
	assert.NotError(files.CreateFolderStructure(filepath.Join(workPath, targetFile + ".abc"))) //create sub-folder with "dangerous" name

	tFile := assert.StringNotError(createFile(filepath.Join(workPath, targetFile), []byte{}))
	similarFile := assert.StringNotError(createFile(filepath.Join(workPath, targetFile + ".abc", targetFile), []byte{}))
	jsonFile  := assert.StringNotError(createFile(filepath.Join(workPath, targetFile + ".json"), []byte{}))
	otherFile := assert.StringNotError(createFile(filepath.Join(workPath, "something.else"), []byte{}))

	assert.TrueNotError("target file was not created")(fileExists(tFile))
	assert.TrueNotError("similar file in sub-folder was not created")(fileExists(similarFile))
	assert.TrueNotError("json file was not created")(fileExists(jsonFile))
	assert.TrueNotError("unrelated file was not created")(fileExists(otherFile))

	clean(commons.Printf, "test clean", job,  workDir)

	assert.FalseNotError("target file was not deleted")(fileExists(targetFile))
	assert.FalseNotError("json file was not deleted")(fileExists(jsonFile))
	assert.TrueNotError("similar file in sub-folder was deleted")(fileExists(similarFile))
	assert.TrueNotError("unrelated file was deleted")(fileExists(otherFile))
}

func createFile(path string, content []byte) (string, error) {
	return path, ioutil.WriteFile(path, content, os.ModePerm)
}

func fileExists(path string) (bool, error) {
	exists, _, err := files.Exists(path)
	return exists, err
}
