package metainfo

import (
	"go-ripper/targetinfo"
	"go-ripper/files"
	"go-ripper/ripper"
)

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
		return ff.metaInfoSource.FetchEpisodeInfo(ti.Id, ti.Season, ti.ItemSeqNo)
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
