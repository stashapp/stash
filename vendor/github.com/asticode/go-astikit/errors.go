package astikit

import (
	"errors"
	"strings"
	"sync"
)

// Errors is an error containing multiple errors
type Errors struct {
	m *sync.Mutex // Locks p
	p []error
}

// NewErrors creates new errors
func NewErrors(errs ...error) *Errors {
	return &Errors{
		m: &sync.Mutex{},
		p: errs,
	}
}

// Add adds a new error
func (errs *Errors) Add(err error) {
	if err == nil {
		return
	}
	errs.m.Lock()
	defer errs.m.Unlock()
	errs.p = append(errs.p, err)
}

// IsNil checks whether the error is nil
func (errs *Errors) IsNil() bool {
	errs.m.Lock()
	defer errs.m.Unlock()
	return len(errs.p) == 0
}

// Loop loops through the errors
func (errs *Errors) Loop(fn func(idx int, err error) bool) {
	errs.m.Lock()
	defer errs.m.Unlock()
	for idx, err := range errs.p {
		if stop := fn(idx, err); stop {
			return
		}
	}
}

// Error implements the error interface
func (errs *Errors) Error() string {
	errs.m.Lock()
	defer errs.m.Unlock()
	var ss []string
	for _, err := range errs.p {
		ss = append(ss, err.Error())
	}
	return strings.Join(ss, " && ")
}

// ErrorCause returns the cause of an error
func ErrorCause(err error) error {
	for {
		if u := errors.Unwrap(err); u != nil {
			err = u
			continue
		}
		return err
	}
}
