package job

import (
	"context"
	"runtime/debug"
	"sync"
	"time"

	"github.com/stashapp/stash/pkg/logger"
)

const maxGraveyardSize = 10
const defaultThrottleLimit = 100 * time.Millisecond

// Manager maintains a queue of jobs. Jobs are executed one at a time.
type Manager struct {
	queue     []*Job
	graveyard []*Job

	mutex    sync.Mutex
	notEmpty *sync.Cond
	stop     chan struct{}

	lastID int

	subscriptions       []*ManagerSubscription
	updateThrottleLimit time.Duration
}

// NewManager initialises and returns a new Manager.
func NewManager() *Manager {
	ret := &Manager{
		stop:                make(chan struct{}),
		updateThrottleLimit: defaultThrottleLimit,
	}

	ret.notEmpty = sync.NewCond(&ret.mutex)

	go ret.dispatcher()

	return ret
}

// Stop is used to stop the dispatcher thread. Once Stop is called, no
// more Jobs will be processed.
func (m *Manager) Stop() {
	m.CancelAll()
	close(m.stop)
}

// Add queues a job.
func (m *Manager) Add(ctx context.Context, description string, e JobExec) int {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	t := time.Now()

	j := Job{
		ID:          m.nextID(),
		Status:      StatusReady,
		Description: description,
		AddTime:     t,
		exec:        e,
		outerCtx:    ctx,
	}

	m.queue = append(m.queue, &j)

	if len(m.queue) == 1 {
		// notify that there is now a job in the queue
		m.notEmpty.Broadcast()
	}

	m.notifyNewJob(&j)

	return j.ID
}

// Start adds a job and starts it immediately, concurrently with any other
// jobs.
func (m *Manager) Start(ctx context.Context, description string, e JobExec) int {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	t := time.Now()

	j := Job{
		ID:          m.nextID(),
		Status:      StatusReady,
		Description: description,
		AddTime:     t,
		exec:        e,
		outerCtx:    ctx,
	}

	m.queue = append(m.queue, &j)

	m.dispatch(ctx, &j)

	return j.ID
}

func (m *Manager) notifyNewJob(j *Job) {
	// assumes lock held
	for _, s := range m.subscriptions {
		// don't block if channel is full
		select {
		case s.newJob <- *j:
		default:
		}
	}
}

func (m *Manager) nextID() int {
	m.lastID++
	return m.lastID
}

func (m *Manager) getReadyJob() *Job {
	// assumes lock held
	for _, j := range m.queue {
		if j.Status == StatusReady {
			return j
		}
	}

	return nil
}

func (m *Manager) dispatcher() {
	m.mutex.Lock()

	for {
		// wait until we have something to process
		j := m.getReadyJob()

		for j == nil {
			m.notEmpty.Wait()

			// it's possible that we have been stopped - check here
			select {
			case <-m.stop:
				m.mutex.Unlock()
				return
			default:
				// keep going
				j = m.getReadyJob()
			}
		}

		done := m.dispatch(j.outerCtx, j)

		// unlock the mutex and wait for the job to finish
		m.mutex.Unlock()
		<-done
		m.mutex.Lock()

		// remove the job from the queue
		m.removeJob(j)

		// process next job
	}
}

func (m *Manager) newProgress(j *Job) *Progress {
	return &Progress{
		updater: &updater{
			m:   m,
			job: j,
		},
		percent: ProgressIndefinite,
	}
}

func (m *Manager) dispatch(ctx context.Context, j *Job) (done chan struct{}) {
	// assumes lock held
	t := time.Now()
	j.StartTime = &t
	j.Status = StatusRunning

	// create a cancellable context for the job that is not canceled by the outer context
	ctx, cancelFunc := context.WithCancel(context.WithoutCancel(ctx))
	j.cancelFunc = cancelFunc

	done = make(chan struct{})
	go m.executeJob(ctx, j, done)

	m.notifyJobUpdate(j)

	return
}

func (m *Manager) executeJob(ctx context.Context, j *Job, done chan struct{}) {
	defer close(done)
	defer m.onJobFinish(j)
	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, log and mark the job as failed
			logger.Errorf("panic in job %d - %s: %v", j.ID, j.Description, p)
			logger.Error(string(debug.Stack()))

			m.mutex.Lock()
			defer m.mutex.Unlock()
			j.Status = StatusFailed
		}
	}()

	progress := m.newProgress(j)
	if err := j.exec.Execute(ctx, progress); err != nil {
		logger.Errorf("task failed due to error: %v", err)
		j.error(err)
	}
}

