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
	path, err := ripper.GetWorkPathForJob(workDir, job)
	if err != nil {
		return result, err
	}

	_, fName := filepath.Split(ripper.GetTargetFileFromJob(job))
	filePattern := filepath.Join(path, fName) + "*"
	printf("cleaning %s related to target \"%s\"\n", desc, ripper.GetTargetFileFromJob(job))
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
