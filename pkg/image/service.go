package image

import (
	"context"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/models"
)

type FinderByFile interface {
	FindByFileID(ctx context.Context, fileID file.ID) ([]*models.Image, error)
	FindByZipFileID(ctx context.Context, zipFileID file.ID) ([]*models.Image, error)
}

type Repository interface {
	FinderByFile
	Destroyer
	models.FileLoader
}

type Service struct {
	File       file.Store
	Repository Repository
}
