package metainfo

import (
	"go-cli/task"
	"go-ripper/ripper"
	"go-ripper/targetinfo"
	"errors"
	"go-ripper/files"
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



func findOrFetch(metaInfo VideoMetaInfoSource, conf *ripper.AppConf, lazy bool) *findOrFetcher {
	return &findOrFetcher{metaInfoSource: metaInfo, conf: conf, lazy: lazy}
}
type findOrFetcher struct {
	metaInfoSource VideoMetaInfoSource
	conf           *ripper.AppConf
	lazy           bool
}

type fetchFunc func() (MetaInfo, error)
func (ff *findOrFetcher) doResolve(metaInfo MetaInfo, metaInfoFileName string, doFetch fetchFunc) (MetaInfo, error) {
	if ff.needToResolve(metaInfoFileName, ff.lazy) {
		 mi, err := doFetch()
		if err == nil {
			err = SaveMetaInfo(metaInfoFileName, mi)
		}
		if err != nil {
			return nil, err
		}
		return mi, nil
	} else {
		err := ReadMetaInfo(metaInfoFileName, metaInfo)
		if err != nil {
			return nil, err
		}
		return metaInfo, nil
	}
}

func (ff * findOrFetcher) movie(ti *targetinfo.Movie) (*MovieMetaInfo, error) {
	mi, err := ff.doResolve(&MovieMetaInfo{}, MovieFileName(ff.conf.MetaInfoRepo, ti.Id), func() (MetaInfo, error) {
		return ff.metaInfoSource.FetchMovieInfo(ti.Id)
	})

	if err != nil {
		return nil, err
	}
	return mi.(*MovieMetaInfo), nil
}

func (ff * findOrFetcher) series(ti *targetinfo.Episode) (*SeriesMetaInfo, error) {
	mi, err := ff.doResolve(&SeriesMetaInfo{}, SeriesFileName(ff.conf.MetaInfoRepo, ti.Id), func() (MetaInfo, error) {
		return ff.metaInfoSource.FetchSeriesInfo(ti.Id)
	})
	if err != nil {
		return nil, err
	}
	return mi.(*SeriesMetaInfo), nil
}

func (ff *findOrFetcher) episode(ti *targetinfo.Episode) (*EpisodeMetaInfo, error) {
	mi, err := ff.doResolve(&EpisodeMetaInfo{}, EpisodeFileName(ff.conf.MetaInfoRepo, ti.Id, ti.Season, ti.Episode), func() (MetaInfo, error) {
		return ff.metaInfoSource.FetchEpisodeInfo(ti.Id, ti.Season, ti.Episode)
	})
	if err != nil {
		return nil, err
	}
	return mi.(*EpisodeMetaInfo), nil
}


func (ff * findOrFetcher) image(id string, imageUri string) error {
	imageFile := ImageFileName(ff. conf.MetaInfoRepo, id, files.Extension(imageUri))
	if !ff.needToResolve(imageFile, ff.lazy) {
		return nil
	}
	imgData, err := ff.metaInfoSource.FetchImage(imageUri)
	if err != nil {
		return err
	}
	return SaveImage(imageFile, imgData)
}

func (ff *findOrFetcher) needToResolve(metaInfFile string, lazy bool) bool {
	if !lazy {
		return true
	}
	alreadyExists, _ := files.Exists(metaInfFile)
	return !alreadyExists
}
