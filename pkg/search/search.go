package search

import (
	"context"
	"time"

	"github.com/stashapp/stash/pkg/event"
	"github.com/stashapp/stash/pkg/logger"
)

// Engine represents a search engine service.
type Engine struct {
	rollUp *rollUp
}

// NewEngine creates a new search engine.
func NewEngine() *Engine {
	return &Engine{
		rollUp: newRollup(),
	}
}

// Start starts the given Engine under a given context, processing events from a given dispatcher.
func (e *Engine) Start(ctx context.Context, d *event.Dispatcher) {
	go func() {
		e.rollUp.start(ctx, d)

		// How often to process batches.
		tick := time.NewTicker(1 * time.Minute)
		for {
			select {
			case <-ctx.Done():
				tick.Stop()
				return
			case <-tick.C:
				// Perform batch insert
				m := e.rollUp.batch()

				if m.hasContent() {
					logger.Infof("processing search batch %v", m)
				}
			}
		}
	}()
}
