package manager

import (
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/stashapp/stash/pkg/hash"
	"github.com/stashapp/stash/pkg/logger"
)

// DownloadStore manages single-use generated files for the UI to download.
type DownloadStore struct {
	m     map[string]*storeFile
	mutex sync.Mutex
}

type storeFile struct {
	path        string
	contentType string
	keep        bool
	wg          sync.WaitGroup
	once        sync.Once
}

func NewDownloadStore() *DownloadStore {
	return &DownloadStore{
		m: make(map[string]*storeFile),
	}
}

func (s *DownloadStore) RegisterFile(fp string, contentType string, keep bool) (string, error) {
	const keyLength = 4
	const attempts = 100

	// keep generating random keys until we get a free one
	// prevent infinite loop by only attempting a finite amount of times
	var h string
	generate := true
	a := 0

	s.mutex.Lock()
	for generate && a < attempts {
		var err error
		h, err = hash.GenerateRandomKey(keyLength)
		if err != nil {
			return "", err
		}
		_, generate = s.m[h]
		a++
	}

	s.m[h] = &storeFile{
		path:        fp,
		contentType: contentType,
		keep:        keep,
	}
	s.mutex.Unlock()

	return h, nil
}

func (s *DownloadStore) Serve(hash string, w http.ResponseWriter, r *http.Request) {
	s.mutex.Lock()
	f, ok := s.m[hash]

	if !ok {
		s.mutex.Unlock()
		http.NotFound(w, r)
		return
	}

	if !f.keep {
		s.waitAndRemoveFile(hash, &w, r)
	}

	s.mutex.Unlock()

	if f.contentType != "" {
		w.Header().Add("Content-Type", f.contentType)
	}
	w.Header().Set("Cache-Control", "no-store")
	http.ServeFile(w, r, f.path)
}

func (s *DownloadStore) waitAndRemoveFile(hash string, w *http.ResponseWriter, r *http.Request) {
	f := s.m[hash]
	notify := r.Context().Done()
	f.wg.Add(1)

	go func() {
		<-notify
		s.mutex.Lock()
		defer s.mutex.Unlock()

		f.wg.Done()
	}()

	go f.once.Do(func() {
		// leave it up for 30 seconds after the first request to allow for multiple requests
		time.Sleep(30 * time.Second)

		f.wg.Wait()
		s.mutex.Lock()
		defer s.mutex.Unlock()

		delete(s.m, hash)
		err := os.Remove(f.path)
		if err != nil {
			logger.Errorf("error removing %s after downloading: %s", f.path, err.Error())
		}
	})
}
