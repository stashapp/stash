package search

import (
	"context"

	"github.com/stashapp/stash/pkg/event"
	"github.com/stashapp/stash/pkg/logger"
)

type eventTrack struct {
	ch chan event.Change
}

func newEventTrack() *eventTrack {
	return &eventTrack{
		ch: make(chan event.Change, 1),
	}
}

func (et *eventTrack) start(ctx context.Context, d *event.Dispatcher) {
	d.Register(et.ch)
	go func() {
		for {
			select {
			case <-ctx.Done():
				d.Unregister(et.ch)
				return
			case ev := <-et.ch:
				logger.Infof("Event: %v", ev)
			}
		}
	}()
}
