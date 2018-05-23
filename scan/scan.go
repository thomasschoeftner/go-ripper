package scan

import (
	"go-ripper/ripper"
	"path/filepath"
	"os"
	"strings"
	"go-ripper/targetinfo"
	"fmt"
	"regexp"
	"strconv"
)

func scan(rootPath string, excludeDirs []string, kind string, conf *ripper.ScanConfig) ([]*targetinfo.TargetInfo, error) {
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

		//discard excluded directories
		folder, file := filepath.Split(path)
		for _, dir := range excludeDirs {
			if strings.Contains(folder, dir) {
				return nil
			}
		}

		id, collection, itemNo, err := dissectPath(path, conf)
		if err != nil {
			return err
		}
		if id != nil {
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

const (
	placeholder_Id = "id"
	placeholder_Collection = "collection"
	placeholder_ItemNo = "itemno"
)

func dissectPath(path string, conf *ripper.ScanConfig) (*string, *int, *int, error) {
	 for _, pattern := range conf.Patterns {
	 	expandedPattern := expandPatterns(pattern, conf.IdPattern, conf.CollectionPattern, conf.ItemNoPattern)
	 	pathTrail := getLastNPathElements(path, strings.Count(expandedPattern, "/") + 1)
	 	re, err := regexp.Compile(expandedPattern)
	 	if err != nil {
	 		return nil, nil, nil, err
		}

		matches := extractParams(re, pathTrail)
		if matches != nil {
			if idVal, isDefined := matches[placeholder_Id]; isDefined {
				id := &idVal
				var col *int
				if colVal, isDefined := matches[placeholder_Collection]; isDefined {
					i, _ := strconv.Atoi(colVal)
					col = &i
				}
				var itemNo *int
				if itemVal, isDefined := matches[placeholder_ItemNo]; isDefined {
					i, _ := strconv.Atoi(itemVal)
					itemNo = &i
				}
				return id, col, itemNo, nil
			}
		}
	 }
	return nil, nil, nil, nil
}

func extractParams(re *regexp.Regexp, path string) map[string]string {
	results := make(map[string]string)
	matches := re.FindStringSubmatch(path)
	if matches == nil {
		return results
	}

	for idx, name := range re.SubexpNames() {
		if idx == 0 || len(name) == 0 {
			continue //ignore first match (contains original string)
		}
		results[name] = matches[idx]
	}
	return results
}

func expandPatterns(pattern string, idPattern string, colPattern string, itemNoPattern string) string {
	expanded := expandPattern(pattern, placeholder_Id, idPattern)
	expanded = expandPattern(expanded, placeholder_Collection, colPattern)
	expanded = expandPattern(expanded, placeholder_ItemNo, itemNoPattern)
	//keep linux file separator!!
	return expanded
}

func expandPattern(pattern string, placeholder string, subPattern string) string {
	replacement := fmt.Sprintf("(?P<%s>%s)", placeholder, subPattern) //e.g. (?P<id>\d*)
	return strings.Replace(pattern, fmt.Sprintf("<%s>",placeholder), replacement, -1)
}

func getLastNPathElements(path string, n int) string {
	pathElems := []string{ }

	path = filepath.Clean(path)
	for i := 0; i < n; i++ {
		last := filepath.Base(path)
		path = filepath.Dir(path)

		pathElems = append([]string{last}, pathElems...)
	}

	//keep/change to linux file separator!!
	return strings.Join(pathElems, "/")
}
