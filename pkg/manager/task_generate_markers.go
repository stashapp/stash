package manager

import (
	"path/filepath"
	"strconv"

	"github.com/remeh/sizedwaitgroup"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stashapp/stash/pkg/utils"
)

type GenerateMarkersTask struct {
	Scene               *models.Scene
	Marker              *models.SceneMarker
	Overwrite           bool
	fileNamingAlgorithm models.HashAlgorithm
}

func (t *GenerateMarkersTask) Start(wg *sizedwaitgroup.SizedWaitGroup) {
	defer wg.Done()

	if t.Scene != nil {
		t.generateSceneMarkers()
	}

	if t.Marker != nil {
		qb := sqlite.NewSceneQueryBuilder()
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
	qb := sqlite.NewSceneMarkerQueryBuilder()
	sceneMarkers, _ := qb.FindBySceneID(t.Scene.ID, nil)
	if len(sceneMarkers) == 0 {
		return
	}

	videoFile, err := ffmpeg.NewVideoFile(instance.FFProbePath, t.Scene.Path)
	if err != nil {
		logger.Errorf("error reading video file: %s", err.Error())
		return
	}

	sceneHash := t.Scene.GetHash(t.fileNamingAlgorithm)

	// Make the folder for the scenes markers
	markersFolder := filepath.Join(instance.Paths.Generated.Markers, sceneHash)
	utils.EnsureDir(markersFolder)

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
	imageExists := t.imageExists(sceneHash, seconds)

	baseFilename := strconv.Itoa(seconds)

	options := ffmpeg.SceneMarkerOptions{
		ScenePath: scene.Path,
		Seconds:   seconds,
		Width:     640,
	}

	encoder := ffmpeg.NewEncoder(instance.FFMPEGPath)

	if t.Overwrite || !videoExists {
		videoFilename := baseFilename + ".mp4"
		videoPath := instance.Paths.SceneMarkers.GetStreamPath(sceneHash, seconds)

		options.OutputPath = instance.Paths.Generated.GetTmpPath(videoFilename) // tmp output in case the process ends abruptly
		if err := encoder.SceneMarkerVideo(*videoFile, options); err != nil {
			logger.Errorf("[generator] failed to generate marker video: %s", err)
		} else {
			_ = utils.SafeMove(options.OutputPath, videoPath)
			logger.Debug("created marker video: ", videoPath)
		}
	}

	if t.Overwrite || !imageExists {
		imageFilename := baseFilename + ".webp"
		imagePath := instance.Paths.SceneMarkers.GetStreamPreviewImagePath(sceneHash, seconds)

		options.OutputPath = instance.Paths.Generated.GetTmpPath(imageFilename) // tmp output in case the process ends abruptly
		if err := encoder.SceneMarkerImage(*videoFile, options); err != nil {
			logger.Errorf("[generator] failed to generate marker image: %s", err)
		} else {
			_ = utils.SafeMove(options.OutputPath, imagePath)
			logger.Debug("created marker image: ", imagePath)
		}
	}
}

func (t *GenerateMarkersTask) isMarkerNeeded() int {
	markers := 0
	qb := sqlite.NewSceneMarkerQueryBuilder()
	sceneMarkers, _ := qb.FindBySceneID(t.Scene.ID, nil)
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

	videoPath := instance.Paths.SceneMarkers.GetStreamPath(sceneChecksum, seconds)
	imagePath := instance.Paths.SceneMarkers.GetStreamPreviewImagePath(sceneChecksum, seconds)
	videoExists, _ := utils.FileExists(videoPath)
	imageExists, _ := utils.FileExists(imagePath)

	return videoExists && imageExists
}

func (t *GenerateMarkersTask) videoExists(sceneChecksum string, seconds int) bool {
	if sceneChecksum == "" {
		return false
	}

	videoPath := instance.Paths.SceneMarkers.GetStreamPath(sceneChecksum, seconds)
	videoExists, _ := utils.FileExists(videoPath)

	return videoExists
}

func (t *GenerateMarkersTask) imageExists(sceneChecksum string, seconds int) bool {
	if sceneChecksum == "" {
		return false
	}

	imagePath := instance.Paths.SceneMarkers.GetStreamPreviewImagePath(sceneChecksum, seconds)
	imageExists, _ := utils.FileExists(imagePath)

	return imageExists
}
