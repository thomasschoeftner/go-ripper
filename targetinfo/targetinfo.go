package targetinfo

import (
	"io/ioutil"
	"encoding/json"
	"os"
	"errors"
	"fmt"
	"path/filepath"
	"go-ripper/ripper"
	"go-ripper/files"
)

const targetinfo_file_extension = "targetinfo"

const (
	TARGETINFO_TYPE_MOVIE   = "movie"
	TARGETINFO_TPYE_EPISODE = "episode"
)

type TargetInfo interface {
	fmt.Stringer
	GetFile() string
	GetFolder() string
	GetType() string
	GetId() string
	GetFullPath() string
}

func fileName(ti TargetInfo) string {
	return files.WithExtension(ti.GetFile(), targetinfo_file_extension)
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

type Movie struct {
	Video
}

type Episode struct {
	//Typed
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

func (v *Video) GetId() string {
	return v.Id
}

func (v *Video) GetFullPath() string {
	return filepath.Join(v.Folder, v.File)
}

func (v *Movie) GetType() string {
	return TARGETINFO_TYPE_MOVIE
}


func (m *Movie) String() string {
	return fmt.Sprintf("movie   (id=%s, file=%s)", m.Id, filepath.Join(m.Folder, m.File))
}

func (e *Episode) GetType() string {
	return TARGETINFO_TPYE_EPISODE
}

func (e *Episode) String() string {
	return fmt.Sprintf("episode (id=%s, season=%-4d, episode=%-4d, item#=%-4d, totalItems=%-4d, file=%s)", e.Id, e.Season, e.Episode, e.ItemSeqNo, e.ItemsTotal, filepath.Join(e.Folder, e.File))
}


func NewMovie(file string, folder string, id string) *Movie {
	return &Movie{Video{Typed: Typed{ Type: TARGETINFO_TYPE_MOVIE}, File: file, Folder: folder, Id: id}}
}

func IsMovie(ti TargetInfo) bool {
	return ti != nil && TARGETINFO_TYPE_MOVIE == ti.GetType()
}

func NewEpisode(file string, folder string, id string, season int, episode int, itemSeqNo int, itemsTotal int) *Episode {
	vid := Video{Typed: Typed{ Type: TARGETINFO_TPYE_EPISODE}, File: file, Folder: folder, Id: id}
	return &Episode{Video: vid, Season: season, Episode: episode, ItemSeqNo: itemSeqNo, ItemsTotal: itemsTotal}
}

func IsEpisode(ti TargetInfo) bool {
	return ti != nil && TARGETINFO_TPYE_EPISODE == ti.GetType()
}


// read TargetInfo for specific target file (input file)
func ForTarget(workDir string, targetPath string) (TargetInfo, error) {
	targetFolder, targetFile := filepath.Split(targetPath)

	workDir, err := ripper.GetWorkPathForTargetFolder(workDir, targetFolder)
	if err != nil {
		return nil, err
	}

	return read(workDir, targetFile)
}

// read TargetInfo with specific filename
func read(workFolder string, targetFileName string) (TargetInfo, error) {
	targetInfoFile := filepath.Join(workFolder, files.WithExtension(targetFileName, targetinfo_file_extension))
	jsonRaw, err := ioutil.ReadFile(targetInfoFile)
	if err != nil {
		return nil, err
	}

	//1. get type
	typed := Typed{}
	err = json.Unmarshal(jsonRaw, &typed)

	var ti TargetInfo
	switch typed.Type {
	case TARGETINFO_TYPE_MOVIE:
		ti = &Movie{}
	case TARGETINFO_TPYE_EPISODE:
		ti = &Episode{}
	default:
		return nil, errors.New(fmt.Sprintf("target info file contains invalid target info type: \"%s\"", typed.Type))
	}

	//2. get complete target-info
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

	err := files.CreateFolderStructure(workFolder)
	if err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(ti, "", "  ")
	if err != nil {
		return err
	}
	targetFile := filepath.Join(workFolder, fileName(ti))
	err = ioutil.WriteFile(targetFile, bytes, os.ModePerm)
	return nil
}
