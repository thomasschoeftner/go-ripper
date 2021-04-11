package metainfo

import (
	"path/filepath"
	"fmt"
	"errors"
	"io/ioutil"
	"github.com/thomasschoeftner/go-ripper/files"
	"os"
)

const (
	SUBDIR_IMAGES = "imgs"
)

type Image []byte

func ImageFileName(repoPath string, id string, extension string) string {
	return filepath.ToSlash(filepath.Join(repoPath, SUBDIR_IMAGES, fmt.Sprintf("%s.%s", id, extension)))
}

func ReadImage(filePath string) (Image, error) {
	if len(filePath) == 0 {
		return nil, errors.New("cannot unmarshall from empty file name")
	}
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func SaveImage(filePath string, img Image) error {
	if img == nil || len(img) == 0 {
		return errors.New("cannot save empty image")
	}

	folder := filepath.Dir(filePath)
	err := files.CreateFolderStructure(folder)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, img, os.ModePerm)
}
