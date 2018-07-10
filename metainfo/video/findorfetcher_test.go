package video

import (
	"go-cli/test"
	"go-ripper/ripper"
	"path/filepath"
	"go-cli/config"
	"go-ripper/files"
	"testing"
	"fmt"
	"bytes"
	"go-ripper/metainfo"
)

func setupFindOrFetcher(assert *test.Assertion, movie *MovieMetaInfo, series *SeriesMetaInfo, episode *EpisodeMetaInfo, images map[string][]byte) (string, *ripper.AppConf) {
		dir := test.MkTempFolder(assert.T)
		conf := &ripper.AppConf{}

		repoDir := filepath.ToSlash(filepath.Join(dir, "repo"))
		workDir := filepath.ToSlash(filepath.Join(dir, "work"))
		assert.NotError(config.FromString(conf, confJson,
			map[string]string {"repodir" : repoDir, "workdir" : workDir}))

		// create meta info files if passed
		if movie != nil {
			assert.NotError(metainfo.SaveMetaInfo(MovieFileName(repoDir, movie.Id), movie))
		}
		if series != nil {
			assert.NotError(metainfo.SaveMetaInfo(SeriesFileName(repoDir, series.Id), series))
		}
		if episode != nil {
			assert.NotError(metainfo.SaveMetaInfo(EpisodeFileName(repoDir, episode.Id, episode.Season, episode.Episode), episode))
		}
		if images != nil {
			for f, image := range images {
				var imgFileName string
				if movie != nil && movie.Poster == f {
					imgFileName = metainfo.ImageFileName(repoDir, movie.Id, files.Extension(movie.Poster))
				} else if series != nil && series.Poster == f {
					imgFileName = metainfo.ImageFileName(repoDir, series.Id, files.Extension(series.Poster))
				} else {
					assert.T.Fatalf("unknown poster name %s matches neither movie, nor series", f)
				}
				assert.NotError(metainfo.SaveImage(imgFileName, image))
			}
		}

		return dir, conf
}

func teardownFindOrFetcher(assert *test.Assertion, dir string) {
	test.RmTempFolder(assert.T, dir)
}

func TestFindOrFetchMovie(t *testing.T) {
	testFindOrFetch := func(lazy bool, preexisingMovie *MovieMetaInfo, noNeedToResolve bool) func(*testing.T) {
		return func(t *testing.T) {
			assert := test.AssertOn(t)
			dir, conf := setupFindOrFetcher(assert, preexisingMovie, nil, nil, nil)
			defer teardownFindOrFetcher(assert, dir)

			miSrc := newVideoMetaInfoSource(&movieMi, &seriesMi, &episodeMi, imageMi)
			fof := findOrFetch(miSrc, conf, lazy)
			gotMovie, err := fof.movie(movieTi)

			assert.NotError(err)
			assert.True("movie meta-info image was fetched")(0 == len(miSrc.imagesFetched))
			assert.False("series was fetched")(miSrc.seriesFetched)
			assert.False("episode meta-info was fetched")(miSrc.episodeFetched)

			if noNeedToResolve {
				assert.False("movie meta-info was unnecessarily fetched")(miSrc.movieFetched)
				assertMoviesEqual(assert, preexisingMovie, gotMovie)
			} else {
				assert.True("movie meta-info was not fetched")(miSrc.movieFetched)
				assertMoviesEqual(assert, &movieMi, gotMovie)
			}
		}
	}
	existingMovie := MovieMetaInfo{IdInfo: metainfo.IdInfo{Id: movieTi.Id}, Title: "an earlier awesome adventure of Sepp", Year: "2008", Poster: "aeaaos.jpg"}


	t.Run("eager without pre-existing meta-info files", testFindOrFetch(false, nil, false))
	t.Run("lazy without pre-existing meta-info files", testFindOrFetch(true, nil, false))
	t.Run("eager with pre-existing meta-info files", testFindOrFetch(false, &existingMovie, false))
	t.Run("lazy with pre-existing meta-info files", testFindOrFetch(true, &existingMovie, true))
}

