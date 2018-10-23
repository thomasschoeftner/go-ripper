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
	tagger, err := NewVideoTagger(conf.Tag, ctx.RunLazy, ctx.Printf.WithIndent(2))
	if err != nil {
		return ripper.ErrorHandler(err)
	}

	return func (job task.Job) ([]task.Job, error) {
		target := ripper.GetTargetFileFromJob(job)
		ctx.Printf("tag video - target %s\n", target)
		ti, err := targetinfo.ForTarget(conf.WorkDirectory, target)
		if err != nil {
			return nil, err
		}

		switch ti.GetType() {
		case targetinfo.TARGETINFO_TYPE_MOVIE:
			err = tagMovie(tagger, ti, conf.WorkDirectory, conf.MetaInfoRepo, conf.Output.Video)
		case targetinfo.TARGETINFO_TPYE_EPISODE:
			err = tagEpisode(tagger, ti, conf.WorkDirectory, conf.MetaInfoRepo, conf.Output.Video)
		default:
			err = fmt.Errorf("unknown type of video target-info found: %s", ti.GetType())
		}
		//NOT do not return nil on error!!!
		return []task.Job{job}, err
	}
}

func tagMovie(tagger VideoTagger, ti targetinfo.TargetInfo, workDir string, metaInfoRepo string, videoOutputExtension string) error {
	fileToTag, outputFile, err := findInputOutputFiles(ti, workDir, videoOutputExtension)
	if err != nil {
		return err
	}

	movieMi := video.MovieMetaInfo{}
	err = metainfo.ReadMetaInfo(video.MovieFileName(metaInfoRepo, ti.GetId()), &movieMi)
	if err != nil {
		return err
	}
	if 0 == len(movieMi.Id) {
		return fmt.Errorf("could not find meta-info for movie: %s\n", ti.String())
	}
	imgFile := metainfo.ImageFileName(metaInfoRepo, movieMi.Id, files.GetExtension(movieMi.Poster))
	//TODO check if missing poster image is actually an error

	return tagger.TagMovie(fileToTag, outputFile, movieMi.Id, movieMi.Title, movieMi.Year, imgFile)
}

func tagEpisode(tagger VideoTagger, ti targetinfo.TargetInfo, workDir string, metaInfoRepo string, videoOutputExtension string) error {
	fileToTag, outputFile, err := findInputOutputFiles(ti, workDir, videoOutputExtension)
	if err != nil {
		return err
	}
	episodeTi := ti.(*targetinfo.Episode)
	episodeMi := video.EpisodeMetaInfo{}
	err = metainfo.ReadMetaInfo(video.EpisodeFileName(metaInfoRepo, episodeTi.Id, episodeTi.Season, episodeTi.Episode), &episodeMi)
	if err != nil {
		return err
	}
	if 0 == len(episodeMi.Id) {
		return fmt.Errorf("could not find meta-info for episode: %s\n", ti.String())
	}

	seriesMi := video.SeriesMetaInfo{}
	err = metainfo.ReadMetaInfo(video.SeriesFileName(metaInfoRepo, episodeTi.Id), &seriesMi)
	if err != nil {
		return err
	}
	if 0 == len(seriesMi.Id) {
		return fmt.Errorf("could not find meta-info for series: %s\n", ti.String())
	}
	imgFile := metainfo.ImageFileName(metaInfoRepo, seriesMi.Id, files.GetExtension(seriesMi.Poster))

	return tagger.TagEpisode(fileToTag, outputFile, seriesMi.Id, seriesMi.Title, episodeMi.Season, episodeMi.Episode, episodeMi.Title, episodeMi.Year, imgFile)
}

func findInputOutputFiles(ti targetinfo.TargetInfo, workDir string, expectedExtension string) (string, string, error) {
	folder, err := ripper.GetWorkPathForTargetFileFolder(workDir, ti.GetFolder())
	if err != nil {
		return "", "", err
	}

	// check work directory for a pre-processed inFile in appropriate format (e.g. a ripped video in .mp4 inFile)
	fName, extension := files.SplitExtension(ti.GetFile())
	preprocessed := filepath.Join(folder, files.WithExtension(fName, expectedExtension))
	if exists, _ := files.Exists(preprocessed); exists {
		return preprocessed, preprocessed, nil
	}

	// if no preprocessed input is available, check if the source inFile can be tagged directly (e.g. if it is an .mp4 video)
	if extension == expectedExtension {
		return filepath.Join(ti.GetFolder(), ti.GetFile()), preprocessed, nil
	} else {
		return "", "", errors.New("unable to find appropriate input inFile for meta-info tagging")
	}
}
