package metainfo

import (
	"go-cli/task"
	"errors"
	"go-ripper/ripper"
	"go-ripper/targetinfo"
)

var tFactory *tokenFactory

func ResolveVideo(omdbTokens []string) task.Handler {
	if tFactory == nil {
		tf, err := newTokenFactory(omdbTokens)
		if err == nil {
			tFactory = tf
		}
	}
	return resolveVideo
}

func resolveVideo(ctx task.Context) task.HandlerFunc {
	if tFactory == nil {
		return errorFunc(errors.New("processing error - omdb-tokens missing"))
	}

	conf := ctx.Config.(*ripper.AppConf)
	return func(job task.Job) ([]task.Job, error) {
		file := job[ripper.JobField_Path]
		id := job[ripper.JobField_TargetId]

		ctx.Printf("process video %s\n", file)
		printf := ctx.Printf.WithIndent(2)
		// read task definition
		ti, err := targetinfo.Read(ripper.GetTempPathFor(job, conf), id)
		if err != nil {
			return nil, err
		}
		printf("recovered target-info: %s\n", ti.String())

		//TODO implement
		//conf.Resolve.Video.Omdb
		return nil, nil
	}
}


func errorFunc(reason error) task.HandlerFunc{
	return func(task.Job) ([]task.Job, error) {
		return nil, reason
	}
}