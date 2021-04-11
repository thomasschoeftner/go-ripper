package scan

import (
	"testing"
	"github.com/thomasschoeftner/go-ripper/ripper"
	"github.com/thomasschoeftner/go-cli/config"
	"github.com/thomasschoeftner/go-cli/test"
	"fmt"
	"path/filepath"
	"github.com/thomasschoeftner/go-cli/commons"
)

func loadConfig(json string) (*ripper.AppConf, error) {
	conf := ripper.AppConf{}

	err := config.FromString(&conf, json, map[string]string {})
	return &conf, err
}

func validateId(t *testing.T, expected string, found string) {
	if expected != found {
		t.Errorf("error matching id - expected \"%s\", but found \"%s\"", expected, found)
	}
}

func validateNumeric(desc string) func (t *testing.T, expected *int, found *int) {
	return func (t *testing.T, expected *int, found *int) {
		if expected == nil {
			if found != nil {
				t.Errorf("error matching %s - expected none, but found \"%d\"", desc, *found)
			}
		} else {
			if found == nil {
				t.Errorf("error matching %s - expected \"%d\", but found none", desc, *expected)
			} else if *expected != *found {
				t.Errorf("error matching %s - expected \"%d\", but found \"%d\"", desc, *expected, *found)
			}
		}
	}
}
var validateCollection = validateNumeric("collection")
var validateItemNo = validateNumeric("itemno")

func dissectPathAndValidate(desc string, t *testing.T,  conf *ripper.AppConf, path string, expectedId string, expectedCol *int, expectedItemNo *int) {
	t.Run(desc , func(t *testing.T) {
		result, err := dissectPath(path, conf.Scan.Video)
		test.CheckError(t, err)
		validateId(t, expectedId, result.Id)
		validateCollection(t, expectedCol, result.Collection)
		validateItemNo(t, expectedItemNo, result.ItemNo)
	})
}

func dissectPathAndValidateNoMatch(desc string, t *testing.T,  conf *ripper.AppConf, path string) {
	t.Run(desc , func(t *testing.T) {
		result, err := dissectPath(path, conf.Scan.Video)
		test.CheckError(t, err)
		if result != nil {
			t.Errorf("error matching path - expected no match, but got %s", result.Id)
		}
	})
}


func TestExtractIdOnly(t *testing.T) {
	conf, err := loadConfig(`
{
  "ignorePrefix" : ".",
  "workDirectory" : "tmp",
  "scan" : {
    "video" : {
      "idPattern" : "tt\\d+",
      "collectionPattern": "\\d+",
      "itemNoPattern" : "\\d+",
      "patterns" : ["<id>.*"]
    }
  }
}`)
	test.CheckError(t, err)

	{
		dissectPathAndValidateNoMatch("id missing", t, conf, "/sepp/hat/gelbe/eier/ttxyz123.abc")
	}

	{
		expectedId := "tt122345"
		path := fmt.Sprintf("/sepp/hat/gelbe/eier/%s.abc", expectedId)
		dissectPathAndValidate("id found", t, conf, path, expectedId, nil, nil)
	}
}

func TestExtractIdItemno(t *testing.T) {
	conf, err := loadConfig(`
{
  "ignorePrefix" : ".",
  "workDirectory" : "tmp",
  "scan" : {
    "video" : {
      "idPattern" : "tt\\d+",
      "collectionPattern": "\\d+",
      "itemNoPattern" : "\\d+",
      "patterns" : ["<id>.*/<itemno>.*"]
    }
  }
}`)
	test.CheckError(t, err)

	{
		expectedId := "tt3453645"
		expectedItemNo := 13
		path := fmt.Sprintf("/sepp/hat/gelbe/eier/%s-name/%d-title", expectedId, expectedItemNo)
		dissectPathAndValidate("id and itemno found", t, conf, path, expectedId, nil, &expectedItemNo)
	}

	{
		dissectPathAndValidateNoMatch("id not found", t, conf, "/sepp/hat/gelbe/eier/ttabcdef/666-title")
	}
}

func TestExtractIdColItemNo(t *testing.T) {
	conf, err := loadConfig(`
{
  "ignorePrefix" : ".",
  "workDirectory" : "tmp",
  "scan" : {
    "video" : {
      "idPattern" : "tt\\d+",
      "collectionPattern": "\\d+",
      "itemNoPattern" : "\\d+",
      "patterns" : [".*<id>.*/s<collection>/e<itemno>.*"]
    }
  }
}`)
	test.CheckError(t, err)

	expectedId := "tt765765"
	expectedCol := 678
	expectedItemNo := 32
	path := fmt.Sprintf("sepp/hat/gelbe/eier/name-%s/s%d/e%d-itemname.txt", expectedId, expectedCol, expectedItemNo)
	dissectPathAndValidate("id, col, and itemno found", t, conf, path, expectedId, &expectedCol, &expectedItemNo)
}

