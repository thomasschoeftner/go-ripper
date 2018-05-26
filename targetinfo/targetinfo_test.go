package targetinfo

import (
	"testing"
	"io/ioutil"
	"os"
	"path/filepath"
	"go-cli/test"
	"fmt"
)

var singleTitle = TargetInfo{"f.g", "/a/b/c", "test", "tt987654321", 0, 0}
var partitionedTitle = TargetInfo{"f.g", "/a/b/c", "test", "tt987654321", 0, 2}
var singleCollectionItem = TargetInfo{"f.g", "/a/b/c", "test", "tt987654321", 3, 12}

func setup(t *testing.T) string {
	systemTempDir := "" //defaults to tmp dir in linux, windows, etc.
	dir, err := ioutil.TempDir(systemTempDir, "go-test")
	test.CheckError(t,err)
	return dir
}

func teardown(t *testing.T, dir string) {
	err := os.RemoveAll(dir)
	test.CheckError(t,err)
}


func TestSaveJson(t *testing.T) {
	t.Run("save single singleTitle", func(t *testing.T) {
		ti := singleTitle
		dir := setup(t)
		defer teardown(t, dir)
		_, err := Save(dir, &ti)
		test.CheckError(t, err)

		fname := ti.Id
		f := filepath.Join(dir, fname)
		info, err := os.Stat(f)
		test.CheckError(t, err)

		if info.Size() == 0 {
			t.Errorf("target file \"%s\" must not be empty", f)
		}
	})

	t.Run("save singleTitle part", func(t *testing.T) {
		ti := partitionedTitle
		dir := setup(t)
		defer teardown(t, dir)
		_, err := Save(dir, &ti)
		test.CheckError(t, err)

		fname := fmt.Sprintf("%s.%d", ti.Id, ti.ItemNo)
		f := filepath.Join(dir, fname)
		info, err := os.Stat(f)
		test.CheckError(t, err)

		if info.Size() == 0 {
			t.Errorf("target file \"%s\" must not be empty", f)
		}
	})

	t.Run("save series singleCollectionItem", func(t *testing.T) {
		ti := singleCollectionItem
		dir := setup(t)
		defer teardown(t, dir)
		_, err := Save(dir, &ti)
		test.CheckError(t, err)

		fname := fmt.Sprintf("%s.%d.%d", ti.Id, ti.Collection, ti.ItemNo)
		f := filepath.Join(dir, fname)
		info, err := os.Stat(f)
		test.CheckError(t, err)

		if info.Size() == 0 {
			t.Errorf("target file \"%s\" must not be empty", f)
		}
	})
}

func TestReadJson(t *testing.T) {
	dir := setup(t)
	defer teardown(t, dir)

	_, err := Save(dir, &singleCollectionItem)
	test.CheckError(t, err)
	read, err := Read(filepath.Join(dir, singleCollectionItem.fileName()))
	test.CheckError(t, err)

	if *read != singleCollectionItem {
		t.Errorf("targetinfo does not match:\n  to json   %v\n  from json %v", singleCollectionItem, *read)
	}
}

func TestOverwriteJsonFile(t *testing.T) {
	dir := setup(t)
	defer teardown(t, dir)

	ti := singleTitle
	_, err := Save(dir, &ti)
	test.CheckError(t, err)
	newId := "tt12345"
	newKind := "changed"
	ti.Id = newId
	ti.Kind = newKind

	//overwrite
	_, err = Save(dir, &ti)
	test.CheckError(t, err)

	read, err := Read(filepath.Join(dir, ti.Id))
	test.CheckError(t, err)

	if read.Id != newId {
		t.Errorf("ids not matching: expected %s, but got %s", newId, read.Id)
	}
	if read.Kind != newKind {
		t.Errorf("kinds not matching: expected %s, but got %s", newKind, read.Kind)
	}
}

func TestSaveNilTargetInfo(t *testing.T) {
	_, err := Save(".", nil)
	if err == nil {
		t.Error("expected error when saving nil TargetInfo")
	}
}
