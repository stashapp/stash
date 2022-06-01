package fsutil

import (
	"context"
	"os/exec"
	"sync"
	"time"
)

type Cancellable interface {
	Cancel()
}

type LockContext struct {
	context.Context
	cancel context.CancelFunc

	cmd *exec.Cmd
}

func (c *LockContext) AttachCommand(cmd *exec.Cmd) {
	c.cmd = cmd
}

func (c *LockContext) Cancel() {
	c.cancel()

	if c.cmd != nil {
		// wait for the process to die before returning
		// don't wait more than a few seconds
		done := make(chan error)
		go func() {
			err := c.cmd.Wait()
			done <- err
		}()

		select {
		case <-done:
			return
		case <-time.After(5 * time.Second):
			return
		}
	}
}

// ReadLockManager manages read locks on file paths.
type ReadLockManager struct {
	readLocks map[string][]*LockContext
	mutex     sync.RWMutex
}

// NewReadLockManager creates a new ReadLockManager.
func NewReadLockManager() *ReadLockManager {
	return &ReadLockManager{
		readLocks: make(map[string][]*LockContext),
	}
}

// ReadLock adds a pending file read lock for fn to its storage, returning a context and cancel function.
// Per standard WithCancel usage, cancel must be called when the lock is freed.
func (m *ReadLockManager) ReadLock(ctx context.Context, fn string) *LockContext {
	retCtx, cancel := context.WithCancel(ctx)

	// if Cancellable, call Cancel() when cancelled
	cancellable, ok := ctx.(Cancellable)
	if ok {
		origCancel := cancel
		cancel = func() {
			origCancel()
			cancellable.Cancel()
		}
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	locks := m.readLocks[fn]

	cc := &LockContext{
		Context: retCtx,
		cancel:  cancel,
	}
	m.readLocks[fn] = append(locks, cc)

	go m.waitAndUnlock(fn, cc)

	return cc
}

func (m *ReadLockManager) waitAndUnlock(fn string, cc *LockContext) {
	<-cc.Done()

	m.mutex.Lock()
	defer m.mutex.Unlock()

	locks := m.readLocks[fn]
	for i, v := range locks {
		if v == cc {
			m.readLocks[fn] = append(locks[:i], locks[i+1:]...)
			return
		}
	}
}

// Cancel cancels all read lock contexts associated with fn.
func (m *ReadLockManager) Cancel(fn string) {
	m.mutex.RLock()
	locks := m.readLocks[fn]
	m.mutex.RUnlock()

	for _, l := range locks {
		l.Cancel()
		<-l.Done()
	}
}
