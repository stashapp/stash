package manager

import (
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type GenerateMarkersTask struct {
	Scene  *models.Scene
	Marker *models.SceneMarker
}

func (t *GenerateMarkersTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	if t.Scene != nil {
		t.generateSceneMarkers()
	}

	if t.Marker != nil {
		qb := models.NewSceneQueryBuilder()
		scene, err := qb.Find(int(t.Marker.SceneID.Int64))
		if err != nil {
			logger.Errorf("error finding scene for marker: %s", err.Error())
			return
		}

		videoFile, err := ffmpeg.NewVideoFile(instance.FFProbePath, t.Scene.Path)
		if err != nil {
			logger.Errorf("error reading video file: %s", err.Error())
			return
		}

		t.generateMarker(videoFile, scene, t.Marker)
	}
}

func (t *GenerateMarkersTask) generateSceneMarkers() {
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

	for i, sceneMarker := range sceneMarkers {
		index := i + 1
		logger.Progressf("[generator] <%s> scene marker %d of %d", t.Scene.Checksum, index, len(sceneMarkers))

		t.generateMarker(videoFile, t.Scene, sceneMarker)
	}
}

func (t *GenerateMarkersTask) generateMarker(videoFile *ffmpeg.VideoFile, scene *models.Scene, sceneMarker *models.SceneMarker) {
	seconds := int(sceneMarker.Seconds)
	baseFilename := strconv.Itoa(seconds)
	videoFilename := baseFilename + ".mp4"
	imageFilename := baseFilename + ".webp"
	videoPath := instance.Paths.SceneMarkers.GetStreamPath(scene.Checksum, seconds)
	imagePath := instance.Paths.SceneMarkers.GetStreamPreviewImagePath(scene.Checksum, seconds)
	videoExists, _ := utils.FileExists(videoPath)
	imageExists, _ := utils.FileExists(imagePath)

	options := ffmpeg.SceneMarkerOptions{
		ScenePath: scene.Path,
		Seconds:   seconds,
		Width:     640,
	}

	encoder := ffmpeg.NewEncoder(instance.FFMPEGPath)

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
