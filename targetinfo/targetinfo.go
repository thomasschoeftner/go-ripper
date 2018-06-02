package targetinfo

import (
	"io/ioutil"
	"encoding/json"
	"os"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

const (
	TARGETINFO_VIDEO   = "v"
	TARGETINFO_EPISODE = "e"
)

type TargetInfo interface {
	fmt.Stringer
	GetFile() string
	GetFolder() string
	GetType() string
	GetId() string
	fileName() string
}

type Video struct {
	File   string `json:"file"`
	Folder string `json:"folder"`
	Id     string `json:"id"`
}

type Episode struct {
	Video
	Season     int `json:"season"`
	Episode    int `json:"episode"`
	ItemSeqNo  int `json:"itemseqno"`
	ItemsTotal int `json:"itemstotal"`
}

func (v *Video) GetFile() string {
	return v.File
}

func (v *Video) GetFolder() string {
	return v.Folder
}

func (v *Video) GetType() string {
	return TARGETINFO_VIDEO
}

func (v *Video) GetId() string {
	return v.Id
}

func (v *Video) fileName() string {
	return fmt.Sprintf("%s%s", v.GetType(), v.Id)
}

func (v *Video) String() string {
	return fmt.Sprintf("video   (id=%s, file=%s)", v.Id, filepath.Join(v.Folder, v.File))
}

func (e *Episode) GetType() string {
	return TARGETINFO_EPISODE
}

func (e *Episode) fileName() string {
	return fmt.Sprintf("%s%s.%d.%d", e.GetType(), e.Id, e.Season, e.Episode)
}

func (e *Episode) String() string {
	return fmt.Sprintf("episode (id=%s, season=%-4d, episode=%-4d, item#=%-4d, totalItems=%-4d, file=%s)", e.Id, e.Season, e.Episode, e.ItemSeqNo, e.ItemsTotal, filepath.Join(e.Folder, e.File))
}


func NewVideo(file string, folder string, id string) *Video {
	return &Video{File: file, Folder: folder, Id: id}
}

func NewEpisode(file string, folder string, id string, season int, episode int, itemSeqNo int, itemsTotal int) *Episode {
	return &Episode{Video: *NewVideo(file, folder, id), Season: season, Episode: episode, ItemSeqNo: itemSeqNo, ItemsTotal: itemsTotal}
}

func Read(targetInfoFile string) (TargetInfo, error) {
	raw, err := ioutil.ReadFile(targetInfoFile)
	if err != nil {
		return nil, err
	}
	jsonString := string(raw[:])

	_, f := filepath.Split(targetInfoFile)
	var ti TargetInfo
	if strings.HasPrefix(f, TARGETINFO_VIDEO) {
		ti = &Video{}
	} else if strings.HasPrefix(f, TARGETINFO_EPISODE) {
		ti = &Episode{}
	} else {
		return nil, errors.New(fmt.Sprintf("target info file suggests invalid target info type: %s ", f))
	}

	err = json.Unmarshal([]byte(jsonString), &ti)
	if err != nil {
		return nil, err
	}

	return ti, nil
}

func Save(tmpFolder string, ti TargetInfo) (*string, error) {
	if ti == nil {
		return nil, errors.New("target info is nil")
	}

	bytes, err := json.Marshal(ti)
	if err != nil {
		return nil, err
	}
	targetFile := filepath.Join(tmpFolder, ti.fileName())
	err = ioutil.WriteFile(targetFile, bytes, os.ModePerm)
	return &targetFile, err
}
