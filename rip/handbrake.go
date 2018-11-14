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

func createHandbrakeRipper(conf *ripper.HandbrakeConfig, lazy bool, printf commons.FormatPrinter) (Ripper, error) {
	timeout, err := time.ParseDuration(conf.Timeout)
	if err != nil {
		return nil, err
	}

	hb := handbrakeRipper{
		lazy:            lazy,
		path:            conf.Path,
		presetsFilePath: conf.PresetsFile,
		presetName:      conf.PresetName,
		timeout :        timeout,
		printf:          printf}
	if conf.ShowErrorOutput {
		hb.errOut = os.Stderr
	}
	if conf.ShowStandardOutput {
		hb.stdOut = os.Stdout
	}
	return &hb, nil
}

type handbrakeRipper struct {
	lazy            bool
	path            string
	presetsFilePath string
	presetName      string
	timeout         time.Duration
	printf          commons.FormatPrinter
	errOut          io.Writer
	stdOut          io.Writer
}

func (hb *handbrakeRipper) process(inFile string, outFile string) error {
	cmd := cli.Command(hb.path, hb.timeout).WithQuotes(" ", '"').
		WithParam(paramImportPreset, hb.presetsFilePath, "").
		WithParam(paramUsePreset, hb.presetName, "").
		WithParam(paramInput, inFile, "").
		WithParam(paramOutput, outFile, "")
	hb.printf(">>>> %s\n", cmd.String())
	return cmd.ExecuteSync(hb.stdOut, hb.errOut)
}
