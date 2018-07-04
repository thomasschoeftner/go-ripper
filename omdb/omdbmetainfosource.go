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
	return &OmdbVideoMetaInfoSource{conf: conf, availableTokens: availableTokens, nextTokenIdx: 0}, nil
}

type OmdbVideoMetaInfoSource struct {
	conf *ripper.OmdbConfig
	availableTokens []string
	nextTokenIdx int
}

// round-robin use of omdb tokens
func (f *OmdbVideoMetaInfoSource) nextToken() string {
	token := f.availableTokens[f.nextTokenIdx]
	noOfTokens := len(f.availableTokens)
	f.nextTokenIdx = (f.nextTokenIdx + 1) % noOfTokens
	return token
}


func (omdb *OmdbVideoMetaInfoSource) FetchMovieInfo(id string) (*video.MovieMetaInfo, error) {
	url := replaceUrlVars(omdb.conf.MovieQuery, map[string]string{urlpattern_omdbtoken : omdb.nextToken(), urlpattern_imdbid : id})
	raw, err := httpGet(url)
	if err != nil {
		return nil, err
	}
	return toMovieMetaInfo(raw)
}

func (omdb *OmdbVideoMetaInfoSource) FetchSeriesInfo(id string) (*video.SeriesMetaInfo, error) {
	url := replaceUrlVars(omdb.conf.SeriesQuery, map[string]string{urlpattern_omdbtoken : omdb.nextToken(), urlpattern_imdbid : id})
	raw, err := httpGet(url)
	if err != nil {
		return nil, err
	}
	return toSeriesMetaInfo(raw)
}

func (omdb *OmdbVideoMetaInfoSource) FetchEpisodeInfo(id string, season int, episode int) (*video.EpisodeMetaInfo, error) {
	url := replaceUrlVars(omdb.conf.EpisodeQuery, map[string]string{
		urlpattern_omdbtoken : omdb.nextToken(),
		urlpattern_imdbid : id,
		urlpattern_season : strconv.Itoa(season),
		urlpattern_episode : strconv.Itoa(episode)})
	raw, err := httpGet(url)
	if err != nil {
		return nil, err
	}
	return toEpisodeMetaInfo(raw)
}

func (omdb *OmdbVideoMetaInfoSource) FetchImage(location string) (metainfo.Image, error) {
	return httpGet(location)
}

func replaceUrlVars(template string, keyVals map[string]string) string {
	result := template
	for variable, value := range keyVals {
		varPlaceholder := fmt.Sprintf("{%s}", variable)
		result = strings.Replace(result, varPlaceholder, value, -1)
	}
	return result
}

func httpGet(url string) ([]byte, error) {
	httpRsp, err:= http.Get(url)
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

