package files

import (
	"os"
)

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, err
}

func CreateFolder(folder string) error {
	return os.Mkdir(folder, os.ModePerm)
}

func CreateFolderStructure(folder string) error {
	return os.MkdirAll(folder, os.ModePerm)
}
