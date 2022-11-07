package astikit

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// Stat names
const (
	StatNameWorkRatio = "astikit.work.ratio"
)

// Chan constants
const (
	// Calling Add() only blocks if the chan has been started and the ctx
	// has not been canceled
	ChanAddStrategyBlockWhenStarted = "block.when.started"
	// Calling Add() never blocks
	ChanAddStrategyNoBlock = "no.block"
	ChanOrderFIFO          = "fifo"
	ChanOrderFILO          = "filo"
)

// Chan is an object capable of executing funcs in a specific order while controlling the conditions
// in which adding new funcs is blocking
// Check out ChanOptions for detailed options
type Chan struct {
	cancel        context.CancelFunc
	c             *sync.Cond
	ctx           context.Context
	fs            []func()
	mc            *sync.Mutex // Locks ctx
	mf            *sync.Mutex // Locks fs
	o             ChanOptions
	running       uint32
	statWorkRatio *DurationPercentageStat
}

// ChanOptions are Chan options
type ChanOptions struct {
	// Determines the conditions in which Add() blocks. See constants with pattern ChanAddStrategy*
	// Default is ChanAddStrategyNoBlock
	AddStrategy string
	// Order in which the funcs will be processed. See constants with pattern ChanOrder*
	// Default is ChanOrderFIFO
	Order string
	// By default the funcs not yet processed when the context is cancelled are dropped.
	// If "ProcessAll" is true,  ALL funcs are processed even after the context is cancelled.
	// However, no funcs can be added after the context is cancelled
	ProcessAll bool
}

// NewChan creates a new Chan
func NewChan(o ChanOptions) *Chan {
	return &Chan{
		c:  sync.NewCond(&sync.Mutex{}),
		mc: &sync.Mutex{},
		mf: &sync.Mutex{},
		o:  o,
	}
}

// Start starts the chan by looping through functions in the buffer and
// executing them if any, or waiting for a new one otherwise
func (c *Chan) Start(ctx context.Context) {
	// Make sure to start only once
	if atomic.CompareAndSwapUint32(&c.running, 0, 1) {
		// Update status
		defer atomic.StoreUint32(&c.running, 0)

		// Create context
		c.mc.Lock()
		c.ctx, c.cancel = context.WithCancel(ctx)
		d := c.ctx.Done()
		c.mc.Unlock()

		// Handle context
		go func() {
			// Wait for context to be done
			<-d

			// Signal
			c.c.L.Lock()
			c.c.Signal()
			c.c.L.Unlock()
		}()

		// Loop
		for {
			// Lock cond here in case a func is added between retrieving l and doing the if on it
			c.c.L.Lock()

			// Get number of funcs in buffer
			c.mf.Lock()
			l := len(c.fs)
			c.mf.Unlock()

			// Only return if context has been cancelled and:
			//   - the user wants to drop funcs that has not yet been processed
			//   - the buffer is empty otherwise
			c.mc.Lock()
			if c.ctx.Err() != nil && (!c.o.ProcessAll || l == 0) {
				c.mc.Unlock()
				c.c.L.Unlock()
				return
			}
			c.mc.Unlock()

			// No funcs in buffer
			if l == 0 {
				c.c.Wait()
				c.c.L.Unlock()
				continue
			}
			c.c.L.Unlock()

			// Get first func
			c.mf.Lock()
			fn := c.fs[0]
			c.mf.Unlock()

			// Execute func
			if c.statWorkRatio != nil {
				c.statWorkRatio.Begin()
			}
			fn()
			if c.statWorkRatio != nil {
				c.statWorkRatio.End()
			}

			// Remove first func
			c.mf.Lock()
			c.fs = c.fs[1:]
			c.mf.Unlock()
		}
	}
}

// Stop stops the chan
func (c *Chan) Stop() {
	c.mc.Lock()
	if c.cancel != nil {
		c.cancel()
	}
	c.mc.Unlock()
}

