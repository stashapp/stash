package search

import (
	"context"
	"fmt"
	"strings"

	"github.com/stashapp/stash/pkg/event"
)

// changeSet is a rollup structure for changes. These are handed off to
// a batch processor when it requests them.
type changeSet struct {
	scenes     map[int]struct{}
	performers map[int]struct{}
}

// newChangemap creates a new initialized empty changeMap.
func newChangeSet() *changeSet {
	return &changeSet{
		scenes:     make(map[int]struct{}),
		performers: make(map[int]struct{}),
	}
}

// track records the given change to the changeMap.
func (s *changeSet) track(e event.Change) {
	switch e.Type {
	case event.Scene:
		s.scenes[e.ID] = struct{}{}
	case event.Performer:
		s.performers[e.ID] = struct{}{}
	default:
		// Ignore changes we don't currently track
	}
}

// cutSceneIds returns a slice of sceneIds from the changeSet.
// The limit argument sets an upper bound on the number of elements
// in the slice. Returns if there's any more elements to cut as well.
func (s *changeSet) cutSceneIds(limit int) ([]int, bool) {
	var ret []int
	for k := range s.scenes {
		if limit == 0 {
			return ret, true
		}

		ret = append(ret, k)
		delete(s.scenes, k)
		limit--
	}

	return ret, false
}

func (s *changeSet) performerIds() []int {
	var ret []int
	for k := range s.performers {
		ret = append(ret, k)
	}

	return ret
}

// hasContent returns true if there are changes to process.
func (s *changeSet) hasContent() bool {
	return len(s.scenes) > 0 || len(s.performers) > 0
}

// String implements the Stringer interface for changeMaps.
func (s *changeSet) String() string {
	var elems []string
	if len(s.scenes) > 0 {
		elems = append(elems, fmt.Sprintf("(%d scenes)", len(s.scenes)))
	}
	if len(s.performers) > 0 {
		elems = append(elems, fmt.Sprintf("(%d performers)", len(s.performers)))
	}

	if len(elems) == 0 {
		return "empty changeset"
	}

	return strings.Join(elems, ", ")
}

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
