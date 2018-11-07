package files

import (
	"testing"
	"go-cli/test"
	"path/filepath"
	"fmt"
	"io/ioutil"
	"os"
	"math"
)

func sizeOf(fileName string) int {
	stats, err := os.Stat(fileName)
	if err != nil {
		return -1
	}
	if stats.Size() > math.MaxInt32 {
		return -2
	}
	return int(stats.Size())
}


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

	t.Run("empty.jpg source file", func(t *testing.T){
		assert := test.AssertOn(t)
		target := filepath.Join(dir, "empty.jpg")
		cnt, err := Copy(filepath.Join(".", testDataFolder, "empty.jpg"), target, false)
		assert.NotError(err)
		assert.IntsEqual(0, int(cnt))

		exists, err := Exists(target)
		assert.NotError(err)
		assert.True("copying empty.jpg file did not create target file")(exists)
	})

	t.Run("small.tiny source file", func(t *testing.T){
		assert := test.AssertOn(t)
		target := filepath.Join(dir, "small.tiny")
		cnt, err := Copy(filepath.Join(".", testDataFolder, "small.tiny"), target, false)
		assert.NotError(err)
		assert.IntsEqual(3, int(cnt))

		exists, err := Exists(target)
		assert.NotError(err)
		assert.True("copying small.tiny file did not create target file")(exists)
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

		cnt, err := Copy(filepath.Join(".", testDataFolder, "small.tiny"), target, false)
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

		cnt, err := Copy(filepath.Join(".", testDataFolder, "small.tiny"), target, true)
		assert.NotError(err)
		assert.IntsEqual(3, int(cnt))
	})
}

func TestReplace(t *testing.T) {
	dir := test.MkTempFolder(t)
	defer test.RmTempFolder(t, dir)

	_cnt := 0
	setup := func(t *testing.T) (*test.Assertion, string, string) {
		_cnt = _cnt + 1
		assert := test.AssertOn(t)
		o := filepath.Join(dir, fmt.Sprintf("original_%d", _cnt))
		r := filepath.Join(dir, fmt.Sprintf("replacement_%d", _cnt))
		assert.AnythingNotError(Copy("./testdata/small.tiny", o, false))
		assert.AnythingNotError(Copy("./testdata/larger.png", r, false))
		return assert, o, r
	}

	t.Run("replace file, keep replacement, but do not keep original", func(t* testing.T) {
		assert, o, r := setup(t)
		originalSize := sizeOf(o)
		assert.TrueNotErrorf("did not replace \"%s\" with \"%s\"", o, r)(Replace(o, r))
		assert.TrueNotErrorf("replacement file \"%s\" does not exist anymore", r)(Exists(r))
		newSize := sizeOf(o)
		assert.IntsEqual(newSize, sizeOf(r))
		assert.False("replaced file is still of same size")(newSize == originalSize)
	})

	t.Run("replace file, but keep original", func(t* testing.T) {
		assert, o, r := setup(t)
		originalSize := sizeOf(o)
		assert.TrueNotErrorf("did not replace \"%s\" with \"%s\"", o, r)(ReplaceButKeepOriginal(o, r, "backup"))
		assert.TrueNotErrorf("replacement file \"%s\" does not exist anymore", r)(Exists(r))
		assert.TrueNotErrorf("replacement file \"%s\" does not exist anymore", r)(Exists(r))
		newSize := sizeOf(o)
		backup := WithExtension(o, "backup")
		assert.TrueNotErrorf("backup file \"%s\" missing", backup)(Exists(backup))
		assert.IntsEqual(newSize, sizeOf(r))
		assert.IntsEqual(originalSize, sizeOf(backup))
	})

	t.Run("check for empty.jpg keep-file extension", func(t* testing.T) {
		assert, o, r := setup(t)
		replaced, err := ReplaceButKeepOriginal(o, r,"")
		assert.ExpectError("expect error when trying to replace file but keep backup with identical name")(err)
		assert.False("should not replace file if keep-extension is missing")(replaced)
	})

	t.Run("fail to replace with non-existent file", func(t* testing.T) {
		assert, o, r := setup(t)
		assert.NotError(os.Remove(r))
		replaced, err := ReplaceButKeepOriginal(o, r,"backup")
		assert.ExpectError("expect error when trying to replace file with non-existing file")(err)
		assert.False("should not replace file with non-existing")(replaced)
	})

	t.Run("silently replace non-existent file", func(t* testing.T) {
		assert, o, r := setup(t)
		assert.NotError(os.Remove(o))
		assert.FalseNotError("should not return true if replaced file did not exist")(ReplaceButKeepOriginal(o, r,"backup"))
		assert.IntsEqual(sizeOf(o), sizeOf(r))
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

	t.Run("return empty.jpg string if no '.'", func (t* testing.T) {
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
