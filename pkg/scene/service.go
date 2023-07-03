package scene

import (
	"context"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/paths"
	"github.com/stashapp/stash/pkg/plugin"
)

type FinderByFile interface {
	FindByFileID(ctx context.Context, fileID file.ID) ([]*models.Scene, error)
}

type FileAssigner interface {
	AssignFiles(ctx context.Context, sceneID int, fileID []file.ID) error
}

type Creator interface {
	Create(ctx context.Context, newScene *models.Scene, fileIDs []file.ID) error
}

type CoverUpdater interface {
	HasCover(ctx context.Context, sceneID int) (bool, error)
	UpdateCover(ctx context.Context, sceneID int, cover []byte) error
}

type Config interface {
	GetVideoFileNamingAlgorithm() models.HashAlgorithm
}

type Repository interface {
	IDFinder
	FinderByFile
	Creator
	PartialUpdater
	Destroyer
	models.VideoFileLoader
	FileAssigner
	CoverUpdater
	models.SceneReader
}

type MarkerRepository interface {
	MarkerFinder
	MarkerDestroyer

	Update(ctx context.Context, updatedObject *models.SceneMarker) error
}

type Service struct {
	File             file.Store
	Repository       Repository
	MarkerRepository MarkerRepository
	PluginCache      *plugin.Cache

	Paths  *paths.Paths
	Config Config
}
