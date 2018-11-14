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
	paramInput              = "--input"              //input file as param
	paramOutput             = "--output"             //mp4 output file as param
	paramImportPreset       = "--preset-import-file" //json file as param
	paramUsePreset          = "--preset"             //selected preset name
	argOptizizeForStreaming = "--optimize"
	argLogToJson            = "--json"
)

func handbrakeRipper(conf *ripper.HandbrakeConfig, lazy bool, printf commons.FormatPrinter) (Ripper, error) {
	timeout, err := time.ParseDuration(conf.Timeout)
	if err != nil {
		return nil, err
	}

	var errOut io.Writer
	if conf.ShowErrorOutput {
		errOut = os.Stderr
	}
	var stdOut io.Writer
	if conf.ShowStandardOutput {
		stdOut = os.Stdout
	}

	return func (inFile string, outFile string) error {
		cmd := cli.Command(conf.Path, timeout).WithQuotes(" ", '"').
		WithParam(paramImportPreset, conf.PresetsFile, "").
		WithParam(paramUsePreset, conf.PresetName, "").
		WithParam(paramInput, inFile, "").
		WithParam(paramOutput, outFile, "")
		//hb.printf(">>>> %s\n", cmd.String())
		return cmd.ExecuteSync(stdOut, errOut)
	}, nil
}
