package file

import (
	"context"
	"io/fs"
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

// FolderGetter provides methods to find Folders.
type FolderGetter interface {
	FindByPath(ctx context.Context, path string) (*Folder, error)
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

// FolderUpdater provides methods to update Folders.
type FolderUpdater interface {
	Update(ctx context.Context, f *Folder) error
}

type FolderDestroyer interface {
	Destroy(ctx context.Context, id FolderID) error
}

// FolderStore provides methods to find, create and update Folders.
type FolderStore interface {
	FolderGetter
	FolderCounter
	FolderCreator
	FolderUpdater
	FolderDestroyer
}
