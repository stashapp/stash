package pkg

import (
	"sync"
	"time"
)

type cacheEntry struct {
	lastModified time.Time
	data         []RemotePackage
}

type repositoryCache struct {
	mu sync.RWMutex
	// cache maps the URL to the last modified time and the data
	cache map[string]cacheEntry
}

func (c *repositoryCache) lastModified(url string) *time.Time {
	if c == nil {
		return nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.cache == nil {
		return nil
	}

	e, found := c.cache[url]

	if !found {
		return nil
	}

	return &e.lastModified
}

func (c *repositoryCache) getPackageList(url string) []RemotePackage {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.cache == nil {
		return nil
	}

	e, found := c.cache[url]

	if !found {
		return nil
	}

	return e.data
}

func (c *repositoryCache) cacheList(url string, lastModified time.Time, data []RemotePackage) {
	if c == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cache == nil {
		c.cache = make(map[string]cacheEntry)
	}

	c.cache[url] = cacheEntry{
		lastModified: lastModified,
		data:         data,
	}
}
