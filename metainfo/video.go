package metainfo

import (
	"path/filepath"
	"fmt"
)

const (
	SUBDIR_MOVIES = "movies"
	SUBDIR_SERIES = "series"
)

var (
	META_INFO_TYPE_MOVIE = "movieTi"
	META_INFO_TYPE_SERIES = "series"
	META_INFO_TYPE_EPISODE = "episodeTi"
)

type VideoMetaInfoSource interface {
	FetchMovieInfo(id string) (*MovieMetaInfo, error)
	FetchSeriesInfo(id string) (*SeriesMetaInfo, error)
	FetchEpisodeInfo(id string, season int, episode int) (*EpisodeMetaInfo, error)
	FetchImage(location string) (Image, error)
}


type MovieMetaInfo struct {
	IdInfo
	Title string
	Year int
	Poster string   //omdb: Poster
}
func (m *MovieMetaInfo) GetType() string {
	return META_INFO_TYPE_MOVIE
}

type SeriesMetaInfo struct {
	IdInfo
	Title string
	Seasons int    //omdb:totalSeasons
	Year int
	Poster string //omdb: Poster
}
func (s *SeriesMetaInfo) GetType() string {
	return META_INFO_TYPE_SERIES
}

type EpisodeMetaInfo struct {
	IdInfo
	Title string
	Episode int //omdb: Episode
	Season int //omdb: Season
	Year int
}
func (e *EpisodeMetaInfo) GetType() string {
	return META_INFO_TYPE_EPISODE
}


func MovieFileName(repoPath string, id string) string {
	return filepath.Join(repoPath, SUBDIR_MOVIES, fmt.Sprintf("%s.%s", id, METAINF_FILE_EXT))
}

func SeriesFileName(repoPath string, id string) string {
	return filepath.Join(repoPath, SUBDIR_SERIES, fmt.Sprintf("%s.%s", id, METAINF_FILE_EXT))
}

func EpisodeFileName(repoPath string, id string, season int, episode int) string {
	return filepath.Join(repoPath, SUBDIR_SERIES, fmt.Sprintf("%s.%d.%d.%s", id, season, episode, METAINF_FILE_EXT))
}
