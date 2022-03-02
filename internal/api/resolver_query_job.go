package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) JobQueue(ctx context.Context) ([]*models.Job, error) {
	queue := manager.GetInstance().JobManager.GetQueue()

	var ret []*models.Job
	for _, j := range queue {
		ret = append(ret, jobToJobModel(j))
	}

	return ret, nil
}

func (r *queryResolver) FindJob(ctx context.Context, input models.FindJobInput) (*models.Job, error) {
	jobID, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, err
	}
	j := manager.GetInstance().JobManager.GetJob(jobID)
	if j == nil {
		return nil, nil
	}

	return jobToJobModel(*j), nil
}

func jobToJobModel(j job.Job) *models.Job {
	ret := &models.Job{
		ID:          strconv.Itoa(j.ID),
		Status:      models.JobStatus(j.Status),
		Description: j.Description,
		SubTasks:    j.Details,
		StartTime:   j.StartTime,
		EndTime:     j.EndTime,
		AddTime:     j.AddTime,
	}

	if j.Progress != -1 {
		ret.Progress = &j.Progress
	}

	return ret
}
