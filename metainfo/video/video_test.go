package video

import (
	"testing"
	"go-cli/test"
	"fmt"
	"go-ripper/metainfo"
)

func TestGetType(t *testing.T) {
	assert := test.AssertOn(t)
	assert.StringsEqual(META_INFO_TYPE_MOVIE, (&MovieMetaInfo{}).GetType())
	assert.StringsEqual(META_INFO_TYPE_SERIES, (&SeriesMetaInfo{}).GetType())
	assert.StringsEqual(META_INFO_TYPE_EPISODE, (&EpisodeMetaInfo{}).GetType())
}

func TestGetFileNames(t *testing.T) {
	dir := "a/b/c"
	id := "tt678"

	t.Run("movie file name", func(t *testing.T) {
		expected := fmt.Sprintf("%s/%s/%s.%s", dir, SUBDIR_MOVIES, id, metainfo.METAINF_FILE_EXT)
		fName := MovieFileName(dir, id)
		test.AssertOn(t).StringsEqual(expected, fName)
	})

	t.Run("series file name", func(t *testing.T) {
		expected := fmt.Sprintf("%s/%s/%s.%s", dir, SUBDIR_SERIES, id, metainfo.METAINF_FILE_EXT)
		fName := SeriesFileName(dir, id)
		test.AssertOn(t).StringsEqual(expected, fName)
	})

	t.Run("episode file name", func(t *testing.T) {
		season := 7
		episode := 21

		expected := fmt.Sprintf("%s/%s/%s.%d.%d.%s", dir, SUBDIR_SERIES, id, season, episode, metainfo.METAINF_FILE_EXT)
		fName := EpisodeFileName(dir, id, season, episode)
		test.AssertOn(t).StringsEqual(expected, fName)
	})
}