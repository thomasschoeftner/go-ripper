package ripper

import (
	"os"
	"go-cli/task"
)

func cleanHandler(ctx tasks.Context, c *tasks.Command) []tasks.Result {
	conf := ctx.Config.(AppConf)
	ctx.Printf("cleaning intermediate data from working folder \"%s\"", conf.WorkDirectory)
	error := os.RemoveAll(conf.WorkDirectory)
	if error != nil {
		ctx.Printf("  ...failed\n")
	} else {
		ctx.Printf("  ...done\n")
	}
	return []tasks.Result{ {c, error} }
}

