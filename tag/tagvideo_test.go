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
	"go-ripper/metainfo"
	"go-ripper/metainfo/video"
	"strconv"
	"go-cli/commons"
	"go-cli/task"
	"go-cli/config"
	"errors"
)

const expectedVideoExtension = "mp4"
func TestFindInputOutputFile(t *testing.T) {
	setup := func(tmpDir string, sourceExtension string, preprocessedExtension string, preprocessedExists bool) (targetinfo.TargetInfo, string, string) {
		sourceDir := "/sepp/hat/gelbe/eier"
		sourceFile := "ripped"
		source := filepath.Join(sourceDir, files.WithExtension(sourceFile, sourceExtension))
		ti := targetinfo.NewMovie(files.WithExtension(sourceFile, sourceExtension), sourceDir, sourceFile)

		workFolder, _ := ripper.GetWorkPathForTargetFolder(tmpDir, ti.GetFolder())
		files.CreateFolderStructure(workFolder)
		preprocessed := filepath.Join(workFolder, files.WithExtension(sourceFile, preprocessedExtension))
		if preprocessedExists {
			//only create preprocessed (ie ripped) inFile if an extension is defined
			ioutil.WriteFile(preprocessed, []byte {1,2,3}, os.ModePerm)
		}
		return ti, source, preprocessed
	}

	t.Run("ripped inFile is available in workdir, source inFile is not a valid output inFile", func(t *testing.T) {
		workDir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, workDir)
		ti, _, preprocessed := setup(workDir, "avi", expectedVideoExtension, true)

		in, out, err := findInputOutputFiles(ti, workDir, expectedVideoExtension)
		if err != nil {
			t.Fatal(err)
		}
		test.AssertOn(t).StringsEqual(preprocessed, in)
		test.AssertOn(t).StringsEqual(preprocessed, out)
	})

	t.Run("ripped inFile is available in workdir, but source inFile is already a valid output inFile", func(t *testing.T) {
		workDir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, workDir)
		ti, _, preprocessed := setup(workDir, expectedVideoExtension, expectedVideoExtension, true)

		in, out, err := findInputOutputFiles(ti, workDir, expectedVideoExtension)
		if err != nil {
			t.Fatal(err)
		}
		test.AssertOn(t).StringsEqual(preprocessed, in)
		test.AssertOn(t).StringsEqual(preprocessed, out)
	})

	t.Run("ripped inFile is missing in workdir, but source inFile is already a valid output inFile", func(t *testing.T) {
		workDir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, workDir)
		ti, source, preprocessed := setup(workDir, expectedVideoExtension, expectedVideoExtension, false)

		in, out, err := findInputOutputFiles(ti, workDir, expectedVideoExtension)
		if err != nil {
			t.Fatal(err)
		}

		test.AssertOn(t).StringsEqual(source, in)
		test.AssertOn(t).StringsEqual(preprocessed, out)
	})

	t.Run("ripped inFile is missing in workdir, source inFile is not a valid output inFile either", func(t *testing.T) {
		workDir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, workDir)
		ti, _, _ := setup(workDir, "avi", expectedVideoExtension, false)

		_, _, err := findInputOutputFiles(ti, workDir, expectedVideoExtension)
		test.AssertOn(t).ExpectError("expected error when finding no suitable inFile - neither source does not have appropriate format, prepocessed inFile is missing")(err)
	})
}


func newTestVideoTagger(raiseError error) func(conf *ripper.TagConfig, lazy bool, printf commons.FormatPrinter) (VideoTagger, error) {
	return func(conf *ripper.TagConfig, lazy bool, printf commons.FormatPrinter) (VideoTagger, error) {
		return &testTagger{raiseError: raiseError}, nil
	}
}

type testTagger struct {
	raiseError error
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
	fmt.Printf("tag movie %s with {id=%s, title=%s, year=%s, image=%s} -> write to %s\n", inFile, id, title, year, posterPath, outFile)
	tagger.inFile = inFile
	tagger.outFile = outFile
	tagger.id = id
	tagger.title = title
	tagger.year = year
	tagger.posterPath = posterPath
	return tagger.raiseError
}

