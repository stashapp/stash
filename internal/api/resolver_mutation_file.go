package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
)

func (r *mutationResolver) MoveFiles(ctx context.Context, input MoveFilesInput) (bool, error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		fileStore := r.repository.File
		folderStore := r.repository.Folder
		mover := file.NewMover(fileStore, folderStore)
		mover.RegisterHooks(ctx)

		var (
			folder   *models.Folder
			basename string
		)

		fileIDs, err := stringslice.StringSliceToIntSlice(input.Ids)
		if err != nil {
			return fmt.Errorf("converting ids: %w", err)
		}

		switch {
		case input.DestinationFolderID != nil:
			var err error

			folderID, err := strconv.Atoi(*input.DestinationFolderID)
			if err != nil {
				return fmt.Errorf("converting destination folder id: %w", err)
			}

			folder, err = folderStore.Find(ctx, models.FolderID(folderID))
			if err != nil {
				return fmt.Errorf("finding destination folder: %w", err)
			}

			if folder == nil {
				return fmt.Errorf("folder with id %d not found", input.DestinationFolderID)
			}

			if folder.ZipFileID != nil {
				return fmt.Errorf("cannot move to %s, is in a zip file", folder.Path)
			}
		case input.DestinationFolder != nil:
			folderPath := *input.DestinationFolder

			// ensure folder path is within the library
			if err := r.validateFolderPath(folderPath); err != nil {
				return err
			}

			// get or create folder hierarchy
			var err error
			folder, err = file.GetOrCreateFolderHierarchy(ctx, folderStore, folderPath)
			if err != nil {
				return fmt.Errorf("getting or creating folder hierarchy: %w", err)
			}
		default:
			return fmt.Errorf("must specify destination folder or path")
		}

		if input.DestinationBasename != nil {
			// ensure only one file was supplied
			if len(input.Ids) != 1 {
				return fmt.Errorf("must specify one file when providing destination path")
			}

			basename = *input.DestinationBasename
		}

		// create the folder hierarchy in the filesystem if needed
		if err := mover.CreateFolderHierarchy(folder.Path); err != nil {
			return fmt.Errorf("creating folder hierarchy %s in filesystem: %w", folder.Path, err)
		}

		for _, fileIDInt := range fileIDs {
			fileID := models.FileID(fileIDInt)
			f, err := fileStore.Find(ctx, fileID)
			if err != nil {
				return fmt.Errorf("finding file %d: %w", fileID, err)
			}

			// ensure that the file extension matches the existing file type
			if basename != "" {
				if err := r.validateFileExtension(f[0].Base().Basename, basename); err != nil {
					return err
				}
			}

			if err := mover.Move(ctx, f[0], folder, basename); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) validateFolderPath(folderPath string) error {
	paths := manager.GetInstance().Config.GetStashPaths()
	if l := paths.GetStashFromDirPath(folderPath); l == nil {
		return fmt.Errorf("folder path %s must be within a stash library path", folderPath)
	}

	return nil
}

func (r *mutationResolver) validateFileExtension(oldBasename, newBasename string) error {
	c := manager.GetInstance().Config
	if err := r.validateFileExtensionList(c.GetVideoExtensions(), oldBasename, newBasename); err != nil {
		return err
	}

	if err := r.validateFileExtensionList(c.GetImageExtensions(), oldBasename, newBasename); err != nil {
		return err
	}

	if err := r.validateFileExtensionList(c.GetGalleryExtensions(), oldBasename, newBasename); err != nil {
		return err
	}

	return nil
}

func (r *mutationResolver) validateFileExtensionList(exts []string, oldBasename, newBasename string) error {
	if fsutil.MatchExtension(oldBasename, exts) && !fsutil.MatchExtension(newBasename, exts) {
		return fmt.Errorf("file extension for %s is inconsistent with old filename %s", newBasename, oldBasename)
	}

	return nil
}

func (r *mutationResolver) DeleteFiles(ctx context.Context, ids []string) (ret bool, err error) {
	fileIDs, err := stringslice.StringSliceToIntSlice(ids)
	if err != nil {
		return false, fmt.Errorf("converting ids: %w", err)
	}

	trashPath := manager.GetInstance().Config.GetDeleteTrashPath()

	fileDeleter := file.NewDeleterWithTrash(trashPath)
	destroyer := &file.ZipDestroyer{
		FileDestroyer:   r.repository.File,
		FolderDestroyer: r.repository.Folder,
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.File

		for _, fileIDInt := range fileIDs {
			fileID := models.FileID(fileIDInt)
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

func (r *mutationResolver) DestroyFiles(ctx context.Context, ids []string) (ret bool, err error) {
	fileIDs, err := stringslice.StringSliceToIntSlice(ids)
	if err != nil {
		return false, fmt.Errorf("converting ids: %w", err)
	}

	destroyer := &file.ZipDestroyer{
		FileDestroyer:   r.repository.File,
		FolderDestroyer: r.repository.Folder,
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.File

		for _, fileIDInt := range fileIDs {
			fileID := models.FileID(fileIDInt)
			f, err := qb.Find(ctx, fileID)
			if err != nil {
				return err
			}

			if len(f) == 0 {
				return fmt.Errorf("file with id %d not found", fileID)
			}

			path := f[0].Base().Path

			// ensure not a primary file
			isPrimary, err := qb.IsPrimary(ctx, fileID)
			if err != nil {
				return fmt.Errorf("checking if file %s is primary: %w", path, err)
			}

			if isPrimary {
				return fmt.Errorf("cannot destroy primary file entry %s", path)
			}

			// destroy DB entries only (no filesystem deletion)
			const deleteFile = false
			if err := destroyer.DestroyZip(ctx, f[0], nil, deleteFile); err != nil {
				return fmt.Errorf("destroying file entry %s: %w", path, err)
			}
		}

		return nil
	}); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) FileSetFingerprints(ctx context.Context, input FileSetFingerprintsInput) (bool, error) {
	fileIDInt, err := strconv.Atoi(input.ID)
	if err != nil {
		return false, fmt.Errorf("converting id: %w", err)
	}

	fileID := models.FileID(fileIDInt)

	// determine what we're doing
	var (
		fingerprints []models.Fingerprint
		toDelete     []string
	)

	for _, i := range input.Fingerprints {
		if i.Type == models.FingerprintTypeMD5 || i.Type == models.FingerprintTypeOshash {
			return false, fmt.Errorf("cannot modify %s fingerprint", i.Type)
		}

		if i.Value == nil {
			toDelete = append(toDelete, i.Type)
		} else {
			// phashes need to be converted from string into uint64
			var v interface{}
			v = *i.Value

			if i.Type == models.FingerprintTypePhash {
				vInt, err := strconv.ParseUint(*i.Value, 16, 64)
				if err != nil {
					return false, fmt.Errorf("converting phash %s: %w", *i.Value, err)
				}

				v = vInt
			}

			fingerprints = append(fingerprints, models.Fingerprint{
				Type:        i.Type,
				Fingerprint: v,
			})
		}
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.File

		if len(fingerprints) > 0 {
			if err := qb.ModifyFingerprints(ctx, fileID, fingerprints); err != nil {
				return fmt.Errorf("modifying fingerprints: %w", err)
			}
		}

		if len(toDelete) > 0 {
			if err := qb.DestroyFingerprints(ctx, fileID, toDelete); err != nil {
				return fmt.Errorf("destroying fingerprints: %w", err)
			}
		}

		return nil
	}); err != nil {
		return false, err
	}

	return true, nil
}
