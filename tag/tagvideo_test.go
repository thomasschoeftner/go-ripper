package tag

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/thomasschoeftner/go-cli/commons"
	"github.com/thomasschoeftner/go-cli/config"
	"github.com/thomasschoeftner/go-cli/task"
	"github.com/thomasschoeftner/go-cli/test"
	"github.com/thomasschoeftner/go-ripper/files"
	"github.com/thomasschoeftner/go-ripper/metainfo"
	"github.com/thomasschoeftner/go-ripper/metainfo/video"
	"github.com/thomasschoeftner/go-ripper/ripper"
	"github.com/thomasschoeftner/go-ripper/targetinfo"
)

const expectedVideoExtension = "mp4"

func newTestVideoTagger(raiseError error) func(conf *ripper.AppConf, lazy bool, printf commons.FormatPrinter) (MovieTagger, EpisodeTagger, error) {
	return func(conf *ripper.AppConf, lazy bool, printf commons.FormatPrinter) (MovieTagger, EpisodeTagger, error) {
		var tagger = &testTagger{conf: conf, raiseError: raiseError}
		return tagger.TagMovie, tagger.TagEpisode, nil
	}
}

type testTagger struct {
	raiseError error
	conf       *ripper.AppConf
	inFile     string
	outFile    string
	id         string
	title      string
	year       string
	posterPath string
	series     string
	season     int
	episode    int
}

func (tagger *testTagger) TagMovie(inFile string, outFile string, id string, title string, year string, posterPath string) error {
	// fmt.Printf("tag movie %s with {id=%s, title=%s, year=%s, image=%s} -> write to %s\n", inFile, id, title, year, posterPath, outFile)
	tagger.inFile = inFile
	tagger.outFile = outFile
	tagger.id = id
	tagger.title = title
	tagger.year = year
	tagger.posterPath = posterPath
	return tagger.raiseError
}

func (tagger *testTagger) TagEpisode(inFile string, outFile string, id string, series string, season int, episode int, title string, year string, posterPath string) error {
	// fmt.Printf("tag episode %s with {id=%s, title=%s, year=%s, image=%s} -> write to %s\n", inFile, id, title, year, posterPath, outFile)
	tagger.inFile = inFile
	tagger.outFile = outFile
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
	emptyConf := &ripper.AppConf{}

	t.Run("expect error when no appropriate input inFile is found", func(t *testing.T) {
		assert := test.AssertOn(t)
		ti := targetinfo.NewMovie(files.WithExtension("movie", "avi"), "/some/dir", "movie-id")
		tagger := testTagger{conf: emptyConf}
		err := tagMovie(tagger.TagMovie, tagger.conf, ti, ti.File)
		assert.ExpectError("expected error when tagging movie without appropriate input inFile, but got none")(err)
	})

	t.Run("expect error when reading movie meta data fails", func(t *testing.T) {
		assert := test.AssertOn(t)
		ti := targetinfo.NewMovie(files.WithExtension("movie", expectedVideoExtension), "/some/dir", "movie-id")
		tagger := testTagger{conf: emptyConf}
		err := tagMovie(tagger.TagMovie, tagger.conf, ti, ti.File)
		assert.ExpectError("expected error when tagging movie without meta-info inFile, but got none")(err)
	})

	t.Run("invoke movie tagger with appropriate params and return tagger error", func(t *testing.T) {
		assert := test.AssertOn(t)

		dir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, dir)
		repoDir := filepath.Join(dir, "repo")

		mi := video.MovieMetaInfo{IdInfo: metainfo.IdInfo{Id: "movie-id"}, Title: "true art", Year: "1966", Poster: "/a/b/c/art.png"}
		metainfo.SaveMetaInfo(video.MovieFileName(repoDir, mi.Id), mi)
		ti := targetinfo.NewMovie(files.WithExtension("movie", expectedVideoExtension), "/some/dir", mi.Id)
		fileToProcess := files.WithExtension("some/other/file", expectedVideoExtension)

		outputDir := filepath.Join(dir, "output")
		conf := &ripper.AppConf{
			MetaInfoRepo: repoDir,
			Output: &ripper.OutputConfig{
				InvalidCharactersInFileName: "",
			},
			OutputDirectory: outputDir,
		}

		tagger := &testTagger{conf: conf}

		err := tagMovie(tagger.TagMovie, tagger.conf, ti, fileToProcess)
		assert.NotError(err)
		assert.StringsEqual(mi.Id, tagger.id)
		assert.StringsEqual(mi.Title, tagger.title)
		assert.StringsEqual(mi.Year, tagger.year)
		assert.StringsEqual(metainfo.ImageFileName(repoDir, mi.Id, files.GetExtension(mi.Poster)), tagger.posterPath)
		assert.StringsEqual(fileToProcess, tagger.inFile)
		assert.StringsEqual(filepath.Join(outputDir, files.WithExtension(mi.Title, expectedVideoExtension)), tagger.outFile)
	})
}

