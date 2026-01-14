// Package gallery provides application logic for managing galleries.
// This functionality is exposed via the [Service] type.
package gallery

import (
	"context"

	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
)

type ImageFinder interface {
	FindByFolderID(ctx context.Context, folder models.FolderID) ([]*models.Image, error)
	FindByZipFileID(ctx context.Context, zipFileID models.FileID) ([]*models.Image, error)
	models.GalleryIDLoader
}

type ImageService interface {
	Destroy(ctx context.Context, i *models.Image, fileDeleter *image.FileDeleter, deleteGenerated, deleteFile, destroyFileEntry bool) error
	DestroyZipImages(ctx context.Context, zipFile models.File, fileDeleter *image.FileDeleter, deleteGenerated bool) ([]*models.Image, error)
	DestroyFolderImages(ctx context.Context, folderID models.FolderID, fileDeleter *image.FileDeleter, deleteGenerated, deleteFile bool) ([]*models.Image, error)
}

type Service struct {
	Repository   models.GalleryReaderWriter
	ImageFinder  ImageFinder
	ImageService ImageService
	File         models.FileReaderWriter
	Folder       models.FolderReaderWriter
}
