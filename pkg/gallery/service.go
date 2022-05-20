package gallery

import (
	"context"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
)

type Repository interface {
	FindByFileID(ctx context.Context, fileID file.ID) ([]*models.Gallery, error)
	Destroy(ctx context.Context, id int) error
}

type ImageFinder interface {
	FindByFolderID(ctx context.Context, folder file.FolderID) ([]*models.Image, error)
	FindByZipFileID(ctx context.Context, zipFileID file.ID) ([]*models.Image, error)
}

type ImageService interface {
	Destroy(ctx context.Context, i *models.Image, fileDeleter *image.FileDeleter, deleteGenerated, deleteFile bool) error
}

type Service struct {
	Repository   Repository
	ImageFinder  ImageFinder
	ImageService ImageService
}
