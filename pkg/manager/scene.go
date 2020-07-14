package manager

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

// DestroyScene deletes a scene and its associated relationships from the
// database.
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

// DeleteGeneratedSceneFiles deletes generated files for the provided scene.
func DeleteGeneratedSceneFiles(scene *models.Scene, useMD5 bool) {
	sceneHash := scene.GetHash(useMD5)

	if sceneHash == "" {
		return
	}

	markersFolder := filepath.Join(GetInstance().Paths.Generated.Markers, sceneHash)

	exists, _ := utils.FileExists(markersFolder)
	if exists {
		err := os.RemoveAll(markersFolder)
		if err != nil {
			logger.Warnf("Could not delete folder %s: %s", markersFolder, err.Error())
		}
	}

	thumbPath := GetInstance().Paths.Scene.GetThumbnailScreenshotPath(sceneHash)
	exists, _ = utils.FileExists(thumbPath)
	if exists {
		err := os.Remove(thumbPath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", thumbPath, err.Error())
		}
	}

	normalPath := GetInstance().Paths.Scene.GetScreenshotPath(sceneHash)
	exists, _ = utils.FileExists(normalPath)
	if exists {
		err := os.Remove(normalPath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", normalPath, err.Error())
		}
	}

	streamPreviewPath := GetInstance().Paths.Scene.GetStreamPreviewPath(sceneHash)
	exists, _ = utils.FileExists(streamPreviewPath)
	if exists {
		err := os.Remove(streamPreviewPath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", streamPreviewPath, err.Error())
		}
	}

	streamPreviewImagePath := GetInstance().Paths.Scene.GetStreamPreviewImagePath(sceneHash)
	exists, _ = utils.FileExists(streamPreviewImagePath)
	if exists {
		err := os.Remove(streamPreviewImagePath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", streamPreviewImagePath, err.Error())
		}
	}

	transcodePath := GetInstance().Paths.Scene.GetTranscodePath(sceneHash)
	exists, _ = utils.FileExists(transcodePath)
	if exists {
		// kill any running streams
		KillRunningStreams(transcodePath)

		err := os.Remove(transcodePath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", transcodePath, err.Error())
		}
	}

	spritePath := GetInstance().Paths.Scene.GetSpriteImageFilePath(sceneHash)
	exists, _ = utils.FileExists(spritePath)
	if exists {
		err := os.Remove(spritePath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", spritePath, err.Error())
		}
	}

	vttPath := GetInstance().Paths.Scene.GetSpriteVttFilePath(sceneHash)
	exists, _ = utils.FileExists(vttPath)
	if exists {
		err := os.Remove(vttPath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", vttPath, err.Error())
		}
	}
}

// DeleteSceneMarkerFiles deletes generated files for a scene marker with the
// provided scene and timestamp.
func DeleteSceneMarkerFiles(scene *models.Scene, seconds int, useMD5 bool) {
	videoPath := GetInstance().Paths.SceneMarkers.GetStreamPath(scene.GetHash(useMD5), seconds)
	imagePath := GetInstance().Paths.SceneMarkers.GetStreamPreviewImagePath(scene.GetHash(useMD5), seconds)

	exists, _ := utils.FileExists(videoPath)
	if exists {
		err := os.Remove(videoPath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", videoPath, err.Error())
		}
	}

	exists, _ = utils.FileExists(imagePath)
	if exists {
		err := os.Remove(imagePath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", videoPath, err.Error())
		}
	}
}

// DeleteSceneFile deletes the scene video file from the filesystem.
func DeleteSceneFile(scene *models.Scene) {
	// kill any running encoders
	KillRunningStreams(scene.Path)

	err := os.Remove(scene.Path)
	if err != nil {
		logger.Warnf("Could not delete file %s: %s", scene.Path, err.Error())
	}
}

func setInitialMD5Config() {
	// if there are no scene files in the database, then default the useMD5
	// and calculateMD5 config settings to false, otherwise set them to true
	// for backwards compatibility purposes
	sqb := models.NewSceneQueryBuilder()
	count, err := sqb.Count()
	if err != nil {
		logger.Errorf("Error while counting scenes: %s", err.Error())
		return
	}

	usingMD5 := count != 0

	viper.SetDefault(config.UseMD5, usingMD5)
	viper.SetDefault(config.CalculateMD5, usingMD5)

	if err := config.Write(); err != nil {
		logger.Errorf("Error while writing configuration file: %s", err.Error())
	}
}