func TestFindOrFetchImage(t *testing.T) {
	testFindOrFetch := func(lazy bool, existingMovie *MovieMetaInfo, existingImages map[string][]byte, noNeedToResolve bool) func(t *testing.T){
		return func (t *testing.T) {
			assert := test.AssertOn(t)
			dir, conf := setupFindOrFetcher(assert, existingMovie, nil, nil, existingImages)
			defer teardownFindOrFetcher(assert, dir)

			miSrc := newVideoMetaInfoSource(&movieMi, &seriesMi, &episodeMi, imageMi)
			fof := findOrFetch(miSrc, conf, lazy)

			imgFileName := metainfo.ImageFileName(conf.MetaInfoRepo, movieMi.Id, files.Extension(movieMi.Poster))

			assert.NotError(fof.image(movieMi.Id, movieMi.Poster))
			img, err := metainfo.ReadImage(imgFileName)
			assert.NotError(err)

			if noNeedToResolve {
				assertImagesEqual(assert, existingImages[existingMovie.Poster], img)
				assert.True("image was unnecessarily fetched")(0 == len(miSrc.imagesFetched))
			} else {
				assertImagesEqual(assert, imageMi[movieMi.Poster], img)
				assert.True("image was not fetched")(1 == len(miSrc.imagesFetched))
				assert.StringsEqual(movieMi.Poster, miSrc.imagesFetched[0])
			}
		}
	}
	existingMovie := MovieMetaInfo{IdInfo: metainfo.IdInfo{Id: movieTi.Id}, Title: "an earlier awesome adventure of Sepp", Year: "2008", Poster: "aeaaos.jpg"}
	existingImages := map[string][]byte{existingMovie.Poster: {12, 13, 14, 15}}

	t.Run("eager without pre-existing image", testFindOrFetch(false, nil, nil, false))
	t.Run("lazy without pre-existing image", testFindOrFetch(true, nil, nil, false))
	t.Run("eager with pre-existing image", testFindOrFetch(false, &existingMovie, existingImages, false))
	t.Run("lazy with pre-existing image", testFindOrFetch(true, &existingMovie, existingImages, true))
}

func TestFindOrFetchSeries(t *testing.T) {
	testFindOrFetch := func(lazy bool, existingSeries *SeriesMetaInfo, noNeedToResolve bool) func(t *testing.T) {
		return func(t *testing.T) {
			assert := test.AssertOn(t)
			dir, conf := setupFindOrFetcher(assert, nil, existingSeries, nil, nil)
			defer teardownFindOrFetcher(assert, dir)

			miSrc := newVideoMetaInfoSource(&movieMi, &seriesMi, &episodeMi, imageMi)
			fof := findOrFetch(miSrc, conf, lazy)

			gotSeries, err := fof.series(episodeTi)
			assert.NotError(err)
			assert.True("series meta-info image was fetched")(0 == len(miSrc.imagesFetched))
			assert.False("episode was fetched")(miSrc.episodeFetched)
			assert.False("movie meta-info was fetched")(miSrc.movieFetched)

			if noNeedToResolve {
				assert.False("series meta-info was unnecessarily fetched")(miSrc.seriesFetched)
				assertSeriesEqual(assert, existingSeries, gotSeries)
			} else {
				assert.True("series meta-info was not fetched")(miSrc.seriesFetched)
				assertSeriesEqual(assert, &seriesMi, gotSeries)
			}
		}
	}
	existingSeries := SeriesMetaInfo {IdInfo: metainfo.IdInfo{episodeTi.Id}, Title: "yet another time waster", Seasons: 7, Year: "2002", Poster: "yatw.png"}

	t.Run("eager without pre-existing image", testFindOrFetch(false, nil, false))
	t.Run("lazy without pre-existing image", testFindOrFetch(true, nil, false))
	t.Run("eager with pre-existing image", testFindOrFetch(false, &existingSeries, false))
	t.Run("lazy with pre-existing image", testFindOrFetch(true, &existingSeries, true))
}

