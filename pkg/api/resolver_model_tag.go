package api

import (
	"context"
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
