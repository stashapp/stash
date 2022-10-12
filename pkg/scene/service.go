package scene

import (
	"context"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/models"
)

type FinderByFile interface {
	FindByFileID(ctx context.Context, fileID file.ID) ([]*models.Scene, error)
}

type FileAssigner interface {
	AssignFiles(ctx context.Context, sceneID int, fileID file.ID) error
}

type Repository interface {
	FinderByFile
	Destroyer
	models.VideoFileLoader
	FileAssigner
}

type Service struct {
	File            file.Store
	Repository      Repository
	MarkerDestroyer MarkerDestroyer
}
