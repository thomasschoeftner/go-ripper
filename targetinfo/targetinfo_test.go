package targetinfo

import (
	"testing"
	"os"
	"path/filepath"
	"go-cli/test"
)

var video = NewMovie("f.g", "/a/b/c", "test")
var episode = NewEpisode("f.g", "/a/b/c", "tt987654321", 3, 12, 12, 24)

func TestSaveJson(t *testing.T) {
	t.Run("save single video", func(t *testing.T) {
		ti := video
		dir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, dir)

		err := Save(dir, ti)
		test.CheckError(t, err)
		f := filepath.Join(dir, fileName(ti))

		info, err := os.Stat(f)
		test.CheckError(t, err)

		if info.Size() == 0 {
			t.Errorf("target file \"%s\" must not be empty", f)
		}
	})


	t.Run("save series episode", func(t *testing.T) {
		ti := episode
		dir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, dir)

		err := Save(dir, ti)
		test.CheckError(t, err)

		f := filepath.Join(dir, fileName(ti))
		info, err := os.Stat(f)
		test.CheckError(t, err)

		if info.Size() == 0 {
			t.Errorf("target file \"%s\" must not be empty", f)
		}
	})
}

func TestReadJson(t *testing.T) {
	dir := test.MkTempFolder(t)
	defer test.RmTempFolder(t, dir)

	err := Save(dir, episode)
	test.CheckError(t, err)

	read, err := read(dir, episode.File)
	test.CheckError(t, err)

	readEpisode := read.(*Episode)
	if *episode != *readEpisode {
		t.Errorf("targetinfo does not match:\n  to json   %v\n  from json %v", *episode, *readEpisode)
	}
}


func TestSaveNilTargetInfo(t *testing.T) {
	err := Save(".", nil)
	if err == nil {
		t.Error("expected error when saving nil TargetInfo")
	}
}

func TestCorrectFileNames(t *testing.T) {
	fname := fileName(video)
	expectedFilName := video.GetFile() + "." + targetinfo_file_extension
	if expectedFilName != fname {
		t.Errorf("video filename %s does not match expected file name %s", fname, expectedFilName)
	}
}
