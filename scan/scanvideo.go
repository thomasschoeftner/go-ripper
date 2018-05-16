package scan

import (
	"go-cli/task"
	"go-ripper/ripper"
	"go-ripper/files"
	"go-ripper/targetinfo"
	"path/filepath"
)

func ScanVideo(ctx task.Context, job task.Job) ([]task.Job, error) {
	conf := ctx.Config.(*ripper.AppConf)

	path := job[ripper.JobField_Path]
	tmp := conf.TempDirectoryName
	out := conf.OutputDirectoryName

	ctx.Printf("  scanning contents of \"%s\" (ignoring temp \"%s\" and output \"%s\")\n", path, tmp, out)
	targets, err := scan(path, tmp, out, "video", conf.Scan.Video)
	if err != nil {
		return nil, err
	}

	jobs := []task.Job{}
	ctx.Printf("    found %d targets:\n", len(targets))
	for _, target := range targets {
		//TODO write target to tmp folder
		tmpFolder := filepath.Join(target.Folder, tmp)
		err = files.CheckOrCreateFolder(tmpFolder)
		if err != nil {
			return nil, err
		}
		targetinfo.Save(tmpFolder, target)
		newJob := job.WithParam(ripper.JobField_Path, filepath.Join(target.Folder, target.File))
		jobs = append(jobs, newJob)
		ctx.Printf("    %s\n", target)
	}
	return jobs, nil
}
