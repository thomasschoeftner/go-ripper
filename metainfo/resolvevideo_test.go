package metainfo

import (
	"testing"
	"go-cli/test"
	"go-ripper/targetinfo"
	"go-cli/task"
	"go-cli/commons"
	"go-ripper/ripper"
	"go-cli/config"
	"errors"
	"path/filepath"
	"go-ripper/files"
	"fmt"
	"bytes"
)

var movie = targetinfo.NewMovie("movie.mp4", "/a/b", "tt123456")
var episode = targetinfo.NewEpisode("episode1.mp4", "/a/b/c", "tt654321", 3, 1, 4, 7)
var movieMi = MovieMetaInfo{BasicVideoMetaInfo: BasicVideoMetaInfo{Id: movie.Id, Title: "The awesome adventures of Sepp", Year: 2018}, Poster: "taaos.jpg"}
var seriesMi = SeriesMetaInfo{BasicVideoMetaInfo: BasicVideoMetaInfo{Id: episode.Id, Title: "a space oddity", Year: 2017}, Seasons: 3, Poster: "aso.png"}
var episodeMi = EpisodeMetaInfo{BasicVideoMetaInfo: BasicVideoMetaInfo{Id: episode.Id, Title: "attack of the raffgrns", Year: 2017}, Episode: 4, Season: 3}
var imageMi = map[string][]byte{movieMi.Poster : []byte{1,2,3,4}, seriesMi.Poster : []byte{5,6,7,8}}


type testVideoMetaInfoSource struct {
	movie *MovieMetaInfo
	series *SeriesMetaInfo
	//season *SeasonMetaInfo
	episode *EpisodeMetaInfo
	images map[string][]byte
	movieFetched bool
	seriesFetched bool
	episodeFetched bool
	imageFetched bool
}

func (f *testVideoMetaInfoSource) FetchMovieInfo(id string) (*MovieMetaInfo, error) {
	var err error
	var m *MovieMetaInfo
	if f.movie == nil {
		err = errors.New("test error - no movies defined")
	} else if f.movie.Id != id {
		err = errors.New("test error - movie not found")
	} else {
		m = f.movie
		f.movieFetched = true
	}
	return m, err
}

func (f *testVideoMetaInfoSource) FetchSeriesInfo(id string) (*SeriesMetaInfo, error) {
	var err error
	var s *SeriesMetaInfo
	if f.series == nil {
		err = errors.New("test error - no series defined")
	} else if f.series.Id != id {
		err = errors.New("test error - series not found")
	} else {
		s = f.series
		f.seriesFetched = true
	}
	return s, err
}

func (f *testVideoMetaInfoSource) FetchEpisodeInfo(id string, season int, episode int) (*EpisodeMetaInfo, error) {
	var e *EpisodeMetaInfo
	var err error
	if f.episode == nil {
		err = errors.New("test error - no episode defined")
	} else if f.episode.Id != id || f.episode.Season != season || f.episode.Episode != episode {
		err = errors.New("test error - episode not found")
	} else {
		e = f.episode
		f.episodeFetched = true
	}
	return e, err

}

func (f *testVideoMetaInfoSource) FetchImage(location string) ([]byte, error) {
	var img []byte
	var err error

	if f.images == nil {
		err = errors.New("test error - no images defined")
	} else {
		_, found := f.images[location]
		if !found {
			err = errors.New("test error - image not found")
		} else {
			img = f.images[location]
		}
		f.imageFetched = true
	}
	return img, err

}

func newVideoMetaInfoSource(movie *MovieMetaInfo, series *SeriesMetaInfo, /*season *SeasonMetaInfo,*/ episode *EpisodeMetaInfo, images map[string][]byte) *testVideoMetaInfoSource {
	return &testVideoMetaInfoSource{movie: movie, series: series, /*season: season,*/ episode: episode, images: images}
}


func TestNilVideoFactory(t *testing.T) {
	_, err := ResolveVideo(nil)
	test.AssertOn(t).ExpectError("expected error on initializing resolve handler with nil movie factory, but got none")(err)
}

const confJson = `
{
  "ignorePrefix" : ".",
  "workDirectory" : "${workdir}",
  "metaInfoRepo": "${repodir}",
  "resolve" : {
    "video" : {
    }
  }
}`


