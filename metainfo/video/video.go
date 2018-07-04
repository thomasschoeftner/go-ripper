package video

import (
	"path/filepath"
	"fmt"
	"go-ripper/metainfo"
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
	FetchImage(location string) (metainfo.Image, error)
}


type MovieMetaInfo struct {
	metainfo.IdInfo
	Title string
	Year int
	Poster string   //omdb: Poster
}
func (m *MovieMetaInfo) GetType() string {
	return META_INFO_TYPE_MOVIE
}

type SeriesMetaInfo struct {
	metainfo.IdInfo
	Title string
	Seasons int    //omdb:totalSeasons
	Year int
	Poster string //omdb: Poster
}
func (s *SeriesMetaInfo) GetType() string {
	return META_INFO_TYPE_SERIES
}

type EpisodeMetaInfo struct {
	metainfo.IdInfo
	Title string
	Episode int //omdb: Episode
	Season int //omdb: Season
	Year int
}
func (e *EpisodeMetaInfo) GetType() string {
	return META_INFO_TYPE_EPISODE
}


func MovieFileName(repoPath string, id string) string {
	return filepath.ToSlash(filepath.Join(repoPath, SUBDIR_MOVIES, fmt.Sprintf("%s.%s", id, metainfo.METAINF_FILE_EXT)))
}

func SeriesFileName(repoPath string, id string) string {
	return filepath.ToSlash(filepath.Join(repoPath, SUBDIR_SERIES, fmt.Sprintf("%s.%s", id, metainfo.METAINF_FILE_EXT)))
}

func EpisodeFileName(repoPath string, id string, season int, episode int) string {
	return filepath.ToSlash(filepath.Join(repoPath, SUBDIR_SERIES, fmt.Sprintf("%s.%d.%d.%s", id, season, episode, metainfo.METAINF_FILE_EXT)))
}
