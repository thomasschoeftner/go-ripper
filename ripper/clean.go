package ripper

import (
	"os"
	"go-cli/task"
)

func cleanHandler(ctx task.Context, c *task.Command) []task.Result {
	conf := ctx.Config.(AppConf)
	ctx.Printf("cleaning intermediate data from working folder \"%s\"", conf.WorkDirectory)
	error := os.RemoveAll(conf.WorkDirectory)
	if error != nil {
		ctx.Printf("  ...failed\n")
	} else {
		ctx.Printf("  ...done\n")
	}
	return []task.Result{ {c, error} }
}

