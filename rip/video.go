package rip

import (
	"go-cli/task"
	"go-ripper/ripper"
	"fmt"
)

func RipVideo(ctx task.Context) task.HandlerFunc {
	conf := ctx.Config.(*ripper.AppConf)
	ripperType := conf.Rip.Video.Ripper

	switch ripperType {
	case CONF_RIPPER_HANDBRAKE:
		handbrake, err := handbrakeRipper(conf.Rip.Video.Handbrake, ctx.RunLazy, ctx.Printf, conf.WorkDirectory)
		if err != nil {
			return ripper.ErrorHandler(err)
		}
		return Rip(ctx, handbrake, conf.Rip.Video.AllowedInputExtensions, conf.Output.Video)
	default:
		return ripper.ErrorHandler(fmt.Errorf("unable to create video ripper of type \"%s\"", ripperType))
	}
}
