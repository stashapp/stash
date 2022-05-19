package file

import (
	"context"
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

// FolderGetter provides methods to find Folders.
type FolderGetter interface {
	FindByPath(ctx context.Context, path string) (*Folder, error)
	FindByZipFileID(ctx context.Context, zipFileID ID) ([]*Folder, error)
}

// FolderCreator provides methods to create Folders.
type FolderCreator interface {
	Create(ctx context.Context, f *Folder) error
}

// FolderUpdater provides methods to update Folders.
type FolderUpdater interface {
	Update(ctx context.Context, f *Folder) error
}

// MissingMarker wraps the MarkMissing method.
type FolderMissingMarker interface {
	FindMissing(ctx context.Context, scanStartTime time.Time, scanPaths []string, page uint, limit uint) ([]*Folder, error)
	MarkMissing(ctx context.Context, scanStartTime time.Time, scanPaths []string) (int, error)
}

// FolderStore provides methods to find, create and update Folders.
type FolderStore interface {
	FolderGetter
	FolderCreator
	FolderUpdater
	FolderMissingMarker
}
