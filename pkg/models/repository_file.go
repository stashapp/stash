package models

import (
	"context"
	"io/fs"
)

// FileGetter provides methods to get files by ID.
type FileGetter interface {
	Find(ctx context.Context, id ...FileID) ([]File, error)
}

// FileFinder provides methods to find files.
type FileFinder interface {
	FileGetter
	FindAllByPath(ctx context.Context, path string, caseSensitive bool) ([]File, error)
	FindAllInPaths(ctx context.Context, p []string, limit, offset int) ([]File, error)
	FindByPath(ctx context.Context, path string, caseSensitive bool) (File, error)
	FindByFingerprint(ctx context.Context, fp Fingerprint) ([]File, error)
	FindByZipFileID(ctx context.Context, zipFileID FileID) ([]File, error)
	FindByFileInfo(ctx context.Context, info fs.FileInfo, size int64) ([]File, error)
}

// FileQueryer provides methods to query files.
type FileQueryer interface {
	Query(ctx context.Context, options FileQueryOptions) (*FileQueryResult, error)
}

// FileCounter provides methods to count files.
type FileCounter interface {
	CountAllInPaths(ctx context.Context, p []string) (int, error)
	CountByFolderID(ctx context.Context, folderID FolderID) (int, error)
}

// FileCreator provides methods to create files.
type FileCreator interface {
	Create(ctx context.Context, f File) error
}

// FileUpdater provides methods to update files.
type FileUpdater interface {
	Update(ctx context.Context, f File) error
}

// FileDestroyer provides methods to destroy files.
type FileDestroyer interface {
	Destroy(ctx context.Context, id FileID) error
}

type FileFinderCreator interface {
	FileFinder
	FileCreator
}

type FileFinderUpdater interface {
	FileFinder
	FileUpdater
}

type FileFinderDestroyer interface {
	FileFinder
	FileDestroyer
}

// FileReader provides all methods to read files.
type FileReader interface {
	FileFinder
	FileQueryer
	FileCounter

	GetCaptions(ctx context.Context, fileID FileID) ([]*VideoCaption, error)
	IsPrimary(ctx context.Context, fileID FileID) (bool, error)
}

type FileFingerprintWriter interface {
	ModifyFingerprints(ctx context.Context, fileID FileID, fingerprints []Fingerprint) error
	DestroyFingerprints(ctx context.Context, fileID FileID, types []string) error
}

// FileWriter provides all methods to modify files.
type FileWriter interface {
	FileCreator
	FileUpdater
	FileDestroyer
	FileFingerprintWriter

	UpdateCaptions(ctx context.Context, fileID FileID, captions []*VideoCaption) error
}

// FileReaderWriter provides all file methods.
type FileReaderWriter interface {
	FileReader
	FileWriter
}
