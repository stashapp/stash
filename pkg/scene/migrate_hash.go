package scene

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models/paths"
)

func MigrateHash(p *paths.Paths, oldHash string, newHash string) {
	oldPath := filepath.Join(p.Generated.Markers, oldHash)
	newPath := filepath.Join(p.Generated.Markers, newHash)
	migrateSceneFiles(oldPath, newPath)

	scenePaths := p.Scene
	oldPath = scenePaths.GetVideoPreviewPath(oldHash)
	newPath = scenePaths.GetVideoPreviewPath(newHash)
	migrateSceneFiles(oldPath, newPath)

	oldPath = scenePaths.GetWebpPreviewPath(oldHash)
	newPath = scenePaths.GetWebpPreviewPath(newHash)
	migrateSceneFiles(oldPath, newPath)

	oldPath = scenePaths.GetTranscodePath(oldHash)
	newPath = scenePaths.GetTranscodePath(newHash)
	migrateSceneFiles(oldPath, newPath)

	oldVttPath := scenePaths.GetSpriteVttFilePath(oldHash)
	newVttPath := scenePaths.GetSpriteVttFilePath(newHash)
	migrateSceneFiles(oldVttPath, newVttPath)

	oldPath = scenePaths.GetSpriteImageFilePath(oldHash)
	newPath = scenePaths.GetSpriteImageFilePath(newHash)
	migrateSceneFiles(oldPath, newPath)
	migrateVttFile(newVttPath, oldPath, newPath)

	oldPath = scenePaths.GetInteractiveHeatmapPath(oldHash)
	newPath = scenePaths.GetInteractiveHeatmapPath(newHash)
	migrateSceneFiles(oldPath, newPath)

	// #3986 - migrate scene marker files
	markerPaths := p.SceneMarkers
	oldPath = markerPaths.GetFolderPath(oldHash)
	newPath = markerPaths.GetFolderPath(newHash)
	migrateSceneFolder(oldPath, newPath)
}

func migrateSceneFiles(oldName, newName string) {
	oldExists, err := fsutil.FileExists(oldName)
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

// #2481: migrate vtt file contents in addition to renaming
func migrateVttFile(vttPath, oldSpritePath, newSpritePath string) {
	// #3356 - don't try to migrate if the file doesn't exist
	exists, err := fsutil.FileExists(vttPath)
	if err != nil && !os.IsNotExist(err) {
		logger.Errorf("Error checking existence of %s: %s", vttPath, err.Error())
		return
	}

	if !exists {
		return
	}

	contents, err := os.ReadFile(vttPath)
	if err != nil {
		logger.Errorf("Error reading %s for vtt migration: %v", vttPath, err)
		return
	}

	oldSpriteBasename := filepath.Base(oldSpritePath)
	newSpriteBasename := filepath.Base(newSpritePath)

	contents = bytes.ReplaceAll(contents, []byte(oldSpriteBasename), []byte(newSpriteBasename))

	if err := os.WriteFile(vttPath, contents, 0644); err != nil {
		logger.Errorf("Error writing %s for vtt migration: %v", vttPath, err)
		return
	}
}

func migrateSceneFolder(oldName, newName string) {
	oldExists, err := fsutil.DirExists(oldName)
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
