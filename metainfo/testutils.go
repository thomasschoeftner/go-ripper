package metainfo

import (
	"errors"
	"fmt"
)

type testVideoMetaInfoSource struct {
	movie *MovieMetaInfo
	series *SeriesMetaInfo
	episode *EpisodeMetaInfo
	images map[string][]byte
	movieFetched bool
	seriesFetched bool
	episodeFetched bool
	imagesFetched []string
}

func (f *testVideoMetaInfoSource) FetchMovieInfo(id string) (*MovieMetaInfo, error) {
	var err error
	var m *MovieMetaInfo
	if f.movie == nil {
		err = errors.New("test error - no movies defined")
	} else if f.movie.Id != id {
		err = errors.New("test error - movieTi not found")
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
		err = errors.New("test error - no episodeTi defined")
	} else if f.episode.Id != id || f.episode.Season != season || f.episode.Episode != episode {
		err = errors.New("test error - episodeTi not found")
	} else {
		e = f.episode
		f.episodeFetched = true
	}
	return e, err

}

func (f *testVideoMetaInfoSource) FetchImage(location string) (Image, error) {
	var img []byte
	var err error

	if f.images == nil {
		err = errors.New("test error - no images defined")
	} else {
		_, found := f.images[location]
		if !found {
			err = fmt.Errorf("test error - image \"%s\" not found", location)
		} else {
			img = f.images[location]
			f.imagesFetched = append(f.imagesFetched, location)
		}
	}
	return img, err

}

func newVideoMetaInfoSource(movie *MovieMetaInfo, series *SeriesMetaInfo, /*season *SeasonMetaInfo,*/ episode *EpisodeMetaInfo, images map[string][]byte) *testVideoMetaInfoSource {
	return &testVideoMetaInfoSource{movie: movie, series: series, /*season: season,*/ episode: episode, images: images}
}
