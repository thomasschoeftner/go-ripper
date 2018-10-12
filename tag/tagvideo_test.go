package tag

import (
	"testing"
	"go-ripper/targetinfo"
	"go-ripper/ripper"
	"go-cli/test"
	"path/filepath"
	"io/ioutil"
	"os"
	"go-ripper/files"
	"fmt"
	"errors"
	"go-ripper/metainfo"
	"go-ripper/metainfo/video"
)

const expectedVideoExtension = "mp4"
func TestChooseInputFile(t *testing.T) {

	setup := func(tmpDir string, sourceExtension string, preprocessedExtension *string) (targetinfo.TargetInfo, string, string) {
		sourceDir := "/sepp/hat/gelbe/eier"
		sourceFile := "ripped"
		source := filepath.Join(sourceDir, files.WithExtension(sourceFile, sourceExtension))
		ti := targetinfo.NewMovie(files.WithExtension(sourceFile, sourceExtension), sourceDir, sourceFile)

		tiPath, _ := ripper.GetWorkPathForTargetFileFolder(tmpDir, ti.GetFolder())
		files.CreateFolderStructure(tiPath)
		var preprocessed string
		if preprocessedExtension != nil {
			//only create preprocessed (ie ripped) file if an extension is defined
			preprocessed = filepath.Join(tiPath, files.WithExtension(sourceFile, *preprocessedExtension))
			ioutil.WriteFile(preprocessed, []byte {1,2,3}, os.ModePerm)
		}
		return ti, source, preprocessed
	}

	t.Run("ripped file is available in workdir, source file is not a valid output file", func(t *testing.T) {
		workDir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, workDir)
		preprocessedExt := expectedVideoExtension
		ti, _, preprocessed := setup(workDir, "avi", &preprocessedExt)

		chosen, err := chooseInputFile(ti, workDir, expectedVideoExtension)
		if err != nil {
			t.Fatal(err)
		}
		test.AssertOn(t).StringsEqual(preprocessed, chosen)
	})

	t.Run("ripped file is available in workdir, but source file is already a valid output file", func(t *testing.T) {
		workDir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, workDir)
		preprocessedExt := expectedVideoExtension
		ti, _, preprocessed := setup(workDir, expectedVideoExtension, &preprocessedExt)

		chosen, err := chooseInputFile(ti, workDir, expectedVideoExtension)
		if err != nil {
			t.Fatal(err)
		}
		test.AssertOn(t).StringsEqual(preprocessed, chosen)
	})

	t.Run("ripped file is missing in workdir, but source file is already a valid output file", func(t *testing.T) {
		workDir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, workDir)
		ti, src, _ := setup(workDir, expectedVideoExtension, nil)

		chosen, err := chooseInputFile(ti, workDir, expectedVideoExtension)
		if err != nil {
			t.Fatal(err)
		}
		test.AssertOn(t).StringsEqual(src, chosen)
	})

	t.Run("ripped file is missing in workdir, source file is not a valid output file either", func(t *testing.T) {
		workDir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, workDir)
		ti, _, _ := setup(workDir, "avi", nil)

		_, err := chooseInputFile(ti, workDir, expectedVideoExtension)
		test.AssertOn(t).ExpectError("expected error when finding no suitable file - neither source does not have appropriate format, prepocessed file is missing")(err)
	})
}



type TestTagger struct {
	raiseError error
	file string
	id string
	title string
	year string
	posterPath string
	series string
	season int
	episode int
}

func (tagger *TestTagger) TagMovie(file string, id string, title string, year string, posterPath string) error {
	fmt.Printf("tag movie %s with {id=%s, title=%s, year=%s, image=%s}\n", file, id, title, year, posterPath)
	tagger.file = file
	tagger.id = id
	tagger.title = title
	tagger.year = year
	tagger.posterPath = posterPath
	return tagger.raiseError
}

func (tagger *TestTagger) TagEpisode(file string, id string, series string, season int, episode int, title string, year string, posterPath string) error {
	fmt.Printf("tag episode %s with {id=%s, title=%s, year=%s, image=%s}\n", file, id, title, year, posterPath)
	tagger.file = file
	tagger.id = id
	tagger.series = series
	tagger.season = season
	tagger.episode = episode
	tagger.title = title
	tagger.year = year
	tagger.posterPath = posterPath
	return tagger.raiseError
}


func TestTagMovie(t *testing.T) {
	t.Run("expect error when no appropriate input file is found", func(t *testing.T) {
		ti := targetinfo.NewMovie(files.WithExtension("movie", "avi"), "/some/dir", "movie-id")
		tagger := &TestTagger{}
		test.AssertOn(t).ExpectError("expected error when tagging movie without appropriate input file, but got none")(tagMovie(tagger, ti, "/work", "/repo", expectedVideoExtension))
	})

	t.Run("expect error when reading movie meta data fails", func(t *testing.T) {
		ti := targetinfo.NewMovie(files.WithExtension("movie", expectedVideoExtension), "/some/dir", "movie-id")
		tagger := &TestTagger{}
		test.AssertOn(t).ExpectError("expected error when tagging movie without meta-info file, but got none")(tagMovie(tagger, ti, "/work", "/repo", expectedVideoExtension))
	})

	t.Run("invoke tagger with appropriate params and return tagger error", func(t *testing.T) {
		dir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, dir)
		workDir := filepath.Join(dir, "work")
		repoDir := filepath.Join(dir, "repo")
		mi := video.MovieMetaInfo{IdInfo: metainfo.IdInfo{Id: "movie-id"}, Title: "true art", Year: "1966", Poster: "/a/b/c/art.png" }
		metainfo.SaveMetaInfo(video.MovieFileName(repoDir, mi.Id), mi)

		ti := targetinfo.NewMovie(files.WithExtension("movie", expectedVideoExtension), "/some/dir", mi.Id)
		expectedErr := errors.New("expected")
		tagger := &TestTagger{raiseError: expectedErr}
		assert := test.AssertOn(t)
		assert.ExpectError("expected error when using test tagger on movie, but got none")(tagMovie(tagger, ti, workDir, repoDir, expectedVideoExtension))
		assert.StringsEqual(mi.Id, tagger.id)
		assert.StringsEqual(mi.Title, tagger.title)
		assert.StringsEqual(mi.Year, tagger.year)
		assert.StringsEqual(metainfo.ImageFileName(repoDir, mi.Id, files.GetExtension(mi.Poster)), tagger.posterPath)
	})
}

