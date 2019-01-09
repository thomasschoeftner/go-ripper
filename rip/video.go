package rip

import (
	"go-cli/task"
	"go-ripper/ripper"
	"fmt"
	"go-ripper/processor"
	"go-cli/commons"
)

type RipperFactory func(conf *ripper.HandbrakeConfig, printf commons.FormatPrinter, workDir string) (processor.Processor, error)
var RipperFactories map[string]RipperFactory

func init() {
	RipperFactories = make(map[string]RipperFactory)
	RipperFactories[CONF_RIPPER_HANDBRAKE] = createHandbrakeRipper
	//TODO initialize more video rippers if required
}

func RipVideo(ctx task.Context) task.HandlerFunc {
	conf := ctx.Config.(*ripper.AppConf)
	ripperType := conf.Rip.Video.Ripper

	var rip processor.Processor
	var err error

	rf := RipperFactories[ripperType]
	if rf == nil {
		err = fmt.Errorf("unknown video ripper configured: \"%s\"", ripperType)
	} else {
		rip, err = rf(conf.Rip.Video.Handbrake, ctx.Printf, conf.WorkDirectory)
	}

	if err != nil {
		return ripper.ErrorHandler(err)
	} else {
		return processor.Process(ctx, rip, ripperType,
			processor.DefaultCheckLazy(ctx.RunLazy, conf.Output.Video),
			processor.DefaultInputFileFor(conf.Rip.Video.AllowedInputExtensions),
			processor.DefaultOutputFileFor(conf.Output.Video))
	}
}
