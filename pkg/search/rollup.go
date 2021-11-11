package search

import (
	"context"

	"github.com/stashapp/stash/pkg/event"
)

// rollUp is the type used by the rollup engine.
type rollUp struct {
	// eventCh is the channel registered in the dispatcher.
	// it is used for receiving events.
	eventCh chan event.Change

	// handoff is an (unbuffered) channel used for batch processing.
	// Used for communication when the batch process is ready to work
	// on a new batch of data.
	handoff chan *changeSet

	// cur is the current changemap
	cur *changeSet
}

// newRollup creates a new rollup service.
func newRollup() *rollUp {
	return &rollUp{
		eventCh: make(chan event.Change, 1),
		handoff: make(chan *changeSet),
		cur:     newChangeSet(),
	}
}

// start starts the given rollup service under a given context.
// It will register on the given event dispatcher.
func (r *rollUp) start(ctx context.Context, d *event.Dispatcher) {
	d.Register(r.eventCh)
	go func() {
		for {
			select {
			case <-ctx.Done():
				d.Unregister(r.eventCh)
				return
			case r.handoff <- r.cur:
				// If we can hand off to a waiting receiver, we can't use
				// the current map anymore. Create a new one for the next
				// batch.
				r.cur = newChangeSet()
			case c := <-r.eventCh:
				r.cur.track(c)
			}
		}
	}()
}

// batch receives a batch from the rollup service.
func (r *rollUp) batch() *changeSet {
	return <-r.handoff
}