func (tagger *testTagger) TagEpisode(inFile string, outFile string, id string, series string, season int, episode int, title string, year string, posterPath string) error {
	fmt.Printf("tag episode %s with {id=%s, title=%s, year=%s, image=%s} -> write to %s\n", inFile, id, title, year, posterPath, outFile)
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
	t.Run("expect error when no appropriate input inFile is found", func(t *testing.T) {
		assert := test.AssertOn(t)
		ti := targetinfo.NewMovie(files.WithExtension("movie", "avi"), "/some/dir", "movie-id")
		tagger := &testTagger{}
		dest, err := tagMovie(tagger, ti,"/repo", ti.File)
		assert.ExpectError("expected error when tagging movie without appropriate input inFile, but got none")(err)
		assert.IntsEqual(0, len(dest))
	})

	t.Run("expect error when reading movie meta data fails", func(t *testing.T) {
		assert := test.AssertOn(t)
		ti := targetinfo.NewMovie(files.WithExtension("movie", expectedVideoExtension), "/some/dir", "movie-id")
		tagger := &testTagger{}
		dest, err := tagMovie(tagger, ti,"/repo", ti.File)
		assert.ExpectError("expected error when tagging movie without meta-info inFile, but got none")(err)
		assert.IntsEqual(0, len(dest))
	})

	t.Run("invoke movie tagger with appropriate params and return tagger error", func(t *testing.T) {
		assert := test.AssertOn(t)

		dir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, dir)
		repoDir := filepath.Join(dir, "repo")
		mi := video.MovieMetaInfo{IdInfo: metainfo.IdInfo{Id: "movie-id"}, Title: "true art", Year: "1966", Poster: "/a/b/c/art.png" }
		metainfo.SaveMetaInfo(video.MovieFileName(repoDir, mi.Id), mi)

		ti := targetinfo.NewMovie(files.WithExtension("movie", expectedVideoExtension), "/some/dir", mi.Id)
		fileToProcess := files.WithExtension("some/other/file", expectedVideoExtension)

		tagger := &testTagger{}
		dest, err := tagMovie(tagger, ti,repoDir, fileToProcess)
		assert.NotError(err)
		assert.StringsEqual(mi.Id, tagger.id)
		assert.StringsEqual(mi.Title, tagger.title)
		assert.StringsEqual(mi.Year, tagger.year)
		assert.StringsEqual(metainfo.ImageFileName(repoDir, mi.Id, files.GetExtension(mi.Poster)), tagger.posterPath)
		assert.StringsEqual(fileToProcess, tagger.inFile)
		assert.StringsEqual(fileToProcess, tagger.outFile)
		assert.StringSlicesEqual([]string{files.WithExtension(mi.Title, expectedVideoExtension)}, dest)
	})
}

