package scan

import (
	"go-cli/task"
	"go-ripper/ripper"
	"go-ripper/files"
	"go-ripper/targetinfo"
	"path/filepath"
)

func ScanVideo(ctx task.Context) task.HandlerFunc {
	conf := ctx.Config.(*ripper.AppConf)
	excludeDirs := []string { conf.TempDirectoryName, conf.OutputDirectoryName}

	return func (job task.Job) ([]task.Job, error) {
		scanPath := job[ripper.JobField_Path]

		ctx.Printf("scanning contents of \"%s\" (ignoring temp \"%s\" and output \"%s\")\n", scanPath, conf.TempDirectoryName, conf.OutputDirectoryName)
		targets, err := scan(scanPath, excludeDirs, "video", conf.Scan.Video)
		if err != nil {
			return nil, err
		}

		jobs := []task.Job{}
		ctx.Printf("found %d targets:\n", len(targets))
		for _, target := range targets {
			//write TargetInfo to tmp folder
			tmpFolder := filepath.Join(target.Folder, conf.TempDirectoryName)
			err = files.CheckOrCreateFolder(tmpFolder)
			if err != nil {
				return nil, err
			}

			fileName, err := targetinfo.Save(tmpFolder, target)
			if err != nil {
				//TODO check if error should be ignored
				return nil, err
			}

			//TODO fix sequenceNo (in case they start with 0)

			//create new Job
			newJob := job.WithParam(ripper.JobField_Path, *fileName)
			jobs = append(jobs, newJob)
			ctx.Printf("  %s\n", target)
		}
		return jobs, nil
	}
}
