package files

import (
	"testing"
	"go-cli/test"
	"path/filepath"
	"strings"
	"fmt"
)


func DummyEvac(from, to string) (*evacuated, error) {
	return &evacuated{original: from, evacuatedTo: to}, nil
}

func TestPrepareEvacuation(t *testing.T) {
	dir := test.MkTempFolder(t)
	defer test.RmTempFolder(t, dir)

	t.Run("check for empty temp folder", func(t *testing.T) {
		assert := test.AssertOn(t)
		evacuate := PrepareEvacuation("", nil)
		evacuated, err := evacuate("./testdata/small").By(DummyEvac)
		assert.ExpectError("expected error due to empty temp folder name")(err)
		assert.True("expected evacuated to be nil after finding empty temp folder name")(nil == evacuated)
	})

	t.Run("create temp folder if not existing", func(t *testing.T) {
		assert := test.AssertOn(t)
		tmp := filepath.Join(dir, "temp1")
		_, err := PrepareEvacuation(tmp, nil).Of("./testdata/small").By(DummyEvac)
		assert.NotError(err)
		assert.TrueNotError("temp folder was not created")(Exists(tmp))
	})

	t.Run("re-use existing temp folder", func(t *testing.T) {
		assert := test.AssertOn(t)
		tmp := filepath.Join(dir, "temp2")
		assert.NotError(CreateFolderStructure(tmp))

		evacuated, err := PrepareEvacuation(tmp, nil).Of("./testdata/small").By(DummyEvac)
		assert.NotError(err)
		assert.True("expected temp folder to be re-used")(strings.HasPrefix(evacuated.evacuatedTo, tmp))
	})

	t.Run("detect missing source file", func(t *testing.T) {
		assert := test.AssertOn(t)
		tmp := filepath.Join(dir, "temp3")
		evacuated, err := PrepareEvacuation(tmp, nil).Of("./testdata/missing").By(DummyEvac)
		assert.ExpectError("expected error due to missing source file, but got none")(err)
		assert.True("expected evacuated to be nil after detecting missing source file")(nil == evacuated)
	})

	t.Run("create unique temporary subfolder", func(t *testing.T) {
		assert := test.AssertOn(t)
		tmp := filepath.Join(dir, "temp4")
		evacuated1, err := PrepareEvacuation(tmp, nil)("./testdata/small").By(DummyEvac)
		assert.NotError(err)
		evacuated2, err := PrepareEvacuation(tmp, nil)("./testdata/small").By(DummyEvac)
		assert.NotError(err)
		assert.True(fmt.Sprintf("expected temp folders to be unique, but got \"%s\" and \"%s\"", evacuated1.evacuatedTo, evacuated2.evacuatedTo))(evacuated1.evacuatedTo != evacuated2.evacuatedTo)
	})

	t.Run("calculate proper filename without replacement", func(t *testing.T) {
		assert := test.AssertOn(t)
		tmp := filepath.Join(dir, "temp5")
		evacuated, err := PrepareEvacuation(tmp, nil)("./testdata/small").By(DummyEvac)
		assert.NotError(err)
		assert.True("expected original file name to be kept")(strings.HasSuffix(evacuated.original,"small"))
		assert.True("expected evacuated file name to match original")(strings.HasSuffix(evacuated.evacuatedTo,"small"))
		assert.False("expected only file name, but not folders, to be used for destination file")(strings.Contains(evacuated.evacuatedTo,"testdata"))
	})

	t.Run("replace characters", func(t *testing.T) {
		assert := test.AssertOn(t)
		tmp := filepath.Join(dir, "temp6")
		evacuated, err := PrepareEvacuation(tmp, map[rune]rune {'t': 'x', 'a' : 'i', 'l' : 't'}).Of("./testdata/small").By(DummyEvac)
		assert.NotError(err)
		assert.True("expected characters to be replaced, but were not")(strings.HasSuffix(evacuated.evacuatedTo, "smitt"))
		assert.True("expected original file name to be kept")(strings.HasSuffix(evacuated.original,"small"))
	})
}

