package scene

import (
	"context"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/models"
)

type FinderByFile interface {
	FindByFileID(ctx context.Context, fileID file.ID) ([]*models.Scene, error)
}

type Repository interface {
	FinderByFile
	Destroyer
}

type Service struct {
	File            file.Store
	Repository      Repository
	MarkerDestroyer MarkerDestroyer
}
