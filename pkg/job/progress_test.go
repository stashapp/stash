package job

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createProgress(m *Manager, j *Job) Progress {
	return Progress{
		updater: &updater{
			m:   m,
			job: j,
		},
		total:     100,
		defined:   true,
		processed: 10,
		percent:   10,
	}
}

func TestProgressIndefinite(t *testing.T) {
	m := NewManager()
	j := &Job{}

	p := createProgress(m, j)

	p.Indefinite()

	assert := assert.New(t)

	// ensure job progress was updated
	assert.Equal(ProgressIndefinite, j.Progress)
}

func TestProgressSetTotal(t *testing.T) {
	m := NewManager()
	j := &Job{}

	p := createProgress(m, j)

	p.SetTotal(50)

	assert := assert.New(t)

	// ensure job progress was updated
	assert.Equal(0.2, j.Progress)

	p.SetTotal(0)
	assert.Equal(ProgressIndefinite, j.Progress)

	p.SetTotal(-10)
	assert.Equal(ProgressIndefinite, j.Progress)

	p.SetTotal(9)
	assert.Equal(float64(1), j.Progress)
}

func TestProgressSetProcessed(t *testing.T) {
	m := NewManager()
	j := &Job{}

	p := createProgress(m, j)

	p.SetProcessed(30)

	assert := assert.New(t)

	// ensure job progress was updated
	assert.Equal(0.3, j.Progress)

	p.SetProcessed(-10)
	assert.Equal(float64(0), j.Progress)

	p.SetProcessed(200)
	assert.Equal(float64(1), j.Progress)
}

func TestProgressSetPercent(t *testing.T) {
	m := NewManager()
	j := &Job{}

	p := createProgress(m, j)

	p.SetPercent(0.3)

	assert := assert.New(t)

	// ensure job progress was updated
	assert.Equal(0.3, j.Progress)

	p.SetPercent(-10)
	assert.Equal(float64(0), j.Progress)

	p.SetPercent(200)
	assert.Equal(float64(1), j.Progress)
}

func TestProgressIncrement(t *testing.T) {
	m := NewManager()
	j := &Job{}

	p := createProgress(m, j)

	p.SetProcessed(49)
	p.Increment()

	assert := assert.New(t)

	// ensure job progress was updated
	assert.Equal(0.5, j.Progress)

	p.SetProcessed(100)
	p.Increment()
	assert.Equal(float64(1), j.Progress)
}

func TestExecuteTask(t *testing.T) {
	m := NewManager()
	j := &Job{}

	p := createProgress(m, j)

	c := make(chan struct{}, 1)
	const taskDesciption = "taskDescription"

	go p.ExecuteTask(taskDesciption, func() {
		<-c
	})

	time.Sleep(sleepTime)

	assert := assert.New(t)

	m.mutex.Lock()
	// ensure task is added to the job details
	assert.Equal(taskDesciption, j.Details[0])
	m.mutex.Unlock()

	// allow task to finish
	close(c)

	time.Sleep(sleepTime)

	m.mutex.Lock()
	// ensure task is removed from the job details
	assert.Len(j.Details, 0)
	m.mutex.Unlock()
}
