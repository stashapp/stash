package scene

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/paths"
	"github.com/stashapp/stash/pkg/plugin"
)

type Config interface {
	GetVideoFileNamingAlgorithm() models.HashAlgorithm
}

type MarkerRepository interface {
	MarkerFinder
	MarkerDestroyer

	Update(ctx context.Context, updatedObject *models.SceneMarker) error
}

type Service struct {
	File             models.FileReaderWriter
	Repository       models.SceneReaderWriter
	MarkerRepository MarkerRepository
	PluginCache      *plugin.Cache

	Paths  *paths.Paths
	Config Config
}
