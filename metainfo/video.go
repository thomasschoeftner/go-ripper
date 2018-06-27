package metainfo

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"fmt"
	"go-ripper/files"
)

const METAINF_FILE_EXT = "json"

const (
	SUBDIR_IMAGES = "imgs"
	SUBDIR_MOVIES = "movies"
	SUBDIR_SERIES = "series"
)

type VideoMetaInfoSource interface {
	FetchMovieInfo(id string) (*MovieMetaInfo, error)
	FetchSeriesInfo(id string) (*SeriesMetaInfo, error)
	//FetchSeasonInfo(id string, season int) (*SeasonMetaInfo, error)
	FetchEpisodeInfo(id string, season int, episode int) (*EpisodeMetaInfo, error)
	FetchImage(location string) ([]byte, error)
}

type BasicVideoMetaInfo struct {
	Id       string   //omdb: imdbID
	Title    string   //omdb: Title
	Year     int      //omdb: Year
}

type MovieMetaInfo struct {
	BasicVideoMetaInfo
	Poster   string   //omdb: Poster
}

type SeriesMetaInfo struct {
	BasicVideoMetaInfo
	Seasons int    //omdb:totalSeasons
	Poster  string //omdb: Poster
}

//type SeasonMetaInfo struct {
//	Id       string //omdb: imdbID
//	Season   int    //omdb: Season
//	Episodes int
//}

type EpisodeMetaInfo struct {
	BasicVideoMetaInfo
	Episode int //omdb: Episode
	Season  int //omdb: Season
}


func MovieFileName(repoPath string, id string) string {
	return filepath.Join(repoPath, SUBDIR_MOVIES, fmt.Sprintf("%s.%s", id, METAINF_FILE_EXT))
}

func SeriesFileName(repoPath string, id string) string {
	return filepath.Join(repoPath, SUBDIR_SERIES, fmt.Sprintf("%s.%s", id, METAINF_FILE_EXT))
}

//func SeasonFileName(repoPath string, id string, season int) string {
//	return filepath.Join(repoPath, SUBDIR_SERIES, fmt.Sprintf("%s.%d.%s", id, season, METAINF_FILE_EXT))
//}

func EpisodeFileName(repoPath string, id string, season int, episode int) string {
	return filepath.Join(repoPath, SUBDIR_SERIES, fmt.Sprintf("%s.%d.%d.%s", id, season, episode, METAINF_FILE_EXT))
}

func ImageFileName(repoPath string, id string, extension string) string {
	return filepath.Join(repoPath, SUBDIR_IMAGES, fmt.Sprintf("%s.%s", id, extension))
}

func ReadMetaInfoFile(filePath string, content interface{}) error {
	if content == nil {
		return errors.New("cannot unmarshall to nil")
	}

	jsonRaw, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonRaw, content)
	return err
}

func ReadMetaInfoImage(filePath string) ([]byte, error) {
	if len(filePath) == 0 {
		return nil, errors.New("cannot unmarshall from empty file name")
	}
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func SaveMetaInfoFile(filePath string, metaInfo interface{}) error {
	if metaInfo == nil {
		return errors.New("cannot save nil json")
	}

	folder := filepath.Dir(filePath)
	err := files.CreateFolderStructure(folder)
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(metaInfo, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, bytes, os.ModePerm)
}

func SaveMetaInfoImage(filePath string, img []byte) error {
	if img == nil || len(img) == 0 {
		return errors.New("cannot save empty image")
	}

	folder := filepath.Dir(filePath)
	err := files.CreateFolderStructure(folder)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, img, os.ModePerm)
}