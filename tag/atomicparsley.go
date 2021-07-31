package tag

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/thomasschoeftner/go-cli/cli"
	"github.com/thomasschoeftner/go-cli/commons"
	"github.com/thomasschoeftner/go-ripper/files"
	"github.com/thomasschoeftner/go-ripper/ripper"
)

const conf_tagger_atomicparsley = "atomicparsley"
const (
	paramOutputFile   = "-o"
	argumentOverwrite = "--overWrite"

	paramTitle       = "--title"
	paramYear        = "--year"
	paramPoster      = "--artwork"
	paramDescription = "--description"

	paramSeriesName  = "--TVShowName"
	paramSeason      = "--TVSeasonNum"
	paramEpisode     = "--TVEpisodeNum"
	paramEpisodeName = "--TVEpisode"

	paramComment = "--comment"
	paramGenre   = "--genre"
	paramAlbum   = "--album"
	paramDisc    = "--disk"
)

func createAtomicParsleyVideoTagger(conf *ripper.AppConf, lazy bool, printf commons.FormatPrinter) (MovieTagger, EpisodeTagger, error) {
	apConf := conf.Tag.Video.AtomicParsley
	workDir := conf.WorkDirectory
	tagCtx := &atomicParsley{}
	var err error

	tagCtx.timeout, err = time.ParseDuration(apConf.Timeout)
	if err != nil {
		return nil, nil, err
	}

	tagCtx.path = apConf.Path

	if apConf.ShowErrorOutput {
		tagCtx.errout = os.Stderr
	}
	if apConf.ShowStandardOutput {
		tagCtx.stdout = os.Stdout
	}

	tagCtx.printf = printf.WithIndent(2)
	tagCtx.tempDir = filepath.Join(workDir, files.TEMP_DIR_NAME)
	return tagCtx.movie, tagCtx.episode, nil
}

type atomicParsley struct {
	path     string
	timeout  time.Duration
	printf   commons.FormatPrinter
	stdout   io.Writer
	errout   io.Writer
	evacuate files.EvacuatorFunc
	tempDir  string
}

func (ap *atomicParsley) movie(inFile string, outFile string, id string, title string, year string, posterPath string) error {
	evacuate := files.PrepareEvacuation(ap.tempDir)
	evacuated, err := evacuate(inFile).By(files.Copying)

	if err == nil {
		defer evacuated.Discard()
		resultFile := evacuated.WithSuffix(".tagged")
		//ap.printf("AtomicParsley tags \"%s\"\n", inFile)
		//ap.printf("using {id=%s, title=%s, year=%s, image=%s}\n", id, title, year, posterPath)
		//ap.printf("-> write to \"%s\"\n", outFile)
		cmd := cli.Command(ap.path, ap.timeout).WithQuotes(" ", '"').
			WithArgument(evacuated.Path()).
			WithParam(paramTitle, title, "").
			WithParam(paramPoster, posterPath, "").
			WithParam(paramYear, year, "")
		cmd = cmd.WithParam(paramOutputFile, resultFile, "")

		//ap.printf(">>>> %s\n", cmd.String())
		err = cmd.ExecuteSync(ap.stdout, ap.errout)
		if err == nil {
			err = moveTo(resultFile, outFile)
		}
	}
	return err
}

func (ap *atomicParsley) episode(inFile string, outFile string, id string, series string, season int, episode int, title string, year string, posterPath string) error {
	evacuate := files.PrepareEvacuation(ap.tempDir)
	evacuated, err := evacuate(inFile).By(files.Copying)

	if err == nil {
		defer evacuated.Discard()
		resultFile := evacuated.WithSuffix(".tagged")

		//ap.printf("AtomicParsley tags \"%s\"\n", inFile)
		//ap.printf("using {id=%s, series=%s, season=%d, episode=%d, title=%s, year=%s, image=%s}\n", id, series, season, episode, title, year, posterPath)
		//ap.printf("-> write to \"%s\"\n", outFile)
		cmd := cli.Command(ap.path, ap.timeout).WithQuotes(" ", '"').
			WithArgument(evacuated.Path()).
			WithParam(paramTitle, title, "").
			WithParam(paramPoster, posterPath, "").
			WithParam(paramYear, year, "").
			WithParam(paramSeriesName, series, "").
			WithParam(paramEpisode, strconv.Itoa(episode), "").
			WithParam(paramSeason, strconv.Itoa(season), "").
			WithParam(paramEpisodeName, title, "")
		cmd = cmd.WithParam(paramOutputFile, resultFile, "")

		//ap.printf(">>>> %s\n", cmd.String())
		err = cmd.ExecuteSync(ap.stdout, ap.errout)
		if err == nil {
			err = moveTo(resultFile, outFile)
		}
	}
	return err
}

func moveTo(intermediateFile string, outFile string) error {
	err := files.CreateFolderStructure(filepath.Dir(outFile))
	if err != nil {
		return err
	}

	err = os.Rename(intermediateFile, outFile)
	return err
}
