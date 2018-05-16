package targetinfo

import (
	"testing"
	"io/ioutil"
	"os"
	"path/filepath"
)

var ti = TargetInfo{"f.g", "/a/b/c", "test", "tt987654321", 1, 0}

func setup(t *testing.T) string {
	systemTempDir := "" //defaults to tmp dir in linux, windows, etc.
	dir, err := ioutil.TempDir(systemTempDir, "go-test")
	if err != nil {
		t.Error(err)
	}
	return dir
}

func teardown(t *testing.T, dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		t.Error(err)
	}
}

func checkError(t *testing.T, e error) {
	if e != nil {
		t.Error(e)
	}
}
func TestSaveJson(t *testing.T) {
	dir := setup(t)
	defer teardown(t, dir)
	checkError(t, Save(dir, &ti))

	f := filepath.Join(dir, ti.Id)
	info, err := os.Stat(f)
	checkError(t, err)

	if info.Size() == 0 {
		t.Errorf("target file \"%s\" must not be empty", f)
	}
}

func TestReadJson(t *testing.T) {
	dir := setup(t)
	defer teardown(t, dir)

	checkError(t, Save(dir, &ti))
	read, err := Read(dir, ti.Id)
	checkError(t, err)

	if *read != ti {
		t.Errorf("targetinfo does not match:\n  to json   %v\n  from json %v", ti, *read)
	}
}

func TestOverwriteJsonFile(t *testing.T) {
	dir := setup(t)
	defer teardown(t, dir)

	checkError(t, Save(dir, &ti))
	newId := "tt12345"
	newKind := "changed"
	ti.Id = newId
	ti.Kind = newKind

	//overwrite
	checkError(t, Save(dir, &ti))

	read, err := Read(dir, ti.Id)
	checkError(t, err)

	if read.Id != newId {
		t.Errorf("ids not matching: expected %s, but got %s", newId, read.Id)
	}
	if read.Kind != newKind {
		t.Errorf("kinds not matching: expected %s, but got %s", newKind, read.Kind)
	}
}

func TestReadUnmatchingFileNameAndId(t *testing.T) {
	dir := setup(t)
	defer teardown(t, dir)

	//rename file to create difference btw. contained target-id and target-id in filename
	checkError(t, Save(dir, &ti))
	newFileId := "tt666"
	checkError(t, os.Rename(filepath.Join(dir, ti.Id), filepath.Join(dir, newFileId)))

	_, err := Read(dir, newFileId)
	if err == nil {
		t.Errorf("unmatching IDs in file-name and -contents not detected during read:\n  file=%s\n  json=%s", newFileId, ti.Id)
	}
}

func TestSaveNilTargetInfo(t *testing.T) {
	err := Save(".", nil)
	if err == nil {
		t.Error("expected error when saving nil TargetInfo")
	}
}
