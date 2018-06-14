package clean

import (
	"os"
	"go-cli/task"
	"go-ripper/ripper"
	"go-cli/commons"
	"path/filepath"
)

func CleanHandler(ctx task.Context) task.HandlerFunc {
	return func (job task.Job) ([]task.Job, error) {
		conf := ctx.Config.(*ripper.AppConf)
		return clean(ctx.Printf, "work data", job, conf.WorkDirectory)
	}
}


func clean(printf commons.FormatPrinter, desc string, job task.Job, workDir string) ([]task.Job, error) {
	result := []task.Job{job}
	path, err := ripper.GetWorkPathFor(workDir, job)
	if err != nil {
		return result, err
	}

	_, fName := filepath.Split(ripper.GetTargetFilePathForm(job))
	filePattern := filepath.Join(path, fName) + "*"
	printf("cleaning %s of \"%s\" (%s)\n", desc, ripper.GetTargetFilePathForm(job), filePattern)
	filesToDelete, err := filepath.Glob(filePattern)
	for _, f := range filesToDelete {
		printf("  deleting file: %s\n", f)
		os.Remove(f)
	}

	if err != nil {
		printf("cleaning failed\n  due to: %s\n", err)
	}
	return result, err
}
