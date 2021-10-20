package manager

import (
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
)

// MigrateHashTask renames generated files between oshash and MD5 based on the
// value of the fileNamingAlgorithm flag.
type MigrateHashTask struct {
	Scene               *models.Scene
	fileNamingAlgorithm models.HashAlgorithm
}

// Start starts the task.
func (t *MigrateHashTask) Start() {
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

	scene.MigrateHash(instance.Paths, oldHash, newHash)
}
