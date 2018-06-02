package targetinfo

import (
	"testing"
	"io/ioutil"
	"os"
	"path/filepath"
	"go-cli/test"
	"fmt"
	"strings"
	"strconv"
)

var video = NewVideo("f.g", "/a/b/c", "test")
var episode = NewEpisode("f.g", "/a/b/c", "tt987654321", 3, 12, 12, 24)

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
	t.Run("save single video", func(t *testing.T) {
		ti := video
		dir := setup(t)
		defer teardown(t, dir)
		_, err := Save(dir, ti)
		test.CheckError(t, err)
		f := filepath.Join(dir, TARGETINFO_VIDEO + ti.Id)
		info, err := os.Stat(f)
		test.CheckError(t, err)

		if info.Size() == 0 {
			t.Errorf("target file \"%s\" must not be empty", f)
		}
	})


	t.Run("save series episode", func(t *testing.T) {
		ti := episode
		dir := setup(t)
		defer teardown(t, dir)
		_, err := Save(dir, ti)
		test.CheckError(t, err)

		fname := fmt.Sprintf("%s%s.%d.%d", TARGETINFO_EPISODE, ti.Id, ti.Season, ti.Episode)
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

	_, err := Save(dir, episode)
	test.CheckError(t, err)

	read, err := Read(filepath.Join(dir, episode.fileName()))
	test.CheckError(t, err)

	readEpisode := read.(*Episode)
	if *episode != *readEpisode {
		t.Errorf("targetinfo does not match:\n  to json   %v\n  from json %v", *episode, *readEpisode)
	}
}


func TestSaveNilTargetInfo(t *testing.T) {
	_, err := Save(".", nil)
	if err == nil {
		t.Error("expected error when saving nil TargetInfo")
	}
}

func TestCorrectFileNames(t *testing.T) {
	t.Run("video file name", func(t *testing.T) {
		media := video
		fname := media.fileName()
		if !strings.HasPrefix(fname, TARGETINFO_VIDEO) {
			t.Errorf("video filename %s is missing prefix %s", fname, TARGETINFO_VIDEO)
		}
		id := strings.Replace(fname, TARGETINFO_VIDEO, "", 1)
		if media.Id != id {
			t.Errorf("video filename %s contains incorrect id %s - expected id %s", fname, id, media.Id)
		}
	})

	t.Run("episode file name", func(t *testing.T) {
		media := episode
		fname := media.fileName()
		if !strings.HasPrefix(fname, TARGETINFO_EPISODE) {
			t.Errorf("episode filename %s is missing prefix %s", fname, TARGETINFO_EPISODE)
		}
		fragments := strings.Split(fname, ".")
		id := strings.Replace(fragments[0], TARGETINFO_EPISODE, "", 1)
		if media.Id != id {
			t.Errorf("episode filename %s contains incorrect id %s - expected %s", fname, id, media.Id)
		}
		season, err := strconv.Atoi(fragments[1])
		test.CheckError(t, err)
		if media.Season != season {
			t.Errorf("episode filename %s contains incorrect season# %d - expected %d", fname, season, media.Season)
		}

		episode, err := strconv.Atoi(fragments[2])
		test.CheckError(t, err)
		if media.Episode != episode {
			t.Errorf("expisode filename %s contains incorrect episode# %d - expected %d", fname, episode, media.Episode)
		}
	})
}