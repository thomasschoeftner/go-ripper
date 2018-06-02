package scan

import (
	"testing"
	"go-cli/task"
	"go-cli/commons"
	"go-cli/test"
	"go-ripper/ripper"
)

func TestScanVideo(t *testing.T) {
	confStr := `
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
