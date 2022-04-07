package migration

import (
	"context"
	"path/filepath"

	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/paths"
)

type sceneCoverReader interface {
	GetCover(sceneID int) ([]byte, error)
}

// MigrateCoversJob moves scene covers from the database and generated directory into
// the data directory.
type MigrateCoversJob struct {
	TxnManager          models.TransactionManager
	Paths               *paths.Paths
	FileNamingAlgorithm models.HashAlgorithm
}

func (j *MigrateCoversJob) Execute(ctx context.Context, progress *job.Progress) {
	logger.Infof("Migrating scene covers to %s directory", j.Paths.Data.Scenes)

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

	_ = j.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		coverReader, ok := r.Scene().(sceneCoverReader)
		if !ok {
			logger.Errorf("internal error: scene reader is not a scene cover reader")
			return nil
		}

		for _, scene := range scenes {
			progress.Increment()
			if job.IsCancelled(ctx) {
				logger.Info("Stopping due to user request")
				return nil
			}

			if scene == nil {
				logger.Errorf("nil scene, skipping migrate")
				continue
			}

			j.migrateScene(coverReader, scene)
		}

		logger.Info("Finished migrating")
		return nil
	})
}

func (j *MigrateCoversJob) migrateScene(reader sceneCoverReader, s *models.Scene) {
	dest := j.Paths.Scene.GetCoverPath(s.ID)

	// move the cover from the generated directory
	oldGeneratePath := j.getOldGeneratedPath(s)
	if exists, _ := fsutil.FileExists(oldGeneratePath); exists {
		if err := fsutil.SafeMove(oldGeneratePath, dest); err != nil {
			logger.Errorf("failed to move cover from %s to %v: %s", oldGeneratePath, dest, err)
		} else {
			logger.Infof("moved %s to %s for scene %d", oldGeneratePath, dest, s.ID)
		}
	}

	// overwrite the cover from the database, if present
	// ignore errors here, as the database table may not be present
	cover, _ := reader.GetCover(s.ID)
	if cover != nil {
		if err := fsutil.WriteFile(dest, cover); err != nil {
			logger.Errorf("failed to write cover to %s: %v", dest, err)
		} else {
			logger.Infof("wrote %s from cover in database for scene %d", dest, s.ID)
		}
	}
}

func (j *MigrateCoversJob) getOldGeneratedPath(s *models.Scene) string {
	checksum := s.GetHash(j.FileNamingAlgorithm)
	return filepath.Join(j.Paths.Generated.Screenshots, checksum+".jpg")
}

// MigrateTruncateCoversJob truncates the scene covers table.
type MigrateDropCoversJob struct{}

func (j *MigrateDropCoversJob) Execute(ctx context.Context, progress *job.Progress) {
	// requires access to database
	_, err := database.DB.Exec("DROP TABLE _scenes_cover_deprecated")
	if err != nil {
		logger.Errorf("failed to drop scene covers table: %v", err)
		return
	}

	_, err = database.DB.Exec("VACUUM")
	if err != nil {
		logger.Warnf("error while performing vacuum: %v", err)
	}
}