func TestFindOrFetchEpisode(t *testing.T) {
	testFindOrFetch := func(lazy bool, existingEpisode *EpisodeMetaInfo, noNeedToResolve bool) func(t *testing.T) {
		return func(t *testing.T) {
			assert := test.AssertOn(t)
			dir, conf := setupFindOrFetcher(assert, nil, nil, existingEpisode, nil)
			defer teardownFindOrFetcher(assert, dir)

			miSrc := newVideoMetaInfoSource(&movieMi, &seriesMi, &episodeMi, imageMi)
			fof := findOrFetch(miSrc, conf, lazy)

			gotEpisode, err := fof.episode(episodeTi)
			assert.NotError(err)
			assert.True("episode meta-info image was fetched")(0 == len(miSrc.imagesFetched))
			assert.False("series was fetched")(miSrc.seriesFetched)
			assert.False("movie meta-info was fetched")(miSrc.movieFetched)
			if noNeedToResolve {
				assert.False("episode meta-info was unnecessarily fetched")(miSrc.episodeFetched)
				assertEpisodesEqual(assert, existingEpisode, gotEpisode)
			} else {
				assert.True("episode meta-info was not fetched")(miSrc.episodeFetched)
				assertEpisodesEqual(assert, &episodeMi, gotEpisode)
			}

		}
	}

	existingEpisode := EpisodeMetaInfo{IdInfo: metainfo.IdInfo{Id: episodeTi.Id}, Title: "an earlier attack of the raffgrns", Year: "2008", Episode: episodeTi.Episode, Season: episodeTi.Season}

	t.Run("eager without pre-existing image", testFindOrFetch(false, nil, false))
	t.Run("lazy without pre-existing image", testFindOrFetch(true, nil, false))
	t.Run("eager with pre-existing image", testFindOrFetch(false, &existingEpisode, false))
	t.Run("lazy with pre-existing image", testFindOrFetch(true, &existingEpisode, true))
}


func assertMoviesEqual(assert *test.Assertion, expected *MovieMetaInfo, got *MovieMetaInfo) {
	if expected == nil || got == nil {
		assert.FailWith(fmt.Sprintf("did not expect nil for movie meta-info (expected %v, got %v)", expected, got))
	}
	assert.StringsEqual(expected.Id, got.Id)
	assert.StringsEqual(expected.Title, got.Title)
	assert.True(fmt.Sprintf("expected year %s, but got %s", expected.Year, got.Year))(expected.Year == got.Year)
	assert.StringsEqual(expected.Poster, got.Poster)
}

func assertImagesEqual(assert *test.Assertion, expected []byte, got []byte) {
	assert.True(fmt.Sprintf("expected image %v, but got %v", expected, got))(bytes.Equal(expected, got))
}

func assertSeriesEqual(assert *test.Assertion, expected *SeriesMetaInfo, got *SeriesMetaInfo) {
	if expected == nil || got == nil {
		assert.FailWith(fmt.Sprintf("did not expect nil for series meta-info (expected %v, got %v", expected, got))
	}
	assert.StringsEqual(expected.Title, got.Title)
	assert.StringsEqual(expected.Id, got.Id)
	assert.StringsEqual(expected.Poster, got.Poster)
	assert.StringsEqual(expected.Year, got.Year)
	assert.IntsEqual(expected.Seasons, got.Seasons)
}

func assertEpisodesEqual(assert *test.Assertion, expected *EpisodeMetaInfo, got *EpisodeMetaInfo) {
	if expected == nil || got == nil {
		assert.FailWith(fmt.Sprintf("did not expect nil for series meta-info (expected %v, got %v", expected, got))
	}
	assert.StringsEqual(expected.Id, got.Id)
	assert.StringsEqual(expected.Title, got.Title)
	assert.StringsEqual(expected.Year, got.Year)
	assert.IntsEqual(expected.Season, got.Season)
	assert.IntsEqual(expected.Episode, got.Episode)
}
