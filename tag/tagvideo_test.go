package tag

import (
	"testing"
	"go-ripper/targetinfo"
	"go-ripper/ripper"
	"go-cli/test"
	"path/filepath"
	"io/ioutil"
	"os"
	"go-ripper/files"
)

func TestChooseInputFile(t *testing.T) {
	const expectedExtension = "mp4"

	setup := func(tmpDir string, sourceExtension string, preprocessedExtension *string) (targetinfo.TargetInfo, string, string) {
		sourceDir := "/sepp/hat/gelbe/eier"
		sourceFile := "ripped"
		source := filepath.Join(sourceDir, sourceFile + "." + sourceExtension)

		ti := targetinfo.NewMovie(sourceFile + "." + sourceExtension, sourceDir, sourceFile)
		tiPath, _ := ripper.GetWorkPathForTargetFileFolder(tmpDir, ti.GetFolder())
		files.CreateFolderStructure(tiPath)
		var preprocessed string
		if preprocessedExtension != nil {
			//only create preprocessed (ie ripped) file if an extension is defined
			preprocessed = filepath.Join(tiPath, sourceFile + "." + *preprocessedExtension )
			ioutil.WriteFile(preprocessed, []byte {1,2,3}, os.ModePerm)
		}
		return ti, source, preprocessed
	}

	t.Run("ripped file is available in workdir, source file is not a valid output file", func(t *testing.T) {
		workDir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, workDir)
		preprocessedExt := expectedExtension
		ti, _, preprocessed := setup(workDir, "avi", &preprocessedExt)

		chosen, err := chooseInputFile(ti, workDir, expectedExtension)
		if err != nil {
			t.Fatal(err)
		}
		test.AssertOn(t).StringsEqual(preprocessed, chosen)
	})

	t.Run("ripped file is available in workdir, but source file is already a valid output file", func(t *testing.T) {
		workDir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, workDir)
		preprocessedExt := expectedExtension
		ti, _, preprocessed := setup(workDir, expectedExtension, &preprocessedExt)

		chosen, err := chooseInputFile(ti, workDir, expectedExtension)
		if err != nil {
			t.Fatal(err)
		}
		test.AssertOn(t).StringsEqual(preprocessed, chosen)
	})

	t.Run("ripped file is missing in workdir, but source file is already a valid output file", func(t *testing.T) {
		workDir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, workDir)
		ti, src, _ := setup(workDir, expectedExtension, nil)

		chosen, err := chooseInputFile(ti, workDir, expectedExtension)
		if err != nil {
			t.Fatal(err)
		}
		test.AssertOn(t).StringsEqual(src, chosen)
	})

	t.Run("ripped file is missing in workdir, source file is not a valid output file either", func(t *testing.T) {
		workDir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, workDir)
		ti, _, _ := setup(workDir, "avi", nil)

		_, err := chooseInputFile(ti, workDir, expectedExtension)
		test.AssertOn(t).ExpectError("expected error when finding no suitable file - neither source does not have appropriate format, prepocessed file is missing")(err)
	})
}
