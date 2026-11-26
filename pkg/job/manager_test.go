package job

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const sleepTime time.Duration = 10 * time.Millisecond

type testExec struct {
	started   chan struct{}
	finish    chan struct{}
	cancelled bool
	progress  *Progress
}

func newTestExec(finish chan struct{}) *testExec {
	return &testExec{
		started: make(chan struct{}),
		finish:  finish,
	}
}

func (e *testExec) Execute(ctx context.Context, p *Progress) error {
	e.progress = p
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

	return nil
}

func TestAdd(t *testing.T) {
	m := NewManager()

	const jobName = "test job"
	exec1 := newTestExec(make(chan struct{}))
	jobID := m.Add(context.Background(), jobName, exec1)

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
	job2ID := m.Add(context.Background(), otherJobName, exec2)

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

func TestStart(t *testing.T) {
	m := NewManager()

	// occupy the queue so a queued job could not run
	blocker := newTestExec(make(chan struct{}))
	m.Add(context.Background(), "blocker", blocker)
	<-blocker.started

	exec := newTestExec(make(chan struct{}))
	jobID := m.Start(context.Background(), "started job", exec)

	// expect the started job to run concurrently with the blocker
	select {
	case <-exec.started:
	case <-time.After(time.Second):
		t.Fatal("started job did not run concurrently")
	}

	close(exec.finish)
	time.Sleep(sleepTime)

	// expect the finished job to be moved from the queue to the graveyard
	m.mutex.Lock()
	_, inQueue := m.getJob(m.queue, jobID)
	_, inGraveyard := m.getJob(m.graveyard, jobID)
	m.mutex.Unlock()

	assert.Nil(t, inQueue)
	assert.NotNil(t, inGraveyard)
	assert.Equal(t, StatusFinished, m.GetJob(jobID).Status)

	close(blocker.finish)
}

// a job queued while a started job is still running must still be dispatched
func TestAddWhileStartedJobRunning(t *testing.T) {
	m := NewManager()

	started := newTestExec(make(chan struct{}))
	m.Start(context.Background(), "started job", started)
	<-started.started

	exec := newTestExec(nil)
	m.Add(context.Background(), "queued job", exec)

	select {
	case <-exec.started:
	case <-time.After(time.Second):
		t.Fatal("queued job was not dispatched while a started job was running")
	}

	close(started.finish)
}

func TestTaskQueueProgress(t *testing.T) {
	m := NewManager()
	j := &Job{}
	p := m.newProgress(j)

	tq := NewTaskQueue(context.Background(), p, 10, make(chan struct{}, 2))
	for range 3 {
		tq.Add("task", func(ctx context.Context) {})
	}
	tq.Close()
	tq.Wait()

	assert.Equal(t, 3, p.total)
	assert.Equal(t, 3, p.processed)
	assert.Equal(t, float64(1), p.percent)
}

// queues sharing slots must not run more than cap(slots) tasks between them
func TestTaskQueueSharedSlots(t *testing.T) {
	m := NewManager()
	slots := make(chan struct{}, 2)

	var mu sync.Mutex
	running, maxRunning := 0, 0
	task := func(ctx context.Context) {
		mu.Lock()
		running++
		maxRunning = max(maxRunning, running)
		mu.Unlock()

		time.Sleep(sleepTime)

		mu.Lock()
		running--
		mu.Unlock()
	}

	var queues []*TaskQueue
	for range 3 {
		tq := NewTaskQueue(context.Background(), m.newProgress(&Job{}), 10, slots)
		for range 4 {
			tq.Add("task", task)
		}
		tq.Close()
		queues = append(queues, tq)
	}

	for _, tq := range queues {
		tq.Wait()
	}

	assert.Equal(t, 2, maxRunning)
}

func TestCancel(t *testing.T) {
	m := NewManager()

	// add two jobs
	const jobName = "test job"
	exec1 := newTestExec(make(chan struct{}))
	jobID := m.Add(context.Background(), jobName, exec1)

	const otherJobName = "other job"
	exec2 := newTestExec(make(chan struct{}))
	job2ID := m.Add(context.Background(), otherJobName, exec2)

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

func TestCancelAll(t *testing.T) {
	m := NewManager()

	// add two jobs
	const jobName = "test job"
	exec1 := newTestExec(make(chan struct{}))
	jobID := m.Add(context.Background(), jobName, exec1)

	const otherJobName = "other job"
	exec2 := newTestExec(make(chan struct{}))
	job2ID := m.Add(context.Background(), otherJobName, exec2)

	// wait a tiny bit
	time.Sleep(sleepTime)

	m.CancelAll()

	// allow first job to finish
	close(exec1.finish)

	// wait a tiny bit
	time.Sleep(sleepTime)

	// expect all jobs to be cancelled
	assert := assert.New(t)
	j := m.GetJob(job2ID)
	assert.Equal(StatusCancelled, j.Status)

	j = m.GetJob(jobID)
	assert.Equal(StatusCancelled, j.Status)

	// expect all jobs to be removed from the queue
	assert.Len(m.GetQueue(), 0)

	// expect job to have not have been started
	select {
	case <-exec2.started:
		t.Error("cancelled exec was started")
	default:
	}
}

// CancelAll must remove every ready job, not every other one
func TestCancelAllReadyJobs(t *testing.T) {
	m := NewManager()

	running := newTestExec(make(chan struct{}))
	m.Add(context.Background(), "running job", running)
	<-running.started

	for range 4 {
		m.Add(context.Background(), "ready job", newTestExec(nil))
	}

	m.CancelAll()

	// only the running job (now stopping) remains
	queue := m.GetQueue()
	if assert.Len(t, queue, 1) {
		assert.Equal(t, StatusStopping, queue[0].Status)
	}

	close(running.finish)
}

func TestSubscribe(t *testing.T) {
	m := NewManager()

	m.updateThrottleLimit = time.Millisecond * 100

	ctx, cancel := context.WithCancel(context.Background())

	s := m.Subscribe(ctx)

	// add a job
	const jobName = "test job"
	exec1 := newTestExec(make(chan struct{}))
	jobID := m.Add(context.Background(), jobName, exec1)

	assert := assert.New(t)

	select {
	case newJob := <-s.NewJob:
		assert.Equal(jobID, newJob.ID)
		assert.Equal(jobName, newJob.Description)
		assert.Equal(StatusReady, newJob.Status)
	case <-time.After(time.Second):
		t.Error("new job was not received")
	}

	// should receive an update when the job begins to run
	select {
	case updatedJob := <-s.UpdatedJob:
		assert.Equal(jobID, updatedJob.ID)
		assert.Equal(jobName, updatedJob.Description)
		assert.Equal(StatusRunning, updatedJob.Status)
	case <-time.After(time.Second):
		t.Error("updated job was not received")
	}

	// wait for it to start
	select {
	case <-exec1.started:
		// ok
	case <-time.After(time.Second):
		t.Error("exec was not started")
	}

	// test update throttling
	exec1.progress.SetPercent(0.1)

	// first update should be immediate
	select {
	case updatedJob := <-s.UpdatedJob:
		assert.Equal(0.1, updatedJob.Progress)
	case <-time.After(m.updateThrottleLimit):
		t.Error("updated job was not received")
	}

	exec1.progress.SetPercent(0.2)
	exec1.progress.SetPercent(0.3)

	// should only receive a single update with the second status
	select {
	case updatedJob := <-s.UpdatedJob:
		assert.Equal(0.3, updatedJob.Progress)
	case <-time.After(time.Second):
		t.Error("updated job was not received")
	}

	select {
	case <-s.UpdatedJob:
		t.Error("received an additional updatedJob")
	default:
	}

	// allow job to finish
	close(exec1.finish)

	select {
	case removedJob := <-s.RemovedJob:
		assert.Equal(jobID, removedJob.ID)
		assert.Equal(jobName, removedJob.Description)
		assert.Equal(StatusFinished, removedJob.Status)
	case <-time.After(time.Second):
		t.Error("removed job was not received")
	}

	// should not receive another update
	select {
	case <-s.UpdatedJob:
		t.Error("updated job was received after update")
	case <-time.After(m.updateThrottleLimit):
	}

	// add another job and cancel it
	exec2 := newTestExec(make(chan struct{}))
	jobID = m.Add(context.Background(), jobName, exec2)

	m.CancelJob(jobID)

	select {
	case removedJob := <-s.RemovedJob:
		assert.Equal(jobID, removedJob.ID)
		assert.Equal(jobName, removedJob.Description)
		assert.Equal(StatusCancelled, removedJob.Status)
	case <-time.After(time.Second):
		t.Error("cancelled job was not received")
	}

	cancel()
}