func TestTagEpisode(t *testing.T) {
	t.Run("expect error when no appropriate input inFile is found", func(t *testing.T) {
		assert := test.AssertOn(t)
		ti := targetinfo.NewEpisode(files.WithExtension("movie", "avi"), "/some/dir", "episode-id", 4, 2, 2, 9)
		tagger := &testTagger{}
		dest, err := tagEpisode(tagger, ti,"/repo", ti.File)
		assert.ExpectError("expected error when tagging episode without appropriate input inFile, but got none")(err)
		assert.IntsEqual(0, len(dest))
	})

	t.Run("expect error when reading episode meta data fails", func(t *testing.T) {
		assert := test.AssertOn(t)
		dir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, dir)
		repoDir := filepath.Join(dir, "repo")

		seriesMi := video.SeriesMetaInfo{IdInfo: metainfo.IdInfo{"series-id"}, Title: "traffic education", Seasons: 9, Year: "2010", Poster: "/pic/of/a/car.jpeg" }
		metainfo.SaveMetaInfo(video.SeriesFileName(repoDir, seriesMi.Id), seriesMi)

		ti := targetinfo.NewEpisode(files.WithExtension("trafficeducation-s4e2", expectedVideoExtension), "/some/dir", seriesMi.Id, 2, 4, 4, 9)

		tagger := &testTagger{}
		dest, err := tagEpisode(tagger, ti,repoDir, ti.File)
		assert.ExpectError("expected error when tagging episode without episode meta-info, but got none")(err)
		assert.IntsEqual(0, len(dest))
	})

	t.Run("expect error when reading series meta data fails", func(t *testing.T) {
		assert := test.AssertOn(t)
		dir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, dir)
		repoDir := filepath.Join(dir, "repo")
		episodeMi := video.EpisodeMetaInfo{IdInfo: metainfo.IdInfo{"episode-id"}, Title: "crash boom", Season: 4, Episode: 2, Year: "2014"}
		metainfo.SaveMetaInfo(video.EpisodeFileName(repoDir, "series-id", episodeMi.Season, episodeMi.Episode), episodeMi)

		ti := targetinfo.NewEpisode(files.WithExtension("trafficeducation-s4e2", expectedVideoExtension), "/some/dir", "series-id", 2, 4, 4, 9)

		tagger := &testTagger{}
		dest, err := tagEpisode(tagger, ti, repoDir, ti.File)
		assert.ExpectError("expected error when tagging episode without series meta-info, but got none")(err)
		assert.IntsEqual(0, len(dest))
	})

	t.Run("invoke tagger with appropriate params and return tagger error", func(t *testing.T) {
		assert := test.AssertOn(t)
		dir := test.MkTempFolder(t)
		defer test.RmTempFolder(t, dir)
		repoDir := filepath.Join(dir, "repo")

		seriesMi := video.SeriesMetaInfo{IdInfo: metainfo.IdInfo{"series-id"}, Title: "traffic education", Seasons: 9, Year: "2010", Poster: "/pic/of/a/car.jpeg" }
		episodeMi := video.EpisodeMetaInfo{IdInfo: metainfo.IdInfo{"episode-id"}, Title: "crash boom", Season: 4, Episode: 2, Year: "2014"}
		metainfo.SaveMetaInfo(video.SeriesFileName(repoDir, seriesMi.Id), seriesMi)
		metainfo.SaveMetaInfo(video.EpisodeFileName(repoDir, seriesMi.Id, episodeMi.Season, episodeMi.Episode), episodeMi)

		ti := targetinfo.NewEpisode(files.WithExtension("trafficeducation-s4e2", expectedVideoExtension), "/some/dir", seriesMi.Id, episodeMi.Season, episodeMi.Episode, episodeMi.Episode, 9)

		tagger := &testTagger{}
		fileToProcess := files.WithExtension("some/other/file", expectedVideoExtension)

		dest, err := tagEpisode(tagger, ti, repoDir, fileToProcess)
		assert.NotError(err)
		assert.StringsEqual(seriesMi.Id, tagger.id)
		assert.StringsEqual(episodeMi.Title, tagger.title)
		assert.StringsEqual(episodeMi.Year, tagger.year)
		assert.IntsEqual(episodeMi.Season, tagger.season)
		assert.IntsEqual(episodeMi.Episode, tagger.episode)
		assert.StringsEqual(seriesMi.Title, tagger.series)
		assert.StringsEqual(metainfo.ImageFileName(repoDir, seriesMi.Id, files.GetExtension(seriesMi.Poster)), tagger.posterPath)
		assert.StringsEqual(fileToProcess, tagger.inFile)
		assert.StringsEqual(fileToProcess, tagger.outFile)
		expectedFileName :=  files.WithExtension(fmt.Sprintf(templateEpisodeFilename, seriesMi.Title, episodeMi.Season, episodeMi.Episode, episodeMi.Title), expectedVideoExtension)
		assert.StringSlicesEqual([]string{seriesMi.Title, strconv.Itoa(episodeMi.Season), expectedFileName}, dest)
	})
}