func TestMoving(t *testing.T) {
	dir := test.MkTempFolder(t)
	defer test.RmTempFolder(t, dir)

	t.Run("successful moving", func(t *testing.T) {
		assert := test.AssertOn(t)
		source := filepath.Join(dir, "source1")
		assert.AnythingNotError(Copy("./testdata/small", source, false))
		destination := filepath.Join(dir, "move1")

		evacuated, err := Moving(source, destination)
		assert.NotError(err)
		assert.FalseNotError("expected move source file not to exist anymore")(Exists(source))
		assert.TrueNotError("expected move destination file to exist")(Exists(destination))
		assert.StringsEqual(source, evacuated.original)
		assert.StringsEqual(destination, evacuated.evacuatedTo)
	})

	t.Run("move of missing source", func(t *testing.T) {
		assert := test.AssertOn(t)
		destination := filepath.Join(dir, "move2")
		evacuated, err := Moving(filepath.Join(dir, ".missing"), destination)
		assert.ExpectError("expected error when moving non-existing file")(err)
		assert.True("expect evacuated to be nil after error")(evacuated == nil)
	})

	t.Run("overwrite during move", func(t *testing.T) {
		assert := test.AssertOn(t)
		source := filepath.Join(dir, "source3")
		assert.AnythingNotError(Copy("./testdata/small", source, false))
		destination := filepath.Join(dir, "move3")
		Copy("./testdata/larger", destination, false) //pre-create file
		originalSize := sizeOf(destination)

		evacuated, err := Moving(source, destination)
		assert.NotError(err)
		assert.FalseNotError("expected move source file not to exist anymore")(Exists(source))
		assert.TrueNotError("expected move destination file to exist")(Exists(evacuated.evacuatedTo))
		assert.StringsEqual(source, evacuated.original)
		assert.StringsEqual(destination, evacuated.evacuatedTo)
		newSize := sizeOf(evacuated.evacuatedTo)
		assert.True("expected destination file to have different size after overwriting ig")(newSize != originalSize)
	})
}

func TestCopying(t *testing.T) {
	dir := test.MkTempFolder(t)
	defer test.RmTempFolder(t, dir)

	t.Run("successful copying", func(t *testing.T) {
		assert := test.AssertOn(t)
		source := "./testdata/small"
		destination := filepath.Join(dir, "copy1")
		evacuated, err := Copying(source, destination)
		assert.NotError(err)
		assert.TrueNotError("expected copy source file to exist")(Exists(source))
		assert.TrueNotError("expected copy destination file to exist")(Exists(destination))
		assert.StringsEqual(source, evacuated.original)
		assert.StringsEqual(destination, evacuated.evacuatedTo)
	})

	t.Run("copy of missing source", func(t *testing.T) {
		assert := test.AssertOn(t)
		destination := filepath.Join(dir, "copy2")
		evacuated, err := Copying(".missing", destination)
		assert.ExpectError("expected error when copying non-existing file")(err)
		assert.True("expect evacuated to be nil after error")(evacuated == nil)
	})

	t.Run("copy to pre-existing destination", func(t *testing.T) {
		assert := test.AssertOn(t)
		source := "./testdata/small"
		destination := filepath.Join(dir, "copy3")
		Copy(source, destination, false) //pre-create file

		evacuated, err := Copying(source, destination)
		assert.ExpectError("expected error when overwriting pre-existing file")(err)
		assert.True("expect evacuated to be nil after error")(evacuated == nil)
	})
}

