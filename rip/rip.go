package rip

import (
	"go-cli/task"
	"go-ripper/ripper"
	"errors"
)

type Ripper interface {
	process(inFile string, outFile string) error
}


func Rip(ctx task.Context, theRipper Ripper) task.HandlerFunc {
	conf := ctx.Config.(*ripper.AppConf)
	println(conf.Rip.Video.Ripper)
	//expectedExtension := conf.Output.Video

	return func(job task.Job) ([]task.Job, error) {
		//TODO implement me
		//return []task.Job{job}, nil
		return []task.Job{}, errors.New("IMPLEMENT ME")
	}
}


func RipError(err error) task.HandlerFunc {
	return func (job task.Job) ([]task.Job, error) {
		return []task.Job{}, err
	}
}