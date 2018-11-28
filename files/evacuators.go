package files

import (
	"errors"
	"fmt"
	"path/filepath"
	"os"
	"go-cli/commons"
	"strconv"
)

const TEMP_DIR_NAME = ".tmp"

type Evacuated struct {
	original    string
	evacuatedTo string
}
func (e *Evacuated) Path() string {
	return e.evacuatedTo
}
func (e *Evacuated) Restore() error {
	return e.MoveTo(e.original)
}
func (e *Evacuated) Discard() error {
	return os.Remove(e.evacuatedTo)
}
func (e *Evacuated) MoveTo(file string) error {
	defer e.Discard()
	return os.Rename(e.evacuatedTo, file)
}

func newTempFileName(tempDir, originalFile string) string {
	hash:= commons.Hash32(originalFile)
	_, ext := SplitExtension(originalFile)
	return filepath.Join(tempDir, WithExtension(strconv.Itoa(int(hash)), ext))
}

func PrepareEvacuation(tempDir string) preparedEvacuation {
	// validation
	if 0 == len(tempDir) {
		return evacuationFailure(errors.New("cannot evacuate - temp-dir is undefined"))
	}
	if CreateFolderStructure(tempDir) != nil {
		return evacuationFailure(fmt.Errorf("unable to create temp folder at \"%s\"", tempDir))
	}

	// return actual preparation function
	return func(originalPath string) evacuationTarget {
		if exits, _ := Exists(originalPath); !exits {
			return evacuationFailure(fmt.Errorf("cannot evacuate non-existing file \"%s\"", originalPath))(originalPath)
		}
		return &readyEvacuationTarget{original: originalPath, evacuated: newTempFileName(tempDir, originalPath)}
	}
}

type preparedEvacuation func(string) evacuationTarget
func (pe preparedEvacuation) Of(file string) evacuationTarget {
	return pe(file)
}

func evacuationFailure(err error) preparedEvacuation {
	return func(string) evacuationTarget {
		return &failedEvacuationTarget{err}
	}
}

type evacuationTarget interface {
	By(evacuator EvacuatorFunc) (*Evacuated, error)
}
type failedEvacuationTarget struct {
	err error
}
func (failed *failedEvacuationTarget) By(EvacuatorFunc) (*Evacuated, error) {
	return nil, failed.err
}
type readyEvacuationTarget struct {
	original string
	evacuated string
}
func (hopeful *readyEvacuationTarget) By(evacuate EvacuatorFunc) (*Evacuated, error) {
	return evacuate(hopeful.original, hopeful.evacuated)
}

type EvacuatorFunc func(from, to string) (*Evacuated, error)
func Moving(from, to string) (*Evacuated, error) {
	if err := os.Rename(from, to); err != nil {
		return nil, err
	}
	return &Evacuated{original: from, evacuatedTo: to}, nil
}
func Copying(from, to string) (*Evacuated, error) {
	if _, err := Copy(from, to, false); err != nil {
		return nil, err
	}

	return &Evacuated{original: from, evacuatedTo: to}, nil
}
