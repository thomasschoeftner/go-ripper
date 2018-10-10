package tag

import (
	"go-ripper/ripper"
	"errors"
)

const CONF_ATOMICPARSLEY_TAGGER = "atomicparsley"

func NewAtomicParsleyVideoTagger(conf *ripper.TagConfig) (VideoTagger, error)  {
	return &AtomicParsleyVideoTagger{path: conf.Video.AtomicParsley.Path}, nil
}

type AtomicParsleyVideoTagger struct {
	path string
}

func (ap *AtomicParsleyVideoTagger) TagMovie(file string, id string, title string, year string, posterPath string) error {
	//fmt.Printf("TODO implement me - äöäöäöäöäöäöäöäöäöäöäöäöä - TagMovie(file=%s, id=%s, title=%s, year=%s, poster)ath=%s\n", file, id, title, year, posterPath)
	//TODO implement me
	return errors.New("implement me - atomicparsley.go")
}

func (ap *AtomicParsleyVideoTagger) TagEpisode(file string, id string, series string, season int, episode int, title string, year string, posterPath string) error {
	//fmt.Printf("TODO implement me - äöäöäöäöäöäöäöäöäöäöäöäöä - TagEpisode(file=%s, id=%s, series=%s, season=%d, episode=%d, title=%s, year=%s, poster)ath=%s\n", file, id, series, season, episode, title, year, posterPath)
	//TODO implement me
	return errors.New("implement me - atomicparsley.go")
}
