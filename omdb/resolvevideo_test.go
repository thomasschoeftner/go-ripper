package omdb

import (
	"testing"
	"go-cli/task"
	"go-cli/commons"
	"go-ripper/ripper"
)

func TestOmdbFactoryCreation(t *testing.T) {
	ctx := task.Context{task.TaskSequence{}, nil, commons.DevNullPrintf}
	job := task.Job{ripper.JobField_Path : "."}

	if tFactory != nil {
		t.Errorf("token factory was expected to be nil initially, but was %+v", *tFactory)
	}

	t.Run("use invalid token factory", func(t *testing.T) {
		createHandler := ResolveVideo([]string{})
		if tFactory != nil {
			t.Errorf("token factory was expected to be nil after invalid setup, but was %+v", *tFactory)
		}
		handler := createHandler(ctx)
		_, err := handler(job)
		if err == nil {
			t.Errorf("expected handler error, but got no error")
		}
	})

	if tFactory != nil {
		t.Errorf("token factory was expected to be nil initially, but was %+v", *tFactory)
	}

	t.Run("use valid token factory", func(t *testing.T) {
		tokens := []string {"aaa", "bbb", "ccc"}
		ResolveVideo(tokens)
		if tFactory == nil {
			t.Error("token factory was expected not to be nil after valid setup, but was nil")
		}
	})
}

//TODO add tests for omdb access
