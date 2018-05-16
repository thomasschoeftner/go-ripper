package files

import (
	"os"
)

func CheckOrCreateFolder(path string) error {
	exists, isFolder, err := Exists(path)
	if err != nil {
		return err
	}

	if !exists {
		err = createFolder(path)
	} else if !isFolder {
		err = os.RemoveAll(path)
		if err != nil {
			return err
		}
		err = createFolder(path)
	}
	return err
}

func Exists(path string) (exists bool, isFolder bool, e error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, false, nil
		}
	}
	return true, info.IsDir(), err
}

func createFolder(folder string) error {
	return os.Mkdir(folder, os.ModePerm)
}