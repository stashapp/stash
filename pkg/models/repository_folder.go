package models

import "context"

// FolderGetter provides methods to get folders by ID.
type FolderGetter interface {
	Find(ctx context.Context, id FolderID) (*Folder, error)
	FindMany(ctx context.Context, id []FolderID) ([]*Folder, error)
}

// FolderFinder provides methods to find folders.
type FolderFinder interface {
	FolderGetter
	FindAllInPaths(ctx context.Context, p []string, limit, offset int) ([]*Folder, error)
	FindByPath(ctx context.Context, path string, caseSensitive bool) (*Folder, error)
	FindByZipFileID(ctx context.Context, zipFileID FileID) ([]*Folder, error)
	FindByParentFolderID(ctx context.Context, parentFolderID FolderID) ([]*Folder, error)
}

type FolderQueryer interface {
	Query(ctx context.Context, options FolderQueryOptions) (*FolderQueryResult, error)
}

type FolderCounter interface {
	CountAllInPaths(ctx context.Context, p []string) (int, error)
}

// FolderCreator provides methods to create folders.
type FolderCreator interface {
	Create(ctx context.Context, f *Folder) error
}

// FolderUpdater provides methods to update folders.
type FolderUpdater interface {
	Update(ctx context.Context, f *Folder) error
}

type FolderDestroyer interface {
	Destroy(ctx context.Context, id FolderID) error
}

type FolderFinderCreator interface {
	FolderFinder
	FolderCreator
}

type FolderFinderDestroyer interface {
	FolderFinder
	FolderDestroyer
}

// FolderReader provides all methods to read folders.
type FolderReader interface {
	FolderFinder
	FolderQueryer
	FolderCounter
}

// FolderWriter provides all methods to modify folders.
type FolderWriter interface {
	FolderCreator
	FolderUpdater
	FolderDestroyer
}

// FolderReaderWriter provides all folder methods.
type FolderReaderWriter interface {
	FolderReader
	FolderWriter
}
