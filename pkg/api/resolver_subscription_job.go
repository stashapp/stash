package api

import (
	"context"
	"sync"
	"time"

	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/models"
)

type throttledUpdate struct {
	id             int
	pendingUpdate  *job.Job
	lastUpdate     time.Time
	broadcastTimer *time.Timer
	killTimer      *time.Timer
}

func (tu *throttledUpdate) broadcast(output chan *models.JobStatusUpdate) {
	tu.lastUpdate = time.Now()
	output <- &models.JobStatusUpdate{
		Type: models.JobStatusUpdateTypeUpdate,
		Job:  jobToJobModel(*tu.pendingUpdate),
	}

	tu.broadcastTimer = nil
	tu.pendingUpdate = nil
}

type updateThrottler struct {
	output chan *models.JobStatusUpdate

	updates []*throttledUpdate
	closed  bool
	mutex   sync.Mutex
}

func (u *updateThrottler) findJob(id int) *throttledUpdate {
	for _, tu := range u.updates {
		if tu.id == id {
			return tu
		}
	}

	return nil
}

func (u *updateThrottler) removeJob(id int) {
	for i, tu := range u.updates {
		if tu.id == id {
			u.updates = append(u.updates[:i], u.updates[i+1:]...)
			return
		}
	}
}

func (u *updateThrottler) makeKillTimer(tu *throttledUpdate) {
	tu.killTimer = time.AfterFunc(time.Second, func() {
		u.mutex.Lock()
		defer u.mutex.Unlock()

		if tu.pendingUpdate != nil || time.Since(tu.lastUpdate) > time.Second {
			u.removeJob(tu.id)
		} else {
			u.makeKillTimer(tu)
		}
	})
}

func (u *updateThrottler) updateJob(j *job.Job) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	tu := u.findJob(j.ID)
	if tu == nil {
		tu := &throttledUpdate{
			id:            j.ID,
			pendingUpdate: j,
		}
		u.updates = append(u.updates, tu)

		tu.broadcast(u.output)
		u.makeKillTimer(tu)
	} else {
		tu.pendingUpdate = j

		if tu.broadcastTimer == nil {
			tu.broadcastTimer = time.AfterFunc(time.Second-time.Since(tu.lastUpdate), func() {
				u.mutex.Lock()
				defer u.mutex.Unlock()

				if !u.closed {
					tu.broadcast(u.output)
				}
			})
		}
	}
}

func (u *updateThrottler) close() {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	u.closed = true
}

func (r *subscriptionResolver) JobsSubscribe(ctx context.Context) (<-chan *models.JobStatusUpdate, error) {
	msg := make(chan *models.JobStatusUpdate, 100)

	subscription := manager.GetInstance().JobManager.Subscribe(ctx)

	throttler := updateThrottler{
		output: msg,
	}

	go func() {
		for {
			select {
			case j := <-subscription.NewJob:
				u := models.JobStatusUpdate{
					Type: models.JobStatusUpdateTypeAdd,
					Job:  jobToJobModel(j),
				}
				msg <- &u
			case j := <-subscription.RemovedJob:
				u := models.JobStatusUpdate{
					Type: models.JobStatusUpdateTypeRemove,
					Job:  jobToJobModel(j),
				}
				msg <- &u
			case j := <-subscription.UpdatedJob:
				// throttle updates to no more than one per second
				throttler.updateJob(&j)
			case <-ctx.Done():
				throttler.close()
				close(msg)
				return
			}
		}
	}()

	return msg, nil
}

func (r *subscriptionResolver) ScanCompleteSubscribe(ctx context.Context) (<-chan bool, error) {
	return manager.GetInstance().ScanSubscribe(ctx), nil
}
