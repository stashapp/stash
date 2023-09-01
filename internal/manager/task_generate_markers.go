package manager

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene/generate"
)

type GenerateMarkersTask struct {
	TxnManager          Repository
	Scene               *models.Scene
	Marker              *models.SceneMarker
	Overwrite           bool
	fileNamingAlgorithm models.HashAlgorithm

	ImagePreview bool
	Screenshot   bool

	generator *generate.Generator
}

func (t *GenerateMarkersTask) GetDescription() string {
	if t.Scene != nil {
		return fmt.Sprintf("Generating markers for %s", t.Scene.Path)
	} else if t.Marker != nil {
		return fmt.Sprintf("Generating marker preview for marker ID %d", t.Marker.ID)
	}

	return "Generating markers"
}

func (t *GenerateMarkersTask) Start(ctx context.Context) {
	if t.Scene != nil {
		t.generateSceneMarkers(ctx)
	}

	if t.Marker != nil {
		var scene *models.Scene
		if err := t.TxnManager.WithReadTxn(ctx, func(ctx context.Context) error {
			var err error
			scene, err = t.TxnManager.Scene.Find(ctx, t.Marker.SceneID)
			if err != nil {
				return err
			}
			if scene == nil {
				return fmt.Errorf("scene with id %d not found", t.Marker.SceneID)
			}

			return scene.LoadPrimaryFile(ctx, t.TxnManager.File)
		}); err != nil {
			logger.Errorf("error finding scene for marker generation: %v", err)
			return
		}

		videoFile := scene.Files.Primary()

		if videoFile == nil {
			// nothing to do
			return
		}

		t.generateMarker(videoFile, scene, t.Marker)
	}
}

func (t *GenerateMarkersTask) generateSceneMarkers(ctx context.Context) {
	var sceneMarkers []*models.SceneMarker
	if err := t.TxnManager.WithReadTxn(ctx, func(ctx context.Context) error {
		var err error
		sceneMarkers, err = t.TxnManager.SceneMarker.FindBySceneID(ctx, t.Scene.ID)
		return err
	}); err != nil {
		logger.Errorf("error getting scene markers: %s", err.Error())
		return
	}

	videoFile := t.Scene.Files.Primary()

	if len(sceneMarkers) == 0 || videoFile == nil {
		return
	}

	sceneHash := t.Scene.GetHash(t.fileNamingAlgorithm)

	// Make the folder for the scenes markers
	markersFolder := filepath.Join(instance.Paths.Generated.Markers, sceneHash)
	if err := fsutil.EnsureDir(markersFolder); err != nil {
		logger.Warnf("could not create the markers folder (%v): %v", markersFolder, err)
	}

	for i, sceneMarker := range sceneMarkers {
		index := i + 1
		logger.Progressf("[generator] <%s> scene marker %d of %d", sceneHash, index, len(sceneMarkers))

		t.generateMarker(videoFile, t.Scene, sceneMarker)
	}
}

func (t *GenerateMarkersTask) generateMarker(videoFile *models.VideoFile, scene *models.Scene, sceneMarker *models.SceneMarker) {
	sceneHash := t.Scene.GetHash(t.fileNamingAlgorithm)
	seconds := int(sceneMarker.Seconds)

	g := t.generator

	if err := g.MarkerPreviewVideo(context.TODO(), videoFile.Path, sceneHash, seconds, instance.Config.GetPreviewAudio()); err != nil {
		logger.Errorf("[generator] failed to generate marker video: %v", err)
		logErrorOutput(err)
	}

	if t.ImagePreview {
		if err := g.SceneMarkerWebp(context.TODO(), videoFile.Path, sceneHash, seconds); err != nil {
			logger.Errorf("[generator] failed to generate marker image: %v", err)
			logErrorOutput(err)
		}
	}

	if t.Screenshot {
		if err := g.SceneMarkerScreenshot(context.TODO(), videoFile.Path, sceneHash, seconds, videoFile.Width); err != nil {
			logger.Errorf("[generator] failed to generate marker screenshot: %v", err)
			logErrorOutput(err)
		}
	}
}

func (t *GenerateMarkersTask) markersNeeded(ctx context.Context) int {
	markers := 0
	sceneMarkers, err := t.TxnManager.SceneMarker.FindBySceneID(ctx, t.Scene.ID)
	if err != nil {
		logger.Errorf("error finding scene markers: %s", err.Error())
		return 0
	}

	if len(sceneMarkers) == 0 || t.Scene.Files.Primary() == nil {
		return 0
	}

	sceneHash := t.Scene.GetHash(t.fileNamingAlgorithm)
	for _, sceneMarker := range sceneMarkers {
		seconds := int(sceneMarker.Seconds)

		if t.Overwrite || !t.markerExists(sceneHash, seconds) {
			markers++
		}
	}

	return markers
}

func (t *GenerateMarkersTask) markerExists(sceneChecksum string, seconds int) bool {
	if sceneChecksum == "" {
		return false
	}

	videoExists := t.videoExists(sceneChecksum, seconds)
	imageExists := !t.ImagePreview || t.imageExists(sceneChecksum, seconds)
	screenshotExists := !t.Screenshot || t.screenshotExists(sceneChecksum, seconds)

	return videoExists && imageExists && screenshotExists
}

func (t *GenerateMarkersTask) videoExists(sceneChecksum string, seconds int) bool {
	if sceneChecksum == "" {
		return false
	}

	videoPath := instance.Paths.SceneMarkers.GetVideoPreviewPath(sceneChecksum, seconds)
	videoExists, _ := fsutil.FileExists(videoPath)

	return videoExists
}

func (t *GenerateMarkersTask) imageExists(sceneChecksum string, seconds int) bool {
	if sceneChecksum == "" {
		return false
	}

	imagePath := instance.Paths.SceneMarkers.GetWebpPreviewPath(sceneChecksum, seconds)
	imageExists, _ := fsutil.FileExists(imagePath)

	return imageExists
}

func (t *GenerateMarkersTask) screenshotExists(sceneChecksum string, seconds int) bool {
	if sceneChecksum == "" {
		return false
	}

	screenshotPath := instance.Paths.SceneMarkers.GetScreenshotPath(sceneChecksum, seconds)
	screenshotExists, _ := fsutil.FileExists(screenshotPath)

	return screenshotExists
}
