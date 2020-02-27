package manager

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/jmoiron/sqlx"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func DestroyScene(sceneID int, tx *sqlx.Tx) error {
	qb := models.NewSceneQueryBuilder()
	jqb := models.NewJoinsQueryBuilder()

	_, err := qb.Find(sceneID)
	if err != nil {
		return err
	}

	if err := jqb.DestroyScenesTags(sceneID, tx); err != nil {
		return err
	}

	if err := jqb.DestroyPerformersScenes(sceneID, tx); err != nil {
		return err
	}

	if err := jqb.DestroyScenesMarkers(sceneID, tx); err != nil {
		return err
	}

	if err := jqb.DestroyScenesGalleries(sceneID, tx); err != nil {
		return err
	}

	if err := qb.Destroy(strconv.Itoa(sceneID), tx); err != nil {
		return err
	}

	return nil
}

func DeleteGeneratedSceneFiles(scene *models.Scene) {
	markersFolder := filepath.Join(GetInstance().Paths.Generated.Markers, scene.Checksum)

	exists, _ := utils.FileExists(markersFolder)
	if exists {
		err := os.RemoveAll(markersFolder)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", scene.Path, err.Error())
		}
	}

	thumbPath := GetInstance().Paths.Scene.GetThumbnailScreenshotPath(scene.Checksum)
	exists, _ = utils.FileExists(thumbPath)
	if exists {
		err := os.Remove(thumbPath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", thumbPath, err.Error())
		}
	}

	normalPath := GetInstance().Paths.Scene.GetScreenshotPath(scene.Checksum)
	exists, _ = utils.FileExists(normalPath)
	if exists {
		err := os.Remove(normalPath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", normalPath, err.Error())
		}
	}

	streamPreviewPath := GetInstance().Paths.Scene.GetStreamPreviewPath(scene.Checksum)
	exists, _ = utils.FileExists(streamPreviewPath)
	if exists {
		err := os.Remove(streamPreviewPath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", streamPreviewPath, err.Error())
		}
	}

	streamPreviewImagePath := GetInstance().Paths.Scene.GetStreamPreviewImagePath(scene.Checksum)
	exists, _ = utils.FileExists(streamPreviewImagePath)
	if exists {
		err := os.Remove(streamPreviewImagePath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", streamPreviewImagePath, err.Error())
		}
	}

	transcodePath := GetInstance().Paths.Scene.GetTranscodePath(scene.Checksum)
	exists, _ = utils.FileExists(transcodePath)
	if exists {
		// kill any running streams
		KillRunningStreams(transcodePath)

		err := os.Remove(transcodePath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", transcodePath, err.Error())
		}
	}

	spritePath := GetInstance().Paths.Scene.GetSpriteImageFilePath(scene.Checksum)
	exists, _ = utils.FileExists(spritePath)
	if exists {
		err := os.Remove(spritePath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", spritePath, err.Error())
		}
	}

	vttPath := GetInstance().Paths.Scene.GetSpriteVttFilePath(scene.Checksum)
	exists, _ = utils.FileExists(vttPath)
	if exists {
		err := os.Remove(vttPath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", vttPath, err.Error())
		}
	}
}

func DeleteSceneFile(scene *models.Scene) {
	// kill any running encoders
	KillRunningStreams(scene.Path)

	err := os.Remove(scene.Path)
	if err != nil {
		logger.Warnf("Could not delete file %s: %s", scene.Path, err.Error())
	}
}
