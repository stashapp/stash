package scene

import (
	"context"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/models"
)

type CoverGenerator interface {
	GenerateCover(ctx context.Context, scene *models.Scene, f *file.VideoFile) error
}
