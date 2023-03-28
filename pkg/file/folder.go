package file

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"strconv"
	"time"
)

// FolderID represents an ID of a folder.
type FolderID int32

// String converts the ID to a string.
func (i FolderID) String() string {
	return strconv.Itoa(int(i))
}

// Folder represents a folder in the file system.
type Folder struct {
	ID FolderID `json:"id"`
	DirEntry
	Path           string    `json:"path"`
	ParentFolderID *FolderID `json:"parent_folder_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (f *Folder) Info(fs FS) (fs.FileInfo, error) {
	return f.info(fs, f.Path)
}

type FolderFinder interface {
	Find(ctx context.Context, id FolderID) (*Folder, error)
}

// FolderPathFinder finds Folders by their path.
type FolderPathFinder interface {
	FindByPath(ctx context.Context, path string) (*Folder, error)
}

// FolderGetter provides methods to find Folders.
type FolderGetter interface {
	FolderFinder
	FolderPathFinder
	FindByZipFileID(ctx context.Context, zipFileID ID) ([]*Folder, error)
	FindAllInPaths(ctx context.Context, p []string, limit, offset int) ([]*Folder, error)
	FindByParentFolderID(ctx context.Context, parentFolderID FolderID) ([]*Folder, error)
}

type FolderCounter interface {
	CountAllInPaths(ctx context.Context, p []string) (int, error)
}

// FolderCreator provides methods to create Folders.
type FolderCreator interface {
	Create(ctx context.Context, f *Folder) error
}

type FolderFinderCreator interface {
	FolderPathFinder
	FolderCreator
}

// FolderUpdater provides methods to update Folders.
type FolderUpdater interface {
	Update(ctx context.Context, f *Folder) error
}

type FolderDestroyer interface {
	Destroy(ctx context.Context, id FolderID) error
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
func GetOrCreateFolderHierarchy(ctx context.Context, fc FolderFinderCreator, path string) (*Folder, error) {
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

		folder = &Folder{
			Path:           path,
			ParentFolderID: &parent.ID,
			DirEntry:       DirEntry{
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
