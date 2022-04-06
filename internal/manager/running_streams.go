package manager

import (
	"net/http"
	"sync"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

var (
	streamingFiles      = make(map[string][]*http.ResponseWriter)
	streamingFilesMutex = sync.RWMutex{}
)

func RegisterStream(filepath string, w *http.ResponseWriter) {
	streamingFilesMutex.Lock()
	streams := streamingFiles[filepath]
	streamingFiles[filepath] = append(streams, w)
	streamingFilesMutex.Unlock()
}

func deregisterStream(filepath string, w *http.ResponseWriter) {
	streamingFilesMutex.Lock()
	defer streamingFilesMutex.Unlock()
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

func KillRunningStreams(scene *models.Scene, fileNamingAlgo models.HashAlgorithm) {
	killRunningStreams(scene.Path)

	sceneHash := scene.GetHash(fileNamingAlgo)

	if sceneHash == "" {
		return
	}

	transcodePath := GetInstance().Paths.Scene.GetTranscodePath(sceneHash)
	killRunningStreams(transcodePath)
}

func killRunningStreams(path string) {
	ffmpeg.KillRunningEncoders(path)

	streamingFilesMutex.RLock()
	streams := streamingFiles[path]
	streamingFilesMutex.RUnlock()

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

type SceneServer struct{}

func (s *SceneServer) StreamSceneDirect(scene *models.Scene, w http.ResponseWriter, r *http.Request) {
	fileNamingAlgo := config.GetInstance().GetVideoFileNamingAlgorithm()

	filepath := GetInstance().Paths.Scene.GetStreamPath(scene.Path, scene.GetHash(fileNamingAlgo))
	RegisterStream(filepath, &w)
	http.ServeFile(w, r, filepath)
	WaitAndDeregisterStream(filepath, &w, r)
}

func (s *SceneServer) ServeScreenshot(scene *models.Scene, w http.ResponseWriter, r *http.Request) {
	filepath := GetInstance().Paths.Scene.GetCoverPath(scene.ID)
	http.ServeFile(w, r, filepath)
}
