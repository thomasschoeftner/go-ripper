package omdb

import (
	"errors"
	"go-ripper/metainfo"
	"go-ripper/metainfo/video"
	"go-ripper/ripper"
	"strings"
	"fmt"
	"strconv"
	"net/http"
	"io/ioutil"
	"time"
)

const (
	urlpattern_omdbtoken = "omdbtoken"
	urlpattern_imdbid    = "imdbid"
	urlpattern_season    = "seasonNo"
	urlpattern_episode   = "episodeNo"
)

func NewOmdbVideoQueryFactory(conf *ripper.OmdbConfig, availableTokens []string) (video.VideoMetaInfoSource, error) {
	if conf == nil {
		return nil, errors.New("cannot initialize omdb query factory without OmdbConfig")
	}
	if len(availableTokens) == 0 {
		return nil, errors.New("cannot initialize omdb query Factory with empty list of tokens")
	}

	httpClient := &http.Client{Timeout: time.Second * time.Duration(conf.Timeout)}
	return &omdbVideoMetaInfoSource{conf: conf, availableTokens: availableTokens, nextTokenIdx: 0, httpClient: httpClient}, nil
}

type omdbVideoMetaInfoSource struct {
	conf *ripper.OmdbConfig
	availableTokens []string
	nextTokenIdx int
	httpClient *http.Client
}

// round-robin use of omdb tokens
func (f *omdbVideoMetaInfoSource) nextToken() string {
	token := f.availableTokens[f.nextTokenIdx]
	noOfTokens := len(f.availableTokens)
	f.nextTokenIdx = (f.nextTokenIdx + 1) % noOfTokens
	return token
}


func (omdb *omdbVideoMetaInfoSource) FetchMovieInfo(id string) (*video.MovieMetaInfo, error) {
	url := replaceUrlVars(omdb.conf.MovieQuery, map[string]string{urlpattern_omdbtoken : omdb.nextToken(), urlpattern_imdbid : id})
	raw, err := omdb.httpGet(url)
	if err != nil {
		return nil, err
	}
	return toMovieMetaInfo(raw)
}

func (omdb *omdbVideoMetaInfoSource) FetchSeriesInfo(id string) (*video.SeriesMetaInfo, error) {
	url := replaceUrlVars(omdb.conf.SeriesQuery, map[string]string{urlpattern_omdbtoken : omdb.nextToken(), urlpattern_imdbid : id})
	raw, err := omdb.httpGet(url)
	if err != nil {
		return nil, err
	}
	return toSeriesMetaInfo(raw)
}

func (omdb *omdbVideoMetaInfoSource) FetchEpisodeInfo(id string, season int, episode int) (*video.EpisodeMetaInfo, error) {
	url := replaceUrlVars(omdb.conf.EpisodeQuery, map[string]string{
		urlpattern_omdbtoken : omdb.nextToken(),
		urlpattern_imdbid : id,
		urlpattern_season : strconv.Itoa(season),
		urlpattern_episode : strconv.Itoa(episode)})
	raw, err := omdb.httpGet(url)
	if err != nil {
		return nil, err
	}
	return toEpisodeMetaInfo(raw)
}

func (omdb *omdbVideoMetaInfoSource) FetchImage(location string) (metainfo.Image, error) {
	return omdb.httpGet(location)
}

func replaceUrlVars(template string, keyVals map[string]string) string {
	result := template
	for variable, value := range keyVals {
		varPlaceholder := fmt.Sprintf("{%s}", variable)
		result = strings.Replace(result, varPlaceholder, value, -1)
	}
	return result
}

func (omdb *omdbVideoMetaInfoSource) httpGet(url string) ([]byte, error) {
	httpRsp, err:= omdb.httpClient.Get(url)
	defer httpRsp.Body.Close()
	if err != nil {
		return nil, err
	}
	raw, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

