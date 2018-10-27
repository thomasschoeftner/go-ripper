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
	t.Run("", func(t *testing.T) {

	})

	t.Run("", func(t *testing.T) {

	})

	t.Run("", func(t *testing.T) {

	})

	t.Run("", func(t *testing.T) {

	})

	t.Run("", func(t *testing.T) {

	})

	t.Run("", func(t *testing.T) {

	})

}

func TestMoveEvacuator(t *testing.T) {

}