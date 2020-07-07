package api

import (
	"context"

	"github.com/stashapp/stash/pkg/api/urlbuilders"
	"github.com/stashapp/stash/pkg/models"
)

func (r *tagResolver) SceneCount(ctx context.Context, obj *models.Tag) (*int, error) {
	qb := models.NewSceneQueryBuilder()
	if obj == nil {
		return nil, nil
	}
	count, err := qb.CountByTagID(obj.ID)
	return &count, err
}

func (r *tagResolver) SceneMarkerCount(ctx context.Context, obj *models.Tag) (*int, error) {
	qb := models.NewSceneMarkerQueryBuilder()
	if obj == nil {
		return nil, nil
	}
	count, err := qb.CountByTagID(obj.ID)
	return &count, err
}

func (r *tagResolver) ImagePath(ctx context.Context, obj *models.Tag) (*string, error) {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	imagePath := urlbuilders.NewTagURLBuilder(baseURL, obj.ID).GetTagImageURL()
	return &imagePath, nil
}
