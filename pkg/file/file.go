package file

import (
	"context"
	"io/fs"

	"github.com/stashapp/stash/pkg/models"
)

type Finder interface {
	Find(ctx context.Context, id ...models.FileID) ([]models.File, error)
}

// Getter provides methods to find Files.
type Getter interface {
	Finder
	FindByPath(ctx context.Context, path string) (models.File, error)
	FindAllByPath(ctx context.Context, path string) ([]models.File, error)
	FindByFingerprint(ctx context.Context, fp models.Fingerprint) ([]models.File, error)
	FindByZipFileID(ctx context.Context, zipFileID models.FileID) ([]models.File, error)
	FindAllInPaths(ctx context.Context, p []string, limit, offset int) ([]models.File, error)
	FindByFileInfo(ctx context.Context, info fs.FileInfo, size int64) ([]models.File, error)
}

type Counter interface {
	CountAllInPaths(ctx context.Context, p []string) (int, error)
	CountByFolderID(ctx context.Context, folderID models.FolderID) (int, error)
}

// Creator provides methods to create Files.
type Creator interface {
	Create(ctx context.Context, f models.File) error
}

// Updater provides methods to update Files.
type Updater interface {
	Update(ctx context.Context, f models.File) error
}

type Destroyer interface {
	Destroy(ctx context.Context, id models.FileID) error
}

type GetterUpdater interface {
	Getter
	Updater
}

type GetterDestroyer interface {
	Getter
	Destroyer
}

// Store provides methods to find, create and update Files.
type Store interface {
	Getter
	Counter
	Creator
	Updater
	Destroyer

	IsPrimary(ctx context.Context, fileID models.FileID) (bool, error)
}

// Decorator wraps the Decorate method to add additional functionality while scanning files.
type Decorator interface {
	Decorate(ctx context.Context, fs models.FS, f models.File) (models.File, error)
	IsMissingMetadata(ctx context.Context, fs models.FS, f models.File) bool
}

type FilteredDecorator struct {
	Decorator
	Filter
}

// Decorate runs the decorator if the filter accepts the file.
func (d *FilteredDecorator) Decorate(ctx context.Context, fs models.FS, f models.File) (models.File, error) {
	if d.Accept(ctx, f) {
		return d.Decorator.Decorate(ctx, fs, f)
	}
	return f, nil
}

func (d *FilteredDecorator) IsMissingMetadata(ctx context.Context, fs models.FS, f models.File) bool {
	if d.Accept(ctx, f) {
		return d.Decorator.IsMissingMetadata(ctx, fs, f)
	}

	return false
}
