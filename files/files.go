package files

import (
	"os"
	"strings"
	"fmt"
	"path/filepath"
	"io"
)

func Copy(from, to string, append bool) (int64, error) {
	if srcFileStat, err := os.Stat(from); err != nil {
		return 0, err
	} else if !srcFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("cannot copy file \"%s\" - not a regular file", from)
	}
	src, err := os.Open(from)
	if err != nil {
		return 0, err
	}
	defer src.Close()

	fileOptions := os.O_RDWR|os.O_CREATE|os.O_TRUNC
	if append {
		fileOptions = fileOptions|os.O_APPEND
	} else {
		fileOptions = fileOptions|os.O_EXCL
	}
	dst, err := os.OpenFile(to, fileOptions, 0666)
	if err != nil {
		return 0, err
	}
	defer dst.Close()

	return io.Copy(dst, src)
}

const defaultKeepExt = "_keep_"
func Replace(file, with string) (wasReplaced bool, err error) {
	defer os.Remove(WithExtension(file, defaultKeepExt))
	return ReplaceButKeepOriginal(file, with, defaultKeepExt)
}

func ReplaceButKeepOriginal(file, with, keepExtension string) (wasReplaced bool, e error) {
	if 0 == len(keepExtension) {
		return false, fmt.Errorf("cannot replace file, but keep original with the same name (keep extension is empty)")
	}

	println("replace", file, "->", with)

	originalExists, err := Exists(file)
	if err != nil {
		return false, err
	}
	if replacementExists, err := Exists(with); err != nil {
		return false, err
	} else if !replacementExists {
		return false, fmt.Errorf("cannot replace file \"%s\" with MISSING file \"%s\"", file, with)
	}

	if originalExists {
		err = os.Rename(file, WithExtension(file, keepExtension))
		if err != nil {
			return false, err
		}
	}

	_, err = Copy(with, file, false)
	return originalExists, err
}

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
