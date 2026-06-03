package audio

import (
	"os"

	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models/paths"
)

func MigrateHash(p *paths.Paths, oldHash string, newHash string) {
	audioPaths := p.Audio

	oldPath := audioPaths.GetTranscodePath(oldHash)
	newPath := audioPaths.GetTranscodePath(newHash)
	migrateAudioFiles(oldPath, newPath)
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
