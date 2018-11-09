package rip

import (
	"go-cli/commons"
	"go-ripper/ripper"
	"io"
	"os"
	"go-cli/cli"
	"time"
)

const CONF_RIPPER_HANDBRAKE = "handbrake"
const (
	//TODO define command line params and args
)

func createHandbrakeRipper(conf *ripper.HandbrakeConfig, lazy bool, printf commons.FormatPrinter) (Ripper, error) {
	timeout, err := time.ParseDuration(conf.Timeout)
	if err != nil {
		return nil, err
	}

	hb := handbrakeRipper{
		lazy: lazy,
		path: conf.Path,
		profilePath: conf.Profile,
		timeout : timeout,
		printf: printf}
	if conf.ShowErrorOutput {
		hb.errOut = os.Stderr
	}
	if conf.ShowStandardOutput {
		hb.stdOut = os.Stdout
	}
	return &hb, nil
}


type handbrakeRipper struct {
	lazy bool
	path string
	profilePath string
	timeout time.Duration
	printf commons.FormatPrinter
	errOut io.Writer
	stdOut io.Writer

}

func (hb *handbrakeRipper) process(inFile string, outFile string) error {
	cmd := cli.Command(hb.path, hb.timeout).WithQuotes(" ", '"')
	finish this!!!
	//TODO implement me
	hb.printf(">>>> %s\n", cmd.String())
	return cmd.ExecuteSync(hb.stdOut, hb.errOut)
}