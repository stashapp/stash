package models

import "context"

type FolderReader interface {
	Find(ctx context.Context, id FolderID) (*Folder, error)
	FindAllInPaths(ctx context.Context, p []string, limit, offset int) ([]*Folder, error)
	FindByPath(ctx context.Context, path string) (*Folder, error)
	FindByZipFileID(ctx context.Context, zipFileID FileID) ([]*Folder, error)
	FindByParentFolderID(ctx context.Context, parentFolderID FolderID) ([]*Folder, error)

	CountAllInPaths(ctx context.Context, p []string) (int, error)
}

type FolderWriter interface {
	Create(ctx context.Context, f *Folder) error
	Update(ctx context.Context, f *Folder) error
	Destroy(ctx context.Context, id FolderID) error
}

type FolderReaderWriter interface {
	FolderReader
	FolderWriter
}
