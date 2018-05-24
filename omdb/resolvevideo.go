package omdb

import (
	"go-cli/task"
	"errors"
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
	//conf := ctx.Config.(*ripper.AppConf)

	if tFactory == nil {
		return errorFunc(errors.New("processing error - omdb-tokens missing"))
	}

	return func(j task.Job) ([]task.Job, error) {
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