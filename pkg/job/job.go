// Package job provides the job execution and management functionality for the application.
package job

import (
	"context"
	"time"
)

type JobExecFn func(ctx context.Context, progress *Progress) error

// JobExec represents the implementation of a Job to be executed.
type JobExec interface {
	Execute(ctx context.Context, progress *Progress) error
}

type jobExecImpl struct {
	fn JobExecFn
}

func (j *jobExecImpl) Execute(ctx context.Context, progress *Progress) error {
	return j.fn(ctx, progress)
}

// MakeJobExec returns a simple JobExec implementation using the provided
// function.
func MakeJobExec(fn JobExecFn) JobExec {
	return &jobExecImpl{
		fn: fn,
	}
}

// Status is the status of a Job
type Status string

const (
	// StatusReady means that the Job is not yet started.
	StatusReady Status = "READY"
	// StatusRunning means that the job is currently running.
	StatusRunning Status = "RUNNING"
	// StatusStopping means that the job is cancelled but is still running.
	StatusStopping Status = "STOPPING"
	// StatusFinished means that the job was completed.
	StatusFinished Status = "FINISHED"
	// StatusCancelled means that the job was cancelled and is now stopped.
	StatusCancelled Status = "CANCELLED"
	// StatusFailed means that the job failed.
	StatusFailed Status = "FAILED"
)

// Job represents the status of a queued or running job.
type Job struct {
	ID     int
	Status Status
	// details of the current operations of the job
	Details     []string
	Description string
	// Progress in terms of 0 - 1.
	Progress  float64
	StartTime *time.Time
	EndTime   *time.Time
	AddTime   time.Time
	Error     *string

	outerCtx   context.Context
	exec       JobExec
	cancelFunc context.CancelFunc
}

// TimeElapsed returns the total time elapsed for the job.
// If the EndTime is set, then it uses this to calculate the elapsed time, otherwise it uses time.Now.
func (j *Job) TimeElapsed() time.Duration {
	var end time.Time
	if j.EndTime != nil {
		end = time.Now()
	} else {
		end = *j.EndTime
	}

	return end.Sub(*j.StartTime)
}

func (j *Job) cancel() {
	if j.Status == StatusReady {
		j.Status = StatusCancelled
	} else if j.Status == StatusRunning {
		j.Status = StatusStopping
	}

	if j.cancelFunc != nil {
		j.cancelFunc()
	}
}

func (j *Job) error(err error) {
	errStr := err.Error()
	j.Error = &errStr
	j.Status = StatusFailed
}

// IsCancelled returns true if cancel has been called on the context.
func IsCancelled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
