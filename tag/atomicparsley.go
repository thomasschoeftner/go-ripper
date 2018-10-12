package tag

import (
	"go-ripper/ripper"
	"errors"
	"go-cli/commons"
)

const CONF_ATOMICPARSLEY_TAGGER = "atomicparsley"

func NewAtomicParsleyVideoTagger(conf *ripper.TagConfig, lazy bool, printf commons.FormatPrinter) (VideoTagger, error)  {
	return &atomicParsleyVideoTagger{path: conf.Video.AtomicParsley.Path, printf: printf}, nil
}

type atomicParsleyVideoTagger struct {
	path string
	lazy bool
	printf commons.FormatPrinter
}

func (ap *atomicParsleyVideoTagger) TagMovie(file string, id string, title string, year string, posterPath string) error {
	ap.printf("AP tags %s with {id=%s, title=%s, year=%s, image=%s}\n", file, id, title, year, posterPath)
	//TODO check laziness
	//TODO implement me
	return errors.New("implement me - atomicparsley.go")
}

func (ap *atomicParsleyVideoTagger) TagEpisode(file string, id string, series string, season int, episode int, title string, year string, posterPath string) error {
	ap.printf("AP tags %s with {id=%s, series=%s, season=%d, episode=%d, title=%s, year=%s, image=%s }\n", file, id, series, season, episode, title, year, posterPath)
	//TODO check laziness
	//TODO implement me
	return errors.New("implement me - atomicparsley.go")
}
