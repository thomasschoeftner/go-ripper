package targetinfo

import (
	"io/ioutil"
	"encoding/json"
	"os"
	"errors"
	"fmt"
	"path/filepath"
)

const targetinfo_filetype = "targetinfo"

const (
	TARGETINFO_TYPE_VIDEO   = "video"
	TARGETINFO_TPYE_EPISODE = "episode"
)

type TargetInfo interface {
	fmt.Stringer
	GetFile() string
	GetFolder() string
	GetType() string
	GetId() string
}

func fileName(ti TargetInfo) string {
	return appendFileExtension(ti.GetFile())
}

func appendFileExtension(targetFileName string) string {
	return fmt.Sprintf("%s.%s", targetFileName, targetinfo_filetype)
}

type Typed struct {
	Type string `json:"type"`
}


type Video struct {
	Typed
	File   string `json:"file"`
	Folder string `json:"folder"`
	Id     string `json:"id"`
}

type Episode struct {
	//Typed
	Video
	Season     int    `json:"season"`
	Episode    int    `json:"episode"`
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
	return TARGETINFO_TYPE_VIDEO
}

func (v *Video) GetId() string {
	return v.Id
}

func (v *Video) String() string {
	return fmt.Sprintf("video   (id=%s, file=%s)", v.Id, filepath.Join(v.Folder, v.File))
}

func (e *Episode) GetType() string {
	return TARGETINFO_TPYE_EPISODE
}

func (e *Episode) String() string {
	return fmt.Sprintf("episode (id=%s, season=%-4d, episode=%-4d, item#=%-4d, totalItems=%-4d, file=%s)", e.Id, e.Season, e.Episode, e.ItemSeqNo, e.ItemsTotal, filepath.Join(e.Folder, e.File))
}


func NewVideo(file string, folder string, id string) *Video {
	return &Video{Typed: Typed{ Type: TARGETINFO_TYPE_VIDEO}, File: file, Folder: folder, Id: id}
}

func NewEpisode(file string, folder string, id string, season int, episode int, itemSeqNo int, itemsTotal int) *Episode {
	vid := NewVideo(file, folder, id)
	vid.Type = TARGETINFO_TPYE_EPISODE
	return &Episode{/*Typed: Typed { Type: TARGETINFO_TPYE_EPISODE},*/ Video: *vid, Season: season, Episode: episode, ItemSeqNo: itemSeqNo, ItemsTotal: itemsTotal}
}

func Read(workFolder string, targetFileNeme string) (TargetInfo, error) {
	targetInfoFile := filepath.Join(workFolder, appendFileExtension(targetFileNeme))
	jsonRaw, err := ioutil.ReadFile(targetInfoFile)
	if err != nil {
		return nil, err
	}
	//jsonString := raw[:]

	//1. get type
	typed := Typed{}
	err = json.Unmarshal(jsonRaw, &typed)

	var ti TargetInfo
	switch typed.Type {
	case TARGETINFO_TYPE_VIDEO:
		ti = &Video{}
	case TARGETINFO_TPYE_EPISODE:
		ti = &Episode{}
	default:
		return nil, errors.New(fmt.Sprintf("target info file contains invalid target info type: \"%s\"", typed.Type))
	}

	err = json.Unmarshal(jsonRaw, &ti)
	if err != nil {
		return nil, err
	}
	return ti, nil
}

func Save(workFolder string, ti TargetInfo) error {
	if ti == nil {
		return errors.New("target info is nil")
	}

	bytes, err := json.MarshalIndent(ti, "", "  ")
	if err != nil {
		return err
	}
	targetFile := filepath.Join(workFolder, fileName(ti))
	err = ioutil.WriteFile(targetFile, bytes, os.ModePerm)
	return nil
}
