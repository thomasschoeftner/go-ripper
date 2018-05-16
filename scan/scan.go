package scan

import (
	"go-ripper/ripper"
	"path/filepath"
	"os"
	"fmt"
	"strconv"
	"strings"
)

func scan(rootPath string, tmp string, out string, conf *ripper.ScanConfig) ([]target, error) {
	targets := []target{}

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

		fmt.Printf("    file at %-72s", path) //TODO remove
		folder, file := filepath.Split(path)
		if strings.Contains(folder, tmp) || strings.Contains(folder, out) {
			fmt.Printf("  -  is tmp / out folder (ignore!)\n") //TODO remove
			return nil
		}

		id, collection, group, err := getIdCollectionGroupForFile(path, conf.PathElemWithIdPattern)
		if err != nil {
			return err
		}

		{ //TODO remove block
			if id != nil {
				fmt.Printf("  -  yields id=%s", *id)
			} else {
				fmt.Printf("  -  is NOT a relevant input file")
			}
			if collection != nil {
				fmt.Printf(", collection=%s", *collection)
			}
			if group != nil {
				fmt.Printf(", group=%d", *group)
			}
			fmt.Printf("\n")
		}

		if id != nil {
			targets = append(targets, *newTarget(folder, file, *id, collection, group))
		}
		return nil
	})

	if err != nil {
		return nil, err
	} else {
		return targets, nil
	}
}

func getIdCollectionGroupForFile(path string, pathElemWithIdPattern string) (*string, *string, *int, error) {
	pathElems := getLastNPathElements(path, 3)
	var id, collection *string
	var group *int

	for idx, pathElem := range pathElems {
		containsId, err := filepath.Match(pathElemWithIdPattern, pathElem)
		if err != nil {
			return nil, nil, nil, err
		}

		if containsId { //set id if found and reset other flags
			id = &pathElems[idx] //TODO revise - only use id part instead of entire string
			collection = nil
			group = nil
		} else if id != nil { //after id was set, set collection next
			if collection == nil {
				collection = &pathElems[idx]
			} else { //set group last after id and collection
				no, _ := strconv.Atoi(pathElems[idx]) //TODO revise - calc index from pathName
				group = &no
			}
		}
	}
	return id, collection, group, nil
}

func getLastNPathElements(path string, n int) []string {
	pathElems := []string{}
	for i := 0; i < n; i++ {
		pathElems = append([]string{filepath.Base(path)}, pathElems...) //prepend to slice
		path = filepath.Dir(path)
	}
	return pathElems
}
