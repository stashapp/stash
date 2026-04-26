package file

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

// GetOrCreateFolderHierarchy gets the folder for the given path, or creates a folder hierarchy for the given path if one if no existing folder is found.
// Creates folder entries for each level of the hierarchy that doesn't already exist, up to the provided root paths.
// Does not create any folders in the file system.
func GetOrCreateFolderHierarchy(ctx context.Context, fc models.FolderFinderCreator, path string, rootPaths []string) (*models.Folder, error) {
	// get or create folder hierarchy
	// assume case sensitive when searching for the folder
	const caseSensitive = true
	folder, err := fc.FindByPath(ctx, path, caseSensitive)
	if err != nil {
		return nil, err
	}

	if folder == nil {
		var parentID *models.FolderID

		if !slices.Contains(rootPaths, path) {
			parentPath := filepath.Dir(path)

			// safety check - don't allow parent path to be the same as the current path,
			// otherwise we could end up in an infinite loop
			if parentPath == path {
				// #6618 - log a warning and return nil for the parent ID,
				// which will cause the folder to be created with no parent
				logger.Warnf("parent path is the same as the current path: %s", path)
				return nil, nil
			}

			parent, err := GetOrCreateFolderHierarchy(ctx, fc, parentPath, rootPaths)
			if err != nil {
				return nil, err
			}

			parentID = &parent.ID
		}

		now := time.Now()

		folder = &models.Folder{
			Path:           path,
			ParentFolderID: parentID,
			DirEntry:       models.DirEntry{
				// leave mod time empty for now - it will be updated when the folder is scanned
			},
			CreatedAt: now,
			UpdatedAt: now,
		}

		logger.Infof("%s doesn't exist. Creating new folder entry...", path)

		if err = fc.Create(ctx, folder); err != nil {
			return nil, fmt.Errorf("creating folder %s: %w", path, err)
		}
	}

	return folder, nil
}

type zipHierarchyMover struct {
	folderStore models.FolderReaderWriter
	files       models.FileFinderUpdater
	rootPaths   []string
}

func (m zipHierarchyMover) transferZipHierarchy(ctx context.Context, zipFileID models.FileID, oldPath string, newPath string) error {
	if err := m.transferZipFolderHierarchy(ctx, zipFileID, oldPath, newPath); err != nil {
		return fmt.Errorf("moving folder hierarchy for file %s: %w", oldPath, err)
	}

	if err := m.transferZipFileEntries(ctx, zipFileID, oldPath, newPath); err != nil {
		return fmt.Errorf("moving zip file contents for file %s: %w", oldPath, err)
	}

	return nil
}

// transferZipFolderHierarchy creates the folder hierarchy for zipFileID under newPath, and removes
// ZipFileID from folders under oldPath.
func (m zipHierarchyMover) transferZipFolderHierarchy(ctx context.Context, zipFileID models.FileID, oldPath string, newPath string) error {
	zipFolders, err := m.folderStore.FindByZipFileID(ctx, zipFileID)
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

		newFolder, err := GetOrCreateFolderHierarchy(ctx, m.folderStore, newZfPath, m.rootPaths)
		if err != nil {
			return err
		}

		// add ZipFileID to new folder
		logger.Debugf("adding zip file %s to folder %s", zipFileID, newFolder.Path)
		newFolder.ZipFileID = &zipFileID
		if err = m.folderStore.Update(ctx, newFolder); err != nil {
			return err
		}

		// remove ZipFileID from old folder
		logger.Debugf("removing zip file %s from folder %s", zipFileID, oldFolder.Path)
		oldFolder.ZipFileID = nil
		if err = m.folderStore.Update(ctx, oldFolder); err != nil {
			return err
		}
	}

	return nil
}

func (m zipHierarchyMover) transferZipFileEntries(ctx context.Context, zipFileID models.FileID, oldPath, newPath string) error {
	// move contained files if file is a zip file
	zipFiles, err := m.files.FindByZipFileID(ctx, zipFileID)
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
		newZfFolder, err := GetOrCreateFolderHierarchy(ctx, m.folderStore, newZfDir, m.rootPaths)
		if err != nil {
			return fmt.Errorf("getting or creating folder hierarchy: %w", err)
		}

		// update file parent folder
		zfBase.ParentFolderID = newZfFolder.ID
		logger.Debugf("moving %s to folder %s", zfBase.Path, newZfFolder.Path)
		if err := m.files.Update(ctx, zf); err != nil {
			return fmt.Errorf("updating file %s: %w", oldZfPath, err)
		}
	}

	return nil
}
