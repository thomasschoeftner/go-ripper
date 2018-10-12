package files

import (
	"os"
	"strings"
	"fmt"
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

func SplitExtension(file string) (string, string) {
	suffix := filepath.Ext(file)
	prefix := strings.TrimSuffix(file, suffix) //cut extension (including leading '.') from path
	if 0 == len(prefix) { //use suffix as filename in case there is no extension and file starts with '.'
		prefix = suffix
		suffix = ""
	}

	if len(suffix) > 0 { //removing leading '.' from extension
		suffix = suffix[1:] //remove leading '.'
	}
	return prefix, suffix
}

func GetExtension(filePath string) string {
	_, ext := SplitExtension(filePath)
	return ext
}

func WithExtension(name string, extension string) string {
	return fmt.Sprintf("%s.%s", name, strings.TrimLeft(extension, "."))
}