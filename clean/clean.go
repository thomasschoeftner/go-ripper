package clean

import (
	"os"
	"go-cli/task"
	"go-ripper/ripper"
	"go-cli/commons"
)

func CleanTmpHandler(ctx task.Context, job task.Job) ([]task.Job, error) {
	conf := ctx.Config.(*ripper.AppConf)
	return clean(ctx.Printf, "temporary data", job, conf.TempDirectoryName)
}

func CleanOutHandler(ctx task.Context, job task.Job) ([]task.Job, error) {
	conf := ctx.Config.(*ripper.AppConf)
	return clean(ctx.Printf, "output data", job, conf.OutputDirectoryName)
}

func clean(printf commons.FormatPrinter, desc string, job task.Job, subDir string) ([]task.Job, error) {
	path := ripper.GetWorkPathFor(job, subDir)
	printf("  cleaning %s from \"%s\"", desc, path)
	error := os.RemoveAll(path)
	if error != nil {
		printf("  ...failed\n  due to: %s\n", error)
	} else {
		printf("  ...done\n")
	}
	return []task.Job{job}, error
}
