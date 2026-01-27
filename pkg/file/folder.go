package file

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

// GetOrCreateFolderHierarchy gets the folder for the given path, or creates a folder hierarchy for the given path if one if no existing folder is found.
// Does not create any folders in the file system
func GetOrCreateFolderHierarchy(ctx context.Context, fc models.FolderFinderCreator, path string) (*models.Folder, error) {
	// get or create folder hierarchy
	// assume case sensitive when searching for the folder
	const caseSensitive = true
	folder, err := fc.FindByPath(ctx, path, caseSensitive)
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

func transferZipHierarchy(ctx context.Context, folderStore models.FolderReaderWriter, files models.FileFinderUpdater, zipFileID models.FileID, oldPath string, newPath string) error {
	if err := transferZipFolderHierarchy(ctx, folderStore, zipFileID, oldPath, newPath); err != nil {
		return fmt.Errorf("moving folder hierarchy for file %s: %w", oldPath, err)
	}

	if err := transferZipFileEntries(ctx, folderStore, files, zipFileID, oldPath, newPath); err != nil {
		return fmt.Errorf("moving zip file contents for file %s: %w", oldPath, err)
	}

	return nil
}

// transferZipFolderHierarchy creates the folder hierarchy for zipFileID under newPath, and removes
// ZipFileID from folders under oldPath.
func transferZipFolderHierarchy(ctx context.Context, folderStore models.FolderReaderWriter, zipFileID models.FileID, oldPath string, newPath string) error {
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
		logger.Debugf("adding zip file %s to folder %s", zipFileID, newFolder.Path)
		newFolder.ZipFileID = &zipFileID
		if err = folderStore.Update(ctx, newFolder); err != nil {
			return err
		}

		// remove ZipFileID from old folder
		logger.Debugf("removing zip file %s from folder %s", zipFileID, oldFolder.Path)
		oldFolder.ZipFileID = nil
		if err = folderStore.Update(ctx, oldFolder); err != nil {
			return err
		}
	}

	return nil
}

func transferZipFileEntries(ctx context.Context, folders models.FolderFinderCreator, files models.FileFinderUpdater, zipFileID models.FileID, oldPath, newPath string) error {
	// move contained files if file is a zip file
	zipFiles, err := files.FindByZipFileID(ctx, zipFileID)
	if err != nil {
		return fmt.Errorf("finding contained files in file %s: %w", oldPath, err)
	}
	for _, zf := range zipFiles {
		zfBase := zf.Base()
		oldZfPath := zfBase.Path
		oldZfDir := filepath.Dir(oldZfPath)

		// sanity check - ignore files which aren't under oldPath
		if !strings.HasPrefix(oldZfPath, oldPath) {
			continue
		}

		relZfDir, err := filepath.Rel(oldPath, oldZfDir)
		if err != nil {
			return fmt.Errorf("moving contained file %s: %w", zfBase.ID, err)
		}
		newZfDir := filepath.Join(newPath, relZfDir)

		// folder should have been created by transferZipFolderHierarchy
		newZfFolder, err := GetOrCreateFolderHierarchy(ctx, folders, newZfDir)
		if err != nil {
			return fmt.Errorf("getting or creating folder hierarchy: %w", err)
		}

		// update file parent folder
		zfBase.ParentFolderID = newZfFolder.ID
		logger.Debugf("moving %s to folder %s", zfBase.Path, newZfFolder.Path)
		if err := files.Update(ctx, zf); err != nil {
			return fmt.Errorf("updating file %s: %w", oldZfPath, err)
		}
	}

	return nil
}
