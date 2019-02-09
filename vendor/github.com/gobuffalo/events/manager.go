package events

import (
	"strings"

	"github.com/markbates/safe"
	"github.com/pkg/errors"
)

type DeleteFn func()

// Manager can be implemented to replace the default
// events manager
type Manager interface {
	Listen(string, Listener) (DeleteFn, error)
	Emit(Event) error
}

// DefaultManager implements a map backed Manager
func DefaultManager() Manager {
	return &manager{
		listeners: listenerMap{},
	}
}

// SetManager allows you to replace the default
// event manager with a custom one
func SetManager(m Manager) {
	boss = m
}

var boss Manager = DefaultManager()
var _ listable = &manager{}

type manager struct {
	listeners listenerMap
}

func (m *manager) Listen(name string, l Listener) (DeleteFn, error) {
	_, ok := m.listeners.Load(name)
	if ok {
		return nil, errors.Errorf("listener named %s is already listening", name)
	}

	m.listeners.Store(name, l)

	df := func() {
		m.listeners.Delete(name)
	}

	return df, nil
}

func (m *manager) Emit(e Event) error {
	if err := e.Validate(); err != nil {
		return errors.WithStack(err)
	}
	e.Kind = strings.ToLower(e.Kind)
	if e.IsError() && e.Error == nil {
		e.Error = errors.New(e.Kind)
	}
	go func(e Event) {
		m.listeners.Range(func(key string, l Listener) bool {
			ex := Event{
				Kind:    e.Kind,
				Error:   e.Error,
				Message: e.Message,
				Payload: Payload{},
			}
			for k, v := range e.Payload {
				ex.Payload[k] = v
			}
			go func(e Event, l Listener) {
				safe.Run(func() {
					l(e)
				})
			}(ex, l)
			return true
		})
	}(e)
	return nil
}

func (m *manager) List() ([]string, error) {
	return m.listeners.Keys(), nil
}
