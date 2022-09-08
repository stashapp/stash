package manager

import (
	"context"
	"net/http"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
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

type SceneCoverGetter interface {
	GetCover(ctx context.Context, sceneID int) ([]byte, error)
}

type SceneServer struct {
	TxnManager       txn.Manager
	SceneCoverGetter SceneCoverGetter
}

func (s *SceneServer) StreamSceneDirect(scene *models.Scene, w http.ResponseWriter, r *http.Request) {
	fileNamingAlgo := config.GetInstance().GetVideoFileNamingAlgorithm()

	filepath := GetInstance().Paths.Scene.GetStreamPath(scene.Path, scene.GetHash(fileNamingAlgo))
	streamRequestCtx := NewStreamRequestContext(w, r)

	// #2579 - hijacking and closing the connection here causes video playback to fail in Safari
	// We trust that the request context will be closed, so we don't need to call Cancel on the
	// returned context here.
	_ = GetInstance().ReadLockManager.ReadLock(streamRequestCtx, filepath)
	http.ServeFile(w, r, filepath)
}

func (s *SceneServer) ServeScreenshot(scene *models.Scene, w http.ResponseWriter, r *http.Request) {
	filepath := GetInstance().Paths.Scene.GetScreenshotPath(scene.GetHash(config.GetInstance().GetVideoFileNamingAlgorithm()))

	// fall back to the scene image blob if the file isn't present
	screenshotExists, _ := fsutil.FileExists(filepath)
	if screenshotExists {
		http.ServeFile(w, r, filepath)
	} else {
		var cover []byte
		err := txn.WithTxn(r.Context(), s.TxnManager, func(ctx context.Context) error {
			cover, _ = s.SceneCoverGetter.GetCover(ctx, scene.ID)
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
