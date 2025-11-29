package database

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

type FolderStore interface {
	CountAllInPaths(ctx context.Context, p []string) (int, error)
	Create(ctx context.Context, f *models.Folder) error
	Destroy(ctx context.Context, id models.FolderID) error
	Find(ctx context.Context, id models.FolderID) (*models.Folder, error)
	FindAllInPaths(ctx context.Context, p []string, limit int, offset int) ([]*models.Folder, error)
	FindByParentFolderID(ctx context.Context, parentFolderID models.FolderID) ([]*models.Folder, error)
	FindByPath(ctx context.Context, p string) (*models.Folder, error)
	FindByZipFileID(ctx context.Context, zipFileID models.FileID) ([]*models.Folder, error)
	Update(ctx context.Context, updatedObject *models.Folder) error
	FindByIDs(ctx context.Context, ids []models.FolderID) ([]*models.Folder, error)
	FindMany(ctx context.Context, ids []models.FolderID) ([]*models.Folder, error)
	Query(ctx context.Context, options models.FolderQueryOptions) (*models.FolderQueryResult, error)
}
