package job

import (
	"sync"
	"time"

	"golang.org/x/net/context"
)

// Manager maintains a queue of jobs. Jobs are executed one at a time.
type Manager struct {
	queue []*Job

	mutex    sync.Mutex
	notEmpty *sync.Cond
	stop     chan struct{}

	lastID int
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

	return j.ID
}

func (m *Manager) nextID() int {
	m.lastID += 1
	return m.lastID
}

func (m *Manager) dispatcher() {
	m.mutex.Lock()

	for {
		// wait until we have something to process
		for len(m.queue) == 0 {
			m.notEmpty.Wait()

			// it's possible that we have been stopped - check here
			select {
			case <-m.stop:
				m.mutex.Unlock()
				return
			default:
				// keep going
			}
		}

		// grab to top job from the queue
		j := m.queue[0]

		if j.Status != StatusCancelled {
			done := m.dispatch(j)

			// unlock the mutex and wait for the job to finish
			m.mutex.Unlock()
			<-done
			m.mutex.Lock()
		}

		// remove the job from the queue
		// TODO - probably still want to store it somewhere
		m.queue = m.queue[1:]

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
		j.exec.Execute(ctx, &updater{
			m:   m,
			job: j,
		})

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

func (m *Manager) getJob(id int) (index int, job *Job) {
	// assumes lock held
	for i, j := range m.queue {
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

	index, j := m.getJob(id)
	if j != nil {
		j.cancel()

		if j.Status == StatusCancelled {
			// remove from the queue
			m.queue = append(m.queue[:index], m.queue[index+1:]...)
		}
	}
}

// CancelAll cancels all of the jobs in the queue. This is the same as
// calling CancelJob on all jobs in the queue.
func (m *Manager) CancelAll() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var newQueue []*Job

	// call cancel on all
	for _, j := range m.queue {
		j.cancel()

		if j.Status != StatusCancelled {
			newQueue = append(newQueue, j)
		}
	}

	m.queue = newQueue
}

// GetJob returns a copy of the Job for the provided id. Returns nil if the job
// does not exist.
func (m *Manager) GetJob(id int) *Job {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, j := m.getJob(id)
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

type updater struct {
	m   *Manager
	job *Job
}

func (u *updater) SetProgress(progress float64) {
	u.m.mutex.Lock()
	defer u.m.mutex.Unlock()

	u.job.Progress = progress

	// TODO need to notify
}

func (u *updater) AddSubTask(subtask string) int {
	u.m.mutex.Lock()
	defer u.m.mutex.Unlock()

	// TODO need to notify

	return u.job.addSubTask(subtask)
}

func (u *updater) RemoveSubTask(subtaskID int) {
	u.m.mutex.Lock()
	defer u.m.mutex.Unlock()

	// TODO need to notify

	u.job.removeSubTask(subtaskID)
}
