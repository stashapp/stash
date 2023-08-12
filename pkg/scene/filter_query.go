package scene

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

type FilterQueryer interface {
	Query(ctx context.Context, sceneFilterFilter *models.SceneFilterFilterType, findFilter *models.FindFilterType) ([]*models.SceneFilter, int, error)
}

type FilterCountQueryer interface {
	QueryCount(ctx context.Context, sceneFilterFilter *models.SceneFilterFilterType, findFilter *models.FindFilterType) (int, error)
}
