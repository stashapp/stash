package api

import (
	"context"
	"github.com/stashapp/stash/pkg/models"
	"strconv"
)

func (r *tagResolver) ID(ctx context.Context, obj *models.Tag) (string, error) {
	return strconv.Itoa(obj.ID), nil
}

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
