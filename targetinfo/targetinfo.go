package targetinfo

import (
	"io/ioutil"
	"encoding/json"
	"os"
	"errors"
	"fmt"
	"path/filepath"
)

type TargetInfo struct {
	File       string `json:"file"` //target file name
	Folder     string `json:"folder"` //folder containing target file
	Kind       string `json:"kind"` //type of target (e.g. audio, video)
	Id         string `json:"id"` //id for name resolve
	Collection int    `json:"collection"` //Season#, CD#, etc.
	ItemNo     int    `json:"itemno"` //track#, singleCollectionItem#
}

func From(file string, folder string, kind string, id string, collection *int, itemNo *int) *TargetInfo {
	c := 0
	if collection != nil {
		c = *collection
	}

	i := 0
	if itemNo != nil {
		i = *itemNo
	}

	return &TargetInfo{file, folder, kind, id, c, i}
}

func Read(targetInfoFile string) (*TargetInfo, error) {
	raw, err := ioutil.ReadFile(targetInfoFile)
	if err != nil {
		return nil, err
	}
	jsonString := string(raw[:])

	ti := TargetInfo{}
	err = json.Unmarshal([]byte(jsonString), &ti)
	if err != nil {
		return nil, err
	}

	return &ti, nil
}

func Save(tmpFolder string, ti *TargetInfo) (*string, error) {
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

func (ti *TargetInfo) fileName() string {
	if ti.Collection != 0 && ti.ItemNo != 0 {
		return fmt.Sprintf("%s.%d.%d", ti.Id, ti.Collection, ti.ItemNo)
	} else if ti.ItemNo != 0 {
		return fmt.Sprintf("%s.%d", ti.Id, ti.ItemNo)
	} else {
		return ti.Id
	}
}

func (ti *TargetInfo) String() string {
	itemNo := "undef"
	if ti.ItemNo != 0 {
		itemNo = fmt.Sprintf("%d", ti.ItemNo)
	}
	collection := "undef"
	if ti.Collection != 0 {
		collection = fmt.Sprintf("%d", ti.Collection)
	}

	return fmt.Sprintf("%s(id=%s, coll=%-5s, itemNo=%-5s, file=%s)", ti.Kind, ti.Id, collection, itemNo, filepath.Join(ti.Folder, ti.File))
}