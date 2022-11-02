package api

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
)

func (r *mutationResolver) DeleteFiles(ctx context.Context, ids []string) (ret bool, err error) {
	fileIDs, err := stringslice.StringSliceToIntSlice(ids)
	if err != nil {
		return false, err
	}

	fileDeleter := file.NewDeleter()
	destroyer := &file.ZipDestroyer{
		FileDestroyer:   r.repository.File,
		FolderDestroyer: r.repository.Folder,
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.File

		for _, fileIDInt := range fileIDs {
			fileID := file.ID(fileIDInt)
			f, err := qb.Find(ctx, fileID)
			if err != nil {
				return err
			}

			path := f[0].Base().Path

			// ensure not a primary file
			isPrimary, err := qb.IsPrimary(ctx, fileID)
			if err != nil {
				return fmt.Errorf("checking if file %s is primary: %w", path, err)
			}

			if isPrimary {
				return fmt.Errorf("cannot delete primary file %s", path)
			}

			// destroy files in zip file
			inZip, err := qb.FindByZipFileID(ctx, fileID)
			if err != nil {
				return fmt.Errorf("finding zip file contents for %s: %w", path, err)
			}

			for _, ff := range inZip {
				const deleteFileInZip = false
				if err := file.Destroy(ctx, qb, ff, fileDeleter, deleteFileInZip); err != nil {
					return fmt.Errorf("destroying file %s: %w", ff.Base().Path, err)
				}
			}

			const deleteFile = true
			if err := destroyer.DestroyZip(ctx, f[0], fileDeleter, deleteFile); err != nil {
				return fmt.Errorf("deleting file %s: %w", path, err)
			}
		}

		return nil
	}); err != nil {
		fileDeleter.Rollback()
		return false, err
	}

	// perform the post-commit actions
	fileDeleter.Commit()

	return true, nil
}
