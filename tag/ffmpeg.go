package tag

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/thomasschoeftner/go-cli/cli"
	"github.com/thomasschoeftner/go-cli/commons"
	"github.com/thomasschoeftner/go-ripper/files"
	"github.com/thomasschoeftner/go-ripper/ripper"
)

const conf_tagger_ffmpeg = "ffmpeg"

const (
	// ffmpeg -i <input-file>.mp4 -i <artwork>.jpg -map 0 -map 1  -metadata title="<title>" -metadata year="<year>" -metadata genre="<genre>" -c copy -disposition:v:1 attached_pic <output>.mp4
	ffmpeg_paramInputFile = "-i"

	ffmpeg_paramMetaData = "-metadata"

	ffmpeg_tagTitleKey       = "title"
	ffmpeg_tagDescriptionKey = "description"
	ffmpeg_tagGenreKey       = "genre"
	ffmpeg_tagYearKey        = "year"

	ffmpeg_tagSeriesNameKey = "show"
	ffmpeg_tagGroupingKey   = "grouping"
	//ffmpeg_tagSeasonKey = "" // TODO - look up
	ffmpeg_tagEpisodeKey = "episode_id"
	//ffmpeg_tagEpisodeNameKey = "--TVEpisode" // TODO - look up

	ffmpeg_tagCommentKey = "comment"
	ffmpeg_tagAlbumKey   = "album"
	ffmpeg_tagTrackKey   = "track"
)

func createFFMPEGVideoTagger(conf *ripper.AppConf, lazy bool, printf commons.FormatPrinter) (MovieTagger, EpisodeTagger, error) {
	apConf := conf.Tag.Video.FFMPEG
	workDir := conf.WorkDirectory
	tagCtx := &ffmpegTagger{}
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

type ffmpegTagger struct {
	path     string
	timeout  time.Duration
	printf   commons.FormatPrinter
	stdout   io.Writer
	errout   io.Writer
	evacuate files.EvacuatorFunc
	tempDir  string
}

func (ffmpeg *ffmpegTagger) movie(inFile string, outFile string, id string, title string, year string, posterPath string) error {
	cmd := cli.Command(ffmpeg.path, ffmpeg.timeout). //WithQuotes(" ", '"').
								WithParam(ffmpeg_paramInputFile, inFile, "").
								WithParam(ffmpeg_paramInputFile, posterPath, "").
								WithParam("-map", "0", "").
								WithParam("-map", "1", "").
								WithParam(ffmpeg_paramMetaData, fmt.Sprintf("%s=%s", ffmpeg_tagTitleKey, title), "").
								WithParam(ffmpeg_paramMetaData, fmt.Sprintf("%s=%s", ffmpeg_tagYearKey, year), "").
								WithParam("-c", "copy", "").                       // do not perform encode step
								WithParam("-disposition:v:1", "attached_pic", ""). // use 2nd input file as artwork
								WithArgument(fmt.Sprintf("%s", outFile))

	ffmpeg.printf(">>>> %s\n", cmd.String()) // TODO - comment out
	return cmd.ExecuteSync(ffmpeg.stdout, ffmpeg.errout)
}

func (ffmpeg *ffmpegTagger) episode(inFile string, outFile string, id string, series string, season int, episode int, title string, year string, posterPath string) error {
	cmd := cli.Command(ffmpeg.path, ffmpeg.timeout). //WithQuotes(" ", '"').
								WithParam(ffmpeg_paramInputFile, inFile, "").
								WithParam(ffmpeg_paramInputFile, posterPath, "").
								WithParam("-map", "0", "").
								WithParam("-map", "1", "").
								WithParam(ffmpeg_paramMetaData, fmt.Sprintf("%s=%s", ffmpeg_tagTitleKey, title), "").
								WithParam(ffmpeg_paramMetaData, fmt.Sprintf("%s=%s", ffmpeg_tagYearKey, year), "").
								WithParam(ffmpeg_paramMetaData, fmt.Sprintf("%s=%s", ffmpeg_tagSeriesNameKey, series), "").
								WithParam(ffmpeg_paramMetaData, fmt.Sprintf("%s=%d", ffmpeg_tagGroupingKey, season), "").
								WithParam(ffmpeg_paramMetaData, fmt.Sprintf("%s=%d", ffmpeg_tagEpisodeKey, episode), "").
								WithParam("-c", "copy", "").                       // do not perform encode step
								WithParam("-disposition:v:1", "attached_pic", ""). // use 2nd input file as artwork
								WithArgument(fmt.Sprintf("%s", outFile))

	ffmpeg.printf(">>>> %s\n", cmd.String())
	return cmd.ExecuteSync(ffmpeg.stdout, ffmpeg.errout)
}
