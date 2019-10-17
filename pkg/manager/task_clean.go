package manager

import (
	"context"
	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type CleanTask struct {
	Scene models.Scene
}

func (t *CleanTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	if t.fileExists(t.Scene.Path) {
		logger.Debugf("Found: %s", t.Scene.Path)
	} else {
		logger.Debugf("Deleting missing file: %s", t.Scene.Path)
		t.deleteScene(t.Scene.ID)
	}
}

func (t *CleanTask) deleteScene(sceneID int) {
	ctx := context.TODO()
	qb := models.NewSceneQueryBuilder()
	jqb := models.NewJoinsQueryBuilder()
	tx := database.DB.MustBeginTx(ctx, nil)
	strSceneID := strconv.Itoa(sceneID)
	defer tx.Commit()

	//check and make sure it still exists. scene is also used to delete generated files
	scene, err := qb.Find(sceneID)
	if err != nil {
		_ = tx.Rollback()
	}

	if err := jqb.DestroyScenesTags(sceneID, tx); err != nil {
		_ = tx.Rollback()
	}

	if err := jqb.DestroyPerformersScenes(sceneID, tx); err != nil {
		_ = tx.Rollback()
	}

	if err := jqb.DestroyScenesMarkers(sceneID, tx); err != nil {
		_ = tx.Rollback()
	}

	if err := jqb.DestroyScenesGalleries(sceneID, tx); err != nil {
		_ = tx.Rollback()
	}

	if err := qb.Destroy(strSceneID, tx); err != nil {
		_ = tx.Rollback()
	}

	t.deleteGeneratedSceneFiles(scene)
}


func (t *CleanTask)  deleteGeneratedSceneFiles(scene *models.Scene) {
	markersFolder := filepath.Join(instance.Paths.Generated.Markers, scene.Checksum)

	exists, _ := utils.FileExists(markersFolder)
	if exists {
		err := os.RemoveAll(markersFolder)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", scene.Path, err.Error())
		}
	}

	thumbPath := instance.Paths.Scene.GetThumbnailScreenshotPath(scene.Checksum)
	exists, _ = utils.FileExists(thumbPath)
	if exists {
		err := os.Remove(thumbPath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", thumbPath, err.Error())
		}
	}

	screenshotPath := instance.Paths.Scene.GetScreenshotPath(scene.Checksum)
	exists, _ = utils.FileExists(screenshotPath)
	if exists {
		err := os.Remove(screenshotPath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", screenshotPath, err.Error())
		}
	}

	streamPreviewPath := instance.Paths.Scene.GetStreamPreviewPath(scene.Checksum)
	exists, _ = utils.FileExists(streamPreviewPath)
	if exists {
		err := os.Remove(streamPreviewPath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", streamPreviewPath, err.Error())
		}
	}

	streamPreviewImagePath := instance.Paths.Scene.GetStreamPreviewImagePath(scene.Checksum)
	exists, _ = utils.FileExists(streamPreviewImagePath)
	if exists {
		err := os.Remove(streamPreviewImagePath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", streamPreviewImagePath, err.Error())
		}
	}

	transcodePath := instance.Paths.Scene.GetTranscodePath(scene.Checksum)
	exists, _ = utils.FileExists(transcodePath)
	if exists {
		err := os.Remove(transcodePath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", transcodePath, err.Error())
		}
	}

	spritePath := instance.Paths.Scene.GetSpriteImageFilePath(scene.Checksum)
	exists, _ = utils.FileExists(spritePath)
	if exists {
		err := os.Remove(spritePath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", spritePath, err.Error())
		}
	}

	vttPath := instance.Paths.Scene.GetSpriteVttFilePath(scene.Checksum)
	exists, _ = utils.FileExists(vttPath)
	if exists {
		err := os.Remove(vttPath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", vttPath, err.Error())
		}
	}
}

func (t *CleanTask) fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
