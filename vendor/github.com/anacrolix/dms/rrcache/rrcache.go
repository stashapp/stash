// Package rrcache implements a random replacement cache. Items are set with
// an associated size. When the capacity is exceeded, items will be randomly
// evicted until it is not.
package rrcache

import (
	"math/rand"
)

type RRCache struct {
	capacity int64
	size     int64

	keys  []interface{}
	table map[interface{}]*entry
}

type entry struct {
	size  int64
	value interface{}
}

func New(capacity int64) *RRCache {
	return &RRCache{
		capacity: capacity,
		table:    make(map[interface{}]*entry),
	}
}

// Returns the sum size of all items currently in the cache.
func (c *RRCache) Size() int64 {
	return c.size
}

func (c *RRCache) Set(key interface{}, value interface{}, size int64) {
	if size > c.capacity {
		return
	}
	_entry := c.table[key]
	if _entry == nil {
		_entry = new(entry)
		c.keys = append(c.keys, key)
		c.table[key] = _entry
	}
	sizeDelta := size - _entry.size
	_entry.value = value
	_entry.size = size
	c.size += sizeDelta
	for c.size > c.capacity {
		i := rand.Intn(len(c.keys))
		key := c.keys[i]
		c.keys[i] = c.keys[len(c.keys)-1]
		c.keys = c.keys[:len(c.keys)-1]
		c.size -= c.table[key].size
		delete(c.table, key)
	}
}

func (c *RRCache) Get(key interface{}) (value interface{}, ok bool) {
	entry, ok := c.table[key]
	if !ok {
		return
	}
	value = entry.value
	return
}

type Item struct {
	Key, Value interface{}
}

// Return all items currently in the cache. This is made available for
// serialization purposes.
func (c *RRCache) Items() (itens []Item) {
	for k, e := range c.table {
		itens = append(itens, Item{k, e.value})
	}
	return
}
