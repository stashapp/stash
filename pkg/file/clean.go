package file

import (
	"context"
	"errors"
	"fmt"
	"io/fs"

	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/txn"
)

// Cleaner scans through stored file and folder instances and removes those that are no longer present on disk.
type Cleaner struct {
	FS         FS
	Repository Repository

	CleanHandlers []CleanHandler
}

type cleanJob struct {
	*Cleaner

	progress *job.Progress
	options  CleanOptions
}

// ScanOptions provides options for scanning files.
type CleanOptions struct {
	Paths []string

	// Do a dry run. Don't delete any files
	DryRun bool

	// PathFilter are used to determine if a file should be included.
	// Excluded files are marked for cleaning.
	PathFilter PathFilter
}

// Clean starts the clean process.
func (s *Cleaner) Clean(ctx context.Context, options CleanOptions, progress *job.Progress) {
	j := &cleanJob{
		Cleaner:  s,
		progress: progress,
		options:  options,
	}

	if err := j.execute(ctx); err != nil {
		logger.Errorf("error cleaning files: %w", err)
		return
	}
}

type deleteSet struct {
	list []ID
	set  map[ID]string
}

func newDeleteSet() deleteSet {
	return deleteSet{
		set: make(map[ID]string),
	}
}

func (s *deleteSet) add(id ID, path string) {
	if _, ok := s.set[id]; !ok {
		s.list = append(s.list, id)
		s.set[id] = path
	}
}

func (s *deleteSet) len() int {
	return len(s.list)
}

func (j *cleanJob) execute(ctx context.Context) error {
	const batchSize = 1000
	offset := 0
	progress := j.progress

	toDelete := newDeleteSet()

	more := true
	if err := txn.WithTxn(ctx, j.Repository, func(ctx context.Context) error {
		for more {
			if job.IsCancelled(ctx) {
				return nil
			}

			files, err := j.Repository.FindAllByPaths(ctx, j.options.Paths, batchSize, offset)
			if err != nil {
				return fmt.Errorf("error querying for files: %w", err)
			}

			for _, f := range files {
				path := f.Base().Path
				err = nil
				progress.ExecuteTask(fmt.Sprintf("Assessing file %s for clean", path), func() {
					if j.shouldClean(ctx, f) {
						// add contained files first
						fileID := f.Base().ID
						var containedFiles []File
						containedFiles, err = j.Repository.FindByZipFileID(ctx, fileID)
						if err != nil {
							err = fmt.Errorf("error finding contained files: %w", err)
							return
						}

						for _, cf := range containedFiles {
							toDelete.add(cf.Base().ID, cf.Base().Path)
						}

						toDelete.add(f.Base().ID, f.Base().Path)
					} else {
						// increment progress, no further processing
						progress.Increment()
					}
				})
				if err != nil {
					return err
				}
			}

			if len(files) != batchSize {
				more = false
			} else {
				offset += batchSize
			}
		}

		return nil
	}); err != nil {
		return err
	}

	if j.options.DryRun && toDelete.len() > 0 {
		// add progress for scenes that would've been deleted
		progress.AddProcessed(toDelete.len())
	}

	if !j.options.DryRun && len(toDelete.list) > 0 {
		progress.ExecuteTask(fmt.Sprintf("Cleaning %d files", toDelete.len()), func() {
			for _, fileID := range toDelete.list {
				if job.IsCancelled(ctx) {
					return
				}

				j.deleteFile(ctx, fileID, toDelete.set[fileID])

				progress.Increment()
			}
		})
	}

	return nil
}

func (j *cleanJob) shouldClean(ctx context.Context, f File) bool {
	path := f.Base().Path

	info, err := f.Base().Info(j.FS)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		logger.Errorf("error getting file info for %q, not cleaning: %w", path, err)
		return false
	}

	if info == nil {
		// info is nil - file not exist
		logger.Infof("File not found. Marking to clean: \"%s\"", path)
		return true
	}

	// run through path filter, if returns false then the file should be cleaned
	filter := j.options.PathFilter

	// don't log anything - assume filter will have logged the reason
	return !filter.Accept(ctx, path, info)
}

func (j *cleanJob) deleteFile(ctx context.Context, fileID ID, fn string) {
	// delete associated objects
	fileDeleter := NewDeleter()
	if err := txn.WithTxn(ctx, j.Repository, func(ctx context.Context) error {
		fileDeleter.RegisterHooks(ctx, j.Repository)

		if err := j.fireHandlers(ctx, fileDeleter, fileID); err != nil {
			return err
		}

		return j.Repository.Destroy(ctx, fileID)
	}); err != nil {
		logger.Errorf("Error deleting file %q from database: %s", fn, err.Error())
		return
	}
}

func (j *cleanJob) fireHandlers(ctx context.Context, fileDeleter *Deleter, fileID ID) error {
	for _, h := range j.CleanHandlers {
		if err := h.Handle(ctx, fileDeleter, fileID); err != nil {
			return err
		}
	}

	return nil
}
