package job

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const sleepTime time.Duration = 1 * time.Millisecond

type testExec struct {
	started   chan struct{}
	finish    chan struct{}
	cancelled bool
}

func newTestExec(finish chan struct{}) *testExec {
	return &testExec{
		started: make(chan struct{}),
		finish:  finish,
	}
}

func (e *testExec) Execute(ctx context.Context, updater StatusUpdater) {
	close(e.started)

	if e.finish != nil {
		<-e.finish

		select {
		case <-ctx.Done():
			e.cancelled = true
		default:
			// fall through
		}
	}
}

func TestAdd(t *testing.T) {
	m := NewManager()

	const jobName = "test job"
	exec1 := newTestExec(make(chan struct{}))
	jobID := m.Add(jobName, exec1)

	// expect jobID to be the first ID
	assert := assert.New(t)
	assert.Equal(1, jobID)

	// wait a tiny bit
	time.Sleep(sleepTime)

	// expect job to have started
	select {
	case <-exec1.started:
		// ok
	default:
		t.Error("exec was not started")
	}

	// expect status to be running
	j := m.GetJob(jobID)

	assert.Equal(StatusRunning, j.Status)

	// expect description to be set
	assert.Equal(jobName, j.Description)

	// expect startTime and addTime to be set
	assert.NotNil(j.StartTime)
	assert.NotNil(j.AddTime)

	// expect endTime to not be set
	assert.Nil(j.EndTime)

	// add another job to the queue
	const otherJobName = "other job name"
	exec2 := newTestExec(make(chan struct{}))
	job2ID := m.Add(otherJobName, exec2)

	// expect status to be ready
	j2 := m.GetJob(job2ID)

	assert.Equal(StatusReady, j2.Status)

	// expect addTime to be set
	assert.NotNil(j2.AddTime)

	// expect startTime and endTime to not be set
	assert.Nil(j2.StartTime)
	assert.Nil(j2.EndTime)

	// allow first job to finish
	close(exec1.finish)

	// wait a tiny bit
	time.Sleep(sleepTime)

	// expect first job to be finished
	j = m.GetJob(jobID)
	assert.Equal(StatusFinished, j.Status)

	// expect end time to be set
	assert.NotNil(j.EndTime)

	// expect second job to have started
	select {
	case <-exec2.started:
		// ok
	default:
		t.Error("exec was not started")
	}

	// expect status to be running
	j2 = m.GetJob(job2ID)

	assert.Equal(StatusRunning, j2.Status)

	// expect startTime to be set
	assert.NotNil(j2.StartTime)
}

func TestCancel(t *testing.T) {
	m := NewManager()

	// add two jobs
	const jobName = "test job"
	exec1 := newTestExec(make(chan struct{}))
	jobID := m.Add(jobName, exec1)

	const otherJobName = "other job"
	exec2 := newTestExec(make(chan struct{}))
	job2ID := m.Add(otherJobName, exec2)

	// wait a tiny bit
	time.Sleep(sleepTime)

	m.CancelJob(job2ID)

	// expect job to be cancelled
	assert := assert.New(t)
	j := m.GetJob(job2ID)
	assert.Equal(StatusCancelled, j.Status)

	// expect end time not to be set
	assert.Nil(j.EndTime)

	// expect job to be removed from the queue
	assert.Len(m.GetQueue(), 1)

	// expect job to have not have been started
	select {
	case <-exec2.started:
		t.Error("cancelled exec was started")
	default:
	}

	// cancel running job
	m.CancelJob(jobID)

	// wait a tiny bit
	time.Sleep(sleepTime)

	// expect status to be stopping
	j = m.GetJob(jobID)
	assert.Equal(StatusStopping, j.Status)

	// expect job to still be in the queue
	assert.Len(m.GetQueue(), 1)

	// allow first job to finish
	close(exec1.finish)

	// wait a tiny bit
	time.Sleep(sleepTime)

	// expect job to be removed from the queue
	assert.Len(m.GetQueue(), 0)

	// expect job to be cancelled
	j = m.GetJob(jobID)
	assert.Equal(StatusCancelled, j.Status)

	// expect endtime to be set
	assert.NotNil(j.EndTime)

	// expect job to have been cancelled via context
	assert.True(exec1.cancelled)
}
