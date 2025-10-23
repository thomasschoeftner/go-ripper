package video

import (
	"errors"

	"github.com/thomasschoeftner/go-cli/task"
	"github.com/thomasschoeftner/go-ripper/ripper"
	"github.com/thomasschoeftner/go-ripper/targetinfo"
)

// needs to be set for successful creation of a video meta-info source
var NewVideoMetaInfoSource func(conf *ripper.VideoResolveConfig) (VideoMetaInfoSource, error)

func ResolveVideo(ctx task.Context) task.HandlerFunc {
	conf := ctx.Config.(*ripper.AppConf)
	if nil == NewVideoMetaInfoSource {
		return ripper.ErrorHandler(errors.New("video meta-info source is undefined"))
	}
	metaInfoSrc, err := NewVideoMetaInfoSource(conf.Resolve.Video)
	if err != nil {
		return ripper.ErrorHandler(err)
	}
	findOrFetcher := findOrFetch(metaInfoSrc, conf, ctx.RunLazy)

	return func(job task.Job) ([]task.Job, error) {
		target := ripper.GetTargetFileFromJob(job)
		ctx.Printf("resolve video - target %s\n", target)

		ti, err := targetinfo.ForTarget(conf.WorkDirectory, target)
		if err != nil {
			return nil, err
		}

		printf := ctx.Printf.WithIndent(2)
		printf("recovered target-info: %s\n", ti.String())

		if targetinfo.IsEpisode(ti) {
			err = resolveEpisode(findOrFetcher, ti.(*targetinfo.Episode))
		} else if targetinfo.IsMovie(ti) {
			err = resolveMovie(findOrFetcher, ti.(*targetinfo.Movie))
		} else {
			//ignore other target-info types (e.g audio)
		}
		if err != nil {
			return nil, err
		}
		return []task.Job{job}, nil
	}
}

func resolveMovie(findOrFetch *findOrFetcher, ti *targetinfo.Movie) error {
	movie, err := findOrFetch.movie(ti)
	if err != nil {
		return err
	}
	return findOrFetch.image(movie.Id, movie.Poster)
}

func resolveEpisode(findOrFetch *findOrFetcher, ti *targetinfo.Episode) error {
	series, err := findOrFetch.series(ti)
	if err != nil {
		return err
	}

	// TODO - consider adding Season Meta-Info

	_, err = findOrFetch.episode(ti)
	if err != nil {
		return err
	}

	return findOrFetch.image(series.Id, series.Poster)
}
