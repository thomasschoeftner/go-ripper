package tag

import (
	"go-cli/task"
	"go-ripper/ripper"
	"errors"
	"go-ripper/targetinfo"
	"go-ripper/metainfo/video"
	"go-ripper/metainfo"
	"fmt"
	"go-ripper/files"
	"go-cli/commons"
	"path/filepath"
	"strconv"
)

type VideoTagger interface {
	TagMovie(inFile string, outFile string, id string, title string, year string, posterPath string) error
	TagEpisode(inFile string, outFile string, id string, series string, season int, episode int, title string, year string, posterPath string) error
}

var NewVideoTagger func(conf *ripper.TagConfig, lazy bool, printf commons.FormatPrinter) (VideoTagger, error)

func TagVideo(ctx task.Context) task.HandlerFunc {
	conf := ctx.Config.(*ripper.AppConf)

	if nil == NewVideoTagger {
		return ripper.ErrorHandler(errors.New("video-tagger is undefined"))
	}
	tagger, err := NewVideoTagger(conf.Tag, ctx.RunLazy, ctx.Printf)
	if err != nil {
		return ripper.ErrorHandler(err)
	}

	expectedExtension := conf.Output.Video
	invalidFileNameChars := conf.Output.InvalidCharactersInFileName
	evacuate := files.PrepareEvacuation(filepath.Join(conf.WorkDirectory, files.TEMP_DIR_NAME)) //replace spaces with underscores

	return func(job task.Job) ([]task.Job, error) {
		target := ripper.GetTargetFileFromJob(job)
		ctx.Printf("tag video - target %s\n", target)
		ti, err := targetinfo.ForTarget(conf.WorkDirectory, target)
		if err != nil {
			return nil, err
		}

		in, inputIsOriginal, err := findInputFile(ti, conf.WorkDirectory, expectedExtension)
		if err != nil {
			return nil, err
		}

		//ctx.Printf("  input=%s\n", in)   //TODO remove
		var evacuated *files.Evacuated
		if inputIsOriginal {
			evacuated, err = evacuate(in).By(files.Copying)
		} else { //already preprocessed in work-directory is input
			evacuated, err = evacuate(in).By(files.Moving)
		}
		if err != nil {
			return nil, err
		}
		defer evacuated.Discard()

		var subPathElems []string
		switch ti.GetType() {
		case targetinfo.TARGETINFO_TYPE_MOVIE:
			subPathElems, err = tagMovie(tagger, ti.(*targetinfo.Movie), conf.MetaInfoRepo, evacuated.Path())
		case targetinfo.TARGETINFO_TPYE_EPISODE:
			subPathElems, err = tagEpisode(tagger, ti.(*targetinfo.Episode), conf.MetaInfoRepo, evacuated.Path())
		default:
			err = fmt.Errorf("unknown type of video target-info found: %s", ti.GetType())
		}
		if err != nil {
			return nil, err
		}
		if 0 == len(subPathElems) {
			return nil, fmt.Errorf("empty output path returned")
		}

		dst := buildDestinationPath(ctx.OutputDir, subPathElems, invalidFileNameChars)

		//move evacuated file to output folder
		if err := files.CreateFolderStructure(filepath.Dir(dst)); err != nil {
			return nil, err
		}
		err = evacuated.MoveTo(dst)
		if err != nil {
			return nil, err
		}
		return []task.Job{job}, nil
	}
}

func buildDestinationPath(outputDir string, pathElems []string, invalidFileNameChars string) string {
	dstPathElems := []string {outputDir}
	for _, pathElem := range pathElems {
		dstPathElems = append(dstPathElems, commons.RemoveCharacters(pathElem, invalidFileNameChars))
	}

	return filepath.Join(dstPathElems...)
}

func tagMovie(tagger VideoTagger, ti *targetinfo.Movie, metaInfoRepo string, inputFile string) ([]string, error) {
	noPath := []string{}
	movieMi := video.MovieMetaInfo{}
	err := metainfo.ReadMetaInfo(video.MovieFileName(metaInfoRepo, ti.GetId()), &movieMi)
	if err != nil {
		return noPath, err
	}
	if 0 == len(movieMi.Id) {
		return noPath, fmt.Errorf("could not find meta-info for movie: %s\n", ti.String())
	}
	imgFile := metainfo.ImageFileName(metaInfoRepo, movieMi.Id, files.GetExtension(movieMi.Poster))
	//TODO check if missing poster image is actually an error

	err = tagger.TagMovie(inputFile, inputFile, movieMi.Id, movieMi.Title, movieMi.Year, imgFile)
	if err != nil {
		return noPath, err
	}

	_, ext := files.SplitExtension(inputFile)
	return []string{files.WithExtension(movieMi.Title, ext)}, nil
}

const templateEpisodeFilename = "%s-s%02de%02d-%s"

func tagEpisode(tagger VideoTagger, ti *targetinfo.Episode, metaInfoRepo string, inputFile string) ([]string, error) {
	noPath := []string{}
	episodeMi := video.EpisodeMetaInfo{}
	err := metainfo.ReadMetaInfo(video.EpisodeFileName(metaInfoRepo, ti.Id, ti.Season, ti.Episode), &episodeMi)
	if err != nil {
		return noPath, err
	}
	if 0 == len(episodeMi.Id) {
		return noPath, fmt.Errorf("could not find meta-info for episode: %s\n", ti.String())
	}

	seriesMi := video.SeriesMetaInfo{}
	err = metainfo.ReadMetaInfo(video.SeriesFileName(metaInfoRepo, ti.Id), &seriesMi)
	if err != nil {
		return noPath, err
	}
	if 0 == len(seriesMi.Id) {
		return noPath, fmt.Errorf("could not find meta-info for series: %s\n", ti.String())
	}
	imgFile := metainfo.ImageFileName(metaInfoRepo, seriesMi.Id, files.GetExtension(seriesMi.Poster))

	err = tagger.TagEpisode(inputFile, inputFile, seriesMi.Id, seriesMi.Title, episodeMi.Season, episodeMi.Episode, episodeMi.Title, episodeMi.Year, imgFile)
	if err != nil {
		return noPath, err
	}

	_, ext := files.SplitExtension(inputFile)
	fName := files.WithExtension(fmt.Sprintf(templateEpisodeFilename, seriesMi.Title, episodeMi.Season, episodeMi.Episode, episodeMi.Title), ext)
	return []string {seriesMi.Title, strconv.Itoa(episodeMi.Season), fName}, nil
}

func findInputFile(ti targetinfo.TargetInfo, workDir string, expectedExtension string) (string, bool, error) {
	// check work directory for a pre-processed inFile in appropriate format (e.g. a ripped video in .mp4 inFile)
	preprocessed, err := ripper.GetProcessingArtifactPathFor(workDir, ti.GetFolder(), ti.GetFile(), expectedExtension)
	if err != nil {
		return "", false, err
	}
	if exists, err := files.Exists(preprocessed); err != nil {
		return "", false, err
	} else if exists {
		return preprocessed, false, nil
	}

	// if no preprocessed input is available, check if the source inFile can be tagged directly (e.g. f
	_, extension := files.SplitExtension(ti.GetFile())
	if extension == expectedExtension {
		return filepath.Join(ti.GetFolder(), ti.GetFile()), true, nil
	} else {
		return "", false, fmt.Errorf("unable to find appropriate input file (\"%s\") for meta-info tagging", expectedExtension)
	}
}
