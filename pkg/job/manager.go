package job

import (
	"sync"
	"time"

	"golang.org/x/net/context"
)

const maxGraveyardSize = 10

// Manager maintains a queue of jobs. Jobs are executed one at a time.
type Manager struct {
	queue     []*Job
	graveyard []*Job

	mutex    sync.Mutex
	notEmpty *sync.Cond
	stop     chan struct{}

	lastID int

	subscriptions []ManagerSubscription
}

// NewManager initialises and returns a new Manager.
func NewManager() *Manager {
	ret := &Manager{}

	ret.notEmpty = sync.NewCond(&ret.mutex)
	ret.stop = make(chan struct{})

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
func (m *Manager) Add(description string, e JobExec) int {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	t := time.Now()

	j := Job{
		ID:          m.nextID(),
		Status:      StatusReady,
		Description: description,
		AddTime:     t,
		exec:        e,
	}

	m.queue = append(m.queue, &j)

	if len(m.queue) == 1 {
		// notify that there is now a job in the queue
		m.notEmpty.Broadcast()
	}

	m.notifyNewJob(&j)

	return j.ID
}

func (m *Manager) Start(description string, e JobExec) int {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	t := time.Now()

	j := Job{
		ID:          m.nextID(),
		Status:      StatusReady,
		Description: description,
		AddTime:     t,
		exec:        e,
	}

	m.queue = append(m.queue, &j)

	m.dispatch(&j)

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
	m.lastID += 1
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

		done := m.dispatch(j)

		// unlock the mutex and wait for the job to finish
		m.mutex.Unlock()
		<-done
		m.mutex.Lock()

		// remove the job from the queue
		m.removeJob(j)

		// process next job
	}
}

func (m *Manager) dispatch(j *Job) (done chan struct{}) {
	// assumes lock held
	t := time.Now()
	j.StartTime = &t
	j.Status = StatusRunning

	ctx, cancelFunc := context.WithCancel(context.Background())
	j.cancelFunc = cancelFunc

	done = make(chan struct{})
	go func() {
		progress := &Progress{
			updater: &updater{
				m:   m,
				job: j,
			},
		}
		j.exec.Execute(ctx, progress)

		m.onJobFinish(j)

		close(done)
	}()

	return
}

func (m *Manager) onJobFinish(job *Job) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if job.Status == StatusStopping {
		job.Status = StatusCancelled
	} else {
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

		if j.Status != StatusCancelled {
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
func (m *Manager) GetQueue() []*Job {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var ret []*Job

	for _, j := range m.queue {
		jCopy := *j
		ret = append(ret, &jCopy)
	}

	return ret
}

func (m *Manager) Subscribe(ctx context.Context) ManagerSubscription {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	ret := newSubscription()

	m.subscriptions = append(m.subscriptions, ret)

	go func() {
		<-ctx.Done()
		m.mutex.Lock()
		defer m.mutex.Unlock()

		ret.close()
	}()

	return ret
}

func (m *Manager) notifyJobUpdate(j *Job) {
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
	m   *Manager
	job *Job
}

func (u *updater) UpdateProgress(progress float64, details []string) {
	u.m.mutex.Lock()
	defer u.m.mutex.Unlock()

	u.job.Progress = progress
	u.job.Details = details

	u.m.notifyJobUpdate(u.job)
}
