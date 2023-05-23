package utils

import (
	"bytes"
	"net/http"
	"time"

	"github.com/stashapp/stash/pkg/hash/md5"
)

// Returns an MD5 hash of data, formatted for use as an HTTP ETag header.
// Intended for use with `http.ServeContent`, to respond to conditional requests.
func GenerateETag(data []byte) string {
	hash := md5.FromBytes(data)
	return `"` + hash + `"`
}

// Serves static content, adding Cache-Control: no-cache and a generated ETag header.
// Responds to conditional requests using the ETag.
func ServeStaticContent(w http.ResponseWriter, r *http.Request, data []byte) {
	if r.URL.Query().Has("t") {
		w.Header().Set("Cache-Control", "private, max-age=31536000, immutable")
	} else {
		w.Header().Set("Cache-Control", "no-cache")
	}
	w.Header().Set("ETag", GenerateETag(data))

	http.ServeContent(w, r, "", time.Time{}, bytes.NewReader(data))
}

// Serves static content at filepath, adding Cache-Control: no-cache.
// Responds to conditional requests using the file modtime.
func ServeStaticFile(w http.ResponseWriter, r *http.Request, filepath string) {
	if r.URL.Query().Has("t") {
		w.Header().Set("Cache-Control", "private, max-age=31536000, immutable")
	} else {
		w.Header().Set("Cache-Control", "no-cache")
	}

	http.ServeFile(w, r, filepath)
}
