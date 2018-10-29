package files

import (
	"errors"
	"io/ioutil"
	"fmt"
	"strings"
	"path/filepath"
	"os"
)

type evacuated struct {
	original    string
	evacuatedTo string
}
func (e *evacuated) Path() string {
	return e.evacuatedTo
}
func (e *evacuated) Restore() error {
	return e.MoveTo(e.original)
}
func (e *evacuated) Discard() error {
	return os.RemoveAll(filepath.Dir(e.evacuatedTo))
}
func (e *evacuated) MoveTo(file string) error {
	defer e.Discard()
	return os.Rename(e.evacuatedTo, file)
}


func PrepareEvacuation(tempDir string, toReplace map[rune]rune) preparedEvacuation {
	// validation
	if 0 == len(tempDir) {
		return evacuationFailure(errors.New("cannot evacuate - temp-dir is undefined"))
	}
	if CreateFolderStructure(tempDir) != nil {
		return evacuationFailure(fmt.Errorf("unable to create temp folder at \"%s\"", tempDir))
	}

	// create temporary directory
	dir, err := ioutil.TempDir(tempDir, ".")
	if err != nil {
		return evacuationFailure(err)
	}

	// return actual preparation function
	return func(originalPath string) evacuationTarget {
		if exits, _ := Exists(originalPath); !exits {
			return evacuationFailure(fmt.Errorf("cannot evacuate non-existing file \"%s\"", originalPath))(originalPath)
		}

		// replace specific characters in filename
		_, file := filepath.Split(originalPath)
		file = strings.Map(func(originalChar rune) rune {
			replacementChar, mustReplace := toReplace[originalChar]
			if mustReplace {
				return replacementChar
			} else {
				return originalChar
			}
		}, file)
		return &readyEvacuationTarget{original: originalPath, evacuated: filepath.Join(dir, file)}
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
	By(evacuator evacuator) (*evacuated, error)
}
type failedEvacuationTarget struct {
	err error
}
func (failed *failedEvacuationTarget) By(evacuator) (*evacuated, error) {
	return nil, failed.err
}
type readyEvacuationTarget struct {
	original string
	evacuated string
}
func (hopeful *readyEvacuationTarget) By(evacuate evacuator) (*evacuated, error) {
	return evacuate(hopeful.original, hopeful.evacuated)
}

type evacuator func(from, to string) (*evacuated, error)
func Moving(from, to string) (*evacuated, error) {
	if err := os.Rename(from, to); err != nil {
		return nil, err
	}
	return &evacuated {original: from, evacuatedTo: to}, nil
}
func Copying(from, to string) (*evacuated, error) {
	if _, err := Copy(from, to, false); err != nil {
		return nil, err
	}

	return &evacuated {original: from, evacuatedTo: to}, nil
}


//
//type EvacuatorFactory func(tempDir string, toReplace map[rune]rune) evacuator
//type evacuator func(filePath string) (evacuated, error)
//
//func prepareEvacuation(tempDir string, toReplace map[rune]rune, path string) (string, error) {
//	// validation
//	if 0 == len(tempDir) {
//		return "", errors.New("cannot evacuate path - temp-dir is undefined")
//	}
//	if CreateFolderStructure(tempDir) != nil {
//		return "", fmt.Errorf("unable to create temp folder at \"%s\"", tempDir)
//	}
//	if exits, _ := Exists(path); !exits {
//		return "", errors.New("cannot evacuate non-existing path")
//	}
//
//	// create temporary directory
//	dir, err := ioutil.TempDir(tempDir, ".")
//	if err != nil {
//		return "", err
//	}
//
//	// calc new filename
//	_, file := filepath.Split(path)
//	file = strings.Map(func(in rune) rune {
//		replacement, mustReplace := toReplace[in]
//		if mustReplace {
//			return replacement
//		} else {
//			return in
//		}
//	}, file)
//
//	return filepath.Join(dir, file), nil
//}
//
//func MovingEvacuator(tempDir string, toReplace map[rune]rune) evacuator {
//	return func(path string) (evacuated, error) {
//		moveTo, err := prepareEvacuation(tempDir, toReplace, path)
//		if err != nil {
//			return nil, err
//		}
//		m := &moved {original: path, movedTo: moveTo}
//		return m, os.Rename(m.original, m.movedTo)
//	}
//}
//
//type moved struct {
//	original string
//	movedTo  string
//}
//func (m *moved) Restore() error {
//	return m.MoveTo(m.original)
//}
//func (m *moved) Discard() error {
//	return os.RemoveAll(filepath.Dir(m.movedTo))
//}
//func (m *moved) MoveTo(file string) error {
//	defer m.Discard()
//	return os.Rename(m.movedTo, file)
//}
//
//
//
//func CopyingEvacuator(tempDir string, toReplace map[rune]rune) evacuator {
//	return func(path string) (evacuated, error) {
//		copiedTo, err := prepareEvacuation(tempDir, toReplace, path)
//		if err != nil {
//			return nil, err
//		}
//
//		// copy file
//		_, err = Copy(path, copiedTo, false)
//		if err != nil {
//			return nil, err
//		}
//		return &copied{original: path, copy: copiedTo}, nil
//	}
//}
//
//type copied struct {
//	original string
//	copy string
//}
//
//func (c *copied) Restore() error {
//	return c.MoveTo(c.original)
//}
//
//func (c *copied) Discard() error {
//	return os.RemoveAll(filepath.Dir(c.copy))
//}
//
//func (c *copied) MoveTo(file string) error {
//	defer c.Discard()
//	return os.Rename(c.copy, file)
//}
