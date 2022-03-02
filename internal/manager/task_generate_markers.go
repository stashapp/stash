package manager

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type GenerateMarkersTask struct {
	TxnManager          models.TransactionManager
	Scene               *models.Scene
	Marker              *models.SceneMarker
	Overwrite           bool
	fileNamingAlgorithm models.HashAlgorithm

	ImagePreview bool
	Screenshot   bool
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
		t.generateSceneMarkers()
	}

	if t.Marker != nil {
		var scene *models.Scene
		if err := t.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
			var err error
			scene, err = r.Scene().Find(int(t.Marker.SceneID.Int64))
			return err
		}); err != nil {
			logger.Errorf("error finding scene for marker: %s", err.Error())
			return
		}

		if scene == nil {
			logger.Errorf("scene not found for id %d", t.Marker.SceneID.Int64)
			return
		}

		ffprobe := instance.FFProbe
		videoFile, err := ffprobe.NewVideoFile(t.Scene.Path, false)
		if err != nil {
			logger.Errorf("error reading video file: %s", err.Error())
			return
		}

		t.generateMarker(videoFile, scene, t.Marker)
	}
}

func (t *GenerateMarkersTask) generateSceneMarkers() {
	var sceneMarkers []*models.SceneMarker
	if err := t.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		var err error
		sceneMarkers, err = r.SceneMarker().FindBySceneID(t.Scene.ID)
		return err
	}); err != nil {
		logger.Errorf("error getting scene markers: %s", err.Error())
		return
	}

	if len(sceneMarkers) == 0 {
		return
	}

	ffprobe := instance.FFProbe
	videoFile, err := ffprobe.NewVideoFile(t.Scene.Path, false)
	if err != nil {
		logger.Errorf("error reading video file: %s", err.Error())
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

func (t *GenerateMarkersTask) generateMarker(videoFile *ffmpeg.VideoFile, scene *models.Scene, sceneMarker *models.SceneMarker) {
	sceneHash := t.Scene.GetHash(t.fileNamingAlgorithm)
	seconds := int(sceneMarker.Seconds)

	videoExists := t.videoExists(sceneHash, seconds)
	imageExists := !t.ImagePreview || t.imageExists(sceneHash, seconds)
	screenshotExists := !t.Screenshot || t.screenshotExists(sceneHash, seconds)

	baseFilename := strconv.Itoa(seconds)

	options := ffmpeg.SceneMarkerOptions{
		ScenePath: scene.Path,
		Seconds:   seconds,
		Width:     640,
		Audio:     instance.Config.GetPreviewAudio(),
	}

	encoder := instance.FFMPEG

	if t.Overwrite || !videoExists {
		videoFilename := baseFilename + ".mp4"
		videoPath := instance.Paths.SceneMarkers.GetStreamPath(sceneHash, seconds)

		options.OutputPath = instance.Paths.Generated.GetTmpPath(videoFilename) // tmp output in case the process ends abruptly
		if err := encoder.SceneMarkerVideo(*videoFile, options); err != nil {
			logger.Errorf("[generator] failed to generate marker video: %s", err)
		} else {
			_ = fsutil.SafeMove(options.OutputPath, videoPath)
			logger.Debug("created marker video: ", videoPath)
		}
	}

	if t.ImagePreview && (t.Overwrite || !imageExists) {
		imageFilename := baseFilename + ".webp"
		imagePath := instance.Paths.SceneMarkers.GetStreamPreviewImagePath(sceneHash, seconds)

		options.OutputPath = instance.Paths.Generated.GetTmpPath(imageFilename) // tmp output in case the process ends abruptly
		if err := encoder.SceneMarkerImage(*videoFile, options); err != nil {
			logger.Errorf("[generator] failed to generate marker image: %s", err)
		} else {
			_ = fsutil.SafeMove(options.OutputPath, imagePath)
			logger.Debug("created marker image: ", imagePath)
		}
	}

	if t.Screenshot && (t.Overwrite || !screenshotExists) {
		screenshotFilename := baseFilename + ".jpg"
		screenshotPath := instance.Paths.SceneMarkers.GetStreamScreenshotPath(sceneHash, seconds)

		screenshotOptions := ffmpeg.ScreenshotOptions{
			OutputPath: instance.Paths.Generated.GetTmpPath(screenshotFilename), // tmp output in case the process ends abruptly
			Quality:    2,
			Width:      videoFile.Width,
			Time:       float64(seconds),
		}
		if err := encoder.Screenshot(*videoFile, screenshotOptions); err != nil {
			logger.Errorf("[generator] failed to generate marker screenshot: %s", err)
		} else {
			_ = fsutil.SafeMove(screenshotOptions.OutputPath, screenshotPath)
			logger.Debug("created marker screenshot: ", screenshotPath)
		}
	}
}

func (t *GenerateMarkersTask) markersNeeded() int {
	markers := 0
	var sceneMarkers []*models.SceneMarker
	if err := t.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		var err error
		sceneMarkers, err = r.SceneMarker().FindBySceneID(t.Scene.ID)
		return err
	}); err != nil {
		logger.Errorf("errror finding scene markers: %s", err.Error())
		return 0
	}

	if len(sceneMarkers) == 0 {
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

	videoPath := instance.Paths.SceneMarkers.GetStreamPath(sceneChecksum, seconds)
	videoExists, _ := fsutil.FileExists(videoPath)

	return videoExists
}

func (t *GenerateMarkersTask) imageExists(sceneChecksum string, seconds int) bool {
	if sceneChecksum == "" {
		return false
	}

	imagePath := instance.Paths.SceneMarkers.GetStreamPreviewImagePath(sceneChecksum, seconds)
	imageExists, _ := fsutil.FileExists(imagePath)

	return imageExists
}

func (t *GenerateMarkersTask) screenshotExists(sceneChecksum string, seconds int) bool {
	if sceneChecksum == "" {
		return false
	}

	screenshotPath := instance.Paths.SceneMarkers.GetStreamScreenshotPath(sceneChecksum, seconds)
	screenshotExists, _ := fsutil.FileExists(screenshotPath)

	return screenshotExists
}
