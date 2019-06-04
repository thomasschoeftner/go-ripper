package tag

import (
	"fmt"
	"path/filepath"
	"strconv"
	"go-cli/task"
	"go-ripper/ripper"
	"go-ripper/targetinfo"
	"go-ripper/metainfo/video"
	"go-ripper/metainfo"
	"go-ripper/files"
	"go-cli/commons"
	"go-ripper/processor"
)

type MovieTagger func(inFile string, outFile string, id string, title string, year string, posterPath string) error
type EpisodeTagger func(inFile string, outFile string, id string, series string, season int, episode int, title string, year string, posterPath string) error


type TaggerFactory func(conf *ripper.AppConf, lazy bool, printf commons.FormatPrinter, workDir string) (MovieTagger, EpisodeTagger, error)
var TaggerFactories map[string]TaggerFactory

func init() {
	TaggerFactories = make(map[string]TaggerFactory)
	TaggerFactories[conf_tagger_atomicparsley] = createAtomicParsleyVideoTagger
}

func TagVideo(ctx task.Context) task.HandlerFunc {
	conf := ctx.Config.(*ripper.AppConf)
	taggerType := conf.Tag.Video.Tagger

	var movieTagger MovieTagger
	var episodeTagger EpisodeTagger
	var err error

	tf := TaggerFactories[taggerType]
	if tf == nil {
		err = fmt.Errorf("unknown video tagger configured: \"%s\"", conf.Tag.Video.Tagger)
	} else {
		movieTagger, episodeTagger, err = createAtomicParsleyVideoTagger(conf, ctx.RunLazy, ctx.Printf, conf.WorkDirectory)
	}

	if err != nil {
		return ripper.ErrorHandler(err)
	} else {
		return processor.Process(ctx, getProcessor(conf, movieTagger, episodeTagger), taggerType,
			processor.NeverLazy(ctx.RunLazy, taggerType, ctx.Printf),
			processor.DefaultInputFileFor([]string{conf.Output.Video}),
			processor.DefaultOutputFileFor(conf.Output.Video))
	}
}

func getProcessor(conf *ripper.AppConf, movieTagger MovieTagger, episodeTagger EpisodeTagger) processor.Processor {
	return func(ti targetinfo.TargetInfo, inFile string, outFile string) error {
		var err error

		switch ti.GetType() {
		case targetinfo.TARGETINFO_TYPE_MOVIE:
			err = tagMovie(movieTagger, conf, ti.(*targetinfo.Movie), inFile)
		case targetinfo.TARGETINFO_TPYE_EPISODE:
			err = tagEpisode(episodeTagger, conf, ti.(*targetinfo.Episode), inFile)
		default:
			err = fmt.Errorf("unknown type of video target-info found: %s", ti.GetType())
		}
		return err
	}
}

func tagMovie(tag MovieTagger, conf *ripper.AppConf, ti *targetinfo.Movie, inputFile string) error {
	movieMi := video.MovieMetaInfo{}
	err := metainfo.ReadMetaInfo(video.MovieFileName(conf.MetaInfoRepo, ti.GetId()), &movieMi)
	if err != nil {
		return err
	}

	if 0 == len(movieMi.Id) {
		return fmt.Errorf("could not find meta-info for movie: %s\n", ti.String())
	}
	imgFile := metainfo.ImageFileName(conf.MetaInfoRepo, movieMi.Id, files.GetExtension(movieMi.Poster))
	//TODO check if missing poster image is actually an error

	ext := files.GetExtension(inputFile)
	outputFile := buildDestinationPath(conf.Output.InvalidCharactersInFileName, conf.OutputDirectory, files.WithExtension(movieMi.Title, ext))

	err = tag(inputFile, outputFile, movieMi.Id, movieMi.Title, movieMi.Year, imgFile)
	return err
}

const templateEpisodeFilename = "%s-s%02de%02d-%s"

func tagEpisode(tag EpisodeTagger, conf *ripper.AppConf, ti *targetinfo.Episode, inputFile string) error {
	episodeMi := video.EpisodeMetaInfo{}
	err := metainfo.ReadMetaInfo(video.EpisodeFileName(conf.MetaInfoRepo, ti.Id, ti.Season, ti.Episode), &episodeMi)
	if err != nil {
		return err
	}
	if 0 == len(episodeMi.Id) {
		return fmt.Errorf("could not find meta-info for episode: %s\n", ti.String())
	}

	seriesMi := video.SeriesMetaInfo{}
	err = metainfo.ReadMetaInfo(video.SeriesFileName(conf.MetaInfoRepo, ti.Id), &seriesMi)
	if err != nil {
		return err
	}
	if 0 == len(seriesMi.Id) {
		return fmt.Errorf("could not find meta-info for series: %s\n", ti.String())
	}
	imgFile := metainfo.ImageFileName(conf.MetaInfoRepo, seriesMi.Id, files.GetExtension(seriesMi.Poster))

	fName := files.WithExtension(fmt.Sprintf(templateEpisodeFilename, seriesMi.Title, episodeMi.Season, episodeMi.Episode, episodeMi.Title), files.GetExtension(inputFile))
	outputFile := buildDestinationPath(conf.Output.InvalidCharactersInFileName, conf.OutputDirectory, seriesMi.Title, strconv.Itoa(episodeMi.Season), fName)

	err = tag(inputFile, outputFile, seriesMi.Id, seriesMi.Title, episodeMi.Season, episodeMi.Episode, episodeMi.Title, episodeMi.Year, imgFile)
	return err
}

func buildDestinationPath(invalidFileNameChars string, outputDir string, pathElems ...string) string {
	dstPathElems := []string {outputDir}
	for _, pathElem := range pathElems {
		dstPathElems = append(dstPathElems, commons.RemoveCharacters(pathElem, invalidFileNameChars))
	}

	return filepath.Join(dstPathElems...)
}
