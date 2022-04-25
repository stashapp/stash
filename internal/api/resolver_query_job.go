package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/job"
)

func (r *queryResolver) JobQueue(ctx context.Context) ([]*Job, error) {
	queue := manager.GetInstance().JobManager.GetQueue()

	var ret []*Job
	for _, j := range queue {
		ret = append(ret, jobToJobModel(j))
	}

	return ret, nil
}

func (r *queryResolver) FindJob(ctx context.Context, input FindJobInput) (*Job, error) {
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

func jobToJobModel(j job.Job) *Job {
	ret := &Job{
		ID:          strconv.Itoa(j.ID),
		Status:      JobStatus(j.Status),
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
