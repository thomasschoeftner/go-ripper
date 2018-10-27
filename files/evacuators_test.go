package files

import (
	"testing"
	"go-cli/test"
	"path/filepath"
	"strings"
	"fmt"
)

func TestPrepareEvacuation(t *testing.T) {
	dir := test.MkTempFolder(t)
	defer test.RmTempFolder(t, dir)

	t.Run("check for empty temp folder", func(t *testing.T) {
		dst, err := prepareEvacuation("", nil, "./testdata/small")
		assert := test.AssertOn(t)
		assert.ExpectError("expected error due to empty temp folder name")(err)
		assert.StringsEqual("", dst)
	})

	t.Run("create temp folder if not existing", func(t *testing.T) {
		tmp := filepath.Join(dir, "temp1")
		_, err := prepareEvacuation(tmp, nil, "./testdata/small")
		assert := test.AssertOn(t)
		assert.NotError(err)
		assert.TrueNotError("temp folder was not created")(Exists(tmp))
	})

	t.Run("re-use existing temp folder", func(t *testing.T) {
		tmp := filepath.Join(dir, "temp2")
		assert := test.AssertOn(t)
		assert.NotError(CreateFolderStructure(tmp))

		dst, err := prepareEvacuation(tmp, nil, "./testdata/small")
		assert.NotError(err)
		assert.True("expected temp folder to be re-used")(strings.HasPrefix(dst, tmp))
	})

	t.Run("detect missing source file", func(t *testing.T) {
		tmp := filepath.Join(dir, "temp3")
		dst, err := prepareEvacuation(tmp, nil, "missing")
		assert := test.AssertOn(t)
		assert.ExpectError("expected error due to missing source file, but got none")(err)
		assert.StringsEqual("", dst)
	})

	t.Run("create unique temporary subfolder", func(t *testing.T) {
		tmp := filepath.Join(dir, "temp4")
		assert := test.AssertOn(t)
		dst1, err := prepareEvacuation(tmp, nil, "./testdata/small")
		assert.NotError(err)
		dst2, err := prepareEvacuation(tmp, nil, "./testdata/small")
		assert.NotError(err)
		assert.True(fmt.Sprintf("expected temp folders to be unique, but got \"%s\" and \"%s\"", dst1, dst2))(dst1 != dst2)
	})

	t.Run("calculate proper filename without replacement", func(t *testing.T) {
		tmp := filepath.Join(dir, "temp5")
		assert := test.AssertOn(t)
		dst, err := prepareEvacuation(tmp, nil, "./testdata/small")
		assert.NotError(err)
		assert.True("expected original file name to be kept")(strings.HasSuffix(dst, "small"))
		assert.False("expected only file name to be used for destination file")(strings.Contains(dst, "testdata"))
	})

	t.Run("replace characters", func(t *testing.T) {
		tmp := filepath.Join(dir, "temp6")
		assert := test.AssertOn(t)
		dst, err := prepareEvacuation(tmp, map[rune]rune {'t': 'x', 'a' : 'i', 'l' : 't'}, "./testdata/small")
		assert.NotError(err)
		assert.True("expected characters to be replaced, but were not")(strings.HasSuffix(dst, "smitt"))
		assert.False("expected only file name to be used for destination file")(strings.Contains(dst, "testdata"))
	})
}


func TestCopyEvacuator(t *testing.T) {
	dir := test.MkTempFolder(t)
	defer test.RmTempFolder(t, dir)
	source := "./testdata/larger"

	t.Run("copy created while original remains", func(t *testing.T) {
		evacuate := CopyingEvacuator(dir, nil)
		evac, err := evacuate(source)

		copy := evac.(*copied).copy
		assert := test.AssertOn(t)
		assert.NotError(err)
		assert.TrueNotError("original not available after copy")(Exists(source))
		assert.TrueNotError("copy was not created")(Exists(copy))
		assert.False("orignal and copy are the same")(source == copy)
	})

	t.Run("restore deletes copy and leaves original", func(t *testing.T) {
		evacuate := CopyingEvacuator(dir, nil)
		evac, _ := evacuate(source)

		copy := evac.(*copied).copy
		assert := test.AssertOn(t)
		assert.TrueNotError("original not available after copy")(Exists(source))
		assert.TrueNotError("copy was not created")(Exists(copy))

		assert.NotError(evac.Restore())
		assert.TrueNotError("original not available after copy and restore")(Exists(source))
		assert.FalseNotError("copy was not deleted during restore")(Exists(copy))
	})

	t.Run("discard deletes copy but leaves original", func(t *testing.T) {
		evacuate := CopyingEvacuator(dir, nil)
		evac, _ := evacuate(source)

		copy := evac.(*copied).copy
		assert := test.AssertOn(t)
		assert.TrueNotError("original not available after copy")(Exists(source))
		assert.TrueNotError("copy was not created")(Exists(copy))

		assert.NotError(evac.Discard())
		assert.TrueNotError("original not available after copy and restore")(Exists(source))
		assert.FalseNotError("copy was not deleted during restore")(Exists(copy))
	})

	t.Run("move leaves original and moves copy", func(t *testing.T) {
		evacuate := CopyingEvacuator(dir, nil)
		evac, _ := evacuate(source)

		copy := evac.(*copied).copy
		assert := test.AssertOn(t)
		assert.TrueNotError("original not available after copy")(Exists(source))
		assert.TrueNotError("copy was not created")(Exists(copy))

		moved := filepath.Join(dir, "moved")
		assert.NotError(evac.MoveTo(moved))
		assert.TrueNotError("original not available after copy and restore")(Exists(source))
		assert.FalseNotError("copy was not deleted during move")(Exists(copy))
		assert.TrueNotError("no new copy at location where it was moved")(Exists(moved))
	})
}

func TestMoveEvacuator(t *testing.T) {
	t.Run("evacuates file to new location - original location is empty", func(t *testing.T) {

	})

	t.Run("restore moves evacuated back to original location", func(t *testing.T) {

	})

	t.Run("discard deletes evacuated without leaving original", func(t *testing.T) {

	})

	t.Run("move moves evacuated to another location", func(t *testing.T) {

	})
}