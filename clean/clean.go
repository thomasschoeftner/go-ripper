package clean

import (
	"go-cli/task"
	"go-ripper/ripper"
	"go-cli/commons"
	"path/filepath"
	"os"
)

func CleanHandler(ctx task.Context) task.HandlerFunc {
	return func (job task.Job) ([]task.Job, error) {
		conf := ctx.Config.(*ripper.AppConf)
		return clean(ctx.Printf, "work data", job, conf.WorkDirectory)
	}
}


func clean(printf commons.FormatPrinter, desc string, job task.Job, workDir string) ([]task.Job, error) {
	result := []task.Job{job}
	workPath, err := ripper.GetWorkPathForJob(workDir, job)
	if err != nil {
		return result, err
	}

	targetPath:= ripper.GetTargetFileFromJob(job)
	_, targetFile := filepath.Split(targetPath)
	filePattern := filepath.Join(workPath, targetFile) + "*"

	printf("cleaning %s related to target \"%s\" (matching %s)\n", desc, targetPath, filePattern)

	filesToDelete, err := filepath.Glob(filePattern)
	if err != nil {
		printf("cleaning failed\n  due to: %s\n", err)
	} else {
		for _, f := range filesToDelete {
			printf("  deleting file: %s\n", f)
			os.Remove(f)
		}
		printf("all artifacts removed\n")
	}
	return result, err
}