func TestColAndItemNoInFilename(t *testing.T) {
	conf, err := loadConfig(`
{
  "ignorePrefix" : ".",
  "workDirectory" : "tmp",
  "scan" : {
    "video" : {
      "idPattern" : "tt\\d+",
      "collectionPattern": "\\d+",
      "itemNoPattern" : "\\d+",
      "patterns" : ["<id>.*/s<collection>e<itemno>.*"]
    }
  }
}`)
	test.CheckError(t, err)
	{
		expectedId := "tt46864"
		expectedCol := 6
		expectedItemNo := 10

		path := fmt.Sprintf("sepp/hat/gelbe/eier/%s-title/s%de%d", expectedId, expectedCol, expectedItemNo)
		dissectPathAndValidate("collection and itemno in filename", t, conf, path, expectedId, &expectedCol, &expectedItemNo)
	}
}

func TestEliminateLeadingZeroes(t *testing.T) {
	conf, err := loadConfig(`
{
  "ignorePrefix" : ".",
  "workDirectory" : "tmp",
  "scan" : {
    "video" : {
      "idPattern" : "tt\\d+",
      "collectionPattern": "\\d+",
      "itemNoPattern" : "\\d+",
      "patterns" : ["<id>.*/<collection>.*/e<itemno>.*"]
    }
  }
}`)
	test.CheckError(t, err)

	expectedId := "tt333444"
	expectedCol := 34
	expectedItemNo := 56

	path := fmt.Sprintf("sepp/hat/gelbe/eier/%s-name/00%d/e000%d-itemname.txt", expectedId, expectedCol, expectedItemNo)
	dissectPathAndValidate("eliminate leading 0s in numbers", t, conf, path, expectedId, &expectedCol, &expectedItemNo)
}

func TestExtractWithMultiplePatternOptions(t *testing.T) {
	conf, err := loadConfig(`
{
  "ignorePrefix" : ".",
  "workDirectory" : "tmp",
  "scan" : {
    "video" : {
      "idPattern" : "tt\\d+",
      "collectionPattern": "\\d+",
      "itemNoPattern" : "\\d+",
      "patterns" : [
        "<id>.*/s<collection>e<itemno>.*",
        "<id>.*/<collection>/<itemno>.*",
        "<id>.*/<collection>/\\D*<itemno>.*",
        "<id>.*/<itemno>.*",
        "<id>.*"]
    }
  }
}`)
	test.CheckError(t, err)
	{
		expectedId := "tt74658"
		dissectPathAndValidate("id found", t, conf, fmt.Sprintf("/sepp/hat/gelbe/eier/%s.abc", expectedId), expectedId, nil, nil)
	}

	{
		expectedId := "tt74659"
		expectedCol := 12
		expectedItemNo := 17
		dissectPathAndValidate("all found", t, conf, fmt.Sprintf("/sepp/hat/gelbe/eier/%s/%d/%d.abc", expectedId, expectedCol, expectedItemNo), expectedId, &expectedCol, &expectedItemNo)

	}

	{
		expectedId := "tt74657"
		expectedCol := 19
		expectedItemNo := 117
		dissectPathAndValidate("all found", t, conf, fmt.Sprintf("/sepp/hat/gelbe/eier/%s/%d/title %d.abc", expectedId, expectedCol, expectedItemNo), expectedId, &expectedCol, &expectedItemNo)
	}
}

func TestGetLastNPathElements(t *testing.T) {
	{
		expected := "c/d/e.fg"
		path := "a/b/" + expected
		last3 := getLastNPathElements(path, 3)
		if last3 != expected {
			t.Errorf("getting last 3.avi path elements from %s - expected %s, but got %s", path, expected, last3)
		}
	}

	{
		path := "a/b.cde"
		expected := "./" + path
		last3 := getLastNPathElements(path, 3)
		if last3 != expected {
			t.Errorf("getting last 3.avi path elements from %s - expected %s, but got %s", path, expected, last3)
		}
	}

	{
		expected := "b/c"
		path := "a/b/c/"
		last2 := getLastNPathElements(path, 2)
		if last2 != expected {
			t.Errorf("getting last 2 path elements from %s - expected %s, but got %s", path, expected, last2)
		}
	}
}




func TestScanSingleTitles(t *testing.T) {
	conf := `
{
  "ignorePrefix" : ".",
  "workDirectory" : "tmp",
  "scan" : {
    "video" : {
      "idPattern" : "tt\\d+",
      "collectionPattern": "\\d+",
      "itemNoPattern" : "\\d+",
      "patterns" : ["<id>.*/.*","<id>.*"],
      "allowSpaces" : false,
      "allowedExtensions" : ["avi"]
    }
  }
}`

	movies := []string{"tt987654321", "tt555", "tt666", "tt34543"}
	testScanVideos(t, movies, nil, conf, "testdata")
}


func TestScanMixedSinglesAndCollections(t *testing.T) {
	conf := `
{
  "ignorePrefix" : ".",
  "workDirectory" : "tmp",
  "scan" : {
    "video" : {
      "idPattern" : "tt\\d+",
      "collectionPattern": "\\d+",
      "itemNoPattern" : "\\d+",
      "patterns" : [
        "<id>.*/season<collection>/<itemno>.*",
        "<id>.*/.*",
        "<id>.*"],
        "allowSpaces" : false,
        "allowedExtensions" : ["avi"]
    }
  }
}`
	movies := []string{"tt987654321", "tt555", "tt666", "tt34543"}
	episodes := map[int]map[string]int {
		2 : {"0.avi":0, "1.avi":1, "2.a.b.c.avi":2, "4.avi":4}, //TODO change to logical numbering?
		4 : {"1.avi":1, "2.a.b.c.avi":2, "3.avi":3}}
	testScanVideos(t, movies, episodes, conf, "testdata")
}


