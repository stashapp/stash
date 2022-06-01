package astikit

import (
	"context"
	"os"
	"os/signal"
	"sync"
)

// Worker represents an object capable of blocking, handling signals and stopping
type Worker struct {
	cancel context.CancelFunc
	ctx    context.Context
	l      SeverityLogger
	os, ow sync.Once
	wg     *sync.WaitGroup
}

// WorkerOptions represents worker options
type WorkerOptions struct {
	Logger StdLogger
}

// NewWorker builds a new worker
func NewWorker(o WorkerOptions) (w *Worker) {
	w = &Worker{
		l:  AdaptStdLogger(o.Logger),
		wg: &sync.WaitGroup{},
	}
	w.ctx, w.cancel = context.WithCancel(context.Background())
	w.wg.Add(1)
	w.l.Info("astikit: starting worker...")
	return
}

// HandleSignals handles signals
func (w *Worker) HandleSignals(hs ...SignalHandler) {
	// Prepend mandatory handler
	hs = append([]SignalHandler{TermSignalHandler(w.Stop)}, hs...)

	// Notify
	ch := make(chan os.Signal, 1)
	signal.Notify(ch)

	// Execute in a task
	w.NewTask().Do(func() {
		for {
			select {
			case s := <-ch:
				// Loop through handlers
				for _, h := range hs {
					h(s)
				}

				// Return
				if isTermSignal(s) {
					return
				}
			case <-w.Context().Done():
				return
			}
		}
	})
}

// Stop stops the Worker
func (w *Worker) Stop() {
	w.os.Do(func() {
		w.l.Info("astikit: stopping worker...")
		w.cancel()
		w.wg.Done()
	})
}

// Wait is a blocking pattern
func (w *Worker) Wait() {
	w.ow.Do(func() {
		w.l.Info("astikit: worker is now waiting...")
		w.wg.Wait()
	})
}

// NewTask creates a new task
func (w *Worker) NewTask() *Task {
	return newTask(w.wg)
}

// Context returns the worker's context
func (w *Worker) Context() context.Context {
	return w.ctx
}

// Logger returns the worker's logger
func (w *Worker) Logger() SeverityLogger {
	return w.l
}

// TaskFunc represents a function that can create a new task
type TaskFunc func() *Task

// Task represents a task
type Task struct {
	od, ow  sync.Once
	wg, pwg *sync.WaitGroup
}

func newTask(parentWg *sync.WaitGroup) (t *Task) {
	t = &Task{
		wg:  &sync.WaitGroup{},
		pwg: parentWg,
	}
	t.pwg.Add(1)
	return
}

// NewSubTask creates a new sub task
func (t *Task) NewSubTask() *Task {
	return newTask(t.wg)
}

// Do executes the task
func (t *Task) Do(f func()) {
	go func() {
		// Make sure to mark the task as done
		defer t.Done()

		// Custom
		f()

		// Wait for first level subtasks to be done
		// Wait() can also be called in f() if something needs to be executed just after Wait()
		t.Wait()
	}()
}

// Done indicates the task is done
func (t *Task) Done() {
	t.od.Do(func() {
		t.pwg.Done()
	})
}

// Wait waits for first level subtasks to be finished
func (t *Task) Wait() {
	t.ow.Do(func() {
		t.wg.Wait()
	})
}
