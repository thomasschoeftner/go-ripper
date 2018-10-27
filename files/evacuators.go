package files

import (
	"errors"
	"io/ioutil"
	"fmt"
	"strings"
	"path/filepath"
	"os"
)

type EvacuatorFactory func(tempDir string, toReplace map[rune]rune) evacuator
type evacuator func(filePath string) (evacuated, error)
type evacuated interface {
	Restore() error
	Discard() error
	MoveTo(file string) error
}

func prepareEvacuation(tempDir string, toReplace map[rune]rune, path string) (string, error) {
	// validation
	if 0 == len(tempDir) {
		return "", errors.New("cannot evacuate path - temp-dir is undefined")
	}
	if CreateFolderStructure(tempDir) != nil {
		return "", fmt.Errorf("unable to create temp folder at \"%s\"", tempDir)
	}
	if exits, _ := Exists(path); !exits {
		return "", errors.New("cannot evacuate non-existing path")
	}

	// create temporary directory
	dir, err := ioutil.TempDir(tempDir, ".")
	if err != nil {
		return "", err
	}

	// calc new filename
	_, file := filepath.Split(path)
	file = strings.Map(func(in rune) rune {
		replacement, mustReplace := toReplace[in]
		if mustReplace {
			return replacement
		} else {
			return in
		}
	}, file)

	return filepath.Join(dir, file), nil
}


func MovingEvacuator(tempDir string, toReplace map[rune]rune) evacuator {
	return func(path string) (evacuated, error) {
		moveTo, err := prepareEvacuation(tempDir, toReplace, path)
		if err != nil {
			return nil, err
		}
		m := &moved {original: path, movedTo: moveTo}
		return m, os.Rename(m.original, m.movedTo)
	}
}

type moved struct {
	original string
	movedTo  string
}
func (m *moved) Restore() error {
	return m.MoveTo(m.original)
}
func (m *moved) Discard() error {
	return os.RemoveAll(filepath.Dir(m.movedTo))
}
func (m *moved) MoveTo(file string) error {
	err := os.Rename(m.movedTo, file)
	if err != nil {
		return err
	}
	return m.Discard()
}



func CopyingEvacuator(tempDir string, toReplace map[rune]rune) evacuator {
	return func(path string) (evacuated, error) {
		copyTo, err := prepareEvacuation(tempDir, toReplace, path)
		if err != nil {
			return nil, err
		}

		// copy file
		_, err = Copy(path, copyTo, false)
		if err != nil {
			return nil, err
		}
		return &copied{copy: copyTo}, nil
	}
}

type copied struct {
	copy string
}

func (c *copied) Restore() error {
	return c.Discard()
}

func (c *copied) Discard() error {
	return os.RemoveAll(filepath.Dir(c.copy))
}

func (c *copied) MoveTo(file string) error {
	err := os.Rename(c.copy, file)
	if err != nil {
		return err
	}
	return c.Discard()
}