func TestTagEpisode(t *testing.T) {
	emptyConf := &ripper.AppConf{}

	t.Run("expect error when no appropriate input inFile is found", func(t *testing.T) {
		assert := test.AssertOn(t)
		ti := targetinfo.NewEpisode(files.WithExtension("movie", "avi"), "/some/dir", "episode-id", 4, 2, 2, 9)
		tagger := testTagger{conf: emptyConf}
		err := tagEpisode(tagger.TagEpisode, tagger.conf, ti, ti.File)
		assert.ExpectError("expected error when tagging episode without appropriate input inFile, but got none")(err)
	})

	t.Run("expect error when reading episode meta data fails", func(t *testing.T) {
		assert := test.AssertOn(t)
		dir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, dir)
		repoDir := filepath.Join(dir, "repo")

		seriesMi := video.SeriesMetaInfo{IdInfo: metainfo.IdInfo{"series-id"}, Title: "traffic education", Seasons: 9, Year: "2010", Poster: "/pic/of/a/car.jpeg"}
		metainfo.SaveMetaInfo(video.SeriesFileName(repoDir, seriesMi.Id), seriesMi)

		ti := targetinfo.NewEpisode(files.WithExtension("trafficeducation-s4e2", expectedVideoExtension), "/some/dir", seriesMi.Id, 2, 4, 4, 9)

		tagger := testTagger{conf: emptyConf}
		err := tagEpisode(tagger.TagEpisode, tagger.conf, ti, ti.File)
		assert.ExpectError("expected error when tagging episode without episode meta-info, but got none")(err)
	})

	t.Run("expect error when reading series meta data fails", func(t *testing.T) {
		assert := test.AssertOn(t)
		dir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, dir)
		repoDir := filepath.Join(dir, "repo")
		episodeMi := video.EpisodeMetaInfo{IdInfo: metainfo.IdInfo{"episode-id"}, Title: "crash boom", Season: 4, Episode: 2, Year: "2014"}
		metainfo.SaveMetaInfo(video.EpisodeFileName(repoDir, "series-id", episodeMi.Season, episodeMi.Episode), episodeMi)

		ti := targetinfo.NewEpisode(files.WithExtension("trafficeducation-s4e2", expectedVideoExtension), "/some/dir", "series-id", 2, 4, 4, 9)

		tagger := testTagger{conf: emptyConf}
		err := tagEpisode(tagger.TagEpisode, tagger.conf, ti, ti.File)
		assert.ExpectError("expected error when tagging episode without series meta-info, but got none")(err)
	})

	t.Run("invoke tagger with appropriate params and return tagger error", func(t *testing.T) {
		assert := test.AssertOn(t)
		dir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, dir)
		repoDir := filepath.Join(dir, "repo")

		seriesMi := video.SeriesMetaInfo{IdInfo: metainfo.IdInfo{"series-id"}, Title: "traffic education", Seasons: 9, Year: "2010", Poster: "/pic/of/a/car.jpeg"}
		episodeMi := video.EpisodeMetaInfo{IdInfo: metainfo.IdInfo{"episode-id"}, Title: "crash boom", Season: 4, Episode: 2, Year: "2014"}
		metainfo.SaveMetaInfo(video.SeriesFileName(repoDir, seriesMi.Id), seriesMi)
		metainfo.SaveMetaInfo(video.EpisodeFileName(repoDir, seriesMi.Id, episodeMi.Season, episodeMi.Episode), episodeMi)
		ti := targetinfo.NewEpisode(files.WithExtension("trafficeducation-s4e2", expectedVideoExtension), "/some/dir", seriesMi.Id, episodeMi.Season, episodeMi.Episode, episodeMi.Episode, 9)
		fileToProcess := files.WithExtension("some/other/file", expectedVideoExtension)

		outputDir := filepath.Join(dir, "output")
		conf := &ripper.AppConf{
			MetaInfoRepo: repoDir,
			Output: &ripper.OutputConfig{
				InvalidCharactersInFileName: "",
			},
			OutputDirectory: outputDir,
		}

		tagger := testTagger{conf: conf}

		err := tagEpisode(tagger.TagEpisode, tagger.conf, ti, fileToProcess)
		assert.NotError(err)
		assert.StringsEqual(seriesMi.Id, tagger.id)
		assert.StringsEqual(episodeMi.Title, tagger.title)
		assert.StringsEqual(episodeMi.Year, tagger.year)
		assert.IntsEqual(episodeMi.Season, tagger.season)
		assert.IntsEqual(episodeMi.Episode, tagger.episode)
		assert.StringsEqual(seriesMi.Title, tagger.series)
		assert.StringsEqual(metainfo.ImageFileName(repoDir, seriesMi.Id, files.GetExtension(seriesMi.Poster)), tagger.posterPath)
		assert.StringsEqual(fileToProcess, tagger.inFile)
		expectedFileName := files.WithExtension(fmt.Sprintf(templateEpisodeFilename, seriesMi.Title, episodeMi.Season, episodeMi.Episode, episodeMi.Title), expectedVideoExtension)
		assert.StringsEqual(filepath.Join(outputDir, seriesMi.Title, strconv.Itoa(episodeMi.Season), expectedFileName), tagger.outFile)
	})
}

