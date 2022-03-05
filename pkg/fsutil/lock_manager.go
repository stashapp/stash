package fsutil

import (
	"context"
	"sync"
)

type cancelContext struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// ReadLockManager manages read locks on file paths.
type ReadLockManager struct {
	readLocks map[string][]*cancelContext
	mutex     sync.RWMutex
}

// NewReadLockManager creates a new ReadLockManager.
func NewReadLockManager() *ReadLockManager {
	return &ReadLockManager{
		readLocks: make(map[string][]*cancelContext),
	}
}

// ReadLock adds a pending file read lock for fn to its storage, returning a context and cancel function.
// Per standard WithCancel usage, cancel must be called when the lock is freed.
func (m *ReadLockManager) ReadLock(ctx context.Context, fn string) (context.Context, context.CancelFunc) {
	retCtx, cancel := context.WithCancel(ctx)

	m.mutex.Lock()
	defer m.mutex.Unlock()

	locks := m.readLocks[fn]

	cc := &cancelContext{retCtx, cancel}
	m.readLocks[fn] = append(locks, cc)

	go m.waitAndUnlock(fn, cc)

	return retCtx, cancel
}

func (m *ReadLockManager) waitAndUnlock(fn string, cc *cancelContext) {
	<-cc.ctx.Done()

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
		l.cancel()
	}
}
