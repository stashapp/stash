package manager

import (
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type GenerateMarkersTask struct {
	Scene models.Scene
}

func (t *GenerateMarkersTask) Start(wg *sync.WaitGroup, markersCh chan<- struct{}, errorCh chan<- struct{}) {
	defer wg.Done()

	qb := models.NewSceneMarkerQueryBuilder()
	sceneMarkers, _ := qb.FindBySceneID(t.Scene.ID, nil)
	if len(sceneMarkers) == 0 {
		return
	}

	videoFile, err := ffmpeg.NewVideoFile(instance.FFProbePath, t.Scene.Path)
	if err != nil {
		logger.Errorf("error reading video file: %s", err.Error())
		return
	}

	// Make the folder for the scenes markers
	markersFolder := filepath.Join(instance.Paths.Generated.Markers, t.Scene.Checksum)
	_ = utils.EnsureDir(markersFolder)

	encoder := ffmpeg.NewEncoder(instance.FFMPEGPath)

	for _, sceneMarker := range sceneMarkers {
		createMarker := false
		errorMarker := false
		seconds := int(sceneMarker.Seconds)
		baseFilename := strconv.Itoa(seconds)
		videoFilename := baseFilename + ".mp4"
		imageFilename := baseFilename + ".webp"
		videoPath := instance.Paths.SceneMarkers.GetStreamPath(t.Scene.Checksum, seconds)
		imagePath := instance.Paths.SceneMarkers.GetStreamPreviewImagePath(t.Scene.Checksum, seconds)
		videoExists, _ := utils.FileExists(videoPath)
		imageExists, _ := utils.FileExists(imagePath)

		options := ffmpeg.SceneMarkerOptions{
			ScenePath: t.Scene.Path,
			Seconds:   seconds,
			Width:     640,
		}
		if !videoExists {
			logger.Infof("[generator] generating scene video markers for %s", t.Scene.Path)
			options.OutputPath = instance.Paths.Generated.GetTmpPath(videoFilename) // tmp output in case the process ends abruptly
			if err := encoder.SceneMarkerVideo(*videoFile, options); err != nil {
				logger.Errorf("[generator] failed to generate marker video: %s", err)
				errorMarker = true
			} else {
				_ = os.Rename(options.OutputPath, videoPath)
				logger.Debug("created marker video: ", videoPath)
				createMarker = true
			}
		}

		if !imageExists {
			logger.Infof("[generator] generating scene image markers for %s", t.Scene.Path)
			options.OutputPath = instance.Paths.Generated.GetTmpPath(imageFilename) // tmp output in case the process ends abruptly
			if err := encoder.SceneMarkerImage(*videoFile, options); err != nil {
				logger.Errorf("[generator] failed to generate marker image: %s", err)
				errorMarker = true
			} else {
				_ = os.Rename(options.OutputPath, imagePath)
				logger.Debug("created marker image: ", videoPath)
				createMarker = true
			}
		}
		if createMarker { // marker created either image or video
			markersCh <- struct{}{}
		} else if errorMarker { // marker creation failed , error found
			errorCh <- struct{}{}
		}
	}
}

func (t *GenerateMarkersTask) isMarkerNeeded() int {

	markers := 0
	qb := models.NewSceneMarkerQueryBuilder()
	sceneMarkers, _ := qb.FindBySceneID(t.Scene.ID, nil)
	if len(sceneMarkers) == 0 {
		return 0
	}

	for _, sceneMarker := range sceneMarkers {
		seconds := int(sceneMarker.Seconds)
		videoPath := instance.Paths.SceneMarkers.GetStreamPath(t.Scene.Checksum, seconds)
		imagePath := instance.Paths.SceneMarkers.GetStreamPreviewImagePath(t.Scene.Checksum, seconds)
		videoExists, _ := utils.FileExists(videoPath)
		imageExists, _ := utils.FileExists(imagePath)

		if (!videoExists) || (!imageExists) {
			markers++
		}

	}
	return markers
}
