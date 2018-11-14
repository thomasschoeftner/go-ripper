package rip

import (
	"go-cli/task"
	"go-ripper/ripper"
	"fmt"
	"go-ripper/targetinfo"
	"go-ripper/files"
	"go-cli/commons"
)


type Ripper func(inFile string, outFile string) error

func Rip(ctx task.Context, rip Ripper, allowedInputExtensions []string, expectedOutputExtension string) task.HandlerFunc {
	conf := ctx.Config.(*ripper.AppConf)
	ripperType := conf.Rip.Video.Ripper

	if nil == rip {
		return ripper.ErrorHandler(fmt.Errorf("failed to initialize ripper module using %s because actual ripper is undefined", ripperType))
	}

	return func(job task.Job) ([]task.Job, error) {
		target := ripper.GetTargetFileFromJob(job)
		ctx.Printf("use %s to rip file - target %s\n", ripperType, target)

		in, out, err := findInputOutputFiles(target, conf.WorkDirectory, allowedInputExtensions, expectedOutputExtension)
		if err != nil {
			return nil, err
		}

		err = rip(in, out)
		if err != nil {
			return []task.Job{}, err
		} else {
			return []task.Job{job}, nil
		}
	}
}

func findInputOutputFiles(targetPath string, workDir string, allowedInputExtensions []string, expectedOutputExtension string) (string, string, error) {
	ti, err := targetinfo.ForTarget(workDir, targetPath)
	if err != nil {
		return "", "", err
	}

	preprocessed, err := ripper.GetProcessingArtifactPathFor(workDir, ti.GetFolder(), ti.GetFile(), expectedOutputExtension)
	if err != nil {
		return "", "", err
	}

	for _, ext := range allowedInputExtensions {
		fName, err := ripper.GetProcessingArtifactPathFor(workDir, ti.GetFolder(), ti.GetFile(), ext)
		if err != nil {
			return "", "", err
		}
		if exists, err := files.Exists(fName); err != nil {
			return "", "", err
		} else if exists {
			return fName, preprocessed, nil
		}
	}

	_, inputExt := files.SplitExtension(targetPath)
	if commons.IsStringAmong(inputExt, allowedInputExtensions) {
		return targetPath, preprocessed, nil
	}
	return "", "", fmt.Errorf("unable to find valid input file related to \"%s\" - valid file types are %v", targetPath, allowedInputExtensions)
}
