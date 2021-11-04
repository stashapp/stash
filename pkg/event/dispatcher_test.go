package event

import (
	"context"
	"math/rand"
	"testing"
	"time"
)

func randomEvent() Change {
	id := rand.Int()
	ty := rand.Intn(3)

	return Change{
		ID:   id,
		Type: ChangeType(ty),
	}
}

func TestDispatcher(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	d := NewDispatcher()
	d.Start(ctx)

	numCh := 4
	chs := make([]chan Change, 0, numCh)
	for i := 0; i < numCh; i++ {
		ch := make(chan Change, 1)
		d.Register(ch)
		chs = append(chs, ch)
	}

	trials := 20
	timeOut := 5 * time.Second
	// Each trial:
	// * Unregisters a random channel
	// * Publishes a random message which should then go to all other channels.
	// * Checks that the unregistered channel doesn't receive the event
	// * Re-registers the unregistered channel
	for i := 0; i < trials; i++ {
		have := randomEvent()
		j := rand.Intn(numCh)
		d.Unregister(chs[j])
		d.Publish(have)
		for k := range chs {
			if k == j {
				select {
				case e := <-chs[k]:
					t.Errorf("received event %v on unregistered channel", e)
				default:
				}
			} else {
				select {
				case got := <-chs[k]:
					if got.ID != have.ID || got.Type != have.Type {
						t.Errorf("chan[%d]: got: %+v; want %+v", k, got, have)
					}
				case <-time.After(timeOut):
					t.Fatal("Did not receive event in time")
				}
			}
		}
		d.Register(chs[j])
	}

	cancel()
}