func TestResolveVideo(t *testing.T) {
	assert := test.AssertOn(t)

	dir := test.MkTempFolder(assert.T)
	defer test.RmTempFolder(t, dir)

	conf := &ripper.AppConf{}

	repoDir := filepath.ToSlash(filepath.Join(dir, "repo"))
	workDir := filepath.ToSlash(filepath.Join(dir, "work"))
	assert.NotError(config.FromString(conf, confJson,
		map[string]string {"repodir" : repoDir, "workdir" : workDir}))
	ctx := task.Context{nil, conf, commons.Printf, false}

	// create target info files
	targetInfos := []targetinfo.TargetInfo{movie, episode}
	for _, ti := range targetInfos {
		workDir, err := ripper.GetWorkPathForTargetFileFolder(conf.WorkDirectory, ti.GetFolder())
		assert.NotError(err)
		assert.NotError(files.CreateFolderStructure(workDir))
		assert.NotError(targetinfo.Save(workDir, ti))
	}

	// create jobs
	movieJob := task.Job{}.WithParam(ripper.JobField_Path, filepath.Join(movie.GetFolder(), movie.GetFile()))
	episodeJob := task.Job{}.WithParam(ripper.JobField_Path, filepath.Join(episode.GetFolder(), episode.GetFile()))

	t.Run("movie", func (t *testing.T) {
		// create resolve task
		miSource := newVideoMetaInfoSource(&movieMi, &seriesMi, &episodeMi, imageMi)
		handler, err := ResolveVideo(miSource)
		assert.NotError(err)
		resolve := handler(ctx)

		assert := test.AssertOn(t)
		resultJobs, err := resolve(movieJob)
		assert.NotError(err)

		assert.True("expected 1 result job")(1 == len(resultJobs))
		assert.StringsEqual(filepath.Join(movie.GetFolder(), movie.GetFile()), resultJobs[0][ripper.JobField_Path])
		assert.True("movie not fetched from meta-info source")(miSource.movieFetched)
		assert.True("image not fetched from meta-info source")(miSource.imageFetched)
	})

	t.Run("episode", func (t *testing.T) {
		// create resolve task
		miSource := newVideoMetaInfoSource(&movieMi, &seriesMi, &episodeMi, imageMi)
		handler, err := ResolveVideo(miSource)
		assert.NotError(err)
		resolve := handler(ctx)

		assert := test.AssertOn(t)
		miSource.episodeFetched = false
		resultJobs, err := resolve(episodeJob)
		assert.NotError(err)

		assert.True("expected 1 result job")(1 == len(resultJobs))
		assert.StringsEqual(filepath.Join(episode.GetFolder(), episode.GetFile()), resultJobs[0][ripper.JobField_Path])
		assert.True("series not fetched from meta-info source")(miSource.seriesFetched)
		assert.True("image not fetched from meta-info source")(miSource.imageFetched)
		assert.True("episode not fetched from meta-info source")(miSource.episodeFetched)
	})
}









func setup(assert *test.Assertion, confJson string, movieMetaInfo *MovieMetaInfo, seriesMetaInfo *SeriesMetaInfo,
	episodeMetaInfo *EpisodeMetaInfo, imageMetaInfo map[string][]byte) (string, *ripper.AppConf) {

	dir := test.MkTempFolder(assert.T)
	conf := &ripper.AppConf{}

	repoDir := filepath.ToSlash(filepath.Join(dir, "repo"))
	workDir := filepath.ToSlash(filepath.Join(dir, "work"))
	assert.NotError(config.FromString(conf, confJson,
		map[string]string {"repodir" : repoDir, "workdir" : workDir}))

	// create meta info files if passed
	if movieMetaInfo != nil {
		assert.NotError(SaveMetaInfoFile(MovieFileName(repoDir, movieMetaInfo.Id), movieMetaInfo))
	}
	if seriesMetaInfo != nil {
		assert.NotError(SaveMetaInfoFile(SeriesFileName(repoDir, seriesMetaInfo.Id), seriesMetaInfo))
	}
	if episodeMetaInfo != nil {
		assert.NotError(SaveMetaInfoFile(EpisodeFileName(repoDir, episodeMetaInfo.Id, episodeMetaInfo.Season, episodeMetaInfo.Episode), episodeMetaInfo))
	}
	if imageMetaInfo != nil {
		for f, image := range imageMetaInfo {
			println(f, image)
			var imgFileName string
			if movieMetaInfo != nil && movieMetaInfo.Poster == f {
				imgFileName = ImageFileName(repoDir, movieMetaInfo.Id, files.Extension(movieMetaInfo.Poster))
			} else if seriesMetaInfo != nil && seriesMetaInfo.Poster == f {
				imgFileName = ImageFileName(repoDir, seriesMetaInfo.Id, files.Extension(seriesMetaInfo.Poster))
			} else {
				assert.T.Fatalf("unknown poster name %s matches neither movie, nor series", f)
			}
			assert.NotError(SaveMetaInfoImage(imgFileName, image))
		}
	}

	return dir, conf
}

