package ripper

import "go-cli/task"

func ErrorHandler(err error) task.HandlerFunc {
	return func(job task.Job) ([]task.Job, error) {
		return nil, err
	}
}