// Add adds a new item to the chan
func (c *Chan) Add(i func()) {
	// Check context
	c.mc.Lock()
	if c.ctx != nil && c.ctx.Err() != nil {
		c.mc.Unlock()
		return
	}
	c.mc.Unlock()

	// Wrap the function
	var fn func()
	var wg *sync.WaitGroup
	if c.o.AddStrategy == ChanAddStrategyBlockWhenStarted {
		wg = &sync.WaitGroup{}
		wg.Add(1)
		fn = func() {
			defer wg.Done()
			i()
		}
	} else {
		fn = i
	}

	// Add func to buffer
	c.mf.Lock()
	if c.o.Order == ChanOrderFILO {
		c.fs = append([]func(){fn}, c.fs...)
	} else {
		c.fs = append(c.fs, fn)
	}
	c.mf.Unlock()

	// Signal
	c.c.L.Lock()
	c.c.Signal()
	c.c.L.Unlock()

	// Wait
	if wg != nil {
		wg.Wait()
	}
}

// Reset resets the chan
func (c *Chan) Reset() {
	c.mf.Lock()
	defer c.mf.Unlock()
	c.fs = []func(){}
}

// Stats returns the chan stats
func (c *Chan) Stats() []StatOptions {
	if c.statWorkRatio == nil {
		c.statWorkRatio = NewDurationPercentageStat()
	}
	return []StatOptions{
		{
			Handler: c.statWorkRatio,
			Metadata: &StatMetadata{
				Description: "Percentage of time doing work",
				Label:       "Work ratio",
				Name:        StatNameWorkRatio,
				Unit:        "%",
			},
		},
	}
}

// BufferPool represents a *bytes.Buffer pool
type BufferPool struct {
	bp *sync.Pool
}

// NewBufferPool creates a new BufferPool
func NewBufferPool() *BufferPool {
	return &BufferPool{bp: &sync.Pool{New: func() interface{} { return &bytes.Buffer{} }}}
}

// New creates a new BufferPoolItem
func (p *BufferPool) New() *BufferPoolItem {
	return newBufferPoolItem(p.bp.Get().(*bytes.Buffer), p.bp)
}

// BufferPoolItem represents a BufferPool item
type BufferPoolItem struct {
	*bytes.Buffer
	bp *sync.Pool
}

func newBufferPoolItem(b *bytes.Buffer, bp *sync.Pool) *BufferPoolItem {
	return &BufferPoolItem{
		Buffer: b,
		bp:     bp,
	}
}

// Close implements the io.Closer interface
func (i *BufferPoolItem) Close() error {
	i.Reset()
	i.bp.Put(i.Buffer)
	return nil
}

// GoroutineLimiter is an object capable of doing several things in parallel while maintaining the
// max number of things running in parallel under a threshold
type GoroutineLimiter struct {
	busy   int
	c      *sync.Cond
	ctx    context.Context
	cancel context.CancelFunc
	o      GoroutineLimiterOptions
}

// GoroutineLimiterOptions represents GoroutineLimiter options
type GoroutineLimiterOptions struct {
	Max int
}

// NewGoroutineLimiter creates a new GoroutineLimiter
func NewGoroutineLimiter(o GoroutineLimiterOptions) (l *GoroutineLimiter) {
	l = &GoroutineLimiter{
		c: sync.NewCond(&sync.Mutex{}),
		o: o,
	}
	if l.o.Max <= 0 {
		l.o.Max = 1
	}
	l.ctx, l.cancel = context.WithCancel(context.Background())
	go l.handleCtx()
	return
}

// Close closes the limiter properly
func (l *GoroutineLimiter) Close() error {
	l.cancel()
	return nil
}

func (l *GoroutineLimiter) handleCtx() {
	<-l.ctx.Done()
	l.c.L.Lock()
	l.c.Broadcast()
	l.c.L.Unlock()
}

// GoroutineLimiterFunc is a GoroutineLimiter func
type GoroutineLimiterFunc func()

// Do executes custom work in a goroutine
func (l *GoroutineLimiter) Do(fn GoroutineLimiterFunc) (err error) {
	// Check context in case the limiter has already been closed
	if err = l.ctx.Err(); err != nil {
		return
	}

	// Lock
	l.c.L.Lock()

	// Wait for a goroutine to be available
	for l.busy >= l.o.Max {
		l.c.Wait()
	}

	// Check context in case the limiter has been closed while waiting
	if err = l.ctx.Err(); err != nil {
		return
	}

	// Increment
	l.busy++

	// Unlock
	l.c.L.Unlock()

	// Execute in a goroutine
	go func() {
		// Decrement
		defer func() {
			l.c.L.Lock()
			l.busy--
			l.c.Signal()
			l.c.L.Unlock()
		}()

		// Execute
		fn()
	}()
	return
}