func teardown(t *testing.T, dir string) {
	test.RmTempFolder(t, dir)
}

func TestResolveMovie(t *testing.T) {
	t.Run("eager without pre-existing meta-info files", func(t *testing.T) {
		const lazy = false
		assert := test.AssertOn(t)

		dir, conf := setup(assert, confJson, nil, nil, nil, nil)
		defer teardown(t, dir)

		src := newVideoMetaInfoSource(&movieMi, &seriesMi, &episodeMi, imageMi)
		gotMi, gotImg := testResolveMovie(assert, src, conf, lazy, movie)
		assert.True("movie was not fetched")(src.movieFetched)
		assert.True("movie image was not fetched")(src.imageFetched)
		assert.False("series was fetched")(src.seriesFetched)
		assert.False("episode was fetched")(src.episodeFetched)

		assertMoviesEqual(assert, &movieMi, gotMi)
		assertImagesEqual(assert, imageMi[movieMi.Poster], gotImg)
	})

	t.Run("eager with pre-existing meta-info files", func(t *testing.T) {
		const lazy = false
		assert := test.AssertOn(t)

		//existingMovie := MovieMetaInfo{BasicVideoMetaInfo: BasicVideoMetaInfo{Id: movie.Id, Title: "an earlier awesome adventure of Sepp", Year: 2008}, Poster: "aeaaos.jpg"}
		//existingImage := map[string][]byte{movieMi.Poster : []byte{12,13,14,15}}
		//dir, conf := setup(assert, confJson, &existingMovie, nil, nil, existingImage)
		dir, conf := setup(assert, confJson, nil, nil, nil, nil)
		defer teardown(t, dir) //TODO activate

		src := newVideoMetaInfoSource(&movieMi, &seriesMi, &episodeMi, imageMi)
		gotMi, gotImg := testResolveMovie(assert, src, conf, lazy, movie)
		assert.True("movie was not fetched")(src.movieFetched)
		assert.True("movie image was not fetched")(src.imageFetched)
		assert.False("series was fetched")(src.seriesFetched)
		assert.False("episode was fetched")(src.episodeFetched)

		assertMoviesEqual(assert, &movieMi, gotMi)
		assertImagesEqual(assert, imageMi[movieMi.Poster], gotImg)
	})

	t.Run("lazy without pre-existing meta-info files", func(t *testing.T) {
		const lazy = true
		assert := test.AssertOn(t)

		dir, conf := setup(assert, confJson, nil, nil, nil, nil)
		defer teardown(t, dir)

		src := newVideoMetaInfoSource(&movieMi, &seriesMi, &episodeMi, imageMi)
		gotMi, gotImg := testResolveMovie(assert, src, conf, lazy, movie)
		assert.True("movie was not fetched")(src.movieFetched)
		assert.True("movie image was not fetched")(src.imageFetched)
		assert.False("series was fetched")(src.seriesFetched)
		assert.False("episode was fetched")(src.episodeFetched)

		assertMoviesEqual(assert, &movieMi, gotMi)
		assertImagesEqual(assert, imageMi[movieMi.Poster], gotImg)
	})

	t.Run("lazy with pre-existing meta-info files", func(t *testing.T) {
		const lazy = true
		assert := test.AssertOn(t)

		existingMovie := MovieMetaInfo{BasicVideoMetaInfo: BasicVideoMetaInfo{Id: movie.Id, Title: "an earlier awesome adventure of Sepp", Year: 2008}, Poster: "aeaaos.jpg"}
		existingImage := map[string][]byte{existingMovie.Poster : []byte{12,13,14,15}}
		dir, conf := setup(assert, confJson, &existingMovie, nil, nil, existingImage)
		defer teardown(t, dir)

		src := newVideoMetaInfoSource(&movieMi, &seriesMi, &episodeMi, imageMi)
		gotMi, gotImg := testResolveMovie(assert, src, conf, lazy, movie)
		assert.False("movie was unnecessarily fetched")(src.movieFetched)
		assert.False("movie image was unnecessarily fetched")(src.imageFetched)
		assert.False("series was fetched")(src.seriesFetched)
		assert.False("episode was fetched")(src.episodeFetched)

		assertMoviesEqual(assert, &existingMovie, gotMi)
		assertImagesEqual(assert, existingImage[existingMovie.Poster], gotImg)
	})

	t.Run("lazy with partially pre-existing meta-info files", func(t *testing.T) {
		const lazy = true
		assert := test.AssertOn(t)

		existingMovie := MovieMetaInfo{BasicVideoMetaInfo: BasicVideoMetaInfo{Id: movie.Id, Title: "an earlier awesome adventure of Sepp", Year: 2008}, Poster: "aeaaos.jpg"}
		missingImage := map[string][]byte{movieMi.Poster : []byte{4, 5, 6, 7}, existingMovie.Poster : []byte{12,13,14,15}}
		dir, conf := setup(assert, confJson, &existingMovie, nil, nil, nil)
		defer teardown(t, dir)

		src := newVideoMetaInfoSource(&movieMi, &seriesMi, &episodeMi, missingImage)
		gotMi, gotImg := testResolveMovie(assert, src, conf, lazy, movie)
		assert.False("movie was unnecessarily fetched")(src.movieFetched)
		assert.True("movie image was not fetched")(src.imageFetched)
		assert.False("series was fetched")(src.seriesFetched)
		assert.False("episode was fetched")(src.episodeFetched)

		assertMoviesEqual(assert, &existingMovie, gotMi)
		assertImagesEqual(assert, missingImage[existingMovie.Poster], gotImg)
	})
}

