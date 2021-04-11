package omdb

import (
	"errors"
	"github.com/thomasschoeftner/go-ripper/metainfo"
	"github.com/thomasschoeftner/go-ripper/metainfo/video"
	"github.com/thomasschoeftner/go-ripper/ripper"
	"strings"
	"fmt"
	"strconv"
	"net/http"
	"io/ioutil"
	"time"
	"encoding/json"
)

const (
	urlpattern_omdbtoken = "omdbtoken"
	urlpattern_imdbid    = "imdbid"
	urlpattern_season    = "seasonNo"
	urlpattern_episode   = "episodeNo"
)

const CONF_OMDB_RESOLVER = "omdb"

func NewOmdbVideoMetaInfoSource(conf *ripper.VideoResolveConfig) (video.VideoMetaInfoSource, error) {
	if conf == nil {
		return nil, errors.New("cannot initialize omdb query factory without OmdbConfig")
	}
	if len(conf.Omdb.OmdbTokens) == 0 {
		return nil, errors.New("cannot initialize omdb query Factory with empty list of tokens")
	}

	httpClient := &http.Client{Timeout: time.Second * time.Duration(conf.Omdb.Timeout)}
	return &omdbVideoMetaInfoSource{conf: conf.Omdb, availableTokens: conf.Omdb.OmdbTokens, nextTokenIdx: 0, httpClient: httpClient}, nil
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
	raw, err := httpGet(omdb.httpClient).WithValidation(validateOmdbResponse).WithRetries(omdb.conf.Retries)(func() string {
		return replaceUrlVars(omdb.conf.MovieQuery, map[string]string{urlpattern_omdbtoken : omdb.nextToken(), urlpattern_imdbid : id})
	})
	if err != nil {
		return nil, err
	}
	return toMovieMetaInfo(raw)
}

func (omdb *omdbVideoMetaInfoSource) FetchSeriesInfo(id string) (*video.SeriesMetaInfo, error) {
	raw, err := httpGet(omdb.httpClient).WithValidation(validateOmdbResponse).WithRetries(omdb.conf.Retries)(func() string {
		return replaceUrlVars(omdb.conf.SeriesQuery, map[string]string{urlpattern_omdbtoken : omdb.nextToken(), urlpattern_imdbid : id})
	})
	if err != nil {
		return nil, err
	}
	return toSeriesMetaInfo(raw)
}

func (omdb *omdbVideoMetaInfoSource) FetchEpisodeInfo(id string, season int, episode int) (*video.EpisodeMetaInfo, error) {
	raw, err := httpGet(omdb.httpClient).WithValidation(validateOmdbResponse).WithRetries(omdb.conf.Retries)(func() string {
		return replaceUrlVars(omdb.conf.EpisodeQuery, map[string]string{
			urlpattern_omdbtoken: omdb.nextToken(),
			urlpattern_imdbid:    id,
			urlpattern_season:    strconv.Itoa(season),
			urlpattern_episode:   strconv.Itoa(episode)})
	})
	if err != nil {
		return nil, err
	}
	return toEpisodeMetaInfo(raw)
}

func (omdb *omdbVideoMetaInfoSource) FetchImage(location string) (metainfo.Image, error) {
	return httpGet(omdb.httpClient).WithRetries(omdb.conf.Retries)(func() string {
		return location
	})
}


func replaceUrlVars(template string, keyVals map[string]string) string {
	result := template
	for variable, value := range keyVals {
		varPlaceholder := fmt.Sprintf("{%s}", variable)
		result = strings.Replace(result, varPlaceholder, value, -1)
	}
	return result
}

type urlBuilder func() string
type httpGetFunc func(urlBuilder) ([]byte, error)
func httpGet(client *http.Client) httpGetFunc {
	return func(buildUrl urlBuilder) ([]byte, error) {
		url := buildUrl()
		httpRsp, err:= http.Get(url)
		if err != nil {
			return nil, err
		}
		defer httpRsp.Body.Close()

		if httpRsp.StatusCode == 401 {
			return nil, fmt.Errorf("invalid OMDB token used for URL: %s", url)
		}
		if httpRsp.StatusCode != 200 {
			return nil, fmt.Errorf("received unexpected response code %d when getting %s", httpRsp.StatusCode, url)
		}
		raw, err := ioutil.ReadAll(httpRsp.Body)
		if err != nil {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		return raw, nil
	}
}

func (hgf httpGetFunc) WithRetries(retries int) httpGetFunc {
	return func(url urlBuilder) ([]byte, error) {
		var errs []error
		v, e := hgf(url)
		if e == nil {
			return v, nil
		}
		errs = append(errs, e)
		for i:=0; i<retries; i++ {
			v, e = hgf(url)
			if e == nil {
				return v, nil
			}
			errs = append(errs, e)
		}
		errMsg := fmt.Sprintf("unable to resolve meta-info after %d tries due to: \n", retries + 1)
		for _, err := range errs {
			errMsg = fmt.Sprintf("%s   -%s\n", errMsg, err.Error())
		}
		return nil, errors.New(errMsg)
	}
}

func (hgf httpGetFunc) WithValidation(validateFunc func(raw []byte) error) httpGetFunc {
	return func(url urlBuilder) ([]byte, error) {
		v, e := hgf(url)
		if e != nil {
			return nil, e
		}
		e = validateFunc(v)
		if e != nil {
			return nil, e
		}
		return v, nil
	}
}

func validateOmdbResponse(raw []byte) error {
	status := basicOmdbResponse{}
	err := json.Unmarshal(raw, &status)
	if err != nil {
		return err
	}

	if strings.ToLower(status.Response) == "false" {
		return fmt.Errorf("%s", status.Error)
	}
	return nil
}

type basicOmdbResponse struct {
	Response string
	Error string
}
