package files

import (
	"os"
)

//func CheckOrCreateFolder(path string) error {
//	exists, isFolder, err := Exists(path)
//	if err != nil {
//		return err
//	}
//
//	if !exists {
//		err = CreateFolderStructure(path)
//	} else if !isFolder {
//		err = errors.New(fmt.Sprintf("unable to create folder structure \"%s\" because file with same name exists", path))
//	}
//	return err
//}

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
