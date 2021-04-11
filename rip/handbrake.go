package rip

import (
	"github.com/thomasschoeftner/go-cli/commons"
	"github.com/thomasschoeftner/go-ripper/ripper"
	"io"
	"os"
	"github.com/thomasschoeftner/go-cli/cli"
	"time"
	"path/filepath"
	"github.com/thomasschoeftner/go-ripper/files"
	"github.com/thomasschoeftner/go-ripper/targetinfo"
	"github.com/thomasschoeftner/go-ripper/processor"
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

func createHandbrakeRipper(conf *ripper.AppConf, printf commons.FormatPrinter, workDir string) (processor.Processor, error) {
	hbConf := conf.Rip.Video.Handbrake
	timeout, err := time.ParseDuration(hbConf.Timeout)
	if err != nil {
		return nil, err
	}

	var errOut io.Writer
	if hbConf.ShowErrorOutput {
		errOut = os.Stderr
	}
	var stdOut io.Writer
	if hbConf.ShowStandardOutput {
		stdOut = os.Stdout
	}

	return func (ti targetinfo.TargetInfo, inFile string, outFile string) error {
		evacuated, err := files.PrepareEvacuation(filepath.Join(workDir, files.TEMP_DIR_NAME)).Of(inFile).By(files.Moving)
		if err != nil {
			return err
		}
		defer evacuated.Restore()

		tmpOut := evacuated.WithSuffix(".ripped")
		cmd := cli.Command(hbConf.Path, timeout).WithQuotes(" ", '\'').
		WithParam(paramImportPreset, filepath.ToSlash(hbConf.PresetsFile), "").
		WithParam(paramUsePreset, hbConf.PresetName, "").
		WithParam(paramInput, filepath.ToSlash(evacuated.Path()), "").
		WithParam(paramOutput, filepath.ToSlash(tmpOut), "")
		//printf(">>>> %s\n", cmd.String())
		err = cmd.ExecuteSync(stdOut, errOut)
		if err != nil {
			return err
		}
		return os.Rename(tmpOut, outFile)
	}, nil
}