func (m *Manager) onJobFinish(job *Job) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if job.Status == StatusStopping {
		job.Status = StatusCancelled
	} else if job.Status != StatusFailed {
		job.Status = StatusFinished
	}
	t := time.Now()
	job.EndTime = &t
}

func (m *Manager) removeJob(job *Job) {
	// assumes lock held
	index, _ := m.getJob(m.queue, job.ID)
	if index == -1 {
		return
	}

	// clear any subtasks
	job.Details = nil

	m.queue = append(m.queue[:index], m.queue[index+1:]...)

	m.graveyard = append(m.graveyard, job)
	if len(m.graveyard) > maxGraveyardSize {
		m.graveyard = m.graveyard[1:]
	}

	// notify job removed
	for _, s := range m.subscriptions {
		// don't block if channel is full
		select {
		case s.removedJob <- *job:
		default:
		}
	}
}

func (m *Manager) getJob(list []*Job, id int) (index int, job *Job) {
	// assumes lock held
	for i, j := range list {
		if j.ID == id {
			index = i
			job = j
			return
		}
	}

	return -1, nil
}

// CancelJob cancels the job with the provided id. Jobs that have been started
// are notified that they are stopping. Jobs that have not yet started are
// removed from the queue. If no job exists with the provided id, then there is
// no effect. Likewise, if the job is already cancelled, there is no effect.
func (m *Manager) CancelJob(id int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, j := m.getJob(m.queue, id)
	if j != nil {
		j.cancel()

		if j.Status == StatusCancelled {
			// remove from the queue
			m.removeJob(j)
		}
	}
}

// CancelAll cancels all of the jobs in the queue. This is the same as
// calling CancelJob on all jobs in the queue.
func (m *Manager) CancelAll() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// call cancel on all
	for _, j := range m.queue {
		j.cancel()

		if j.Status == StatusCancelled {
			// add to graveyard
			m.removeJob(j)
		}
	}
}

// GetJob returns a copy of the Job for the provided id. Returns nil if the job
// does not exist.
func (m *Manager) GetJob(id int) *Job {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// get from the queue or graveyard
	_, j := m.getJob(append(m.queue, m.graveyard...), id)
	if j != nil {
		// make a copy of the job and return the pointer
		jCopy := *j
		return &jCopy
	}

	return nil
}

// GetQueue returns a copy of the current job queue.
func (m *Manager) GetQueue() []Job {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var ret []Job

	for _, j := range m.queue {
		jCopy := *j
		ret = append(ret, jCopy)
	}

	return ret
}

// Subscribe subscribes to changes to jobs in the manager queue.
func (m *Manager) Subscribe(ctx context.Context) *ManagerSubscription {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	ret := newSubscription()

	m.subscriptions = append(m.subscriptions, ret)

	go func() {
		<-ctx.Done()
		m.mutex.Lock()
		defer m.mutex.Unlock()

		ret.close()

		// remove from the list
		for i, s := range m.subscriptions {
			if s == ret {
				m.subscriptions = append(m.subscriptions[:i], m.subscriptions[i+1:]...)
				break
			}
		}
	}()

	return ret
}

func (m *Manager) notifyJobUpdate(j *Job) {
	// don't update if job is finished or cancelled - these are handled
	// by removeJob
	if j.Status == StatusCancelled || j.Status == StatusFinished {
		return
	}

	// assumes lock held
	for _, s := range m.subscriptions {
		// don't block if channel is full
		select {
		case s.updatedJob <- *j:
		default:
		}
	}
}

type updater struct {
	m           *Manager
	job         *Job
	lastUpdate  time.Time
	updateTimer *time.Timer
}

func (u *updater) notifyUpdate() {
	// assumes lock held
	u.m.notifyJobUpdate(u.job)
	u.lastUpdate = time.Now()
	u.updateTimer = nil
}

func (u *updater) updateProgress(progress float64, details []string) {
	u.m.mutex.Lock()
	defer u.m.mutex.Unlock()

	u.job.Progress = progress
	u.job.Details = details

	if time.Since(u.lastUpdate) < u.m.updateThrottleLimit {
		if u.updateTimer == nil {
			u.updateTimer = time.AfterFunc(u.m.updateThrottleLimit-time.Since(u.lastUpdate), func() {
				u.m.mutex.Lock()
				defer u.m.mutex.Unlock()

				u.notifyUpdate()
			})
		}
	} else {
		u.notifyUpdate()
	}
}
