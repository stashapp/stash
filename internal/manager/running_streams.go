package manager

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/internal/static"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
	"github.com/stashapp/stash/pkg/utils"
)

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
	streamRequestCtx := ffmpeg.NewStreamRequestContext(w, r)

	// #2579 - hijacking and closing the connection here causes video playback to fail in Safari
	// We trust that the request context will be closed, so we don't need to call Cancel on the
	// returned context here.
	_ = GetInstance().ReadLockManager.ReadLock(streamRequestCtx, filepath)
	http.ServeFile(w, r, filepath)
}

func (s *SceneServer) ServeScreenshot(scene *models.Scene, w http.ResponseWriter, r *http.Request) {
	const defaultSceneImage = "scene/scene.svg"

	var cover []byte
	readTxnErr := txn.WithReadTxn(r.Context(), s.TxnManager, func(ctx context.Context) error {
		cover, _ = s.SceneCoverGetter.GetCover(ctx, scene.ID)
		return nil
	})
	if errors.Is(readTxnErr, context.Canceled) {
		return
	}
	if readTxnErr != nil {
		logger.Warnf("read transaction error on fetch screenshot: %v", readTxnErr)
	}

	if cover == nil {
		// fallback to legacy image if present
		if scene.Path != "" {
			filepath := GetInstance().Paths.Scene.GetLegacyScreenshotPath(scene.GetHash(config.GetInstance().GetVideoFileNamingAlgorithm()))

			// fall back to the scene image blob if the file isn't present
			screenshotExists, _ := fsutil.FileExists(filepath)
			if screenshotExists {
				http.ServeFile(w, r, filepath)
				return
			}
		}

		// fallback to default cover if none found
		// should always be there
		f, _ := static.Scene.Open(defaultSceneImage)
		defer f.Close()
		stat, _ := f.Stat()
		http.ServeContent(w, r, "scene.svg", stat.ModTime(), f.(io.ReadSeeker))
	}

	if err := utils.ServeImage(cover, w, r); err != nil {
		logger.Warnf("error serving screenshot image: %v", err)
	}
}
