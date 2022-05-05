package manager

import (
	"context"
	"net/http"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type StreamRequestContext struct {
	context.Context
	ResponseWriter http.ResponseWriter
}

func NewStreamRequestContext(w http.ResponseWriter, r *http.Request) *StreamRequestContext {
	return &StreamRequestContext{
		Context:        r.Context(),
		ResponseWriter: w,
	}
}

func (c *StreamRequestContext) Cancel() {
	hj, ok := (c.ResponseWriter).(http.Hijacker)
	if !ok {
		return
	}

	// hijack and close the connection
	conn, _, _ := hj.Hijack()
	if conn != nil {
		conn.Close()
	}
}

func KillRunningStreams(scene *models.Scene, fileNamingAlgo models.HashAlgorithm) {
	instance.ReadLockManager.Cancel(scene.Path)

	sceneHash := scene.GetHash(fileNamingAlgo)

	if sceneHash == "" {
		return
	}

	transcodePath := GetInstance().Paths.Scene.GetTranscodePath(sceneHash)
	instance.ReadLockManager.Cancel(transcodePath)
}

type SceneServer struct {
	TXNManager models.TransactionManager
}

func (s *SceneServer) StreamSceneDirect(scene *models.Scene, w http.ResponseWriter, r *http.Request) {
	fileNamingAlgo := config.GetInstance().GetVideoFileNamingAlgorithm()

	filepath := GetInstance().Paths.Scene.GetStreamPath(scene.Path, scene.GetHash(fileNamingAlgo))
	streamRequestCtx := NewStreamRequestContext(w, r)
	lockCtx := GetInstance().ReadLockManager.ReadLock(streamRequestCtx, filepath)
	defer lockCtx.Cancel()

	setContentType(scene, w)
	http.ServeFile(w, r, filepath)
}

func setContentType(scene *models.Scene, w http.ResponseWriter) {
	// it's possible that the file extension does not match the required content type
	// set it here. If no mime type is matched, it will defer to original behaviour
	container := ffmpeg.Container(scene.Format.String)
	mimeType := container.MimeType()

	if mimeType != "" {
		w.Header().Set("Content-Type", mimeType)
	}
}

func (s *SceneServer) ServeScreenshot(scene *models.Scene, w http.ResponseWriter, r *http.Request) {
	filepath := GetInstance().Paths.Scene.GetScreenshotPath(scene.GetHash(config.GetInstance().GetVideoFileNamingAlgorithm()))

	// fall back to the scene image blob if the file isn't present
	screenshotExists, _ := fsutil.FileExists(filepath)
	if screenshotExists {
		http.ServeFile(w, r, filepath)
	} else {
		var cover []byte
		err := s.TXNManager.WithReadTxn(r.Context(), func(repo models.ReaderRepository) error {
			cover, _ = repo.Scene().GetCover(scene.ID)
			return nil
		})
		if err != nil {
			logger.Warnf("read transaction failed while serving screenshot: %v", err)
		}

		if err = utils.ServeImage(cover, w, r); err != nil {
			logger.Warnf("unable to serve screenshot image: %v", err)
		}
	}
}
