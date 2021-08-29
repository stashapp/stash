package manager

import (
	"os"
	"path/filepath"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/utils"
)

func MigrateHash(oldHash string, newHash string) {
	oldPath := filepath.Join(instance.Paths.Generated.Markers, oldHash)
	newPath := filepath.Join(instance.Paths.Generated.Markers, newHash)
	migrate(oldPath, newPath)

	scenePaths := GetInstance().Paths.Scene
	oldPath = scenePaths.GetThumbnailScreenshotPath(oldHash)
	newPath = scenePaths.GetThumbnailScreenshotPath(newHash)
	migrate(oldPath, newPath)

	oldPath = scenePaths.GetScreenshotPath(oldHash)
	newPath = scenePaths.GetScreenshotPath(newHash)
	migrate(oldPath, newPath)

	oldPath = scenePaths.GetHeatmapPath(oldHash)
	newPath = scenePaths.GetHeatmapPath(newHash)
	migrate(oldPath, newPath)

	oldPath = scenePaths.GetStreamPreviewPath(oldHash)
	newPath = scenePaths.GetStreamPreviewPath(newHash)
	migrate(oldPath, newPath)

	oldPath = scenePaths.GetStreamPreviewImagePath(oldHash)
	newPath = scenePaths.GetStreamPreviewImagePath(newHash)
	migrate(oldPath, newPath)

	oldPath = scenePaths.GetTranscodePath(oldHash)
	newPath = scenePaths.GetTranscodePath(newHash)
	migrate(oldPath, newPath)

	oldPath = scenePaths.GetSpriteVttFilePath(oldHash)
	newPath = scenePaths.GetSpriteVttFilePath(newHash)
	migrate(oldPath, newPath)

	oldPath = scenePaths.GetSpriteImageFilePath(oldHash)
	newPath = scenePaths.GetSpriteImageFilePath(newHash)
	migrate(oldPath, newPath)
}

func migrate(oldName, newName string) {
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