func TestScanDeep(t *testing.T) {
	conf := `
{
  "ignorePrefix" : ".",
  "workDirectory" : "tmp",
  "scan" : {
    "video" : {
      "idPattern" : "tt\\d+",
      "collectionPattern": "\\d+",
      "itemNoPattern" : "\\d+",
      "patterns" : [
        "<id>.*/season<collection>/<itemno>.*",
        "<id>.*/season<collection>/.*/.*/.*/<itemno>.*",
        "<id>.*/.*",
        "<id>.*"],
        "allowSpaces" : false,
        "allowedExtensions" : ["avi"]
    }
  }
}`

	movies := []string{"tt987654321", "tt555", "tt666", "tt34543"}
	episodes := map[int]map[string]int {
		2 : {"0.avi":0, "1.avi":1, "2.a.b.c.avi":2, "4.avi":4}, //TODO change to logical numbering
		4 : {"1.avi":1, "2.a.b.c.avi":2, "3.avi":3},
		7 : {"6.avi":6, "7.a.b.c.avi":7, "8.avi":8}}

	testScanVideos(t, movies, episodes, conf, "testdata")
}

func testScanVideos(t *testing.T, expectedMovies []string, expectedEpisodes map[int]map[string]int, confStr string, testDataFolder string) {
	conf, err := loadConfig(confStr)
	test.CheckError(t, err)

	path, _ := filepath.Abs(filepath.Join(".", testDataFolder))
	results, err := scan(path, conf.IgnorePrefix, conf.Scan.Video, commons.Printf)
	test.CheckError(t, err)

	for _, result := range results  {
		  fmt.Printf("found %v\n", *result)
	}

	expectedNoOfMatches := len(expectedMovies)
	for _, season := range expectedEpisodes {
		expectedNoOfMatches = expectedNoOfMatches + len(season)
	}

	if len(results) == expectedNoOfMatches+ 1 {
		t.Errorf("scaning %s yielded unexpected number of relevant files - .hidden folder should not be searched", path)
	} else if len(results) != expectedNoOfMatches {
		t.Errorf("scaning %s yielded unexpected number of relevant files - expected %d, but got %d", path, expectedNoOfMatches, len(results))
	}

	for _, result := range results {
		s := 0
		if result.Collection !=  nil {
			s = *result.Collection
		}

		if season, seasonFound := expectedEpisodes[s]; s == 0 || !seasonFound {
			if !idFoundIn(result.Id, expectedMovies) {
				t.Errorf("unexpected video %s extracted from path %s", result.Id, filepath.Join(result.Folder, result.File))
			}
		} else {
			if episode, episodeFound := season[result.File]; !episodeFound {
				t.Errorf("unexpected file %s extracted from path %s", result.File, filepath.Join(result.Folder, result.File))
			} else {
				if episode != *result.ItemNo {
					t.Errorf("extracted itemNo %d from path %s, but should be itemNo %d", *result.ItemNo, filepath.Join(result.Folder, result.File), episode)
				}
			}
		}
	}
}

func TestExclusion(t *testing.T) {
	const ignorePrefix = "."
	t.Run("do not ignore ordinary files with valid extension", func(t *testing.T) {
		path := "a/b/c.valid"
		if shouldIgnore(path, ignorePrefix, []string{}, []string{"valid"}) {
			t.Errorf("expected file \"%s\" not to be ignored, but is", path)
		}
	})

	t.Run("ignore ordinary files with invalid extension", func(t *testing.T) {
		path := "a/b/c.invalid"
		if !shouldIgnore(path, ignorePrefix, []string{}, []string{"valid"}) {
			t.Errorf("expected file \"%s\" to be ignored, but is not", path)
		}
	})

	t.Run("ignore file", func(t *testing.T) {
		path := "a/b/.c.valid"
		if !shouldIgnore(path, ignorePrefix, []string{}, []string {"valid"}) {
			t.Errorf("expected file \"%s\" to be ignored, but is not", path)
		}
	})

	t.Run("ignore folder", func(t *testing.T) {
		path := "a/.b/c.valid"
		if !shouldIgnore(path, ignorePrefix, []string{}, []string {"valid"}) {
			t.Errorf("expected folder content \"%s\" to be ignored, but is not", path)
		}
	})

	t.Run("ignore ignored sub-folder", func(t *testing.T) {
		path := "a/.b/c/d/e.valid"
		if !shouldIgnore(path, ignorePrefix, []string{"a/.b"}, []string {"valid"}) {
			t.Errorf("expected folder content \"%s\" to be ignored, but is not", path)
		}
	})
}

func idFoundIn(id string, ids []string) bool {
	for _, m := range ids {
		if m == id {
			return true
		}
	}
	return false
}
