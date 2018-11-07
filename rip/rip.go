package rip

import (
	"go-ripper/targetinfo"
	"go-ripper/ripper"
	"go-ripper/files"
	"path/filepath"
	"errors"
)

func findInputOutputFiles(ti targetinfo.TargetInfo, workDir string, expectedExtension string) (string, string, error) {
	folder, err := ripper.GetWorkPathForTargetFolder(workDir, ti.GetFolder())
	if err != nil {
		return "", "", err
	}

	// check work directory for a pre-processed inFile in appropriate format (e.g. a ripped video in .mp4 inFile)
	fName, extension := files.SplitExtension(ti.GetFile())
	preprocessed := filepath.Join(folder, files.WithExtension(fName, expectedExtension))
	if exists, _ := files.Exists(preprocessed); exists {
		return preprocessed, preprocessed, nil
	}

	// if no preprocessed input is available, check if the source inFile can be tagged directly (e.g. if it is an .mp4 video)
	if extension == expectedExtension {
		return filepath.Join(ti.GetFolder(), ti.GetFile()), preprocessed, nil
	} else {
		return "", "", errors.New("unable to find appropriate input inFile for meta-info tagging")
	}
}