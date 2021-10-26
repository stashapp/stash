package utils

// MutexManager manages access to mutexes using a mutex type and key.
type MutexManager struct {
	mapChan chan map[string]<-chan struct{}
}

// NewMutexManager returns a new instance of MutexManager.
func NewMutexManager() *MutexManager {
	ret := &MutexManager{
		mapChan: make(chan map[string]<-chan struct{}, 1),
	}

	initial := make(map[string]<-chan struct{})
	ret.mapChan <- initial

	return ret
}

// Claim blocks until the mutex for the mutexType and key pair is available.
// The mutex is then claimed by the calling code until the provided done
// channel is closed.
func (csm *MutexManager) Claim(mutexType string, key string, done <-chan struct{}) {
	mapKey := mutexType + "_" + key
	success := false

	var existing <-chan struct{}
	for !success {
		// grab the map
		m := <-csm.mapChan

		// get the entry for the given key
		newEntry := m[mapKey]

		// if its the existing entry or nil, then it's available, add our channel
		if newEntry == nil || newEntry == existing {
			m[mapKey] = done
			success = true
		}

		// return the map
		csm.mapChan <- m

		// if there is an existing entry, now we can wait for it to
		// finish, then repeat the process
		if newEntry != nil {
			existing = newEntry
			<-newEntry
		}
	}

	// add to goroutine to remove from the map only
	go func() {
		<-done

		m := <-csm.mapChan

		if m[mapKey] == done {
			delete(m, mapKey)
		}

		csm.mapChan <- m
	}()
}
