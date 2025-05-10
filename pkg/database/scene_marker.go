package database

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

type SceneMarkerStore interface {
	All(ctx context.Context) ([]*models.SceneMarker, error)
	Count(ctx context.Context) (int, error)
	CountByTagID(ctx context.Context, tagID int) (int, error)
	Create(ctx context.Context, newObject *models.SceneMarker) error
	Destroy(ctx context.Context, id int) error
	Find(ctx context.Context, id int) (*models.SceneMarker, error)
	FindBySceneID(ctx context.Context, sceneID int) ([]*models.SceneMarker, error)
	FindMany(ctx context.Context, ids []int) ([]*models.SceneMarker, error)
	GetMarkerStrings(ctx context.Context, q *string, sort *string) ([]*models.MarkerStringsResultType, error)
	GetTagIDs(ctx context.Context, id int) ([]int, error)
	Query(ctx context.Context, sceneMarkerFilter *models.SceneMarkerFilterType, findFilter *models.FindFilterType) ([]*models.SceneMarker, int, error)
	QueryCount(ctx context.Context, sceneMarkerFilter *models.SceneMarkerFilterType, findFilter *models.FindFilterType) (int, error)
	Update(ctx context.Context, updatedObject *models.SceneMarker) error
	UpdatePartial(ctx context.Context, id int, partial models.SceneMarkerPartial) (*models.SceneMarker, error)
	UpdateTags(ctx context.Context, id int, tagIDs []int) error
	Wall(ctx context.Context, q *string) ([]*models.SceneMarker, error)
}