func TestTagVideo(t *testing.T) {
	workDir := test.MkTempFolder(t)
	defer test.RmTempFolder(t, workDir)

	setup := func(t *testing.T, inputFile string) (*test.Assertion, task.Context, task.Job, *ripper.AppConf) {
		assert := test.AssertOn(t)

		TaggerFactories = make(map[string]TaggerFactory)

		conf := ripper.AppConf{}
		test.AssertOn(t).NotError(config.FromFile(&conf, "../go-ripper.conf", nil))
		conf.WorkDirectory = workDir
		conf.MetaInfoRepo = "./testdata/meta"
		conf.Tag.Video.Tagger = "test-tagger"
		conf.OutputDirectory = filepath.Join(workDir, "out")
		conf.Output.Video = "mp4"

		ctx := task.Context{
			Config:  &conf,
			Printf:  commons.Printf,
			RunLazy: true}
		job := map[string]string{ripper.JobField_Path: inputFile}
		return assert, ctx, job, &conf
	}

	movieTi := targetinfo.NewMovie("flick.mp4", "./testdata/in", "some-flick")
	movieFile := filepath.Join(movieTi.Folder, movieTi.File)
	mWorkDir, _ := ripper.GetWorkPathForTargetFolder(workDir, movieTi.GetFolder())
	targetinfo.Save(mWorkDir, movieTi)
	episodeTi := targetinfo.NewEpisode("part1.mp4", "./testdata/in", "part1", 3, 1, 1, 3)
	episodeFile := filepath.Join(episodeTi.Folder, episodeTi.File)
	eWorkDir, _ := ripper.GetWorkPathForTargetFolder(workDir, episodeTi.GetFolder())
	targetinfo.Save(eWorkDir, episodeTi)

	t.Run("return error when video tagger factory is unset/nil", func(t *testing.T) {
		assert, ctx, job, _ := setup(t, movieFile)
		TaggerFactories["non-existing-tagger"] = nil

		handlerFunc := TagVideo(ctx)
		jobs, err := handlerFunc(job)
		assert.ExpectError("expected error when calling TagVideo without Video-Tagger-Factory, but got none")(err)
		assert.IntsEqual(0, len(jobs))
	})

	t.Run("return error if target-info not found", func(t *testing.T) {
		assert, ctx, job, _ := setup(t, "./testdata/in/unknown.mp4")
		TaggerFactories["test-tagger"] = newTestVideoTagger(nil)

		handlerFunc := TagVideo(ctx)
		jobs, err := handlerFunc(job)
		assert.ExpectError("expected error when tagging video with missing targetinfo")(err)
		assert.IntsEqual(0, len(jobs))
	})

	t.Run("return error if no valid input file is found", func(t *testing.T) {
		assert, ctx, job, _ := setup(t, episodeFile)
		TaggerFactories["test-tagger"] = newTestVideoTagger(nil)

		ctx.Config.(*ripper.AppConf).Output.Video = "ogg" //expect ogg files as tagger input
		handlerFunc := TagVideo(ctx)
		jobs, err := handlerFunc(job)
		assert.ExpectError("expected error when no appropriate input file is found")(err)
		assert.IntsEqual(0, len(jobs))
	})

	t.Run("detect tagging error and discard evacuated file", func(t *testing.T) {
		assert, ctx, job, _ := setup(t, episodeFile)
		TaggerFactories["test-tagger"] = newTestVideoTagger(errors.New("expected test error"))

		handlerFunc := TagVideo(ctx)
		jobs, err := handlerFunc(job)
		assert.ExpectError("expected intentional test error")(err)
		assert.IntsEqual(0, len(jobs))
	})

	t.Run("invoke tagMovie for Movie", func(t *testing.T) {
		assert, ctx, job, conf := setup(t, movieFile)
		TaggerFactories["test-tagger"] = newTestVideoTagger(nil)
		conf.OutputDirectory += "-movies"

		handlerFunc := TagVideo(ctx)
		jobs, err := handlerFunc(job)
		assert.NotError(err)
		assert.IntsEqual(1, len(jobs))
	})

	t.Run("invoke tagEpisode for Episode", func(t *testing.T) {
		assert, ctx, job, conf := setup(t, episodeFile)
		TaggerFactories["test-tagger"] = newTestVideoTagger(nil)
		conf.OutputDirectory += "-episodes"

		handlerFunc := TagVideo(ctx)
		jobs, err := handlerFunc(job)
		assert.NotError(err)
		assert.IntsEqual(1, len(jobs))
	})
}
