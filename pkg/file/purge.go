package file

import (
	"context"

	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

const purgeBatchSize = 1000

// purger deletes files and folders, along with their associated objects, from the database.
type purger struct {
	Repository    Repository
	PurgeHandlers []PurgeHandler
	TrashPath     string
}

// deleteFile returns false if the file could not be deleted.
func (p purger) deleteFile(ctx context.Context, fileID models.FileID, fn string) bool {
	return p.delete(ctx, "file", fn, func(ctx context.Context, d *Deleter) error {
		for _, h := range p.PurgeHandlers {
			if err := h.HandleFile(ctx, d, fileID); err != nil {
				return err
			}
		}
		return p.Repository.File.Destroy(ctx, fileID)
	})
}

// deleteFolder returns false if the folder could not be deleted.
func (p purger) deleteFolder(ctx context.Context, folderID models.FolderID, fn string) bool {
	return p.delete(ctx, "folder", fn, func(ctx context.Context, d *Deleter) error {
		for _, h := range p.PurgeHandlers {
			if err := h.HandleFolder(ctx, d, folderID); err != nil {
				return err
			}
		}
		return p.Repository.Folder.Destroy(ctx, folderID)
	})
}

func (p purger) delete(ctx context.Context, kind, fn string, destroy func(ctx context.Context, d *Deleter) error) bool {
	fileDeleter := NewDeleterWithTrash(p.TrashPath)
	if err := p.Repository.WithTxn(ctx, func(ctx context.Context) error {
		fileDeleter.RegisterHooks(ctx)
		return destroy(ctx, fileDeleter)
	}); err != nil {
		logger.Errorf("Error deleting %s %q from database: %s", kind, fn, err.Error())
		return false
	}

	return true
}

// forEachBatch fetches items in batches of purgeBatchSize until a short batch is returned or the job is cancelled.
// visit returns how far to advance the offset for the next fetch.
func forEachBatch[T any](ctx context.Context, r Repository, fetch func(ctx context.Context, offset int) ([]T, error), visit func(batch []T) (advance int, err error)) error {
	offset := 0
	for {
		if job.IsCancelled(ctx) {
			return nil
		}

		var batch []T
		if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
			var err error
			batch, err = fetch(ctx, offset)
			return err
		}); err != nil {
			return err
		}

		advance, err := visit(batch)
		if err != nil {
			return err
		}

		if len(batch) != purgeBatchSize {
			return nil
		}
		offset += advance
	}
}
