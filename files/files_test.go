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

func TestCopy(t *testing.T) {
	dir := test.MkTempFolder(t)
	defer test.RmTempFolder(t, dir)
	const testDataFolder = "testdata"

	t.Run("empty source file", func(t *testing.T){
		assert := test.AssertOn(t)
		target := filepath.Join(dir, "empty")
		cnt, err := Copy(filepath.Join(".", testDataFolder, "empty"), target, false)
		assert.NotError(err)
		assert.IntsEqual(0, int(cnt))

		exists, err := Exists(target)
		assert.NotError(err)
		assert.True("copying empty file did not create target file")(exists)
	})

	t.Run("small source file", func(t *testing.T){
		assert := test.AssertOn(t)
		target := filepath.Join(dir, "small")
		cnt, err := Copy(filepath.Join(".", testDataFolder, "small"), target, false)
		assert.NotError(err)
		assert.IntsEqual(3, int(cnt))

		exists, err := Exists(target)
		assert.NotError(err)
		assert.True("copying small file did not create target file")(exists)
	})

	t.Run("missing source file", func(t *testing.T){
		assert := test.AssertOn(t)
		target := filepath.Join(dir, "missing")
		cnt, err := Copy(filepath.Join(".", testDataFolder, "missing"), target, false)
		assert.ExpectError("expected error when copying missing file, but got none")(err)
		assert.IntsEqual(0, int(cnt))

		exists, err := Exists(target)
		assert.NotError(err)
		assert.False("copying missing file did create target file")(exists)
	})

	t.Run("pre-existing destination file without truncate", func(t *testing.T){
		assert := test.AssertOn(t)
		target := filepath.Join(dir, "preexists1")
		f, _ := os.Create(target)
		f.Close()
		exists, err := Exists(target)
		assert.True("preexisting target file exists")(exists)

		cnt, err := Copy(filepath.Join(".", testDataFolder, "small"), target, false)
		assert.ExpectError("expected error when copying to preexisting file with truncation disabled")(err)
		assert.IntsEqual(0, int(cnt))
	})

	t.Run("pre-existing destination file with truncate", func(t *testing.T) {
		assert := test.AssertOn(t)
		target := filepath.Join(dir, "preexists2")
		f, _ := os.Create(target)
		f.Close()
		exists, err := Exists(target)
		assert.True("preexisting target file exists")(exists)

		cnt, err := Copy(filepath.Join(".", testDataFolder, "small"), target, true)
		assert.NotError(err)
		assert.IntsEqual(3, int(cnt))
	})
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

func TestGetExtension(t *testing.T) {
	t.Run("remove leading '.'", func (t *testing.T) {
		test.AssertOn(t).StringsEqual("def", GetExtension("s/b/c.def"))
	})

	t.Run("remove multiple leading '.'", func (t *testing.T) {
		test.AssertOn(t).StringsEqual("def", GetExtension("s/b/c....def"))
	})

	t.Run("return empty string if no '.'", func (t* testing.T) {
		test.AssertOn(t).StringsEqual("", GetExtension("s/b/cdef"))
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

	t.Run("complete folders including '.'s in folder names", func (t* testing.T) {
		folder := "/a/b/.c/..d/...e/f./g../.h./"
		name := "sepp"
		extension := "rtf"

		rest, ext:= SplitExtension(folder + name + "." + extension)
		assert := test.AssertOn(t)
		assert.StringsEqual(folder + name, rest)
		assert.StringsEqual(extension, ext)
	})

	t.Run("files with leading '.' with extension", func (t* testing.T) {
		name, ext := SplitExtension(".abc.def")
		assert := test.AssertOn(t)
		assert.StringsEqual(".abc", name)
		assert.StringsEqual("def", ext)
	})

	t.Run("files with leading '.' without extension", func (t* testing.T) {
		name, ext := SplitExtension(".abc")
		assert := test.AssertOn(t)
		assert.StringsEqual(".abc", name)
		assert.StringsEqual("", ext)
	})

	t.Run("files with multiple '.' before extension", func (t* testing.T) {
		name, ext := SplitExtension(".abc...xyz")
		assert := test.AssertOn(t)
		assert.StringsEqual(".abc..", name)
		assert.StringsEqual("xyz", ext)
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
