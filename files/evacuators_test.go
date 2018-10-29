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


//func TestCopyEvacuator(t *testing.T) {
//	dir := test.MkTempFolder(t)
//	defer test.RmTempFolder(t, dir)
//	largeFile := filepath.Join(dir, "larger")
//	emptyFile := filepath.Join(dir, "empty")
//	Copy("./testdata/larger", largeFile, false)
//	Copy("./testdata/empty", emptyFile, false)
//
//	t.Run("copy created while original remains", func(t *testing.T) {
//		evacuate := CopyingEvacuator(dir, nil)
//		evac, err := evacuate(largeFile)
//
//		copy := evac.(*copied).copy
//		assert := test.AssertOn(t)
//		assert.NotError(err)
//		assert.TrueNotError("original not available after copy")(Exists(largeFile))
//		assert.TrueNotError("copy was not created")(Exists(copy))
//		assert.False("original and copy are the same")(largeFile == copy)
//	})
//
//	t.Run("restore replaces original with copy", func(t *testing.T) {
//		evacuate := CopyingEvacuator(dir, nil)
//		evac, _ := evacuate(largeFile)
//		originalSize := sizeOf(largeFile)
//
//		copy := evac.(*copied).copy
//		assert := test.AssertOn(t)
//		assert.TrueNotError("original not available after copy")(Exists(largeFile))
//		assert.TrueNotError("copy was not created")(Exists(copy))
//
//		//modify (aka replace) evacuated file
//		assert.TrueNotError("failed to replace (modify) evacuated file")(Replace(copy, emptyFile))
//
//		assert.NotError(evac.Restore())
//		assert.TrueNotError("original not available after copy and restore")(Exists(largeFile))
//		assert.FalseNotError("copy was not deleted during restore")(Exists(copy))
//		newSize := sizeOf(largeFile)
//		assert.Falsef("file size must differ after evacuation (original=%d, evacuated=%d", originalSize, newSize)(originalSize == newSize)
//	})
//
//	t.Run("discard deletes copy but leaves original", func(t *testing.T) {
//		evacuate := CopyingEvacuator(dir, nil)
//		evac, _ := evacuate(largeFile)
//
//		copy := evac.(*copied).copy
//		assert := test.AssertOn(t)
//		assert.TrueNotError("original not available after copy")(Exists(largeFile))
//		assert.TrueNotError("copy was not created")(Exists(copy))
//
//		assert.NotError(evac.Discard())
//		assert.TrueNotError("original not available after copy and restore")(Exists(largeFile))
//		assert.FalseNotError("copy was not deleted during restore")(Exists(copy))
//	})
//
//	t.Run("move leaves original and moves copy", func(t *testing.T) {
//		evacuate := CopyingEvacuator(dir, nil)
//		evac, _ := evacuate(largeFile)
//
//		copy := evac.(*copied).copy
//		assert := test.AssertOn(t)
//		assert.TrueNotError("original not available after copy")(Exists(largeFile))
//		assert.TrueNotError("copy was not created")(Exists(copy))
//
//		moved := filepath.Join(dir, "moved")
//		assert.NotError(evac.MoveTo(moved))
//		assert.TrueNotError("original not available after copy and restore")(Exists(largeFile))
//		assert.FalseNotError("copy was not deleted during move")(Exists(copy))
//		assert.TrueNotError("no new copy at location where it was moved")(Exists(moved))
//	})
//}
//
//func TestMoveEvacuator(t *testing.T) {
//	dir := test.MkTempFolder(t)
//	defer test.RmTempFolder(t, dir)
//	largeFile := filepath.Join(dir, "larger")
//	emptyFile := filepath.Join(dir, "empty")
//	Copy("./testdata/larger", largeFile, false)
//	Copy("./testdata/empty", emptyFile, false)
//
//	t.Run("evacuates file to new location - original location is empty", func(t *testing.T) {
//		src := filepath.Join(dir, "file1")
//		Copy(largeFile, src, false)
//
//		evacuate := MovingEvacuator(dir, nil)
//		evac, err := evacuate(src)
//		m := evac.(*moved)
//
//		assert := test.AssertOn(t)
//		assert.NotError(err)
//		assert.FalseNotError("original still available after move")(Exists(src))
//		assert.TrueNotError("moved file not found")(Exists(m.movedTo))
//		assert.StringsEqual(src, m.original)
//	})
//
//	t.Run("restore moves evacuated back to original location", func(t *testing.T) {
//		src := filepath.Join(dir, "file2")
//		Copy(largeFile, src, false)
//
//		evacuate := MovingEvacuator(dir, nil)
//		evac, _ := evacuate(src)
//		m := evac.(*moved)
//
//		assert := test.AssertOn(t)
//		assert.FalseNotError("original still available after move")(Exists(src))
//		assert.TrueNotError("moved file not found")(Exists(m.movedTo))
//
//		assert.NotError(evac.Restore())
//		assert.TrueNotError("original not available after move and restore")(Exists(src))
//		assert.FalseNotError("moved file still exists")(Exists(m.movedTo))
//
//	})
//
//	t.Run("discard deletes evacuated without leaving original", func(t *testing.T) {
//		src := filepath.Join(dir, "file3")
//		Copy(largeFile, src, false)
//
//		evacuate := MovingEvacuator(dir, nil)
//		evac, _ := evacuate(src)
//
//		m := evac.(*moved)
//		assert := test.AssertOn(t)
//		assert.FalseNotError("original still available after move")(Exists(src))
//		assert.TrueNotError("moved file not found")(Exists(m.movedTo))
//
//		assert.NotError(evac.Discard())
//		assert.FalseNotError("original available despite moved was discarded")(Exists(src))
//		assert.FalseNotError("moved file still exists")(Exists(m.movedTo))
//	})
//
//	t.Run("move moves evacuated to another location", func(t *testing.T) {
//		src := filepath.Join(dir, "file4")
//		Copy(largeFile, src, false)
//
//		evacuate := MovingEvacuator(dir, nil)
//		evac, _ := evacuate(src)
//
//		m := evac.(*moved)
//		assert := test.AssertOn(t)
//		assert.FalseNotError("original still available after move")(Exists(src))
//		assert.TrueNotError("moved file not found")(Exists(m.movedTo))
//
//		moved := filepath.Join(dir, "moved-somewhere-else")
//		assert.NotError(evac.MoveTo(moved))
//		assert.FalseNotError("original available despite moved was discarded")(Exists(src))
//		assert.FalseNotError("moved file still exists")(Exists(m.movedTo))
//		assert.TrueNotError("file is not available at location where it was moved")(Exists(moved))
//	})
//}
