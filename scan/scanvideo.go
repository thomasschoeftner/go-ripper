package scan

import (
	"go-cli/task"
	"go-ripper/ripper"
	"go-ripper/targetinfo"
	"path/filepath"
	"errors"
	"fmt"
)

func ScanVideo(ctx task.Context) task.HandlerFunc {
	conf := ctx.Config.(*ripper.AppConf)

	return func (job task.Job) ([]task.Job, error) {
		scanPath := job[ripper.JobField_Path]

		ctx.Printf("scanning contents of \"%s\"", scanPath)
		scanResults, err := scan(scanPath, conf.IgnorePrefix, conf.Scan.Video)
		if err != nil {
			return nil, err
		}

		//convert scanResults to TargetInfos
		targets, err := toTargetInfos(scanResults)
		if err != nil {
			return nil, err
		}

		jobs := []task.Job{}
		ctx.Printf("found %d targets:\n", len(targets))
		for _, target := range targets {
			//write TargetInfo to work folder
			workDir, err := ripper.GetWorkPathForTargetFileFolder(conf.WorkDirectory, target.GetFolder())
			if err != nil {
				return nil, err
			}
			
			err = targetinfo.Save(workDir, target)
			if err != nil {
				//TODO check if error should be ignored - in a worst case the target file will be missing
				return nil, err
			}

			//create new Job
			newJob := job.WithParam(ripper.JobField_Path, filepath.Join(target.GetFolder(), target.GetFile()))
			jobs = append(jobs, newJob)
			ctx.Printf("  %s\n", target)
		}
		return jobs, nil
	}
}


func toTargetInfos (results []*scanResult) ([]targetinfo.TargetInfo, error) {
	var targetInfos []targetinfo.TargetInfo
	episodeCount := map[string]map[int]int{}

	for _, r := range results {
		if r.Collection != nil {
			series := r.Id
			path := filepath.Join(r.Folder, r.File)
			season := *r.Collection
			if r.ItemNo == nil {
				return nil, errors.New(fmt.Sprintf("invalid episode found - season# is set (%d), but episode# is missing in file %s", season, path))
			}
			episode := *r.ItemNo

			seasons := episodeCount[series]
			if seasons == nil {
				seasons = map[int]int {}
				episodeCount[series] = seasons
			}
			cnt := seasons[season] + 1
			seasons[season] = cnt

			episodeInfo := targetinfo.NewEpisode(r.File, r.Folder, r.Id, season, episode, cnt, 0)
			targetInfos = append(targetInfos, episodeInfo)
		} else { //single video
			targetInfos = append(targetInfos, targetinfo.NewMovie(r.File, r.Folder, r.Id))
		}
	}

	//finally update all total # of episodes for all episodes
	for _, ti := range targetInfos {
		if targetinfo.TARGETINFO_TPYE_EPISODE == ti.GetType() {
			e := ti.(*targetinfo.Episode)
			e.ItemsTotal = episodeCount[e.Id][e.Season]
		}
	}

	return targetInfos, nil
}