func TestEvacuated(t *testing.T) {
	setup := func(t *testing.T) (*test.Assertion, string, *evacuated) {
		assert := test.AssertOn(t)
		tempDir := test.MkTempFolder(t)
		evacDir := filepath.Join(tempDir, "evac")

		e := &evacuated{filepath.Join(tempDir, "original"), filepath.Join(evacDir, "evacuated")}
		CreateFolderStructure(evacDir)

		assert.AnythingNotError(Copy("./testdata/small", e.original, false))
		assert.AnythingNotError(Copy("./testdata/larger", e.evacuatedTo, false))
		return assert, tempDir, e
	}

	t.Run("return correct path", func(t *testing.T) {
		assert := test.AssertOn(t)
		e := &evacuated{"orignate/from", "evacuated/to"}
		assert.StringsEqual(e.evacuatedTo, e.Path())
	})

	t.Run("successfully restore", func(t *testing.T) {
		assert, dir, evacuated := setup(t)
		defer test.RmTempFolder(t, dir)
		originalSize := sizeOf(evacuated.original)
		evacuatedSize := sizeOf(evacuated.evacuatedTo)
		assert.True("original and evacuated need different size for this test")(evacuatedSize != originalSize)
		assert.NotError(evacuated.Restore())
		assert.IntsEqual(evacuatedSize, sizeOf(evacuated.original))
		assert.FalseNotErrorf("expected evacuated file \"%s\" to be gone after restoring back to \"%s\"", evacuated.evacuatedTo, evacuated.original)(Exists(evacuated.evacuatedTo))
	})

	t.Run("delete all files in evacuation folder and folder itself", func(t *testing.T) {
		assert, dir, evacuated := setup(t)
		defer test.RmTempFolder(t, dir)

		//create extra dummz files in evacuation folder
		evacDir := filepath.Dir(evacuated.evacuatedTo)
		extraFile1 := filepath.Join(evacDir, "file1")
		extraFile2 := filepath.Join(evacDir, "file2")
		assert.AnythingNotError(Copy("./testdata/empty", extraFile1, false))
		assert.AnythingNotError(Copy("./testdata/small", extraFile2, false))

		// check preconditions
		assert.TrueNotErrorf("expected evacuated file \"%s\" to exist before discard", evacuated.evacuatedTo)(Exists(evacuated.evacuatedTo))
		assert.TrueNotErrorf("expected extra file \"%s\" to exist before discard", extraFile1)(Exists(extraFile1))
		assert.TrueNotErrorf("expected extra file \"%s\" to exist before discard", extraFile2)(Exists(extraFile2))

		assert.NotError(evacuated.Discard())

		// validate impact of discard operation
		assert.FalseNotErrorf("expected evacuation directory \"%s\" to be gone after discard", evacDir)(Exists(evacDir))
		assert.TrueNotErrorf("expected original file \"%s\" to still exist after discard", evacuated.original)(Exists(evacuated.original))
	})

	t.Run("move evacuated file and discard evac folder", func(t *testing.T) {
		assert, dir, evacuated := setup(t)
		defer test.RmTempFolder(t, dir)

		//create extra dummz files in evacuation folder
		evacDir := filepath.Dir(evacuated.evacuatedTo)
		extraFile1 := filepath.Join(evacDir, "file1")
		extraFile2 := filepath.Join(evacDir, "file2")
		assert.AnythingNotError(Copy("./testdata/empty", extraFile1, false))
		assert.AnythingNotError(Copy("./testdata/small", extraFile2, false))

		// check preconditions
		assert.TrueNotErrorf("expected evacuated file \"%s\" to exist before discard", evacuated.evacuatedTo)(Exists(evacuated.evacuatedTo))
		assert.TrueNotErrorf("expected extra file \"%s\" to exist before discard", extraFile1)(Exists(extraFile1))
		assert.TrueNotErrorf("expected extra file \"%s\" to exist before discard", extraFile2)(Exists(extraFile2))

		movedTo := filepath.Join(dir, "moved")
		assert.NotError(evacuated.MoveTo(movedTo))

		// validate impact of discard operation
		assert.FalseNotErrorf("expected evacuation directory \"%s\" to be gone after move", evacDir)(Exists(evacDir))
		assert.TrueNotErrorf("expected original file \"%s\" to still exist after move", evacuated.original)(Exists(evacuated.original))
		assert.TrueNotErrorf("expected moved file \"%s\" to exist after move", movedTo)(Exists(movedTo))
	})
}