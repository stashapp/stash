package utils

import (
	"sync"
	"testing"
)

// should be run with -race
func TestMutexManager(t *testing.T) {
	m := NewMutexManager()

	map1 := make(map[string]bool)
	map2 := make(map[string]bool)
	map3 := make(map[string]bool)
	maps := []map[string]bool{
		map1,
		map2,
		map3,
	}

	types := []string{
		"foo",
		"foo",
		"bar",
	}

	const key = "baz"

	const workers = 8
	const loops = 300
	var wg sync.WaitGroup
	for k := 0; k < workers; k++ {
		wg.Add(1)
		go func(wk int) {
			defer wg.Done()
			for l := 0; l < loops; l++ {
				func(l int) {
					c := make(chan struct{})
					defer close(c)

					m.Claim(types[l%3], key, c)

					maps[l%3][key] = true
				}(l)
			}
		}(k)
	}

	wg.Wait()
}
