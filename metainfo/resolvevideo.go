package metainfo

import (
	"go-cli/task"
	"go-ripper/ripper"
	"go-ripper/targetinfo"
	"errors"
	"go-ripper/files"
)

func ResolveVideo(metaInfo VideoMetaInfoSource) (task.Handler, error) {
	if metaInfo == nil {
		return nil, errors.New("unable to create ResolveVideo Handler without movie metainfo fetcher (nil)")
	}

	return func (ctx task.Context) task.HandlerFunc {
		conf := ctx.Config.(*ripper.AppConf)

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
				err = resolveEpisode(metaInfo, ti.(*targetinfo.Episode), conf, ctx.RunLazy)
			} else if targetinfo.IsMovie(ti) {
				err = resolveMovie(metaInfo, ti.(*targetinfo.Movie), conf, ctx.RunLazy)
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

func resolveMovie(metaInfo VideoMetaInfoSource, movie *targetinfo.Movie, conf *ripper.AppConf, lazy bool) error {
	var info *MovieMetaInfo
	metaInfFile := MovieFileName(conf.MetaInfoRepo, movie.Id)
	if needToResolve(metaInfFile, lazy) {
		mmi, err := metaInfo.FetchMovieInfo(movie.Id)
		if err != nil {
			return err
		}
		err = SaveMetaInfoFile(metaInfFile, mmi)
		if err != nil {
			return err
		}
		info = mmi
	} else {
		mmi := &MovieMetaInfo{}
		err := ReadMetaInfoFile(metaInfFile, mmi)
		if err != nil {
			return err
		}
		info = mmi
	}

	return findOrFetchImage(metaInfo, info.Id, info.Poster, conf, lazy)
}


func resolveEpisode(metaInfo VideoMetaInfoSource, ti *targetinfo.Episode, conf *ripper.AppConf, lazy bool) error {
	series, err := findOrFetchSeries(metaInfo, ti, conf, lazy)
	if err != nil {
		return err
	}

	//for now - assume each season is complete!!!
	//for later: TODO correct episode numbering
	//_, err = findOrFetchSeason(metaInfo, ti, conf, lazy)
	//if err != nil {
	//	return err
	//}
	ti.Episode = ti.ItemSeqNo
	targetinfo.Save(conf.WorkDirectory, ti)

	err = findOrFetchEpisode(metaInfo, ti, conf, lazy)
	if err != nil {
		return err
	}

	return findOrFetchImage(metaInfo, series.Id, series.Poster, conf, lazy)
}

func findOrFetchImage(metaInfo VideoMetaInfoSource, id string, imageUri string, conf *ripper.AppConf, lazy bool) error {
	imageFile := ImageFileName(conf.MetaInfoRepo, id, files.Extension(imageUri))
	if !needToResolve(imageFile, lazy) {
		return nil
	}
	imgData, err := metaInfo.FetchImage(imageUri)
	if err != nil {
		return err
	}
	return SaveMetaInfoImage(imageFile, imgData)
}

func findOrFetchSeries(metaInfo VideoMetaInfoSource, ti *targetinfo.Episode, conf *ripper.AppConf, lazy bool) (*SeriesMetaInfo, error) {
	seriesFile := SeriesFileName(conf.MetaInfoRepo, ti.Id)
	if needToResolve(seriesFile, lazy) {
		series, err := metaInfo.FetchSeriesInfo(ti.Id)
		if err != nil {
			return nil, err
		}
		err = SaveMetaInfoFile(seriesFile, series)
		if err != nil {
			return nil, err
		}
		return series, nil
	} else {
		series := &SeriesMetaInfo{}
		err := ReadMetaInfoFile(seriesFile, series)
		if err != nil {
			return nil, err
		}
		return series, nil
	}
}

func findOrFetchEpisode(metaInfo VideoMetaInfoSource, ti *targetinfo.Episode, conf *ripper.AppConf, lazy bool) error {
	episodeFile := EpisodeFileName(conf.MetaInfoRepo, ti.Id, ti.Season, ti.Episode)
	if needToResolve(episodeFile, lazy) {
		episode, err := metaInfo.FetchEpisodeInfo(ti.Id, ti.Season, ti.Episode)
		if err == nil {
			err = SaveMetaInfoFile(episodeFile, episode)
		}
		return err
	}
	return nil
}

func needToResolve(metaInfFile string, lazy bool) bool {
	alreadyExists, _ := files.Exists(metaInfFile)
	return !(lazy && alreadyExists)
}
