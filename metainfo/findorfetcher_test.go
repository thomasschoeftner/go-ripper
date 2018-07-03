package metainfo

import (
	"go-cli/test"
	"go-ripper/ripper"
	"path/filepath"
	"go-cli/config"
	"go-ripper/files"
	"testing"
)

func setupFindOrFetcher(assert *test.Assertion,  movie *MovieMetaInfo, series *SeriesMetaInfo, episode *EpisodeMetaInfo, images map[string][]byte) (string, *ripper.AppConf) {
		dir := test.MkTempFolder(assert.T)
		conf := &ripper.AppConf{}

		repoDir := filepath.ToSlash(filepath.Join(dir, "repo"))
		workDir := filepath.ToSlash(filepath.Join(dir, "work"))
		assert.NotError(config.FromString(conf, confJson,
			map[string]string {"repodir" : repoDir, "workdir" : workDir}))

		// create meta info files if passed
		if movie != nil {
			assert.NotError(SaveMetaInfo(MovieFileName(repoDir, movie.Id), movie))
		}
		if series != nil {
			assert.NotError(SaveMetaInfo(SeriesFileName(repoDir, series.Id), series))
		}
		if episode != nil {
			assert.NotError(SaveMetaInfo(EpisodeFileName(repoDir, episode.Id, episode.Season, episode.Episode), episode))
		}
		if images != nil {
			for f, image := range images {
				println(f, image)
				var imgFileName string
				if movie != nil && movie.Poster == f {
					imgFileName = ImageFileName(repoDir, movie.Id, files.Extension(movie.Poster))
				} else if series != nil && series.Poster == f {
					imgFileName = ImageFileName(repoDir, series.Id, files.Extension(series.Poster))
				} else {
					assert.T.Fatalf("unknown poster name %s matches neither movie, nor series", f)
				}
				assert.NotError(SaveImage(imgFileName, image))
			}
		}

		return dir, conf
}

func teardownFindOrFetcher(assert *test.Assertion, dir string) {
	test.RmTempFolder(assert.T, dir)
}



