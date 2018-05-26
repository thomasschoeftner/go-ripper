package omdb

import (
	"errors"
	"go-ripper/metainfo"
	"go-ripper/ripper"
	"strings"
	"fmt"
	"strconv"
)

const (
	urlpattern_omdbtoken = "omdbtoken"
	urlpattern_imdbid    = "imdbid"
	urlpattern_season    = "seasonNo"
	urlpattern_episode   = "episodeNo"
)

func NewOmdbVideoQueryFactory(conf *ripper.OmdbConfig, availableTokens []string) (metainfo.VideoMetaInfoQueryFactory, error) {
	if conf == nil {
		return nil, errors.New("cannot initialize omdb query factory without OmdbConfig")
	}
	if len(availableTokens) == 0 {
		return nil, errors.New("cannot initialize omdb query Factory with empty list of tokens")
	}

	return &OmdbVideoQueryFactory{conf: conf, availableTokens: availableTokens, nextTokenIdx: 0}, nil
}

type OmdbVideoQueryFactory struct {
	conf *ripper.OmdbConfig
	availableTokens []string
	nextTokenIdx int
}

// round-robin use of omdb tokens
func (f *OmdbVideoQueryFactory) nextToken() string {
	token := f.availableTokens[f.nextTokenIdx]
	noOfTokens := len(f.availableTokens)
	f.nextTokenIdx = (f.nextTokenIdx + 1) % noOfTokens
	return token
}


func (omdb *OmdbVideoQueryFactory) NewTitleQuery(id string) metainfo.MetaInfoQuery {
	url := replaceVars(omdb.conf.TitleQuery, map[string]string{urlpattern_omdbtoken : omdb.nextToken(), urlpattern_imdbid : id})
	return newQuery(url)
}

func (omdb *OmdbVideoQueryFactory) NewEpisodeQuery(id string, season int, episode int) metainfo.MetaInfoQuery {
	url := replaceVars(omdb.conf.TitleQuery, map[string]string{
		urlpattern_omdbtoken : omdb.nextToken(),
		urlpattern_imdbid : id,
		urlpattern_season : strconv.Itoa(season),
		urlpattern_episode : strconv.Itoa(episode)})
	return newQuery(url)
}

func replaceVars(template string, keyVals map[string]string) string {
	result := template
	for variable, value := range keyVals {
		result = replaceVar(result, variable, value)
	}
	return result
}

func replaceVar(template string, variable string, value string) string {
	return strings.Replace(template, fmt.Sprintf("{%s}", variable), value, -1)
}
