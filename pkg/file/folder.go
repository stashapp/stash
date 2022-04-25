package file

import (
	"context"
	"strconv"
	"time"
)

type FolderID int32

func (i FolderID) String() string {
	return strconv.Itoa(int(i))
}

type Folder struct {
	ID FolderID `json:"id"`
	DirEntry

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FolderGetter interface {
	GetByPath(ctx context.Context, path string) (*Folder, error)
}

type FolderCreator interface {
	Create(ctx context.Context, f Folder) (*Folder, error)
}

type FolderUpdater interface {
	Update(ctx context.Context, f Folder) (*Folder, error)
}

type FolderStore interface {
	FolderGetter
	FolderCreator
	FolderUpdater
}
