package metainfo

import (
	"go-cli/task"
	"go-ripper/ripper"
	"go-ripper/targetinfo"
)

func ResolveVideo(queryFactory VideoMetaInfoQueryFactory) task.Handler {
	//todo handle queryFactory == nil??

	return func (ctx task.Context) task.HandlerFunc {
		//conf := ctx.Config.(*ripper.AppConf)

		return func(job task.Job) ([]task.Job, error) {
			targetInfofile := job[ripper.JobField_Path]
			//tmpPath := ripper.GetTempPathFor(job, conf)


			ctx.Printf("process video - targetinfo %s\n", targetInfofile)

			printf := ctx.Printf.WithIndent(2)
			ti, err := targetinfo.Read(targetInfofile)
			if err != nil {
				return nil, err
			}
			printf("recovered target-info: %s\n", ti.String())

			//todo check if already handled (lazy)
			//todo issue proper queries
			metaInfo, err := Get(queryFactory.NewTitleQuery(ti.GetId()))
			println(metaInfo)

			//TODO persist meta-info in file
			return nil, nil
		}

	}
}
