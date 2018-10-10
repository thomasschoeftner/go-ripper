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
	TagMovie(file string, id string, title string, year string, posterPath string) error
	TagEpisode(file string, id string, series string, season int, episode int, title string, year string, posterPath string) error
}

var NewVideoTagger func(conf *ripper.TagConfig) (VideoTagger, error)


func TagVideo(ctx task.Context) task.HandlerFunc {
	conf := ctx.Config.(*ripper.AppConf)

	if nil == NewVideoTagger {
		return ripper.ErrorHandler(errors.New("video-tagger is undefined"))
	}
	tagger, err := NewVideoTagger(conf.Tag)
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
		printf := ctx.Printf.WithIndent(2)

		switch ti.GetType() {
		case targetinfo.TARGETINFO_TYPE_MOVIE:
			err = tagMovie(tagger, ti, conf, printf)
		case targetinfo.TARGETINFO_TPYE_EPISODE:
			err = tagEpisode(tagger, ti, conf, printf)
		default:
			err = fmt.Errorf("unknown type of video target-info found: %s", ti.GetType())
		}
		//NOT do not return nil on error!!!
		return []task.Job{job}, err
	}
}

func tagMovie(tagger VideoTagger, ti targetinfo.TargetInfo, conf *ripper.AppConf, printf commons.FormatPrinter) error {
	fileToTag, err := chooseInputFile(ti, conf.WorkDirectory, conf.Output.Video)
	if err != nil {
		return err
	}

	movieMi := video.MovieMetaInfo{}
	metainfo.ReadMetaInfo(video.MovieFileName(conf.MetaInfoRepo, ti.GetId()), &movieMi)
	if 0 == len(movieMi.Id) {
		return fmt.Errorf("could not find meta-info for movie: %s\n", ti.String())
	}
	imgFile := metainfo.ImageFileName(conf.MetaInfoRepo, movieMi.Id, files.Extension(movieMi.Poster))
	//TODO check if missing poster image is actually an error

	return tagger.TagMovie(fileToTag, movieMi.Id, movieMi.Title, movieMi.Year, imgFile)
}

func tagEpisode(tagger VideoTagger, ti targetinfo.TargetInfo, conf *ripper.AppConf, printf commons.FormatPrinter) error {
	fileToTag, err := chooseInputFile(ti, conf.WorkDirectory, conf.Output.Video)
	if err != nil {
		return err
	}

	episodeTi := ti.(*targetinfo.Episode)
	episodeMi := video.EpisodeMetaInfo{}
	metainfo.ReadMetaInfo(video.EpisodeFileName(conf.MetaInfoRepo, episodeTi.Id, episodeTi.Season, episodeTi.Episode), &episodeMi)
	if 0 == len(episodeMi.Id) {
		return fmt.Errorf("could not find meta-info for episode: %s\n", ti.String())
	}
	seriesMi := video.SeriesMetaInfo{}
	metainfo.ReadMetaInfo(video.SeriesFileName(conf.MetaInfoRepo, episodeTi.Id), &seriesMi)
	if 0 == len(seriesMi.Id) {
		return fmt.Errorf("could not find meta-info for series: %s\n", ti.String())
	}
	imgFile := metainfo.ImageFileName(conf.MetaInfoRepo, seriesMi.Id, files.Extension(seriesMi.Poster))

	return tagger.TagEpisode(fileToTag, episodeMi.Id, seriesMi.Title, episodeMi.Season, episodeMi.Episode, episodeMi.Title, episodeMi.Year, imgFile)
}

func chooseInputFile(ti targetinfo.TargetInfo, workDir string, expectedExtension string) (string, error) {
	folder, err := ripper.GetWorkPathForTargetFileFolder(workDir, ti.GetFolder())
	if err != nil {
		return "", err
	}
	fName, extension := files.SplitExtension(ti.GetFile())
	preprocessed := filepath.Join(folder, fName + "." + expectedExtension)
	fmt.Printf("expect=%s, fname=%s, ext=%s, preproc=%s\n", expectedExtension, fName, extension, preprocessed)
	if exists, _ := files.Exists(preprocessed); exists {
		return preprocessed, nil
	}

	if extension == expectedExtension {
		return filepath.Join(ti.GetFolder(), ti.GetFile()), nil
	} else {
		return "", errors.New("unable to find appropriate input file for meta-info tagging")
	}
}