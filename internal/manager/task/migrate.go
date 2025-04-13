package task

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/sqlite"
)

type migrateJobConfig interface {
	GetBackupDirectoryPath() string
	GetBackupDirectoryPathOrDefault() string
}

type MigrateJob struct {
	BackupPath string
	Config     migrateJobConfig
	Database   *sqlite.Database
}

type databaseSchemaInfo struct {
	CurrentSchemaVersion  uint
	RequiredSchemaVersion uint
	StepsRequired         uint
}

func (s *MigrateJob) Execute(ctx context.Context, progress *job.Progress) error {
	schemaInfo, err := s.required()
	if err != nil {
		return err
	}

	if schemaInfo.StepsRequired == 0 {
		logger.Infof("database is already at the latest schema version")
		return nil
	}

	logger.Infof("Migrating database from %d to %d", schemaInfo.CurrentSchemaVersion, schemaInfo.RequiredSchemaVersion)

	// set the number of tasks = backup + required steps + optimise
	progress.SetTotal(int(schemaInfo.StepsRequired + 2))

	database := s.Database

	// always backup so that we can roll back to the previous version if
	// migration fails
	backupPath := s.BackupPath
	if backupPath == "" {
		backupPath = database.DatabaseBackupPath(s.Config.GetBackupDirectoryPath())
	} else {
		// check if backup path is a filename or path
		// filename goes into backup directory, path is kept as is
		filename := filepath.Base(backupPath)
		if backupPath == filename {
			backupPath = filepath.Join(s.Config.GetBackupDirectoryPathOrDefault(), filename)
		}
	}

	progress.ExecuteTask("Backing up database", func() {
		defer progress.Increment()

		// perform database backup
		err = database.Backup(backupPath)
	})

	if err != nil {
		return fmt.Errorf("error backing up database: %s", err)
	}

	err = s.runMigrations(ctx, progress)

	if err != nil {
		errStr := fmt.Sprintf("error performing migration: %s", err)

		// roll back to the backed up version
		restoreErr := database.RestoreFromBackup(backupPath)
		if restoreErr != nil {
			errStr = fmt.Sprintf("ERROR: unable to restore database from backup after migration failure: %s\n%s", restoreErr.Error(), errStr)
		} else {
			errStr = "An error occurred migrating the database to the latest schema version. The backup database file was automatically renamed to restore the database.\n" + errStr
		}

		return errors.New(errStr)
	}

	// if no backup path was provided, then delete the created backup
	if s.BackupPath == "" {
		if err := os.Remove(backupPath); err != nil {
			logger.Warnf("error removing unwanted database backup (%s): %s", backupPath, err.Error())
		}
	}

	// reinitialise the database
	if err := database.ReInitialise(); err != nil {
		return fmt.Errorf("error reinitialising database: %s", err)
	}

	logger.Infof("Database migration complete")

	return nil
}

func (s *MigrateJob) required() (ret databaseSchemaInfo, err error) {
	database := s.Database

	m, err := sqlite.NewMigrator(database)
	if err != nil {
		return
	}

	defer m.Close()

	ret.CurrentSchemaVersion = m.CurrentSchemaVersion()
	ret.RequiredSchemaVersion = m.RequiredSchemaVersion()

	if ret.RequiredSchemaVersion < ret.CurrentSchemaVersion {
		// shouldn't happen
		return
	}

	ret.StepsRequired = ret.RequiredSchemaVersion - ret.CurrentSchemaVersion
	return
}

func (s *MigrateJob) runMigrations(ctx context.Context, progress *job.Progress) error {
	database := s.Database

	m, err := sqlite.NewMigrator(database)
	if err != nil {
		return err
	}

	defer m.Close()

	logger.Info("Running migrations")

	for {
		currentSchemaVersion := m.CurrentSchemaVersion()
		targetSchemaVersion := m.RequiredSchemaVersion()

		if currentSchemaVersion >= targetSchemaVersion {
			break
		}

		var err error
		progress.ExecuteTask(fmt.Sprintf("Migrating database to schema version %d", currentSchemaVersion+1), func() {
			err = m.RunMigration(ctx, currentSchemaVersion+1)
		})

		if err != nil {
			return fmt.Errorf("error running migration for schema %d: %s", currentSchemaVersion+1, err)
		}

		progress.Increment()
	}

	// perform post-migrate analyze using the migrator connection
	progress.ExecuteTask("Optimising database", func() {
		err = m.PostMigrate(ctx)
		progress.Increment()
	})

	if err != nil {
		return fmt.Errorf("error optimising database: %s", err)
	}

	return nil
}