func TestFindOrFetchMovie(t *testing.T) {
	t.Run("eager without pre-existing meta-info files", func(t *testing.T) {
		const lazy= false
		assert := test.AssertOn(t)
		dir, conf := setupFindOrFetcher(assert, nil, nil, nil, nil)
		defer teardownFindOrFetcher(assert, dir)

		miSrc := newVideoMetaInfoSource(&movieMi, &seriesMi, &episodeMi, imageMi)
		fof := findOrFetch(miSrc, conf, lazy)

		gotMovie, err := fof.movie(movieTi)

		assert.NotError(err)
		assert.True("movie meta-info was not fetched")(miSrc.movieFetched)
		assert.True("movie meta-info image was fetched")(0 == len(miSrc.imagesFetched))
		assert.False("series was fetched")(miSrc.seriesFetched)
		assert.False("episode meta-info was fetched")(miSrc.episodeFetched)

		assert.StringsEqual(movieMi.Id, gotMovie.Id)
		assert.StringsEqual(movieMi.Poster, gotMovie.Poster)
	})

	t.Run("lazy without pre-existing meta-info files", func(t *testing.T) {
		const lazy= true
		assert := test.AssertOn(t)
		dir, conf := setupFindOrFetcher(assert, nil, nil, nil, nil)
		defer teardownFindOrFetcher(assert, dir)

		miSrc := newVideoMetaInfoSource(&movieMi, &seriesMi, &episodeMi, imageMi)
		fof := findOrFetch(miSrc, conf, lazy)

		gotMovie, err := fof.movie(movieTi)

		assert.NotError(err)
		assert.True("movie meta-info was not fetched")(miSrc.movieFetched)
		assert.True("movie meta-info image was fetched")(0 == len(miSrc.imagesFetched))
		assert.False("series was fetched")(miSrc.seriesFetched)
		assert.False("episode meta-info was fetched")(miSrc.episodeFetched)

		assert.StringsEqual(movieMi.Id, gotMovie.Id)
		assert.StringsEqual(movieMi.Poster, gotMovie.Poster)
	})

	t.Run("eager with pre-existing meta-info files", func(t *testing.T) {
		const lazy= false
		assert := test.AssertOn(t)

		existingMovie := MovieMetaInfo{IdInfo: IdInfo{Id: movieTi.Id}, Title: "an earlier awesome adventure of Sepp", Year: 2008, Poster: "aeaaos.jpg"}
		existingImages := map[string][]byte{existingMovie.Poster: []byte{12, 13, 14, 15}}
		dir, conf := setupFindOrFetcher(assert, &existingMovie, nil, nil, existingImages)
		defer teardownFindOrFetcher(assert, dir)

		miSrc := newVideoMetaInfoSource(&movieMi, &seriesMi, &episodeMi, imageMi)
		fof := findOrFetch(miSrc, conf, lazy)

		gotMovie, err := fof.movie(movieTi)
		assert.NotError(err)
		assert.True("movie was not fetched")(miSrc.movieFetched)
		assert.True("movie image was fetched")(0 == len(miSrc.imagesFetched))
		assert.False("series was fetched")(miSrc.seriesFetched)
		assert.False("episode was fetched")(miSrc.episodeFetched)

		assert.StringsEqual(movieMi.Id, gotMovie.Id)
		assert.StringsEqual(movieMi.Poster, gotMovie.Poster)
	})

	t.Run("lazy with pre-existing meta-info files", func(t *testing.T) {
		const lazy= true
		assert := test.AssertOn(t)

		existingMovie := MovieMetaInfo{IdInfo: IdInfo{Id: movieTi.Id}, Title: "an earlier awesome adventure of Sepp", Year: 2008, Poster: "aeaaos.jpg"}
		existingImages := map[string][]byte{existingMovie.Poster: []byte{12, 13, 14, 15}}
		dir, conf := setupFindOrFetcher(assert, &existingMovie, nil, nil, existingImages)
		defer teardownFindOrFetcher(assert, dir)

		miSrc := newVideoMetaInfoSource(&movieMi, &seriesMi, &episodeMi, imageMi)
		fof := findOrFetch(miSrc, conf, lazy)

		gotMovie, err := fof.movie(movieTi)
		assert.NotError(err)
		assert.False("movie was unnecessarily fetched")(miSrc.movieFetched)
		assert.True("movie image was fetched")(0 == len(miSrc.imagesFetched))
		assert.False("series was fetched")(miSrc.seriesFetched)
		assert.False("episode was fetched")(miSrc.episodeFetched)

		assert.StringsEqual(existingMovie.Id, gotMovie.Id)
		assert.StringsEqual(existingMovie.Poster, gotMovie.Poster)
	})

	//	t.Run("lazy with partially pre-existing meta-info files", func(t *testing.T) {
	//		const lazy = true
	//		assert := test.AssertOn(t)
	//
	//		existingMovie := MovieMetaInfo{BasicVideoMetaInfo: BasicVideoMetaInfo{Id: movieTi.Id, Title: "an earlier awesome adventure of Sepp", Year: 2008}, Poster: "aeaaos.jpg"}
	//		missingImage := map[string][]byte{movieMi.Poster : []byte{4, 5, 6, 7}, existingMovie.Poster : []byte{12,13,14,15}}
	//		dir, conf := setup(assert, confJson, &existingMovie, nil, nil, nil)
	//		defer teardown(t, dir)
	//
	//		src := newVideoMetaInfoSource(&movieMi, &seriesMi, &episodeMi, missingImage)
	//		gotMi, gotImg := testResolveMovie(assert, src, conf, lazy, movieTi)
	//		assert.False("movieTi was unnecessarily fetched")(src.movieFetched)
	//		assert.True("movieTi image was not fetched")(src.imageFetched)
	//		assert.False("series was fetched")(src.seriesFetched)
	//		assert.False("episodeTi was fetched")(src.episodeFetched)
	//
	//		assertMoviesEqual(assert, &existingMovie, gotMi)
	//		assertImagesEqual(assert, missingImage[existingMovie.Poster], gotImg)
	//	})
	//}
	//
	//func assertMoviesEqual(assert *test.Assertion, expected *MovieMetaInfo, got *MovieMetaInfo) {
	//	if expected == nil || got == nil {
	//		assert.FailWith(fmt.Sprintf("did not expect nil for movieTi metainfo (expected %v, got %v)", expected, got))
	//	}
	//	assert.StringsEqual(expected.Id, got.Id)
	//	assert.StringsEqual(expected.Title, got.Title)
	//	assert.True(fmt.Sprintf("expected year %d, but got %d", expected.Year, got.Year))(expected.Year == got.Year)
	//	assert.StringsEqual(expected.Poster, got.Poster)
	//}
	//
	//func assertImagesEqual(assert *test.Assertion, expected []byte, got []byte) {
	//	assert.True(fmt.Sprintf("expected image %v, but got %v", expected, got))(bytes.Equal(expected, got))
	//}
	//
	//
	//func testResolve____TODO___(assert *test.Assertion, miSource VideoMetaInfoSource, conf *ripper.AppConf, lazy bool, ti *targetinfo.Episode) *EpisodeMetaInfo {
	//	assert.NotError(resolveEpisode(miSource, ti, conf, lazy))
	//	gotMi := &EpisodeMetaInfo{}
	//
	//	//read meta-info and resolve
	//	assert.NotError(ReadMetaInfoFile(MovieFileName(conf.MetaInfoRepo, ti.Id), gotMi))
	//
	//	return gotMi
	//}
	//
	//
	//func TestFindOrFetchImage(t *testing.T) {
	//
	//}
	//
	//func TestFindOrFetchSeries(t *testing.T) {
	//
	//}
	//
	//func TestFindOrFetchEpisode(t *testing.T) {
	//
	//}

}