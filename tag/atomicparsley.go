package tag

import (
	"go-ripper/ripper"
	"errors"
	"go-cli/commons"
)

const CONF_ATOMICPARSLEY_TAGGER = "atomicparsley"
const (
	paramOutput = "output"

	paramTitle = "title"
	paramYear  = "year"
	paramPoster = "poster"
	paramDescription = "desc"

	paramSeriesName = "seriesName"
	paramSeason = "season"
	paramEpisode = "episode"
	paramEpisodeName = "episodeName"

	paramComment = "comment"
	paramGenre = "genre"
	paramAlbum = "album"
	paramDisc = "disc"
)

var params = map[string]string {
	paramOutput: "-o",

	paramTitle: "--title",
	paramYear: "--year",
	paramPoster: "--artwork",
	paramDescription: "--description",

	paramSeriesName: "--TVShowName",
	paramSeason: "--TVSeasonNum",
	paramEpisode: "--TVEpisodeNum",
	paramEpisodeName: "",

	paramComment: "--comment",
	paramGenre: "--genre",
	paramAlbum: "--album",
	paramDisc: "--disk",
}

func NewAtomicParsleyVideoTagger(conf *ripper.TagConfig, lazy bool, printf commons.FormatPrinter) (VideoTagger, error)  {
	if lazy {
		printf("WARNING - video tagging via AtomicParsley is NOT lazy - ie files will always be tagged/written!\n")
	}
	return &atomicParsleyVideoTagger{conf: conf.Video.AtomicParsley, lazy: lazy, printf: printf}, nil
}

type atomicParsleyVideoTagger struct {
	conf *ripper.AtomicParsleyConfig
	lazy bool
	printf commons.FormatPrinter
}

func (ap *atomicParsleyVideoTagger) TagMovie(inFile string, outFile string, id string, title string, year string, posterPath string) error {
	ap.printf("AP tags %s with {id=%s, title=%s, year=%s, image=%s} -> write to %s\n", inFile, id, title, year, posterPath, outFile)
	//TODO implement me
	return errors.New("implement me - atomicparsley.go")
}

func (ap *atomicParsleyVideoTagger) TagEpisode(inFile string, outFile string, id string, series string, season int, episode int, title string, year string, posterPath string) error {
	ap.printf("AP tags %s with {id=%s, series=%s, season=%d, episode=%d, title=%s, year=%s, image=%s} -> write to %s\n", inFile, id, series, season, episode, title, year, posterPath, outFile)
	//TODO implement me
	return errors.New("implement me - atomicparsley.go")
}
