package file

import (
	"context"
	"fmt"
	"time"

	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

// MissingPurger purges files and folders marked as missing and their associated objects from the database.
type MissingPurger struct {
	Repository Repository

	PurgeHandlers []PurgeHandler
	TrashPath     string
}

type purgeMissingJob struct {
	*MissingPurger

	progress *job.Progress
	options  PurgeMissingOptions

	missingFileHandler   missingFileFn
	missingFolderHandler missingFolderFn
}

// PurgeMissingOptions provides options for purging missing files.
type PurgeMissingOptions struct {
	Paths []string

	// DryRun indicates if this is a dry run. A dry run will not make any changes and will only log what would have been done.
	DryRun bool

	// MissingSinceBefore is an optional filter to only purge files and folders that have been marked as missing since before the specified time.
	MissingSinceBefore *time.Time
}

// PurgeMissing starts the purge missing process.
func (s *MissingPurger) PurgeMissing(ctx context.Context, options PurgeMissingOptions, progress *job.Progress) {
	j := &purgeMissingJob{
		MissingPurger: s,
		progress:      progress,
		options:       options,
	}

	if err := j.execute(ctx); err != nil {
		logger.Errorf("error purging missing files: %v", err)
		return
	}
}

func (j *purgeMissingJob) execute(ctx context.Context) error {
	progress := j.progress

	var (
		fileCount   int
		folderCount int
	)

	r := j.Repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		var err error
		fileCount, err = r.File.CountMissingInPaths(ctx, j.options.Paths, j.options.MissingSinceBefore)
		if err != nil {
			return err
		}

		folderCount, err = r.Folder.CountMissingInPaths(ctx, j.options.Paths, j.options.MissingSinceBefore)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	progress.AddTotal(fileCount + folderCount)
	progress.Definite()

	if err := j.purgeMissingFiles(ctx, progress); err != nil {
		return err
	}

	if err := j.purgeMissingFolders(ctx); err != nil {
		return err
	}

	return nil
}

func (j *purgeMissingJob) purgeMissingFiles(ctx context.Context, progress *job.Progress) error {
	const batchSize = 1000

	offset := 0
	more := true
	r := j.Repository

	for more {
		var files []models.File

		if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
			if job.IsCancelled(ctx) {
				return nil
			}

			var err error
			files, err = r.File.FindMissingInPaths(ctx, j.options.Paths, j.options.MissingSinceBefore, batchSize, offset)
			if err != nil {
				return fmt.Errorf("error querying for files: %w", err)
			}

			return nil
		}); err != nil {
			return err
		}

		for _, f := range files {
			if j.options.DryRun {
				logger.Infof("Would delete file %q from database", f.Base().Path)
			} else {
				logger.Infof("Deleting file %q from database", f.Base().Path)
				j.deleteFile(ctx, f.Base().ID, f.Base().Path)
			}
			progress.Increment()
		}

		if len(files) != batchSize {
			more = false
		} else if j.options.DryRun {
			// when not in dry run, we should be continuing until there's none left
			// in dry run, we can just increment the offset and continue to the next batch
			offset += batchSize
		}
	}

	return nil
}

func (j *purgeMissingJob) purgeMissingFolders(ctx context.Context) error {
	const batchSize = 1000
	offset := 0
	progress := j.progress

	more := true
	r := j.Repository
	for more {
		if job.IsCancelled(ctx) {
			return nil
		}

		var folders []*models.Folder

		if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
			var err error
			folders, err = r.Folder.FindMissingInPaths(ctx, j.options.Paths, j.options.MissingSinceBefore, batchSize, offset)
			if err != nil {
				return fmt.Errorf("error querying for folders: %w", err)
			}

			return nil
		}); err != nil {
			return err
		}

		for _, f := range folders {
			if j.options.DryRun {
				logger.Infof("Would delete folder %q from database", f.Path)
			} else {
				logger.Infof("Deleting folder %q from database", f.Path)
				j.deleteFolder(ctx, f.ID, f.Path)
			}
			progress.Increment()
		}

		if len(folders) != batchSize {
			more = false
		} else if j.options.DryRun {
			// when not in dry run, we should be continuing until there's none left
			// in dry run, we can just increment the offset and continue to the next batch
			offset += batchSize
		}
	}

	return nil
}

func (j *purgeMissingJob) deleteFile(ctx context.Context, fileID models.FileID, fn string) {
	// delete associated objects
	fileDeleter := NewDeleterWithTrash(j.TrashPath)
	r := j.Repository
	if err := r.WithTxn(ctx, func(ctx context.Context) error {
		fileDeleter.RegisterHooks(ctx)

		if err := j.fireHandlers(ctx, fileDeleter, fileID); err != nil {
			return err
		}

		return r.File.Destroy(ctx, fileID)
	}); err != nil {
		logger.Errorf("Error deleting file %q from database: %s", fn, err.Error())
		return
	}
}

func (j *purgeMissingJob) deleteFolder(ctx context.Context, folderID models.FolderID, fn string) {
	// delete associated objects
	fileDeleter := NewDeleterWithTrash(j.TrashPath)
	r := j.Repository
	if err := r.WithTxn(ctx, func(ctx context.Context) error {
		fileDeleter.RegisterHooks(ctx)

		if err := j.fireFolderHandlers(ctx, fileDeleter, folderID); err != nil {
			return err
		}

		return r.Folder.Destroy(ctx, folderID)
	}); err != nil {
		logger.Errorf("Error deleting folder %q from database: %s", fn, err.Error())
		return
	}
}

func (j *purgeMissingJob) fireHandlers(ctx context.Context, fileDeleter *Deleter, fileID models.FileID) error {
	for _, h := range j.PurgeHandlers {
		if err := h.HandleFile(ctx, fileDeleter, fileID); err != nil {
			return err
		}
	}

	return nil
}

func (j *purgeMissingJob) fireFolderHandlers(ctx context.Context, fileDeleter *Deleter, folderID models.FolderID) error {
	for _, h := range j.PurgeHandlers {
		if err := h.HandleFolder(ctx, fileDeleter, folderID); err != nil {
			return err
		}
	}

	return nil
}
