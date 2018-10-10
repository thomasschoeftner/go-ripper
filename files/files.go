package files

import (
	"os"
	"strings"
	"path/filepath"
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

func Extension(filePath string) string {
	return strings.Replace(filepath.Ext(filePath), ".", "", 1)
}

func SplitExtension(file string) (string, string) {
	idx := strings.LastIndex(file, ".")
	if idx > 0 {
		return file[:idx], file[idx+1:len(file)]
	} else {
		return file, ""
	}
}