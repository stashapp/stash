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
	purger

	progress *job.Progress
	options  PurgeMissingOptions
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
		purger:   purger{Repository: s.Repository, PurgeHandlers: s.PurgeHandlers, TrashPath: s.TrashPath},
		progress: progress,
		options:  options,
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

	if err := j.purgeMissingFiles(ctx); err != nil {
		return err
	}

	if err := j.purgeMissingFolders(ctx); err != nil {
		return err
	}

	return nil
}

func (j *purgeMissingJob) purgeMissingFiles(ctx context.Context) error {
	r := j.Repository
	return forEachBatch(ctx, r, func(ctx context.Context, offset int) ([]models.File, error) {
		files, err := r.File.FindMissingInPaths(ctx, j.options.Paths, j.options.MissingSinceBefore, purgeBatchSize, offset)
		if err != nil {
			return nil, fmt.Errorf("error querying for files: %w", err)
		}
		return files, nil
	}, func(files []models.File) (int, error) {
		failed := 0
		for _, f := range files {
			if !j.purge("file", f.Base().Path, func() bool { return j.deleteFile(ctx, f.Base().ID, f.Base().Path) }) {
				failed++
			}
			j.progress.Increment()
		}
		return j.nextOffset(failed), nil
	})
}

func (j *purgeMissingJob) purgeMissingFolders(ctx context.Context) error {
	r := j.Repository
	return forEachBatch(ctx, r, func(ctx context.Context, offset int) ([]*models.Folder, error) {
		folders, err := r.Folder.FindMissingInPaths(ctx, j.options.Paths, j.options.MissingSinceBefore, purgeBatchSize, offset)
		if err != nil {
			return nil, fmt.Errorf("error querying for folders: %w", err)
		}
		return folders, nil
	}, func(folders []*models.Folder) (int, error) {
		failed := 0
		for _, f := range folders {
			if !j.purge("folder", f.Path, func() bool { return j.deleteFolder(ctx, f.ID, f.Path) }) {
				failed++
			}
			j.progress.Increment()
		}
		return j.nextOffset(failed), nil
	})
}

// purge logs and runs del, or only logs in dry run. Returns false if the delete failed.
func (j *purgeMissingJob) purge(kind, path string, del func() bool) bool {
	if j.options.DryRun {
		logger.Infof("Would delete %s %q from database", kind, path)
		return true
	}

	logger.Infof("Deleting %s %q from database", kind, path)
	return del()
}

// nextOffset returns how far to advance after a full batch.
func (j *purgeMissingJob) nextOffset(failed int) int {
	if j.options.DryRun {
		// nothing is deleted, so page through the results
		return purgeBatchSize
	}
	// deleted rows drop out of the results; skip past rows that failed to delete
	return failed
}