func TestTagEpisode(t *testing.T) {
	t.Run("expect error when no appropriate input file is found", func(t *testing.T) {
		ti := targetinfo.NewEpisode(files.WithExtension("movie", "avi"), "/some/dir", "episode-id", 4, 2, 2, 9)
		tagger := &TestTagger{}
		test.AssertOn(t).ExpectError("expected error when tagging episode without appropriate input file, but got none")(tagEpisode(tagger, ti, "/work", "/repo", expectedVideoExtension))
	})

	t.Run("expect error when reading episode meta data fails", func(t *testing.T) {
		dir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, dir)
		workDir := filepath.Join(dir, "work")
		repoDir := filepath.Join(dir, "repo")

		seriesMi := video.SeriesMetaInfo{IdInfo: metainfo.IdInfo{"series-id"}, Title: "traffic education", Seasons: 9, Year: "2010", Poster: "/pic/of/a/car.jpeg" }
		metainfo.SaveMetaInfo(video.SeriesFileName(repoDir, seriesMi.Id), seriesMi)

		ti := targetinfo.NewEpisode(files.WithExtension("trafficeducation-s4e2", expectedVideoExtension), "/some/dir", seriesMi.Id, 2, 4, 4, 9)

		tagger := &TestTagger{}
		test.AssertOn(t).ExpectError("expected error when tagging episode without episode meta-info, but got none")(tagEpisode(tagger, ti, workDir, repoDir, expectedVideoExtension))
	})

	t.Run("expect error when reading series meta data fails", func(t *testing.T) {
		dir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, dir)
		workDir := filepath.Join(dir, "work")
		repoDir := filepath.Join(dir, "repo")
		episodeMi := video.EpisodeMetaInfo{IdInfo: metainfo.IdInfo{"episode-id"}, Title: "crash boom", Season: 4, Episode: 2, Year: "2014"}
		metainfo.SaveMetaInfo(video.EpisodeFileName(repoDir, "series-id", episodeMi.Season, episodeMi.Episode), episodeMi)

		ti := targetinfo.NewEpisode(files.WithExtension("trafficeducation-s4e2", expectedVideoExtension), "/some/dir", "series-id", 2, 4, 4, 9)

		tagger := &TestTagger{}
		test.AssertOn(t).ExpectError("expected error when tagging episode without series meta-info, but got none")(tagEpisode(tagger, ti, workDir, repoDir, expectedVideoExtension))
	})

	t.Run("invoke tagger with appropriate params and return tagger error", func(t *testing.T) {
		dir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, dir)
		workDir := filepath.Join(dir, "work")
		repoDir := filepath.Join(dir, "repo")

		seriesMi := video.SeriesMetaInfo{IdInfo: metainfo.IdInfo{"series-id"}, Title: "traffic education", Seasons: 9, Year: "2010", Poster: "/pic/of/a/car.jpeg" }
		episodeMi := video.EpisodeMetaInfo{IdInfo: metainfo.IdInfo{"episode-id"}, Title: "crash boom", Season: 4, Episode: 2, Year: "2014"}
		metainfo.SaveMetaInfo(video.SeriesFileName(repoDir, seriesMi.Id), seriesMi)
		metainfo.SaveMetaInfo(video.EpisodeFileName(repoDir, seriesMi.Id, episodeMi.Season, episodeMi.Episode), episodeMi)

		ti := targetinfo.NewEpisode(files.WithExtension("trafficeducation-s4e2", expectedVideoExtension), "/some/dir", seriesMi.Id, episodeMi.Season, episodeMi.Episode, episodeMi.Episode, 9)

		expectedErr := errors.New("expected")
		tagger := &TestTagger{raiseError: expectedErr}

		assert := test.AssertOn(t)
		assert.ExpectError("expected error when using test tagger on episode, but got none")(tagEpisode(tagger, ti, workDir, repoDir, expectedVideoExtension))
		assert.StringsEqual(seriesMi.Id, tagger.id)
		assert.StringsEqual(episodeMi.Title, tagger.title)
		assert.StringsEqual(episodeMi.Year, tagger.year)
		assert.IntsEqual(episodeMi.Season, tagger.season)
		assert.IntsEqual(episodeMi.Episode, tagger.episode)
		assert.StringsEqual(seriesMi.Title, tagger.series)
		assert.StringsEqual(metainfo.ImageFileName(repoDir, seriesMi.Id, files.GetExtension(seriesMi.Poster)), tagger.posterPath)
	})
}