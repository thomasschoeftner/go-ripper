package rip

import (
	"go-cli/commons"
	"go-ripper/ripper"
	"io"
	"os"
)

const CONF_RIPPER_HANDBRAKE = "handbrake"
const (
	//command line params
)

func createHandbrakeRipper(conf *ripper.HandbrakeConfig, lazy bool, printf commons.FormatPrinter) Ripper {
	hb := handbrakeRipper{
		lazy: lazy,
		path: conf.Path,
		profilePath: conf.Profile,
		printf: printf}
	if conf.ShowErrorOutput {
		hb.errOut = os.Stderr
	}
	if conf.ShowStandardOutput {
		hb.stdOut = os.Stdout
	}
	return &hb
}


type handbrakeRipper struct {
	lazy bool
	path string
	profilePath string
	printf commons.FormatPrinter
	errOut io.Writer
	stdOut io.Writer

}

func (hb *handbrakeRipper) process(inFile string, outFile string) error {
	//TODO implement me
	return nil
}