package file

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/models"
)

// GetOrCreateFolderHierarchy gets the folder for the given path, or creates a folder hierarchy for the given path if one if no existing folder is found.
// Does not create any folders in the file system
func GetOrCreateFolderHierarchy(ctx context.Context, fc models.FolderFinderCreator, path string) (*models.Folder, error) {
	// get or create folder hierarchy
	folder, err := fc.FindByPath(ctx, path)
	if err != nil {
		return nil, err
	}

	if folder == nil {
		parentPath := filepath.Dir(path)
		parent, err := GetOrCreateFolderHierarchy(ctx, fc, parentPath)
		if err != nil {
			return nil, err
		}

		now := time.Now()

		folder = &models.Folder{
			Path:           path,
			ParentFolderID: &parent.ID,
			DirEntry:       models.DirEntry{
				// leave mod time empty for now - it will be updated when the folder is scanned
			},
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err = fc.Create(ctx, folder); err != nil {
			return nil, fmt.Errorf("creating folder %s: %w", path, err)
		}
	}

	return folder, nil
}

// TransferZipFolderHierarchy creates the folder hierarchy for zipFileID under newPath, and removes
// ZipFileID from folders under oldPath.
func TransferZipFolderHierarchy(ctx context.Context, folderStore models.FolderReaderWriter, zipFileID models.FileID, oldPath string, newPath string) error {
	zipFolders, err := folderStore.FindByZipFileID(ctx, zipFileID)
	if err != nil {
		return err
	}

	for _, oldFolder := range zipFolders {
		oldZfPath := oldFolder.Path

		// sanity check - ignore folders which aren't under oldPath
		if !strings.HasPrefix(oldZfPath, oldPath) {
			continue
		}

		relZfPath, err := filepath.Rel(oldPath, oldZfPath)
		if err != nil {
			return err
		}
		newZfPath := filepath.Join(newPath, relZfPath)

		newFolder, err := GetOrCreateFolderHierarchy(ctx, folderStore, newZfPath)
		if err != nil {
			return err
		}

		// add ZipFileID to new folder
		newFolder.ZipFileID = &zipFileID
		if err = folderStore.Update(ctx, newFolder); err != nil {
			return err
		}

		// remove ZipFileID from old folder
		oldFolder.ZipFileID = nil
		if err = folderStore.Update(ctx, oldFolder); err != nil {
			return err
		}
	}

	return nil
}
