package processor

import (
	"go-ripper/files"
	"go-ripper/ripper"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"go-ripper/targetinfo"
	"go-cli/test"
	"go-cli/commons"
)


var ti targetinfo.TargetInfo = targetinfo.NewMovie("c.foo", "a/b", "id")

func TestDefaultCheckLazy(t *testing.T) {
	t.Run("recommend not lazy if lazy is off", func(t *testing.T) {
		checkLazy := DefaultCheckLazy(false, "foo")
		test.AssertOn(t).False("expected checklazy to be false")(checkLazy(ti))
	})

	t.Run("recommend not lazy if lazy is on but extension does not match", func(t *testing.T) {
		checkLazy := DefaultCheckLazy(true, "bar")
		test.AssertOn(t).False("expected checklazy to be false")(checkLazy(ti))
	})

	t.Run("recommend lazy if lazy is on and extension matches", func(t *testing.T) {
		checkLazy := DefaultCheckLazy(true, "foo")
		test.AssertOn(t).True("expected checklazy to be true")(checkLazy(ti))
	})
}

func TestNeverLazy(t *testing.T) {
	t.Run("recommend not lazy if lazy is off", func(t *testing.T) {
		checkLazy := NeverLazy(false, "procName", commons.Printf)
		test.AssertOn(t).False("expected checklazy to return false when never-lazy")(checkLazy(ti))
	})

	t.Run("recommend not lazy if lazy is on", func(t *testing.T) {
		checkLazy := NeverLazy(true, "procName", commons.Printf)
		test.AssertOn(t).False("expected checklazy to return false when never-lazy")(checkLazy(ti))
	})
}

func TestDefaultInputFile(t *testing.T) {
	const expectedFileExtension = "xyz"

	setup := func(tmpDir string, sourceExtension string, preprocessedExtension string, preprocessedExists bool) (targetinfo.TargetInfo, string, string) {
		sourceDir := "/sepp/hat/gelbe/eier"
		sourceFile := "ripped"
		source := filepath.Join(sourceDir, files.WithExtension(sourceFile, sourceExtension))
		ti := targetinfo.NewMovie(files.WithExtension(sourceFile, sourceExtension), sourceDir, sourceFile)

		workFolder, _ := ripper.GetWorkPathForTargetFolder(tmpDir, ti.GetFolder())
		files.CreateFolderStructure(workFolder)
		preprocessed := filepath.Join(workFolder, files.WithExtension(sourceFile, preprocessedExtension))
		if preprocessedExists {
			//only create preprocessed (ie ripped) inFile if an extension is defined
			ioutil.WriteFile(preprocessed, []byte {1,2,3}, os.ModePerm)
		}
		return ti, source, preprocessed
	}

	test.Run(t,"ripped inFile is available in workdir, source inFile is not a valid output inFile", func(assert *test.Assertion) {
		workDir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, workDir)
		ti, _, preprocessed := setup(workDir, "avi", expectedFileExtension, true)

		in, err := DefaultInputFileFor([]string {expectedFileExtension})(ti, workDir)
		if err != nil {
			t.Fatal(err)
		}
		assert.StringsEqual(preprocessed, in)
	})

	test.Run(t,"source inFile is already a valid output inFile", func(assert *test.Assertion) {
		workDir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, workDir)
		ti, _, preprocessed := setup(workDir, expectedFileExtension, expectedFileExtension, true)

		in, err := DefaultInputFileFor([]string {expectedFileExtension})(ti, workDir)
		if err != nil {
			t.Fatal(err)
		}
		assert.StringsEqual(preprocessed, in)
	})

	test.Run(t,"ripped inFile is missing in workdir, but source inFile is already a valid output inFile", func(assert *test.Assertion) {
		workDir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, workDir)
		ti, source, _ := setup(workDir, expectedFileExtension, expectedFileExtension, false)

		in, err := DefaultInputFileFor([]string {expectedFileExtension})(ti, workDir)
		if err != nil {
			t.Fatal(err)
		}
		assert.StringsEqual(source, in)
	})

	test.Run(t,"ripped inFile is missing in workdir, source inFile is not a valid output inFile either", func(assert *test.Assertion) {
		workDir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, workDir)
		ti, _, _ := setup(workDir, "avi", expectedFileExtension, false)

		_, err := DefaultInputFileFor([]string {expectedFileExtension})(ti, workDir)
		assert.ExpectError("expected error when finding no suitable input file - neither source does not have appropriate format, prepocessed inFile is missing")(err)
	})

}

func TestGetDefaultOutputFileFor(t *testing.T) {
	assert := test.AssertOn(t)
	const workDirBase = "/work"
	const expectedOutputExtension = "engelbert"
	const sourceDir = "/x/y/z"
	const sourceFile = "source"
	const sourceExtension = "hubert"
	ti := targetinfo.NewMovie(files.WithExtension(sourceFile, sourceExtension), sourceDir, sourceFile)

	workDir, err := ripper.GetWorkPathForTargetFolder(workDirBase, sourceDir)
	assert.NotError(err)
	outputFile, err := DefaultOutputFileFor(expectedOutputExtension)(ti, workDirBase)
	assert.NotError(err)
	assert.StringsEqual(filepath.Join(workDir, files.WithExtension(sourceFile, expectedOutputExtension)), outputFile)
}