package tag

import (
	"go-ripper/ripper"
	"go-cli/commons"
	"go-cli/cli"
	"time"
	"go-ripper/files"
	"fmt"
	"strconv"
	"os"
	"io"
)

const CONF_ATOMICPARSLEY_TAGGER = "atomicparsley"
const (
	paramOutput = "-o"

	paramTitle = "--title"
	paramYear = "--year"
	paramPoster = "--artwork"
	paramDescription = "--description"

	paramSeriesName = "--TVShowName"
	paramSeason = "--TVSeasonNum"
	paramEpisode = "--TVEpisodeNum"
	paramEpisodeName = "--TVEpisode"

	paramComment = "--comment"
	paramGenre = "--genre"
	paramAlbum = "--album"
	paramDisc = "--disk"
)

func NewAtomicParsleyVideoTagger(conf *ripper.TagConfig, lazy bool, printf commons.FormatPrinter) (VideoTagger, error)  {
	if lazy {
		printf("WARNING - video tagging via AtomicParsley is NOT lazy - ie files will always be tagged/written!\n")
	}

	timeout, err := time.ParseDuration(conf.Video.AtomicParsley.Timeout)
	if err != nil {
		return nil, err
	}

	apConf := conf.Video.AtomicParsley
	path := apConf.Path
	exists, err := files.Exists(path)
	if err != nil {
		return nil, err
	} else if !exists {
		return nil, fmt.Errorf("unable to find AtomicParsley binary at \"%s\"\n", path)
	}

	ap := &atomicParsleyVideoTagger{timeout: timeout, path: path, printf: printf.WithIndent(2)}
	if apConf.ShowErrorOutput {
		ap.errOut = os.Stderr
	}
	if apConf.ShowStandardOutput {
		ap.stdOut = os.Stdout
	}
	return ap, nil
}

type atomicParsleyVideoTagger struct {
	timeout time.Duration
	path string
	printf commons.FormatPrinter
	errOut io.Writer
	stdOut io.Writer
}

func (ap *atomicParsleyVideoTagger) TagMovie(inFile string, outFile string, id string, title string, year string, posterPath string) error {
	ap.printf("AtomicParsley tags {id=%s, title=%s, year=%s, image=%s}\n", id, title, year, posterPath)
	ap.printf("-> write to \"%s\"\n", outFile)
	cmd := cli.Command(ap.path, ap.timeout).WithQuotes(" ", '"').
		WithArgument(inFile).
		WithParam(paramTitle, title, "").
		WithParam(paramPoster, posterPath, "").
		WithParam(paramYear, year, "").
		WithParam(paramOutput, outFile, "")
	//ap.printf(">>>> %s\n", cmd.String())
	return cmd.ExecuteSync(ap.stdOut, ap.errOut)
}

func (ap *atomicParsleyVideoTagger) TagEpisode(inFile string, outFile string, id string, series string, season int, episode int, title string, year string, posterPath string) error {
	ap.printf("AtomicParsley tags {id=%s, series=%s, season=%d, episode=%d, title=%s, year=%s, image=%s}\n", id, series, season, episode, title, year, posterPath)
	ap.printf("-> write to \"%s\"\n", outFile)
	cmd := cli.Command(ap.path, ap.timeout).WithQuotes(" ", '"').
		WithArgument(inFile).
		WithParam(paramTitle, title, "").
		WithParam(paramPoster, posterPath, "").
		WithParam(paramYear, year, "").
		WithParam(paramSeriesName, series, "").
		WithParam(paramEpisode, strconv.Itoa(episode), "").
		WithParam(paramSeason, strconv.Itoa(season), "").
		WithParam(paramEpisodeName, title, "").
		WithParam(paramOutput, outFile, "")
	//ap.printf(">>>> %s\n", cmd.String())
	return cmd.ExecuteSync(ap.stdOut, ap.errOut)
}
