package files

import (
	"testing"
	"go-cli/test"
	"path/filepath"
	"strings"
	"fmt"
	"go-cli/commons"
	"strconv"
)


func DummyEvac(from, to string) (*Evacuated, error) {
	return &Evacuated{original: from, evacuatedTo: to}, nil
}

func TestPrepareEvacuation(t *testing.T) {
	dir := test.MkTempFolder(t)
	defer test.RmTempFolder(t, dir)

	t.Run("check for empty.jpg temp folder", func(t *testing.T) {
		assert := test.AssertOn(t)
		evacuate := PrepareEvacuation("")
		evacuated, err := evacuate("./testdata/small.tiny").By(DummyEvac)
		assert.ExpectError("expected error due to empty.jpg temp folder name")(err)
		assert.True("expected Evacuated to be nil after finding empty.jpg temp folder name")(nil == evacuated)
	})

	t.Run("create temp folder if not existing", func(t *testing.T) {
		assert := test.AssertOn(t)
		tmp := filepath.Join(dir, "temp1")
		_, err := PrepareEvacuation(tmp).Of("./testdata/small.tiny").By(DummyEvac)
		assert.NotError(err)
		assert.TrueNotError("temp folder was not created")(Exists(tmp))
	})

	t.Run("re-use existing temp folder", func(t *testing.T) {
		assert := test.AssertOn(t)
		tmp := filepath.Join(dir, "temp2")
		assert.NotError(CreateFolderStructure(tmp))

		evacuated, err := PrepareEvacuation(tmp).Of("./testdata/small.tiny").By(DummyEvac)
		assert.NotError(err)
		assert.True("expected temp folder to be re-used")(strings.HasPrefix(evacuated.evacuatedTo, tmp))
	})

	t.Run("detect missing source file", func(t *testing.T) {
		assert := test.AssertOn(t)
		tmp := filepath.Join(dir, "temp3")
		evacuated, err := PrepareEvacuation(tmp).Of("./testdata/missing").By(DummyEvac)
		assert.ExpectError("expected error due to missing source file, but got none")(err)
		assert.True("expected Evacuated to be nil after detecting missing source file")(nil == evacuated)
	})

	t.Run("create unique temporary evacuation files", func(t *testing.T) {
		assert := test.AssertOn(t)
		tmp := filepath.Join(dir, "temp4")
		evacuated1, err := PrepareEvacuation(tmp)("./testdata/small.tiny").By(DummyEvac)
		assert.NotError(err)
		evacuated2, err := PrepareEvacuation(tmp)("./testdata/larger.png").By(DummyEvac)
		assert.NotError(err)
		assert.True(fmt.Sprintf("expected evacuation files to be unique, but got \"%s\" and \"%s\"", evacuated1.evacuatedTo, evacuated2.evacuatedTo))(evacuated1.evacuatedTo != evacuated2.evacuatedTo)
	})

	t.Run("calculate proper filename", func(t *testing.T) {
		assert := test.AssertOn(t)
		tmp := filepath.Join(dir, "temp5")
		evacuated, err := PrepareEvacuation(tmp)("./testdata/small.tiny").By(DummyEvac)
		assert.NotError(err)
		assert.StringsEqual("./testdata/small.tiny", evacuated.original) //assert original stays untouched
		assert.StringsEqual(filepath.Join(tmp, strconv.Itoa(int(commons.Hash32(evacuated.original))) + ".tiny"), evacuated.evacuatedTo)
		assert.False("expected only file name, but not folders, to be used for destination file")(strings.Contains(evacuated.evacuatedTo,"testdata"))
	})
}

func TestMoving(t *testing.T) {
	dir := test.MkTempFolder(t)
	defer test.RmTempFolder(t, dir)

	t.Run("successful moving", func(t *testing.T) {
		assert := test.AssertOn(t)
		source := filepath.Join(dir, "source1")
		assert.AnythingNotError(Copy("./testdata/small.tiny", source, false))
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
		assert.True("expect Evacuated to be nil after error")(evacuated == nil)
	})

	t.Run("overwrite during move", func(t *testing.T) {
		assert := test.AssertOn(t)
		source := filepath.Join(dir, "source3")
		assert.AnythingNotError(Copy("./testdata/small.tiny", source, false))
		destination := filepath.Join(dir, "move3")
		Copy("./testdata/larger.png", destination, false) //pre-create file
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
		source := "./testdata/small.tiny"
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
		assert.True("expect Evacuated to be nil after error")(evacuated == nil)
	})

	t.Run("copy to pre-existing destination", func(t *testing.T) {
		assert := test.AssertOn(t)
		source := "./testdata/small.tiny"
		destination := filepath.Join(dir, "copy3")
		Copy(source, destination, false) //pre-create file

		evacuated, err := Copying(source, destination)
		assert.ExpectError("expected error when overwriting pre-existing file")(err)
		assert.True("expect Evacuated to be nil after error")(evacuated == nil)
	})
}

