package scan

import (
	"go-cli/task"
	"go-ripper/ripper"
)

func ScanVideo(ctx task.Context, job task.Job) ([]task.Job, error) {
	conf := ctx.Config.(*ripper.AppConf)

	path := job[ripper.JobField_Location]
	tmp := conf.TempDirectoryName
	out := conf.OutputDirectoryName

	ctx.Printf("  scanning contents of \"%s\"   (ignoring temp \"%s\" and output \"%s\")\n", path, tmp, out)
	targets, err := scan(path, tmp, out, conf.Scan.Video)
	if err != nil {
		return nil, err
	}

	jobs := []task.Job{}
	for _, target := range targets {
		//TODO add other fields from target to job!!!
		jobs = append(jobs, job.WithParam(ripper.JobField_Location, target.File)) //copy existing job with updated params
	}
	return jobs, nil
}
