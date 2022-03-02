package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/internal/manager"
)

func (r *mutationResolver) StopJob(ctx context.Context, jobID string) (bool, error) {
	idInt, err := strconv.Atoi(jobID)
	if err != nil {
		return false, err
	}
	manager.GetInstance().JobManager.CancelJob(idInt)

	return true, nil
}

func (r *mutationResolver) StopAllJobs(ctx context.Context) (bool, error) {
	manager.GetInstance().JobManager.CancelAll()
	return true, nil
}
