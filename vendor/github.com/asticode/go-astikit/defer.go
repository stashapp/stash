package astikit

import (
	"sync"
)

type CloseFunc func()
type CloseFuncWithError func() error

// Closer is an object that can close several things
type Closer struct {
	closed bool
	fs     []CloseFuncWithError
	// We need to split into 2 mutexes to allow using .Add() in .Do()
	mc       *sync.Mutex // Locks .Close()
	mf       *sync.Mutex // Locks fs
	onClosed CloserOnClosed
}

type CloserOnClosed func(err error)

// NewCloser creates a new closer
func NewCloser() *Closer {
	return &Closer{
		mc: &sync.Mutex{},
		mf: &sync.Mutex{},
	}
}

// Close implements the io.Closer interface
func (c *Closer) Close() error {
	// Lock
	c.mc.Lock()
	defer c.mc.Unlock()

	// Get funcs
	c.mf.Lock()
	fs := c.fs
	c.mf.Unlock()

	// Loop through closers
	err := NewErrors()
	for _, f := range fs {
		err.Add(f())
	}

	// Reset closers
	c.fs = []CloseFuncWithError{}

	// Update attribute
	c.closed = true

	// Callback
	if c.onClosed != nil {
		c.onClosed(err)
	}

	// Return
	if err.IsNil() {
		return nil
	}
	return err
}

func (c *Closer) Add(f CloseFunc) {
	c.AddWithError(func() error {
		f()
		return nil
	})
}

func (c *Closer) AddWithError(f CloseFuncWithError) {
	// Lock
	c.mf.Lock()
	defer c.mf.Unlock()

	// Append
	c.fs = append([]CloseFuncWithError{f}, c.fs...)
}

func (c *Closer) Append(dst *Closer) {
	// Lock
	c.mf.Lock()
	dst.mf.Lock()
	defer c.mf.Unlock()
	defer dst.mf.Unlock()

	// Append
	c.fs = append(c.fs, dst.fs...)
}

// NewChild creates a new child closer
func (c *Closer) NewChild() (child *Closer) {
	child = NewCloser()
	c.AddWithError(child.Close)
	return
}

// Do executes a callback while ensuring :
//   - closer hasn't been closed before
//   - closer can't be closed in between
func (c *Closer) Do(fn func()) {
	// Lock
	c.mc.Lock()
	defer c.mc.Unlock()

	// Closer already closed
	if c.closed {
		return
	}

	// Callback
	fn()
}

func (c *Closer) OnClosed(fn CloserOnClosed) {
	c.mc.Lock()
	defer c.mc.Unlock()
	c.onClosed = fn
}

func (c *Closer) IsClosed() bool {
	c.mc.Lock()
	defer c.mc.Unlock()
	return c.closed
}
