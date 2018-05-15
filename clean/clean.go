package clean

import (
	"os"
	"go-cli/task"
	"go-ripper/ripper"
)

func CleanHandler(ctx task.Context, job task.Job) ([]task.Job, error) {
	conf := ctx.Config.(*ripper.AppConf)

	path := ripper.GetTempPathFor(job, conf)
	ctx.Printf("cleaning data from \"%s\"", path)
	error := os.RemoveAll(path)
	if error != nil {
		ctx.Printf("  ...failed\n  due to: %s", error)
	} else {
		ctx.Printf("  ...done\n")
	}
	return []task.Job{job}, error
}
