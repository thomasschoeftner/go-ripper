package omdb

import (
	"encoding/json"
	"io"
	"net/http"
	"go-ripper/metainfo"
	"errors"
)


func omdbMapper(raw []byte) (map[string]string, error) {
	var results map[string]string
	err := json.Unmarshal(raw, results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func newQuery(url string) *omdbQuery {
	return &omdbQuery{url, httpGet}
}
type omdbQuery struct {
	url string
	getter func(string) (io.ReadCloser, error)
}

func (oq *omdbQuery) Invoke() (io.ReadCloser, error) {
	return oq.getter(oq.url)
}

func (oq *omdbQuery) Convert(raw []byte) (metainfo.MetaInfo, error) {
	//TODO implement me
	return nil, errors.New("implement me")
}

func httpGet(url string) (io.ReadCloser, error) {
	httpRsp, err:= http.Get(url)
	if err != nil {
		return nil, err
	}
	return httpRsp.Body, nil
}
