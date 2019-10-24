package manager

import (
	"net/http"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
)

var streamingFiles = make(map[string][]*http.ResponseWriter)

func RegisterStream(filepath string, w *http.ResponseWriter) {
	streams := streamingFiles[filepath]

	streamingFiles[filepath] = append(streams, w)
}

func deregisterStream(filepath string, w *http.ResponseWriter) {
	streams := streamingFiles[filepath]

	for i, v := range streams {
		if v == w {
			streamingFiles[filepath] = append(streams[:i], streams[i+1:]...)
			return
		}
	}
}

func WaitAndDeregisterStream(filepath string, w *http.ResponseWriter, r *http.Request) {
	notify := r.Context().Done()
	go func() {
		<-notify
		deregisterStream(filepath, w)
	}()
}

func KillRunningStreams(path string) {
	ffmpeg.KillRunningEncoders(path)

	streams := streamingFiles[path]

	for _, w := range streams {
		hj, ok := (*w).(http.Hijacker)
		if !ok {
			// if we can't close the connection can't really do anything else
			logger.Warnf("cannot close running stream for: %s", path)
			return
		}

		// hijack and close the connection
		conn, _, err := hj.Hijack()
		if err != nil {
			logger.Errorf("cannot close running stream for '%s' due to error: %s", path, err.Error())
		} else {
			conn.Close()
		}
	}
}
