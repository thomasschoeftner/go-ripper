package scan

import (
	"github.com/thomasschoeftner/go-ripper/ripper"
	"path/filepath"
	"os"
	"strings"
	"fmt"
	"regexp"
	"strconv"
	"github.com/thomasschoeftner/go-ripper/files"
	"github.com/thomasschoeftner/go-cli/commons"
)

type scanResult struct {
	Folder     string
	File       string
	Id         string
	Collection *int
	ItemNo     *int
}

func scan(rootPath string, ignorePrefix string, conf *ripper.ScanConfig, printf commons.FormatPrinter) ([]*scanResult, error) {
	results := []*scanResult{}
	ignoredFolders := []string{}

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		//ignore spaces if required
		if !conf.AllowSpaces && strings.Contains(path, " ") {
			printf("WARNING - ignore file \"%s\" due to spaces in path\n", path)
			return nil
		}

		if shouldIgnore(path, ignorePrefix, ignoredFolders, conf.AllowedExtensions) {
			return nil
		}

		result, err := dissectPath(path, conf)

		if err != nil {
			return err
		}
		if result != nil {
			results = append(results, result)
		}

		return nil
	})

	if err != nil {
		return nil, err
	} else {
		return results, nil
	}
}

func shouldIgnore(path string, ignorePrefix string, ignoredFolders []string, allowedExtensions []string) bool {
	folder, file := filepath.Split(path)
	//discard excluded files
	if strings.HasPrefix(file, ignorePrefix) {
		return true
	}

	//discard non-matching file extensions
	_, ext := files.SplitExtension(file)
	if !commons.IsStringAmong(ext, allowedExtensions) {
		return true
	}

	// discard files in ignore folders, and keep list of ignored folders
	folderName := filepath.Base(folder)
	if strings.HasPrefix(folderName, ignorePrefix) {
		//TODO optimize by storing an ignored folder only once?
		ignoredFolders = append(ignoredFolders, folder)
		return true
	}

	//discard sub-directories of excluded directories
	for _, ignored := range ignoredFolders {
		if strings.HasPrefix(folder, ignored) {
			return true
		}
	}

	return false
}


const (
	placeholder_Id = "id"
	placeholder_Collection = "collection"
	placeholder_ItemNo = "itemno"
)

func dissectPath(path string, conf *ripper.ScanConfig) (*scanResult, error) {
	 for _, pattern := range conf.Patterns {
	 	expandedPattern := expandPatterns(pattern, conf.IdPattern, conf.CollectionPattern, conf.ItemNoPattern)
	 	pathTrail := getLastNPathElements(path, strings.Count(expandedPattern, "/") + 1) //folder depth + file nameee
	 	re, err := regexp.Compile(expandedPattern)
	 	if err != nil {
	 		return nil, err
		}

		matches := extractParams(re, pathTrail)
		if matches != nil {
			if id, isDefined := matches[placeholder_Id]; isDefined {
				var collection *int
				if collectionVal, isDefined := matches[placeholder_Collection]; isDefined {
					i, _ := strconv.Atoi(collectionVal)
					collection = &i
				}
				var itemNo *int
				if itemVal, isDefined := matches[placeholder_ItemNo]; isDefined {
					i, _ := strconv.Atoi(itemVal)
					itemNo = &i
				}
				folder, file := filepath.Split(path)
				return &scanResult{Folder: folder, File: file, Id: id, Collection: collection, ItemNo: itemNo}, nil

			}
		}
	 }
	return nil, nil
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
