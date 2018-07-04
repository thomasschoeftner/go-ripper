package video

import (
	"go-cli/task"
	"go-ripper/ripper"
	"go-ripper/targetinfo"
	"errors"
)

func ResolveVideo(metaInfoSrc VideoMetaInfoSource) (task.Handler, error) {
	if metaInfoSrc == nil {
		return nil, errors.New("unable to create ResolveVideo Handler without movieTi metainfo fetcher (nil)")
	}

	return func (ctx task.Context) task.HandlerFunc {
		conf := ctx.Config.(*ripper.AppConf)
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

	}, nil
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

	//for now - assume each season is complete!!!
	//for later: TODO correct episodeTi numbering
	//_, err = findOrFetchSeason(metaInfoSource, ti, conf, lazy)
	//if err != nil {
	//	return err
	//}
	ti.Episode = ti.ItemSeqNo
	targetinfo.Save(findOrFetch.conf.WorkDirectory, ti)

	_, err = findOrFetch.episode(ti)
	if err != nil {
		return err
	}

	return findOrFetch.image(series.Id, series.Poster)
}
