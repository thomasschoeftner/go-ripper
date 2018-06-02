package scan

import (
	"testing"
	"go-ripper/ripper"
	"go-cli/config"
	"go-cli/test"
	"fmt"
	"path/filepath"
)

func loadConfig(json string) (*ripper.AppConf, error) {
	conf := ripper.AppConf{}

	err := config.FromString(&conf, json, map[string]string {})
	return &conf, err
}


func validateId(t *testing.T, expected *string, found *string) {
	if expected == nil {
		if found != nil {
			t.Errorf("error matching id - expected none, but found \"%s\"", *found)
		}
	} else {
		if found == nil {
			t.Errorf("error matching id - expected \"%s\", but found none", *expected)
		} else if *expected != *found {
			t.Errorf("error matching id - expected \"%s\", but found \"%s\"", *expected, *found)
		}

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

func dissectPathAndValidate(desc string, t *testing.T,  conf *ripper.AppConf, path string, expectedId *string, expectedCol *int, expectedItemNo *int) {
	t.Run(desc , func(t *testing.T) {
		id, col, item, err := dissectPath(path, conf.Scan.Video)
		test.CheckError(t, err)
		validateId(t, expectedId, id)
		validateCollection(t, expectedCol, col)
		validateItemNo(t, expectedItemNo, item)
	})
}

func TestExtractIdOnly(t *testing.T) {
	conf, err := loadConfig(`
{
  "ignoreFolderPrefix" : ".",
  "tempDirectoryName" : "tmp",
  "outputDirectoryName" : "out",
  "scan" : {
    "video" : {
      "ignoreFolderPrefix" : ".",
      "idPattern" : "tt\\d+",
      "collectionPattern": "\\d+",
      "itemNoPattern" : "\\d+",
      "patterns" : ["<id>.*"]
    }
  }
}`)
	test.CheckError(t, err)

	{
		dissectPathAndValidate("id missing", t, conf, "/sepp/hat/gelbe/eier/ttxyz123.abc", nil, nil, nil)
	}

	{
		expectedId := "tt122345"
		path := fmt.Sprintf("/sepp/hat/gelbe/eier/%s.abc", expectedId)
		dissectPathAndValidate("id found", t, conf, path, &expectedId, nil, nil)
	}
}

func TestExtractIdItemno(t *testing.T) {
	conf, err := loadConfig(`
{
  "ignoreFolderPrefix" : ".",
  "tempDirectoryName" : "tmp",
  "outputDirectoryName" : "out",
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
		dissectPathAndValidate("id and itemno found", t, conf, path, &expectedId, nil, &expectedItemNo)
	}

	{
		dissectPathAndValidate("id not found", t, conf, "/sepp/hat/gelbe/eier/ttabcdef/666-title", nil, nil, nil)
	}
}

func TestExtractIdColItemNo(t *testing.T) {
	conf, err := loadConfig(`
{
  "ignoreFolderPrefix" : ".",
  "tempDirectoryName" : "tmp",
  "outputDirectoryName" : "out",
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
	dissectPathAndValidate("id, col, and itemno found", t, conf, path, &expectedId, &expectedCol, &expectedItemNo)
}

func TestColAndItemNoInFilename(t *testing.T) {
	conf, err := loadConfig(`
{
  "ignoreFolderPrefix" : ".",
  "tempDirectoryName" : "tmp",
  "outputDirectoryName" : "out",
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
		dissectPathAndValidate("collection and itemno in filename", t, conf, path, &expectedId, &expectedCol, &expectedItemNo)
	}
}

func TestEliminateLeadingZeroes(t *testing.T) {
	conf, err := loadConfig(`
{
  "ignoreFolderPrefix" : ".",
  "tempDirectoryName" : "tmp",
  "outputDirectoryName" : "out",
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
	dissectPathAndValidate("eliminate leading 0s in numbers", t, conf, path, &expectedId, &expectedCol, &expectedItemNo)
}

func TestExtractWithMultiplePatternOptions(t *testing.T) {
	conf, err := loadConfig(`
{
  "ignoreFolderPrefix" : ".",
  "tempDirectoryName" : "tmp",
  "outputDirectoryName" : "out",
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
		dissectPathAndValidate("id found", t, conf, fmt.Sprintf("/sepp/hat/gelbe/eier/%s.abc", expectedId), &expectedId, nil, nil)
	}

	{
		expectedId := "tt74659"
		expectedCol := 12
		expectedItemNo := 17
		dissectPathAndValidate("all found", t, conf, fmt.Sprintf("/sepp/hat/gelbe/eier/%s/%d/%d.abc", expectedId, expectedCol, expectedItemNo), &expectedId, &expectedCol, &expectedItemNo)

	}

	{
		expectedId := "tt74657"
		expectedCol := 19
		expectedItemNo := 117
		dissectPathAndValidate("all found", t, conf, fmt.Sprintf("/sepp/hat/gelbe/eier/%s/%d/title %d.abc", expectedId, expectedCol, expectedItemNo), &expectedId, &expectedCol, &expectedItemNo)
	}
}

func TestGetLastNPathElements(t *testing.T) {
	{
		expected := "c/d/e.fg"
		path := "a/b/" + expected
		last3 := getLastNPathElements(path, 3)
		if last3 != expected {
			t.Errorf("getting last 3 path elements from %s - expected %s, but got %s", path, expected, last3)
		}
	}

	{
		path := "a/b.cde"
		expected := "./" + path
		last3 := getLastNPathElements(path, 3)
		if last3 != expected {
			t.Errorf("getting last 3 path elements from %s - expected %s, but got %s", path, expected, last3)
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
  "ignoreFolderPrefix" : ".",
  "tempDirectoryName" : "tmp",
  "outputDirectoryName" : "out",
  "scan" : {
    "video" : {
      "idPattern" : "tt\\d+",
      "collectionPattern": "\\d+",
      "itemNoPattern" : "\\d+",
      "patterns" : ["<id>.*/.*","<id>.*"]
    }
  }
}`

	movies := []string{"tt987654321", "tt555", "tt666", "tt34543"}
	testScanVideos(t, movies, nil, conf, "testdata")
}


func TestScanMixedSinglesAndCollections(t *testing.T) {
	conf := `
{
  "ignoreFolderPrefix" : ".",
  "tempDirectoryName" : "tmp",
  "outputDirectoryName" : "out",
  "scan" : {
    "video" : {
      "idPattern" : "tt\\d+",
      "collectionPattern": "\\d+",
      "itemNoPattern" : "\\d+",
      "patterns" : [
        "<id>.*/season\\s<collection>/<itemno>.*",
        "<id>.*/.*",
        "<id>.*"]
    }
  }
}`
	movies := []string{"tt987654321", "tt555", "tt666", "tt34543"}
	episodes := map[int]map[string]int {
		2 : {"0.mkv":0, "1.txt":1, "2.a.b.c":2, "4":4}, //TODO change to logical numbering?
		4 : {"1.txt":1, "2.a.b.c":2, "3":3}}
	testScanVideos(t, movies, episodes, conf, "testdata")
}


func TestScanDeep(t *testing.T) {
	conf := `
{
  "ignoreFolderPrefix" : ".",
  "tempDirectoryName" : "tmp",
  "outputDirectoryName" : "out",
  "scan" : {
    "video" : {
      "idPattern" : "tt\\d+",
      "collectionPattern": "\\d+",
      "itemNoPattern" : "\\d+",
      "patterns" : [
        "<id>.*/season\\s<collection>/<itemno>.*",
        "<id>.*/season\\s<collection>/.*/.*/.*/<itemno>.*",
        "<id>.*/.*",
        "<id>.*"]
    }
  }
}`

	movies := []string{"tt987654321", "tt555", "tt666", "tt34543"}
	episodes := map[int]map[string]int {
		2 : {"0.mkv":0, "1.txt":1, "2.a.b.c":2, "4":4}, //TODO change to logical numbering
		4 : {"1.txt":1, "2.a.b.c":2, "3":3},
		7 : {"6.mkv":6, "7.a.b.c":7, "8":8}}

	testScanVideos(t, movies, episodes, conf, "testdata")
}


func testScanVideos(t *testing.T, expectedMovies []string, expectedEpisodes map[int]map[string]int, confStr string, testDataFolder string) {
	conf, err := loadConfig(confStr)
	test.CheckError(t, err)

	path, _ := filepath.Abs(filepath.Join(".", testDataFolder))
	targetInfos, err := scan(path, "video", conf.IgnoreFolderPrefix, conf.Scan.Video)
	test.CheckError(t, err)

	for _, ti := range targetInfos {
		fmt.Printf("%v\n", *ti)
	}

	expectedNoOfItems := len(expectedMovies)
	for _, season := range expectedEpisodes {
		expectedNoOfItems = expectedNoOfItems + len(season)
	}

	if len(targetInfos) == expectedNoOfItems + 1 {
		t.Errorf("scaning %s yielded unexpected number of relevant files - .hidden folder should not be searched", path)
	} else if len(targetInfos) != expectedNoOfItems {
		t.Errorf("scaning %s yielded unexpected number of relevant files - expected %d, but got %d", path, expectedNoOfItems, len(targetInfos))
	}

	for _, ti := range targetInfos {

		if collection, isEpisode := expectedEpisodes[ti.Collection]; ti.Collection == 0 || !isEpisode {
			if !idFoundIn(ti.Id, expectedMovies) {
				t.Errorf("unexpected video %d extracted from path %s", ti.Collection, filepath.Join(ti.Folder, ti.File))
			}
		} else {
			if itemNo, found := collection[ti.File]; !found {
				t.Errorf("unexpected file %s extracted from path %s", ti.File, filepath.Join(ti.Folder, ti.File))
			} else {
				if ti.ItemNo != itemNo {
					t.Errorf("extracted itemNo %d from path %s, but should be itemNo %d", ti.ItemNo, filepath.Join(ti.Folder, ti.File), itemNo)
				}
			}

		}
	}

}

func idFoundIn(id string, ids []string) bool {
	for _, m := range ids {
		if m == id {
			return true
		}
	}
	return false
}
