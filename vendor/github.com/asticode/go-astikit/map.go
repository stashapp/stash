package astikit

import (
	"fmt"
	"sync"
)

// BiMap represents a bidirectional map
type BiMap struct {
	forward map[interface{}]interface{}
	inverse map[interface{}]interface{}
	m       *sync.Mutex
}

// NewBiMap creates a new BiMap
func NewBiMap() *BiMap {
	return &BiMap{
		forward: make(map[interface{}]interface{}),
		inverse: make(map[interface{}]interface{}),
		m:       &sync.Mutex{},
	}
}

func (m *BiMap) get(k interface{}, i map[interface{}]interface{}) (v interface{}, ok bool) {
	m.m.Lock()
	defer m.m.Unlock()
	v, ok = i[k]
	return
}

// Get gets the value in the forward map based on the provided key
func (m *BiMap) Get(k interface{}) (interface{}, bool) { return m.get(k, m.forward) }

// GetInverse gets the value in the inverse map based on the provided key
func (m *BiMap) GetInverse(k interface{}) (interface{}, bool) { return m.get(k, m.inverse) }

// MustGet gets the value in the forward map based on the provided key and panics if key is not found
func (m *BiMap) MustGet(k interface{}) interface{} {
	v, ok := m.get(k, m.forward)
	if !ok {
		panic(fmt.Sprintf("astikit: key %+v not found in foward map", k))
	}
	return v
}

// MustGetInverse gets the value in the inverse map based on the provided key and panics if key is not found
func (m *BiMap) MustGetInverse(k interface{}) interface{} {
	v, ok := m.get(k, m.inverse)
	if !ok {
		panic(fmt.Sprintf("astikit: key %+v not found in inverse map", k))
	}
	return v
}

func (m *BiMap) set(k, v interface{}, f, i map[interface{}]interface{}) *BiMap {
	m.m.Lock()
	defer m.m.Unlock()
	f[k] = v
	i[v] = k
	return m
}

// Set sets the value in the forward and inverse map for the provided forward key
func (m *BiMap) Set(k, v interface{}) *BiMap { return m.set(k, v, m.forward, m.inverse) }

// SetInverse sets the value in the forward and inverse map for the provided inverse key
func (m *BiMap) SetInverse(k, v interface{}) *BiMap { return m.set(k, v, m.inverse, m.forward) }
