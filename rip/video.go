package rip

import (
	"go-cli/task"
	"go-ripper/ripper"
	"fmt"
	"go-ripper/processor"
)

func RipVideo(ctx task.Context) task.HandlerFunc {
	conf := ctx.Config.(*ripper.AppConf)
	ripperType := conf.Rip.Video.Ripper

	var rip processor.Processor
	var err error

	switch ripperType {
	case CONF_RIPPER_HANDBRAKE:
		rip, err = handbrakeRipper(conf.Rip.Video.Handbrake, ctx.Printf, conf.WorkDirectory)
	default:
		err = fmt.Errorf("unable to create video ripper of type \"%s\"", ripperType)
	}

	if err != nil {
		return ripper.ErrorHandler(err)
	}
	return processor.Process(ctx, rip, ripperType,
		processor.DefaultCheckLazy(ctx.RunLazy, conf.Output.Video),
		processor.DefaultInputFileFor(conf.WorkDirectory, conf.Rip.Video.AllowedInputExtensions),
		processor.DefaultOutputFileFor(conf.WorkDirectory, conf.Output.Video))
}
