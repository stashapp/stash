package image

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

type FinderByFile interface {
	FindByFileID(ctx context.Context, fileID models.FileID) ([]*models.Image, error)
	FindByZipFileID(ctx context.Context, zipFileID models.FileID) ([]*models.Image, error)
}

type Repository interface {
	FinderByFile
	Destroyer
	models.FileLoader
}

type Service struct {
	File       models.FileReaderWriter
	Repository Repository
}
