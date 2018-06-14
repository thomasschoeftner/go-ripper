package scan

import (
	"testing"
	"go-cli/task"
	"go-cli/commons"
	"go-cli/test"
	"go-ripper/ripper"
	"go-ripper/targetinfo"
)

func TestScanVideo(t *testing.T) {
	confStr := `
{
  "ignorePrefix" : ".",
  "workDirectory" : "tmp",
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

	conf, err := loadConfig(confStr)
	test.CheckError(t, err)
	ctx := task.Context{nil, conf, commons.Printf, false}
	handler := ScanVideo(ctx)
	job := task.Job{ripper.JobField_Path : "./testdata"}
	results, err := handler(job)
	test.CheckError(t, err)
	expectedNoOfSearchResults := 14
	if len(results) != expectedNoOfSearchResults {
		t.Errorf("found %d number of search results, but expected %d", len(results), expectedNoOfSearchResults)
	}
}

func TestToTargetInfos(t *testing.T) {
	t.Run("nil scan results", func(t *testing.T) {
		ti, err := toTargetInfos(nil)
		test.CheckError(t, err)
		if len(ti) != 0 {
			t.Errorf("expected empty target info list, but got %v", ti)
		}
	})

	t.Run("empty scan results", func(t *testing.T) {
		ti, err := toTargetInfos([]*scanResult{})
		test.CheckError(t, err)
		if len(ti) != 0 {
			t.Errorf("expected empty target info list, but got %v", ti)
		}
	})

	t.Run("count total number of episodes", func(t *testing.T) {
		sr := []*scanResult {
			newScanResult("a/3", "1",  "a", 3, 1),
			newScanResult("a/3", "2",  "a", 3, 2),
			newScanResult("a/3", "3",  "a", 3, 3)}
		targetInfos, err := toTargetInfos(sr)
		test.CheckError(t, err)
		if len(targetInfos) != len(sr) {
			t.Errorf("got %d target infos, but expected %d", len(targetInfos), len(sr))
		}
		ep := targetInfos[0].(*targetinfo.Episode)
		if ep.ItemsTotal != len(sr) {
			t.Errorf("got total # of episodes %d, but expected %d", ep.ItemsTotal, len(sr))
		}
	})

}

func newScanResult(folder string, file string, id string, season int, episode int) *scanResult {
	return &scanResult{Folder: folder, File: file, Id: id, Collection: &season, ItemNo: &episode}
}