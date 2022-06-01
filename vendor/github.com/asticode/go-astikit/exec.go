package astikit

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

// Statuses
const (
	ExecStatusCrashed = "crashed"
	ExecStatusRunning = "running"
	ExecStatusStopped = "stopped"
)

// ExecHandler represents an object capable of handling the execution of a cmd
type ExecHandler struct {
	cancel  context.CancelFunc
	ctx     context.Context
	err     error
	o       sync.Once
	stopped bool
}

// Status returns the cmd status
func (h *ExecHandler) Status() string {
	if h.ctx.Err() != nil {
		if h.stopped || h.err == nil {
			return ExecStatusStopped
		}
		return ExecStatusCrashed
	}
	return ExecStatusRunning
}

// Stop stops the cmd
func (h *ExecHandler) Stop() {
	h.o.Do(func() {
		h.cancel()
		h.stopped = true
	})
}

// ExecCmdOptions represents exec options
type ExecCmdOptions struct {
	Args       []string
	CmdAdapter func(cmd *exec.Cmd, h *ExecHandler) error
	Name       string
	StopFunc   func(cmd *exec.Cmd) error
}

// ExecCmd executes a cmd
// The process will be stopped when the worker stops
func ExecCmd(w *Worker, o ExecCmdOptions) (h *ExecHandler, err error) {
	// Create handler
	h = &ExecHandler{}
	h.ctx, h.cancel = context.WithCancel(w.Context())

	// Create command
	cmd := exec.Command(o.Name, o.Args...)

	// Adapt command
	if o.CmdAdapter != nil {
		if err = o.CmdAdapter(cmd, h); err != nil {
			err = fmt.Errorf("astikit: adapting cmd failed: %w", err)
			return
		}
	}

	// Start
	w.Logger().Infof("astikit: starting %s", strings.Join(cmd.Args, " "))
	if err = cmd.Start(); err != nil {
		err = fmt.Errorf("astikit: executing %s: %w", strings.Join(cmd.Args, " "), err)
		return
	}

	// Handle context
	go func() {
		// Wait for context to be done
		<-h.ctx.Done()

		// Get stop func
		f := func() error { return cmd.Process.Kill() }
		if o.StopFunc != nil {
			f = func() error { return o.StopFunc(cmd) }
		}

		// Stop
		if err = f(); err != nil {
			w.Logger().Error(fmt.Errorf("astikit: stopping cmd failed: %w", err))
			return
		}
	}()

	// Execute in a task
	w.NewTask().Do(func() {
		h.err = cmd.Wait()
		h.cancel()
		w.Logger().Infof("astikit: status is now %s for %s", h.Status(), strings.Join(cmd.Args, " "))
	})
	return
}
