package manager

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

// MigrateHashTask renames generated files between oshash and MD5 based on the
// value of the fileNamingAlgorithm flag.
type MigrateHashTask struct {
	Scene               *models.Scene
	fileNamingAlgorithm models.HashAlgorithm
}

// Start starts the task.
func (t *MigrateHashTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	if !t.Scene.OSHash.Valid || !t.Scene.Checksum.Valid {
		// nothing to do
		return
	}

	oshash := t.Scene.OSHash.String
	checksum := t.Scene.Checksum.String

	oldHash := oshash
	newHash := checksum
	if t.fileNamingAlgorithm == models.HashAlgorithmOshash {
		oldHash = checksum
		newHash = oshash
	}

	oldPath := filepath.Join(instance.Paths.Generated.Markers, oldHash)
	newPath := filepath.Join(instance.Paths.Generated.Markers, newHash)
	t.migrate(oldPath, newPath)

	scenePaths := GetInstance().Paths.Scene
	oldPath = scenePaths.GetThumbnailScreenshotPath(oldHash)
	newPath = scenePaths.GetThumbnailScreenshotPath(newHash)
	t.migrate(oldPath, newPath)

	oldPath = scenePaths.GetScreenshotPath(oldHash)
	newPath = scenePaths.GetScreenshotPath(newHash)
	t.migrate(oldPath, newPath)

	oldPath = scenePaths.GetStreamPreviewPath(oldHash)
	newPath = scenePaths.GetStreamPreviewPath(newHash)
	t.migrate(oldPath, newPath)

	oldPath = scenePaths.GetStreamPreviewImagePath(oldHash)
	newPath = scenePaths.GetStreamPreviewImagePath(newHash)
	t.migrate(oldPath, newPath)

	oldPath = scenePaths.GetTranscodePath(oldHash)
	newPath = scenePaths.GetTranscodePath(newHash)
	t.migrate(oldPath, newPath)

	oldPath = scenePaths.GetSpriteVttFilePath(oldHash)
	newPath = scenePaths.GetSpriteVttFilePath(newHash)
	t.migrate(oldPath, newPath)

	oldPath = scenePaths.GetSpriteImageFilePath(oldHash)
	newPath = scenePaths.GetSpriteImageFilePath(newHash)
	t.migrate(oldPath, newPath)
}

func (t *MigrateHashTask) migrate(oldName, newName string) {
	oldExists, err := utils.FileExists(oldName)
	if err != nil && !os.IsNotExist(err) {
		logger.Errorf("Error checking existence of %s: %s", oldName, err.Error())
		return
	}

	if oldExists {
		logger.Infof("renaming %s to %s", oldName, newName)
		if err := os.Rename(oldName, newName); err != nil {
			logger.Errorf("error renaming %s to %s: %s", oldName, newName, err.Error())
		}
	}
}
