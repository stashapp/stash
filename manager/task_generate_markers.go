package manager

import (
	"github.com/stashapp/stash/ffmpeg"
	"github.com/stashapp/stash/logger"
	"github.com/stashapp/stash/models"
	"github.com/stashapp/stash/utils"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type GenerateMarkersTask struct {
	Scene models.Scene
}

func (t *GenerateMarkersTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	instance.Paths.Generated.EmptyTmpDir()
	qb := models.NewSceneMarkerQueryBuilder()
	sceneMarkers, _ := qb.FindBySceneID(t.Scene.ID, nil)
	if len(sceneMarkers) == 0 {
		return
	}

	videoFile, err := ffmpeg.NewVideoFile(instance.Paths.FixedPaths.FFProbe, t.Scene.Path)
	if err != nil {
		logger.Errorf("error reading video file: %s", err.Error())
		return
	}

	// Make the folder for the scenes markers
	markersFolder := filepath.Join(instance.Paths.Generated.Markers, t.Scene.Checksum)
	_ = utils.EnsureDir(markersFolder)

	encoder := ffmpeg.NewEncoder(instance.Paths.FixedPaths.FFMPEG)
	for i, sceneMarker := range sceneMarkers {
		index := i + 1
		logger.Progressf("[generator] <%s> scene marker %d of %d", t.Scene.Checksum, index, len(sceneMarkers))

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
			Seconds: seconds,
			Width: 640,
		}
		if !videoExists {
			options.OutputPath = instance.Paths.Generated.GetTmpPath(videoFilename) // tmp output in case the process ends abruptly
			if err := encoder.SceneMarkerVideo(*videoFile, options); err != nil {
				logger.Errorf("[generator] failed to generate marker video: %s", err)
			} else {
				_ = os.Rename(options.OutputPath, videoPath)
				logger.Debug("created marker video: ", videoPath)
			}
		}

		if !imageExists {
			options.OutputPath = instance.Paths.Generated.GetTmpPath(imageFilename) // tmp output in case the process ends abruptly
			if err := encoder.SceneMarkerImage(*videoFile, options); err != nil {
				logger.Errorf("[generator] failed to generate marker image: %s", err)
			} else {
				_ = os.Rename(options.OutputPath, imagePath)
				logger.Debug("created marker image: ", videoPath)
			}
		}
	}
}
