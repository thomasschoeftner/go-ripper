package metainfo

import (
	"errors"
	"io/ioutil"
	"encoding/json"
	"path/filepath"
	"github.com/thomasschoeftner/go-ripper/files"
	"os"
)

const (
	METAINF_FILE_EXT = "json"
)

type MetaInfo interface {
	GetId() string
	GetType() string
}

type IdInfo struct {
	Id string
}
func (ii *IdInfo) GetId() string {
	return ii.Id
}

func Is(metaInfo MetaInfo, kind string) bool {
	if metaInfo != nil {
		return kind == metaInfo.GetType()
	}
	return false
}


func ReadMetaInfo(filePath string, content interface{}) error {
	if content == nil {
		return errors.New("cannot unmarshall to nil")
	}

	jsonRaw, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonRaw, content)
	return err
}


func SaveMetaInfo(filePath string, metaInfo interface{}) error {
	if metaInfo == nil {
		return errors.New("cannot save nil json")
	}

	folder := filepath.Dir(filePath)
	err := files.CreateFolderStructure(folder)
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(metaInfo, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, bytes, os.ModePerm)
}
