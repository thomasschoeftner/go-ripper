package video

import (
	"testing"
	"go-cli/test"
	"go-ripper/targetinfo"
	"go-cli/task"
	"go-cli/commons"
	"go-ripper/ripper"
	"go-cli/config"
	"path/filepath"
	"go-ripper/files"
	"go-ripper/metainfo"
)

var movieTi = targetinfo.NewMovie("movieTi.mp4", "/a/b", "tt123456")
var episodeTi = targetinfo.NewEpisode("episode1.mp4", "/a/b/c", "tt654321", 3, 2, 4, 7)

var movieMi = MovieMetaInfo{IdInfo: metainfo.IdInfo{Id: movieTi.Id}, Title: "The awesome adventures of Sepp", Year: "2018", Poster: "taaos.jpg"}
var seriesMi = SeriesMetaInfo{IdInfo: metainfo.IdInfo{Id: episodeTi.Id}, Title: "a space oddity", Year: "2017", Seasons: 3, Poster: "aso.png"}
var episodeMi = EpisodeMetaInfo{IdInfo: metainfo.IdInfo{Id: episodeTi.Id}, Title: "attack of the raffgrns", Year: "2017", Episode: 4, Season: 3}
var imageMi = map[string][]byte{movieMi.Poster : []byte{1,2,3,4}, seriesMi.Poster : []byte{5,6,7,8}}

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


//integration-test
func TestResolveVideo(t *testing.T) {
	assert := test.AssertOn(t)

	dir := test.MkTempFolder(assert.T)
	defer test.RmTempFolder(t, dir)

	conf := &ripper.AppConf{}

	repoDir := filepath.ToSlash(filepath.Join(dir, "repo"))
	workDir := filepath.ToSlash(filepath.Join(dir, "work"))
	assert.NotError(config.FromString(conf, confJson,
		map[string]string {"repodir" : repoDir, "workdir" : workDir}))
	ctx := task.Context{nil, conf, commons.Printf, false, ""}

	// create target info files
	targetInfos := []targetinfo.TargetInfo{movieTi, episodeTi}
	for _, ti := range targetInfos {
		workDir, err := ripper.GetWorkPathForTargetFileFolder(conf.WorkDirectory, ti.GetFolder())
		assert.NotError(err)
		assert.NotError(files.CreateFolderStructure(workDir))
		assert.NotError(targetinfo.Save(workDir, ti))
	}

	// create jobs
	movieJob := task.Job{}.WithParam(ripper.JobField_Path, filepath.Join(movieTi.GetFolder(), movieTi.GetFile()))
	episodeJob := task.Job{}.WithParam(ripper.JobField_Path, filepath.Join(episodeTi.GetFolder(), episodeTi.GetFile()))

	t.Run("movieTi", func (t *testing.T) {
		miSource := newVideoMetaInfoSource(&movieMi, &seriesMi, &episodeMi, imageMi)
		NewVideoMetaInfoSource = func(conf *ripper.VideoResolveConfig) (VideoMetaInfoSource, error) {
			return miSource, nil
		}
		resolve := ResolveVideo(ctx)
		NewVideoMetaInfoSource = nil

		assert := test.AssertOn(t)
		resultJobs, err := resolve(movieJob)
		assert.NotError(err)

		assert.True("expected 1 result job")(1 == len(resultJobs))
		assert.StringsEqual(filepath.Join(movieTi.GetFolder(), movieTi.GetFile()), resultJobs[0][ripper.JobField_Path])
		assert.True("movieTi not fetched from meta-info source")(miSource.movieFetched)
		assert.True("image not fetched from meta-info source")(1 == len(miSource.imagesFetched))
	})

	t.Run("episodeTi", func (t *testing.T) {
		miSource := newVideoMetaInfoSource(&movieMi, &seriesMi, &episodeMi, imageMi)
		NewVideoMetaInfoSource = func(conf *ripper.VideoResolveConfig) (VideoMetaInfoSource, error) {
			return miSource, nil
		}
		resolve := ResolveVideo(ctx)
		NewVideoMetaInfoSource = nil

		assert := test.AssertOn(t)
		miSource.episodeFetched = false
		resultJobs, err := resolve(episodeJob)
		assert.NotError(err)

		assert.True("expected 1 result job")(1 == len(resultJobs))
		assert.StringsEqual(filepath.Join(episodeTi.GetFolder(), episodeTi.GetFile()), resultJobs[0][ripper.JobField_Path])
		assert.True("series not fetched from meta-info source")(miSource.seriesFetched)
		assert.True("image not fetched from meta-info source")(1 == len(miSource.imagesFetched))
		assert.True("episodeTi not fetched from meta-info source")(miSource.episodeFetched)
	})
}

func setupResolver(t *testing.T, lazy bool) (string, *findOrFetcher){
	dir := test.MkTempFolder(t)
	assert := test.AssertOn(t)
	conf := &ripper.AppConf{}
	assert.NotError(config.FromString(conf, confJson, map[string]string{
		"repodir": filepath.ToSlash(filepath.Join(dir, "repo")),
		"workdir": filepath.ToSlash(filepath.Join(dir, "work"))}))
	miSource := newVideoMetaInfoSource(&movieMi, &seriesMi, &episodeMi, imageMi)
	testFindOrFetcher := findOrFetch(miSource, conf, false)
	return dir, testFindOrFetcher
}

func teardownResolver(t *testing.T, dir string) {
	test.RmTempFolder(t, dir)
}

func TestResolveMovie(t *testing.T) {
	dir, testFindOrFetcher := setupResolver(t, false)
	defer teardownResolver(t, dir)
	assert := test.AssertOn(t)
	miSource := testFindOrFetcher.metaInfoSource.(*testVideoMetaInfoSource)

	assert.NotError(resolveMovie(testFindOrFetcher, movieTi))
	assert.True("movie meta-info not invoked")(miSource.movieFetched)
	assert.True("image resolve missing")(1 == len(miSource.imagesFetched))
	assert.StringsEqual(movieMi.Poster, miSource.imagesFetched[0])
	assert.False("series meta-info unnecessarily fetched")(miSource.seriesFetched)
	assert.False("episode meta-info unnecessarily fetched")(miSource.episodeFetched)
}

func TestResolveEpisode(t *testing.T) {
	dir, testFindOrFetcher := setupResolver(t, false)
	defer teardownResolver(t, dir)
	assert := test.AssertOn(t)
	miSource := testFindOrFetcher.metaInfoSource.(*testVideoMetaInfoSource)

	assert.NotError(resolveEpisode(testFindOrFetcher, episodeTi))
	assert.True("series meta-info not fetched")(miSource.seriesFetched)
	assert.True("episode meta-info not invoked")(miSource.episodeFetched)
	assert.True("image resolve missing")(1 == len(miSource.imagesFetched))
	assert.StringsEqual(seriesMi.Poster, miSource.imagesFetched[0])
	assert.False("movie meta-info unnecessarily fetched")(miSource.movieFetched)
}
