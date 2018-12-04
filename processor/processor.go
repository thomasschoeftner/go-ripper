package processor

import (
	"go-cli/task"
	"go-ripper/ripper"
	"fmt"
	"go-ripper/targetinfo"
	"go-ripper/files"
	"strings"
	"go-cli/commons"
)


type Processor func(ti targetinfo.TargetInfo, inFile string, outFile string) error
type CheckLazy func(targetInfo targetinfo.TargetInfo) bool
type InputFile func(targetInfo targetinfo.TargetInfo, workDir string) (string, error)
type OutputFile func(targetInfo targetinfo.TargetInfo, workDir string) (string, error)

func Process(ctx task.Context, process Processor, processorName string, checkLazy CheckLazy, inputFile InputFile, outputFile OutputFile) task.HandlerFunc {
	conf := ctx.Config.(*ripper.AppConf)
	workDir := conf.WorkDirectory
	if nil == process {
		return ripper.ErrorHandler(fmt.Errorf("failed to initialize %s processor because actual processor instance is undefined", processorName))
	}

	return func(job task.Job) ([]task.Job, error) {
		target := ripper.GetTargetFileFromJob(job)
		ti, err := targetinfo.ForTarget(workDir, target)
		if err != nil {
			return nil, err
		}

		in, err := inputFile(ti, workDir)
		if err != nil {
			return nil, err
		}
		out, err := outputFile(ti, workDir)
		if err != nil {
			return nil, err
		}

		if checkLazy(ti) {
			ctx.Printf("input file appears just right -> reuse %s\n", target)
			_, err = files.Copy(in, out, false)
		} else {
			ctx.Printf("use %s to process file %s\n", processorName, target)
			err = process(ti, in, out)
		}

		if err != nil {
			return []task.Job{}, err
		} else {
			return []task.Job{job}, nil
		}
	}
}


func DefaultCheckLazy(lazyEnabled bool, expectedExtension string) CheckLazy {
	return func(targetInfo targetinfo.TargetInfo) bool {
		return lazyEnabled && strings.ToLower(files.GetExtension(targetInfo.GetFile())) == strings.ToLower(expectedExtension)
	}
}

func DefaultInputFileFor(workDir string, allowedInputExtensions []string) InputFile {
	return func(ti targetinfo.TargetInfo, workDir string) (string, error) {
		if ti == nil {
			return "", fmt.Errorf("target-info is undefined")
		}
		originalInputFile := ti.GetFullPath()

		// 1. check among pre-processed artifacts
		for _, ext := range allowedInputExtensions {
			fName, err := ripper.GetProcessingArtifactPathFor(workDir, ti.GetFolder(), ti.GetFile(), ext)
			if err != nil {
				return "", err
			}
			if exists, err := files.Exists(fName); err != nil {
				return "", err
			} else if exists {
				return fName, nil
			}
		}

		//2. check if original input file can be used
		_, inputExt := files.SplitExtension(originalInputFile)
		if commons.IsStringAmong(inputExt, allowedInputExtensions) {
			return ti.GetFullPath(), nil
		}

		return "", fmt.Errorf("unable to find valid input file related to \"%s\" - valid file types are %v", originalInputFile, allowedInputExtensions)
	}
}

func DefaultOutputFileFor(workDir string, expectedOutputExtension string) OutputFile {
	return func(ti targetinfo.TargetInfo, workDir string) (string, error) {
		return ripper.GetProcessingArtifactPathFor(workDir, ti.GetFolder(), ti.GetFile(), expectedOutputExtension)
	}
}