package file

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/models"
)

type FolderFinder interface {
	Find(ctx context.Context, id models.FolderID) (*models.Folder, error)
}

// FolderPathFinder finds Folders by their path.
type FolderPathFinder interface {
	FindByPath(ctx context.Context, path string) (*models.Folder, error)
}

// FolderGetter provides methods to find Folders.
type FolderGetter interface {
	FolderFinder
	FolderPathFinder
	FindByZipFileID(ctx context.Context, zipFileID models.FileID) ([]*models.Folder, error)
	FindAllInPaths(ctx context.Context, p []string, limit, offset int) ([]*models.Folder, error)
	FindByParentFolderID(ctx context.Context, parentFolderID models.FolderID) ([]*models.Folder, error)
}

type FolderCounter interface {
	CountAllInPaths(ctx context.Context, p []string) (int, error)
}

// FolderCreator provides methods to create Folders.
type FolderCreator interface {
	Create(ctx context.Context, f *models.Folder) error
}

type FolderFinderCreator interface {
	FolderPathFinder
	FolderCreator
}

// FolderUpdater provides methods to update Folders.
type FolderUpdater interface {
	Update(ctx context.Context, f *models.Folder) error
}

type FolderDestroyer interface {
	Destroy(ctx context.Context, id models.FolderID) error
}

type FolderGetterDestroyer interface {
	FolderGetter
	FolderDestroyer
}

// FolderStore provides methods to find, create and update Folders.
type FolderStore interface {
	FolderGetter
	FolderCounter
	FolderCreator
	FolderUpdater
	FolderDestroyer
}

// GetOrCreateFolderHierarchy gets the folder for the given path, or creates a folder hierarchy for the given path if one if no existing folder is found.
// Does not create any folders in the file system
func GetOrCreateFolderHierarchy(ctx context.Context, fc FolderFinderCreator, path string) (*models.Folder, error) {
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
func TransferZipFolderHierarchy(ctx context.Context, folderStore FolderStore, zipFileID models.FileID, oldPath string, newPath string) error {
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
