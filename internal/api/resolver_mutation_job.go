package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/internal/manager"
)

func (r *mutationResolver) StopJob(ctx context.Context, jobID string) (bool, error) {
	id, err := strconv.Atoi(jobID)
	if err != nil {
		return false, fmt.Errorf("converting id: %w", err)
	}
	manager.GetInstance().JobManager.CancelJob(id)

	return true, nil
}

func (r *mutationResolver) StopAllJobs(ctx context.Context) (bool, error) {
	manager.GetInstance().JobManager.CancelAll()
	return true, nil
}
