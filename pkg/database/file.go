package database

import (
	"context"
	"io/fs"

	"github.com/stashapp/stash/pkg/models"
)

type FileStore interface {
	CountAllInPaths(ctx context.Context, p []string) (int, error)
	CountByFolderID(ctx context.Context, folderID models.FolderID) (int, error)
	Create(ctx context.Context, f models.File) error
	Destroy(ctx context.Context, id models.FileID) error
	DestroyFingerprints(ctx context.Context, fileID models.FileID, types []string) error
	Find(ctx context.Context, ids ...models.FileID) ([]models.File, error)
	FindAllByPath(ctx context.Context, p string) ([]models.File, error)
	FindAllInPaths(ctx context.Context, p []string, limit int, offset int) ([]models.File, error)
	FindByFileInfo(ctx context.Context, info fs.FileInfo, size int64) ([]models.File, error)
	FindByFingerprint(ctx context.Context, fp models.Fingerprint) ([]models.File, error)
	FindByPath(ctx context.Context, p string) (models.File, error)
	FindByZipFileID(ctx context.Context, zipFileID models.FileID) ([]models.File, error)
	GetCaptions(ctx context.Context, fileID models.FileID) ([]*models.VideoCaption, error)
	IsPrimary(ctx context.Context, fileID models.FileID) (bool, error)
	ModifyFingerprints(ctx context.Context, fileID models.FileID, fingerprints []models.Fingerprint) error
	Query(ctx context.Context, options models.FileQueryOptions) (*models.FileQueryResult, error)
	Update(ctx context.Context, f models.File) error
	UpdateCaptions(ctx context.Context, fileID models.FileID, captions []*models.VideoCaption) error
}