func assertMoviesEqual(assert *test.Assertion, expected *MovieMetaInfo, got *MovieMetaInfo) {
	if expected == nil || got == nil {
		assert.FailWith(fmt.Sprintf("did not expect nil for movie metainfo (expected %v, got %v)", expected, got))
	}
	assert.StringsEqual(expected.Id, got.Id)
	assert.StringsEqual(expected.Title, got.Title)
	assert.True(fmt.Sprintf("expected year %d, but got %d", expected.Year, got.Year))(expected.Year == got.Year)
	assert.StringsEqual(expected.Poster, got.Poster)
}

func assertImagesEqual(assert *test.Assertion, expected []byte, got []byte) {
	assert.True(fmt.Sprintf("expected image %v, but got %v", expected, got))(bytes.Equal(expected, got))
}

func testResolveMovie(assert *test.Assertion, miSource VideoMetaInfoSource, conf *ripper.AppConf, lazy bool, m *targetinfo.Movie) (*MovieMetaInfo, []byte) {
	assert.NotError(resolveMovie(miSource, m, conf, lazy))

	//read meta-info and resolve
	gotMi := &MovieMetaInfo{}
	assert.NotError(ReadMetaInfoFile(MovieFileName(conf.MetaInfoRepo, m.Id), gotMi))

	img, err := ReadMetaInfoImage(ImageFileName(conf.MetaInfoRepo, gotMi.Id, files.Extension(gotMi.Poster)))
	assert.NotError(err)
	return gotMi, img
}

func testResolve____TODO___(assert *test.Assertion, miSource VideoMetaInfoSource, conf *ripper.AppConf, lazy bool, ti *targetinfo.Episode) *EpisodeMetaInfo {
	assert.NotError(resolveEpisode(miSource, ti, conf, lazy))
	gotMi := &EpisodeMetaInfo{}

	//read meta-info and resolve
	assert.NotError(ReadMetaInfoFile(MovieFileName(conf.MetaInfoRepo, ti.Id), gotMi))

	return gotMi
}


func TestFindOrFetchImage(t *testing.T) {

}

func TestFindOrFetchSeries(t *testing.T) {

}

func TestFindOrFetchEpisode(t *testing.T) {

}