// Eventer represents an object that can dispatch simple events (name + payload)
type Eventer struct {
	c  *Chan
	hs map[string][]EventerHandler
	mh *sync.Mutex
}

// EventerOptions represents Eventer options
type EventerOptions struct {
	Chan ChanOptions
}

// EventerHandler represents a function that can handle the payload of an event
type EventerHandler func(payload interface{})

// NewEventer creates a new eventer
func NewEventer(o EventerOptions) *Eventer {
	return &Eventer{
		c:  NewChan(o.Chan),
		hs: make(map[string][]EventerHandler),
		mh: &sync.Mutex{},
	}
}

// On adds an handler for a specific name
func (e *Eventer) On(name string, h EventerHandler) {
	// Lock
	e.mh.Lock()
	defer e.mh.Unlock()

	// Add handler
	e.hs[name] = append(e.hs[name], h)
}

// Dispatch dispatches a payload for a specific name
func (e *Eventer) Dispatch(name string, payload interface{}) {
	// Lock
	e.mh.Lock()
	defer e.mh.Unlock()

	// No handlers
	hs, ok := e.hs[name]
	if !ok {
		return
	}

	// Loop through handlers
	for _, h := range hs {
		func(h EventerHandler) {
			// Add to chan
			e.c.Add(func() {
				h(payload)
			})
		}(h)
	}
}

// Start starts the eventer. It is blocking
func (e *Eventer) Start(ctx context.Context) {
	e.c.Start(ctx)
}

// Stop stops the eventer
func (e *Eventer) Stop() {
	e.c.Stop()
}

// Reset resets the eventer
func (e *Eventer) Reset() {
	e.c.Reset()
}

// RWMutex represents a RWMutex capable of logging its actions to ease deadlock debugging
type RWMutex struct {
	c string // Last successful caller
	l SeverityLogger
	m *sync.RWMutex
	n string // Name
}

// RWMutexOptions represents RWMutex options
type RWMutexOptions struct {
	Logger StdLogger
	Name   string
}

// NewRWMutex creates a new RWMutex
func NewRWMutex(o RWMutexOptions) *RWMutex {
	return &RWMutex{
		l: AdaptStdLogger(o.Logger),
		m: &sync.RWMutex{},
		n: o.Name,
	}
}

func (m *RWMutex) caller() (o string) {
	if _, file, line, ok := runtime.Caller(2); ok {
		o = fmt.Sprintf("%s:%d", file, line)
	}
	return
}

// Lock write locks the mutex
func (m *RWMutex) Lock() {
	c := m.caller()
	m.l.Debugf("astikit: requesting lock for %s at %s", m.n, c)
	m.m.Lock()
	m.l.Debugf("astikit: lock acquired for %s at %s", m.n, c)
	m.c = c
}

// Unlock write unlocks the mutex
func (m *RWMutex) Unlock() {
	m.m.Unlock()
	m.l.Debugf("astikit: unlock executed for %s", m.n)
}

// RLock read locks the mutex
func (m *RWMutex) RLock() {
	c := m.caller()
	m.l.Debugf("astikit: requesting rlock for %s at %s", m.n, c)
	m.m.RLock()
	m.l.Debugf("astikit: rlock acquired for %s at %s", m.n, c)
	m.c = c
}

// RUnlock read unlocks the mutex
func (m *RWMutex) RUnlock() {
	m.m.RUnlock()
	m.l.Debugf("astikit: unlock executed for %s", m.n)
}

// IsDeadlocked checks whether the mutex is deadlocked with a given timeout
// and returns the last caller
func (m *RWMutex) IsDeadlocked(timeout time.Duration) (bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	go func() {
		m.m.Lock()
		cancel()
		m.m.Unlock()
	}()
	<-ctx.Done()
	return errors.Is(ctx.Err(), context.DeadlineExceeded), m.c
}
