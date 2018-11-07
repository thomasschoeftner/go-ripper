package rip

import (
	"testing"
	"go-ripper/targetinfo"
	"path/filepath"
	"go-ripper/files"
	"go-ripper/ripper"
	"io/ioutil"
	"os"
	"go-cli/test"
)

const expectedVideoExtension = "mp4"

func TestFindInputOutputFile(t *testing.T) {
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

	t.Run("ripped inFile is available in workdir, source inFile is not a valid output inFile", func(t *testing.T) {
		workDir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, workDir)
		ti, _, preprocessed := setup(workDir, "avi", expectedVideoExtension, true)

		in, out, err := findInputOutputFiles(ti, workDir, expectedVideoExtension)
		if err != nil {
			t.Fatal(err)
		}
		test.AssertOn(t).StringsEqual(preprocessed, in)
		test.AssertOn(t).StringsEqual(preprocessed, out)
	})

	t.Run("ripped inFile is available in workdir, but source inFile is already a valid output inFile", func(t *testing.T) {
		workDir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, workDir)
		ti, _, preprocessed := setup(workDir, expectedVideoExtension, expectedVideoExtension, true)

		in, out, err := findInputOutputFiles(ti, workDir, expectedVideoExtension)
		if err != nil {
			t.Fatal(err)
		}
		test.AssertOn(t).StringsEqual(preprocessed, in)
		test.AssertOn(t).StringsEqual(preprocessed, out)
	})

	t.Run("ripped inFile is missing in workdir, but source inFile is already a valid output inFile", func(t *testing.T) {
		workDir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, workDir)
		ti, source, preprocessed := setup(workDir, expectedVideoExtension, expectedVideoExtension, false)

		in, out, err := findInputOutputFiles(ti, workDir, expectedVideoExtension)
		if err != nil {
			t.Fatal(err)
		}

		test.AssertOn(t).StringsEqual(source, in)
		test.AssertOn(t).StringsEqual(preprocessed, out)
	})

	t.Run("ripped inFile is missing in workdir, source inFile is not a valid output inFile either", func(t *testing.T) {
		workDir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, workDir)
		ti, _, _ := setup(workDir, "avi", expectedVideoExtension, false)

		_, _, err := findInputOutputFiles(ti, workDir, expectedVideoExtension)
		test.AssertOn(t).ExpectError("expected error when finding no suitable inFile - neither source does not have appropriate format, prepocessed inFile is missing")(err)
	})
}

