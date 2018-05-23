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
	ItemNo     int    `json:"itemno"` //track#, episode#
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

func Read(tmpFolder string, id string) (*TargetInfo, error) {
	targetFile := filepath.Join(tmpFolder, id)
	raw, err := ioutil.ReadFile(targetFile)
	if err != nil {
		return nil, err
	}
	jsonString := string(raw[:])

	ti := TargetInfo{}
	err = json.Unmarshal([]byte(jsonString), &ti)
	if err != nil {
		return nil, err
	}

	if id != ti.Id {
		return nil, errors.New(fmt.Sprintf("read error: id in filename (%s) and json (%s) do not match", id, ti.Id))
	}

	return &ti, nil
}

func Save(tmpFolder string, ti *TargetInfo) error {
	if ti == nil {
		return errors.New("target info is nil")
	}

	bytes, err := json.Marshal(ti)
	if err != nil {
		return err
	}
	targetFile := filepath.Join(tmpFolder, ti.Id)
	err = ioutil.WriteFile(targetFile, bytes, os.ModePerm)
	return err
}

func (t *TargetInfo) String() string {
	itemNo := "undef"
	if t.ItemNo != 0 {
		itemNo = fmt.Sprintf("%d", t.ItemNo)
	}
	collection := "undef"
	if t.Collection != 0 {
		collection = fmt.Sprintf("%d", t.Collection)
	}

	return fmt.Sprintf("%s(id=%s, coll=%-5s, itemNo=%-5s, file=%-80s)", t.Kind, t.Id, collection, itemNo, filepath.Join(t.Folder, t.File))
}