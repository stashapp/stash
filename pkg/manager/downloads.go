package manager

import (
	"net/http"
	"os"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/utils"
)

// DownloadStore manages single-use generated files for the UI to download.
type DownloadStore map[string]storeFile

type storeFile struct {
	path        string
	contentType string
	keep        bool
}

func NewDownloadStore() DownloadStore {
	return make(DownloadStore)
}

func (s *DownloadStore) RegisterFile(fp string, contentType string, keep bool) string {
	const keyLength = 4
	const attempts = 100

	// keep generating random keys until we get a free one
	// prevent infinite loop by only attempting a finite amount of times
	var hash string
	generate := true
	a := 0

	for generate && a < attempts {
		hash = utils.GenerateRandomKey(keyLength)
		_, generate = (*s)[hash]
		a = a + 1
	}

	(*s)[hash] = storeFile{
		path:        fp,
		contentType: contentType,
		keep:        keep,
	}

	return hash
}

func (s *DownloadStore) Serve(hash string, w http.ResponseWriter, r *http.Request) {
	f, ok := (*s)[hash]
	delete(*s, hash)

	if !ok {
		http.NotFound(w, r)
		return
	}

	if f.contentType != "" {
		w.Header().Add("Content-Type", f.contentType)
	}
	http.ServeFile(w, r, f.path)

	if !f.keep {
		s.waitAndRemoveFile(f.path, &w, r)
	}
}

func (s *DownloadStore) waitAndRemoveFile(filepath string, w *http.ResponseWriter, r *http.Request) {
	notify := r.Context().Done()
	go func() {
		<-notify
		err := os.Remove(filepath)
		if err != nil {
			logger.Errorf("error removing %s after downloading: %s", filepath, err.Error())
		}
	}()
}