func TestTagVideo(t *testing.T) {
	workDir := test.MkTempFolder(t)
	defer test.RmTempFolder(t, workDir)

	setup := func(t *testing.T, inputFile string) (*test.Assertion, task.Context, task.Job) {
		assert := test.AssertOn(t)
		conf := ripper.AppConf{}
		test.AssertOn(t).NotError(config.FromFile(&conf, "../go-ripper.conf", nil))
		conf.WorkDirectory = workDir
		conf.MetaInfoRepo = "./testdata/meta"
		ctx := task.Context{
			Config: &conf,
			Printf: commons.Printf,
			RunLazy: true,
			OutputDir: filepath.Join(workDir, "out")}
		job := map[string]string {ripper.JobField_Path : inputFile}
		return assert, ctx, job
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
		NewVideoTagger = nil
		assert, ctx, job := setup(t, movieFile)
		handlerFunc := TagVideo(ctx)
		jobs, err := handlerFunc(job)
		assert.ExpectError("expected error when calling TagVideo without Video-Tagger-Factory, but got none")(err)
		assert.IntsEqual(0, len(jobs))
	})

	t.Run("return error if target-info not found", func(t *testing.T) {
		NewVideoTagger = newTestVideoTagger(nil)
		assert, ctx, job := setup(t,"./testdata/in/unknown.mp4")
		handlerFunc := TagVideo(ctx)
		jobs, err := handlerFunc(job)
		assert.ExpectError("expected error when tagging video with missing targetinfo")(err)
		assert.IntsEqual(0, len(jobs))
	})

	t.Run("return error if no valid input file is found", func(t *testing.T) {
		NewVideoTagger = newTestVideoTagger(nil)
		assert, ctx, job := setup(t, episodeFile)
		ctx.Config.(*ripper.AppConf).Output.Video = "ogg" //expect ogg files as tagger input
		handlerFunc := TagVideo(ctx)
		jobs, err := handlerFunc(job)
		assert.ExpectError("expected error when no appropriate input file is found")(err)
		assert.IntsEqual(0, len(jobs))
	})

	t.Run("detetect tagging error and discard evacuated file", func(t *testing.T) {
		NewVideoTagger = newTestVideoTagger(errors.New("expected test error"))
		assert, ctx, job := setup(t, episodeFile)
		handlerFunc := TagVideo(ctx)
		jobs, err := handlerFunc(job)
		assert.ExpectError("expected intentional test error")(err)
		assert.IntsEqual(0, len(jobs))
		tmpDir := filepath.Join(workDir, tagTempFolder)
		assert.TrueNotErrorf("expected \"%s\" folder to exist", tagTempFolder)(files.Exists(tmpDir))

		containedFiles, err := files.GetDirectoryContents(tmpDir)
		assert.NotError(err)
		assert.IntsEqual(0, len(containedFiles))
	})

	t.Run("invoke tagMovie for Movie", func(t *testing.T) {
		NewVideoTagger = newTestVideoTagger(nil)
		assert, ctx, job := setup(t, movieFile)
		ctx.OutputDir = ctx.OutputDir + "-movies"
		handlerFunc := TagVideo(ctx)
		jobs, err := handlerFunc(job)
		assert.NotError(err)
		assert.IntsEqual(1, len(jobs))
		assert.TrueNotErrorf("output directory \"%s\" should exist", ctx.OutputDir)(files.Exists(ctx.OutputDir))
		containedFiles, _ := files.GetDirectoryContents(ctx.OutputDir)
		assert.IntsEqual(1, len(containedFiles))
		assert.StringsEqual("some flick.mp4", containedFiles[0])
	})

	t.Run("invoke tagEpisode for Episode", func(t *testing.T) {
		NewVideoTagger = newTestVideoTagger(nil)
		assert, ctx, job := setup(t, episodeFile)
		ctx.OutputDir = ctx.OutputDir + "-episodes"
		handlerFunc := TagVideo(ctx)
		jobs, err := handlerFunc(job)
		assert.NotError(err)
		assert.IntsEqual(1, len(jobs))

		outputDir := filepath.Join(ctx.OutputDir, "in many parts", "3")
		assert.TrueNotErrorf("output directory \"%s\" should exist", outputDir)(files.Exists(outputDir))
		containedFiles, _ := files.GetDirectoryContents(outputDir)
		assert.IntsEqual(1, len(containedFiles))
		assert.StringsEqual("in many parts-s03e01-the real 1st part.mp4", containedFiles[0])
	})
}
