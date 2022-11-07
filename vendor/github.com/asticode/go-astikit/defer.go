package astikit

import (
	"sync"
)

// CloseFunc is a method that closes something
type CloseFunc func() error

// Closer is an object that can close several things
type Closer struct {
	fs []CloseFunc
	m  *sync.Mutex
}

// NewCloser creates a new closer
func NewCloser() *Closer {
	return &Closer{
		m: &sync.Mutex{},
	}
}

// Close implements the io.Closer interface
func (c *Closer) Close() error {
	// Lock
	c.m.Lock()
	defer c.m.Unlock()

	// Loop through closers
	err := NewErrors()
	for _, f := range c.fs {
		err.Add(f())
	}

	// Reset closers
	c.fs = []CloseFunc{}

	// Return
	if err.IsNil() {
		return nil
	}
	return err
}

// Add adds a close func at the beginning of the list
func (c *Closer) Add(f CloseFunc) {
	c.m.Lock()
	defer c.m.Unlock()
	c.fs = append([]CloseFunc{f}, c.fs...)
}

// NewChild creates a new child closer
func (c *Closer) NewChild() (child *Closer) {
	child = NewCloser()
	c.Add(child.Close)
	return
}