func TestEvacuated(t *testing.T) {
	setup := func(t *testing.T) (*test.Assertion, string, *Evacuated) {
		assert := test.AssertOn(t)
		tempDir := test.MkTempFolder(t)
		evacDir := filepath.Join(tempDir, "evac")

		e := &Evacuated{filepath.Join(tempDir, "original"), filepath.Join(evacDir, "Evacuated")}
		CreateFolderStructure(evacDir)

		assert.AnythingNotError(Copy("./testdata/small.tiny", e.original, false))
		assert.AnythingNotError(Copy("./testdata/larger.png", e.evacuatedTo, false))
		return assert, tempDir, e
	}

	t.Run("return correct path", func(t *testing.T) {
		assert := test.AssertOn(t)
		e := &Evacuated{"originate/from", "Evacuated/to"}
		assert.StringsEqual(e.evacuatedTo, e.Path())
	})

	t.Run("return correct path with suffix", func(t *testing.T) {
		assert := test.AssertOn(t)
		e := &Evacuated{"originate/from.xyz", "evacuated/to.xyz"}
		assert.StringsEqual("evacuated/to.suffix.xyz", e.WithSuffix(".suffix"))
	})

	t.Run("successfully restore", func(t *testing.T) {
		assert, dir, evacuated := setup(t)
		defer test.RmTempFolder(t, dir)

		originalSize := sizeOf(evacuated.original)
		evacuatedSize := sizeOf(evacuated.evacuatedTo)
		assert.True("original and Evacuated need different size for this test")(evacuatedSize != originalSize)
		assert.NotError(evacuated.Restore())
		assert.IntsEqual(evacuatedSize, sizeOf(evacuated.original))
		assert.FalseNotErrorf("expected Evacuated file \"%s\" to be gone after restoring back to \"%s\"", evacuated.evacuatedTo, evacuated.original)(Exists(evacuated.evacuatedTo))
	})

	t.Run("delete evacuated file but leave temp-dir", func(t *testing.T) {
		assert, dir, evacuated := setup(t)
		defer test.RmTempFolder(t, dir)

		//create extra dummz files in evacuation folder
		evacDir := filepath.Dir(evacuated.evacuatedTo)
		extraFile := filepath.Join(evacDir, "file1")
		assert.AnythingNotError(Copy("./testdata/empty.jpg", extraFile, false))

		// check preconditions
		assert.TrueNotErrorf("expected Evacuated file \"%s\" to exist before discard", evacuated.evacuatedTo)(Exists(evacuated.evacuatedTo))
		assert.TrueNotErrorf("expected extra file \"%s\" to exist before discard", extraFile)(Exists(extraFile))

		assert.NotError(evacuated.Discard())

		// validate impact of discard operation
		assert.TrueNotErrorf("expected evacuation directory \"%s\" to still exist after discard", evacDir)(Exists(evacDir))
		assert.FalseNotErrorf("expect evacuated file \"%s\" to be gone after discard", evacuated.evacuatedTo)(Exists(evacuated.evacuatedTo))
		assert.TrueNotErrorf("expected original file \"%s\" to still exist after discard", evacuated.original)(Exists(evacuated.original))
		assert.TrueNotErrorf("expected extra file \"%s\" to still exist after discard", extraFile)(Exists(extraFile))
	})

	t.Run("move Evacuated file and discard evac folder", func(t *testing.T) {
		assert, dir, evacuated := setup(t)
		defer test.RmTempFolder(t, dir)

		//create extra dummz files in evacuation folder
		evacDir := filepath.Dir(evacuated.evacuatedTo)

		// check preconditions
		assert.TrueNotErrorf("expected Evacuated file \"%s\" to exist before discard", evacuated.evacuatedTo)(Exists(evacuated.evacuatedTo))

		movedTo := filepath.Join(dir, "moved")
		assert.NotError(evacuated.MoveTo(movedTo))

		// validate impact of discard operation

		assert.TrueNotErrorf("expected evacuation directory \"%s\" to still exist after discard", evacDir)(Exists(evacDir))
		assert.FalseNotErrorf("expect evacuated file \"%s\" to be gone after discard", evacuated.evacuatedTo)(Exists(evacuated.evacuatedTo))
		assert.TrueNotErrorf("expected original file \"%s\" to still exist after discard", evacuated.original)(Exists(evacuated.original))
		assert.TrueNotErrorf("expected moved file \"%s\" to exist after move", movedTo)(Exists(movedTo))
	})
}
