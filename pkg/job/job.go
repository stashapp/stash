package job

import (
	"context"
	"time"
)

// JobExec represents the implementation of a Job to be executed.
type JobExec interface {
	Execute(ctx context.Context, updater StatusUpdater)
}

// StatusUpdater is used by JobExec objects to communicate their progress
// to the job manager.
type StatusUpdater interface {
	SetProgress(progress float64)
	AddSubTask(subtask string) int
	RemoveSubTask(subtaskID int)
}

// Status is the status of a Job
type Status string

const (
	// StatusReady means that the Job is not yet started.
	StatusReady   Status = "READY"
	StatusRunning Status = "RUNNING"
	// StatusStopping means that the job is cancelled but is still running.
	StatusStopping Status = "STOPPING"
	// StatusFinished means that the job was completed.
	StatusFinished Status = "FINISHED"
	// StatusCancelled means that the job was cancelled and is now stopped.
	StatusCancelled Status = "CANCELLED"
)

// Job represents the status of a queued or running job.
type Job struct {
	ID          int
	Status      Status
	SubTasks    []SubTask
	Description string
	Progress    float64
	StartTime   *time.Time
	EndTime     *time.Time
	AddTime     time.Time

	exec          JobExec
	cancelFunc    context.CancelFunc
	lastSubtaskID int
}

func (j *Job) nextID() int {
	j.lastSubtaskID += 1
	return j.lastSubtaskID
}

func (j *Job) addSubTask(subtask string) int {
	s := SubTask{
		id:          j.nextID(),
		Description: subtask,
	}

	j.SubTasks = append(j.SubTasks, s)
	return s.id
}

func (j *Job) removeSubTask(id int) {
	for i, subtask := range j.SubTasks {
		if subtask.id == id {
			j.SubTasks = append(j.SubTasks[:i], j.SubTasks[i+1:]...)
			return
		}
	}
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

// SubTask represents a unit of work within a Job.
type SubTask struct {
	id          int
	Description string
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
