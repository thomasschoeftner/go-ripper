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
	excludeDirs := []string { conf.TempDirectoryName, conf.OutputDirectoryName}

	ctx.Printf("  scanning contents of \"%s\" (ignoring temp \"%s\" and output \"%s\")\n", path, conf.TempDirectoryName, conf.OutputDirectoryName)
	targets, err := scan(path, excludeDirs, "video", conf.Scan.Video)
	if err != nil {
		return nil, err
	}

	jobs := []task.Job{}
	ctx.Printf("    found %d targets:\n", len(targets))
	for _, target := range targets {
		//write TargetInfo to tmp folder
		tmpFolder := filepath.Join(target.Folder, conf.TempDirectoryName)
		err = files.CheckOrCreateFolder(tmpFolder)
		if err != nil {
			return nil, err
		}
		targetinfo.Save(tmpFolder, target)

		//create new Job
		newJob := job.WithParam(ripper.JobField_Path, filepath.Join(target.Folder, target.File))
		jobs = append(jobs, newJob)
		ctx.Printf("    %s\n", target)
	}
	return jobs, nil
}
