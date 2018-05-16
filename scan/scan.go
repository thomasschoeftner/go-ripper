package scan

import (
	"go-ripper/ripper"
	"path/filepath"
	"os"
	"strconv"
	"strings"
	"go-ripper/targetinfo"
)

func scan(rootPath string, tmp string, out string, kind string, conf *ripper.ScanConfig) ([]*targetinfo.TargetInfo, error) {
	targets := []*targetinfo.TargetInfo{}

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		path, err = filepath.Abs(path)
		if err != nil {
			return err
		}

		folder, file := filepath.Split(path)
		if strings.Contains(folder, tmp) || strings.Contains(folder, out) { //exclude tmp and out folders
			return nil
		}

		id, collection, itemNo, err := getIdCollectionItemNoForFile(path, conf.PathElemWithIdPattern)
		if err != nil {
			return err
		}

		if id != nil {
			//targets = append(targets, *newTarget(folder, file, *id, collection, itemNo))
			targets = append(targets, targetinfo.From(file, folder, kind, *id, collection,itemNo))
		}

		return nil
	})

	if err != nil {
		return nil, err
	} else {
		return targets, nil
	}
}

func getIdCollectionItemNoForFile(path string, pathElemWithIdPattern string) (*string, *int, *int, error) {
	pathElems := getLastNPathElements(path, 3)
	var id *string      //required!!!
	var collection *int //not required, id only, or id + itemNo is valid
	var itemNo *int     //sequence number of title/track/episode

	for idx, pathElem := range pathElems {
		containsId, err := filepath.Match(pathElemWithIdPattern, pathElem)
		if err != nil {
			return nil, nil, nil, err
		}

		if containsId { //set id if found and reset other flags
			id = &pathElems[idx] //TODO revise - only use id part instead of entire string
			itemNo = nil
			collection = nil
		} else if id != nil { //after id was set, set itemNo next
			if itemNo == nil {
				no, _ := strconv.Atoi(pathElems[idx]) //TODO revise - calc index from pathName, error handling
				itemNo = &no
			} else { //if a 3rd param is specified - shift use 2nd as collection, and 3rd as itemNo
				collection = itemNo
				no, _ := strconv.Atoi(pathElems[idx]) //TODO revise - calc index from pathName, error handling
				itemNo = &no
			}
		}
	}
	return id, collection, itemNo, nil
}

func getLastNPathElements(path string, n int) []string {
	pathElems := []string{}
	for i := 0; i < n; i++ {
		pathElems = append([]string{filepath.Base(path)}, pathElems...) //prepend to slice
		path = filepath.Dir(path)
	}
	return pathElems
}
