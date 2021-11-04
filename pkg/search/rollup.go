package search

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/event"
)

// changeMap is a rollup structure for changes. These are handed off to
// a batch processor when it requests them.
type changeMap struct {
	scenes map[int]struct{}
}

// newChangemap creates a new initialized empty changeMap.
func newChangeMap() *changeMap {
	return &changeMap{
		scenes: make(map[int]struct{}),
	}
}

// track records the given change to the changeMap.
func (m *changeMap) track(e event.Change) {
	switch e.Type {
	case event.Scene:
		m.scenes[e.ID] = struct{}{}
	default:
		// Ignore changes we don't currently track
	}
}

// hasContent returns true if there are changes to process.
func (m *changeMap) hasContent() bool {
	return len(m.scenes) > 0
}

// String implements the Stringer interface for changeMaps.
func (m *changeMap) String() string {
	return fmt.Sprintf("(%d scenes)", len(m.scenes))
}

// rollUp is the type used by the rollup engine.
type rollUp struct {
	// eventCh is the channel registered in the dispatcher.
	// it is used for receiving events.
	eventCh chan event.Change

	// handoff is an (unbuffered) channel used for batch processing.
	// Used for communication when the batch process is ready to work
	// on a new batch of data.
	handoff chan *changeMap

	// cur is the current changemap
	cur *changeMap
}

// newRollup creates a new rollup service.
func newRollup() *rollUp {
	return &rollUp{
		eventCh: make(chan event.Change, 1),
		handoff: make(chan *changeMap),
		cur:     newChangeMap(),
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
				r.cur = newChangeMap()
			case c := <-r.eventCh:
				r.cur.track(c)
			}
		}
	}()
}

// batch receives a batch from the rollup service.
func (r *rollUp) batch() *changeMap {
	return <-r.handoff
}
