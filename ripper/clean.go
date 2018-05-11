package ripper

import (
	"os"
	"go-cli/task"
)

func cleanHandler(ctx task.Context, job task.Job) ([]task.Job, error) {
	conf := ctx.Config.(AppConf)
	path := GetTempPathFor(job, conf)
	ctx.Printf("cleaning data from \"%s\"", path)
	error := os.RemoveAll(path)
	if error != nil {
		ctx.Printf("  ...failed\n  due to: %s", error)
	} else {
		ctx.Printf("  ...done\n")
	}
	return []task.Job{job}, error
}
