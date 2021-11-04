// package event dispatches change-events in stash
package event

import (
	"context"
	"sync"
)

// ChangeType defines what has changed
type ChangeType int

const (
	SceneMarker ChangeType = iota
	Scene
	Image
	Gallery
	Movie
	Performer
	Studio
	Tag
)

// Change represents a pair of a Type and the ID which has changed
type Change struct {
	Type ChangeType
	ID   int
}

// Dispatcher represents a single event dispatcher
type Dispatcher struct {
	incoming chan Change
	mu       sync.Mutex
	chans    []chan Change
}

// NewDispatcher creates a new even dispatcher
func NewDispatcher() *Dispatcher {
	incoming := make(chan Change, 1)
	return &Dispatcher{
		incoming: incoming,
	}
}

// Start starts the dispatcher goroutine under the given context
func (d *Dispatcher) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case change := <-d.incoming:
				d.broadcast(change)
			}
		}
	}()
}

// Register registers chan for receiving events. It is up to the caller to ensure
// the channel doesn't block. I.e., events must be lifted from chan quickly.
func (d *Dispatcher) Register(c chan Change) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.chans = append(d.chans, c)
}

// Unregister removes a channel for dispatches
func (d *Dispatcher) Unregister(c chan Change) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Find the index at which the chan occur
	idx := -1
	for i := range d.chans {
		if d.chans[i] == c {
			idx = i
			break
		}
	}

	// If already unregistered, ignore
	if idx == -1 {
		return
	}

	// Swap chan to last element, cut it off
	last := len(d.chans) - 1
	d.chans[idx] = d.chans[last]
	d.chans = d.chans[:last]
}

// Publish broadcasts a change
func (d *Dispatcher) Publish(c Change) {
	d.incoming <- c
}

// broadcast fans-out a change to all channels registered
func (d *Dispatcher) broadcast(change Change) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, ch := range d.chans {
		ch <- change
	}
}
