package scene

import (
	"os"
	"path/filepath"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/paths"
	"github.com/stashapp/stash/pkg/utils"
)

func MigrateHash(p *paths.Paths, oldHash string, newHash string) {
	oldPath := filepath.Join(p.Generated.Markers, oldHash)
	newPath := filepath.Join(p.Generated.Markers, newHash)
	migrateSceneFiles(oldPath, newPath)

	scenePaths := p.Scene
	oldPath = scenePaths.GetThumbnailScreenshotPath(oldHash)
	newPath = scenePaths.GetThumbnailScreenshotPath(newHash)
	migrateSceneFiles(oldPath, newPath)

	oldPath = scenePaths.GetScreenshotPath(oldHash)
	newPath = scenePaths.GetScreenshotPath(newHash)
	migrateSceneFiles(oldPath, newPath)

	oldPath = scenePaths.GetStreamPreviewPath(oldHash)
	newPath = scenePaths.GetStreamPreviewPath(newHash)
	migrateSceneFiles(oldPath, newPath)

	oldPath = scenePaths.GetStreamPreviewImagePath(oldHash)
	newPath = scenePaths.GetStreamPreviewImagePath(newHash)
	migrateSceneFiles(oldPath, newPath)

	oldPath = scenePaths.GetTranscodePath(oldHash)
	newPath = scenePaths.GetTranscodePath(newHash)
	migrateSceneFiles(oldPath, newPath)

	oldPath = scenePaths.GetSpriteVttFilePath(oldHash)
	newPath = scenePaths.GetSpriteVttFilePath(newHash)
	migrateSceneFiles(oldPath, newPath)

	oldPath = scenePaths.GetSpriteImageFilePath(oldHash)
	newPath = scenePaths.GetSpriteImageFilePath(newHash)
	migrateSceneFiles(oldPath, newPath)

	oldPath = scenePaths.GetInteractiveHeatmapPath(oldHash)
	newPath = scenePaths.GetInteractiveHeatmapPath(newHash)
	migrateSceneFiles(oldPath, newPath)
}

func migrateSceneFiles(oldName, newName string) {
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
