// TODO(audio): update this file

package audio

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
	migrateAudioFiles(oldPath, newPath)

	audioPaths := p.Audio
	oldPath = audioPaths.GetVideoPreviewPath(oldHash)
	newPath = audioPaths.GetVideoPreviewPath(newHash)
	migrateAudioFiles(oldPath, newPath)

	oldPath = audioPaths.GetWebpPreviewPath(oldHash)
	newPath = audioPaths.GetWebpPreviewPath(newHash)
	migrateAudioFiles(oldPath, newPath)

	oldPath = audioPaths.GetTranscodePath(oldHash)
	newPath = audioPaths.GetTranscodePath(newHash)
	migrateAudioFiles(oldPath, newPath)

	oldVttPath := audioPaths.GetSpriteVttFilePath(oldHash)
	newVttPath := audioPaths.GetSpriteVttFilePath(newHash)
	migrateAudioFiles(oldVttPath, newVttPath)

	oldPath = audioPaths.GetSpriteImageFilePath(oldHash)
	newPath = audioPaths.GetSpriteImageFilePath(newHash)
	migrateAudioFiles(oldPath, newPath)
	migrateVttFile(newVttPath, oldPath, newPath)

	oldPath = audioPaths.GetInteractiveHeatmapPath(oldHash)
	newPath = audioPaths.GetInteractiveHeatmapPath(newHash)
	migrateAudioFiles(oldPath, newPath)

	// #3986 - migrate audio marker files
	markerPaths := p.AudioMarkers
	oldPath = markerPaths.GetFolderPath(oldHash)
	newPath = markerPaths.GetFolderPath(newHash)
	migrateAudioFolder(oldPath, newPath)
}

func migrateAudioFiles(oldName, newName string) {
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

func migrateAudioFolder(oldName, newName string) {
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
