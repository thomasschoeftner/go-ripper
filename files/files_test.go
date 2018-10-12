package files

import (
	"testing"
	"go-cli/test"
	"path/filepath"
	"fmt"
	"io/ioutil"
	"os"
)

func TestCreateFolder(t *testing.T) {
	dir := test.MkTempFolder(t)
	defer test.RmTempFolder(t, dir)

	t.Run("new folder", func(t *testing.T) {
		folder := filepath.Join(dir, "fritz")
		assert := test.AssertOn(t)
		assert.NotError(CreateFolder(folder))

		assert.ExpectError("recreate duplicate folder should fail, but did not")(CreateFolder(folder))
	})

	t.Run("folder structure", func(t *testing.T) {
		folder := filepath.Join(dir, "a", "b", "c")
		test.AssertOn(t).ExpectError("creating complex folder structure is expected to fail, but did not")(CreateFolder(folder))
	})
}

func TestCreateFolderStructure(t *testing.T) {
	dir := test.MkTempFolder(t)
	defer test.RmTempFolder(t, dir)

	folder := filepath.Join(dir, "x", "y", "z")
	assert := test.AssertOn(t)
	assert.NotError(CreateFolderStructure(folder))

	//succeed on recreating same folder
	assert.NotError(CreateFolderStructure(folder))

	subFolder := filepath.Join(folder, "sub")
	assert.NotError(CreateFolderStructure(subFolder))
}

func TestExists(t *testing.T) {
	dir := test.MkTempFolder(t)
	defer test.RmTempFolder(t, dir)

	t.Run("folder", func(t *testing.T) {
		assert := test.AssertOn(t)
		folder := filepath.Join(dir, "x", "y", "z")
		assert.FalseNotError(fmt.Sprintf("expected folder %s not to exist yet, but did", folder))(Exists(folder))
		CreateFolderStructure(folder)
		assert.TrueNotError(fmt.Sprintf("expected folder %s to exist, but did not", folder))(Exists(folder))
	})

	t.Run("file", func(t *testing.T) {
		assert := test.AssertOn(t)
		file := filepath.Join(dir, "a.b")
		assert.FalseNotError(fmt.Sprintf("expected file %s not to exist yet, but did", file))(Exists(file))
		ioutil.WriteFile(file, []byte{}, os.ModePerm)
		assert.TrueNotError(fmt.Sprintf("expected file %s to exist, but did not", file))(Exists(file))
	})
}

func TestSplitExtension(t *testing.T) {
	t.Run("common file name", func(t *testing.T) {
		file := "abc.de"
		name, ext := SplitExtension(file)
		assert := test.AssertOn(t)
		assert.StringsEqual("abc", name)
		assert.StringsEqual("de", ext)
	})

	t.Run("no extension", func(t *testing.T) {
		file := "abcde"
		name, ext := SplitExtension(file)
		assert := test.AssertOn(t)
		assert.StringsEqual("abcde", name)
		assert.StringsEqual("", ext)
	})

	t.Run("multiple periods", func(t *testing.T) {
		file := "abc.de.fg"
		name, ext := SplitExtension(file)
		assert := test.AssertOn(t)
		assert.StringsEqual("abc.de", name)
		assert.StringsEqual("fg", ext)
	})
}

func TestWithExtension(t *testing.T) {
	t.Run("add extension without leading '.'", func (t *testing.T) {
		name := "frank"
		ext := "txt"
		test.AssertOn(t).StringsEqual(name + "." + ext, WithExtension(name, ext))
	})

	t.Run("add extension with leading '.'", func (t *testing.T) {
		name := "frank"
		ext := ".txt"
		test.AssertOn(t).StringsEqual(name + ext, WithExtension(name, ext))
	})

}