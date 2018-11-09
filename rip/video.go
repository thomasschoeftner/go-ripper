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
		handbrake, err := createHandbrakeRipper(conf.Rip.Video.Handbrake, ctx.RunLazy, ctx.Printf)
		if err != nil {
			return RipError(err)
		}
		return Rip(ctx, handbrake)
	default:
		return RipError(fmt.Errorf("unable to create video ripper of type \"%s\"", ripperType))
	}
}
