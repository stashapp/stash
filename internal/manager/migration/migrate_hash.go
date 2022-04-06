package migration

import (
	"context"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/paths"
	"github.com/stashapp/stash/pkg/scene"
)

// MigrateHashJob renames generated files between oshash and MD5 based on the
// value of the fileNamingAlgorithm flag.
type MigrateHashJob struct {
	TxnManager          models.TransactionManager
	Paths               *paths.Paths
	FileNamingAlgorithm models.HashAlgorithm
}

func (j *MigrateHashJob) Execute(ctx context.Context, progress *job.Progress) {
	fileNamingAlgo := config.GetInstance().GetVideoFileNamingAlgorithm()
	logger.Infof("Migrating generated files for %s naming hash", fileNamingAlgo.String())

	var scenes []*models.Scene
	if err := j.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		var err error
		scenes, err = r.Scene().All()
		return err
	}); err != nil {
		logger.Errorf("failed to fetch list of scenes for migration: %s", err.Error())
		return
	}

	total := len(scenes)
	progress.SetTotal(total)

	for _, scene := range scenes {
		progress.Increment()
		if job.IsCancelled(ctx) {
			logger.Info("Stopping due to user request")
			return
		}

		if scene == nil {
			logger.Errorf("nil scene, skipping migrate")
			continue
		}

		j.migrateScene(scene)

	}

	logger.Info("Finished migrating")
}

func (j *MigrateHashJob) migrateScene(s *models.Scene) {
	if !s.OSHash.Valid || !s.Checksum.Valid {
		// nothing to do
		return
	}

	oshash := s.OSHash.String
	checksum := s.Checksum.String

	oldHash := oshash
	newHash := checksum
	if j.FileNamingAlgorithm == models.HashAlgorithmOshash {
		oldHash = checksum
		newHash = oshash
	}

	scene.MigrateHash(j.Paths, oldHash, newHash)
}
