package api

import (
	"context"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/models"
)

func makeJobStatusUpdate(t models.JobStatusUpdateType, j job.Job) *models.JobStatusUpdate {
	return &models.JobStatusUpdate{
		Type: t,
		Job:  jobToJobModel(j),
	}
}

func (r *subscriptionResolver) JobsSubscribe(ctx context.Context) (<-chan *models.JobStatusUpdate, error) {
	msg := make(chan *models.JobStatusUpdate, 100)

	subscription := manager.GetInstance().JobManager.Subscribe(ctx)

	go func() {
		for {
			select {
			case j := <-subscription.NewJob:
				msg <- makeJobStatusUpdate(models.JobStatusUpdateTypeAdd, j)
			case j := <-subscription.RemovedJob:
				msg <- makeJobStatusUpdate(models.JobStatusUpdateTypeRemove, j)
			case j := <-subscription.UpdatedJob:
				msg <- makeJobStatusUpdate(models.JobStatusUpdateTypeUpdate, j)
			case <-ctx.Done():
				close(msg)
				return
			}
		}
	}()

	return msg, nil
}

func (r *subscriptionResolver) ScanCompleteSubscribe(ctx context.Context) (<-chan bool, error) {
	return manager.GetInstance().ScanSubscribe(ctx), nil
}
