package api

import (
	"context"
	"github.com/stashapp/stash/api/urlbuilders"
	"github.com/stashapp/stash/models"
	"strconv"
)

func (r *studioResolver) ID(ctx context.Context, obj *models.Studio) (string, error) {
	return strconv.Itoa(obj.ID), nil
}

func (r *studioResolver) Name(ctx context.Context, obj *models.Studio) (string, error) {
	if obj.Name.Valid {
		return obj.Name.String, nil
	}
	panic("null name") // TODO make name required
}

func (r *studioResolver) URL(ctx context.Context, obj *models.Studio) (*string, error) {
	if obj.Url.Valid {
		return &obj.Url.String, nil
	}
	return nil, nil
}

func (r *studioResolver) ImagePath(ctx context.Context, obj *models.Studio) (*string, error) {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	imagePath := urlbuilders.NewStudioURLBuilder(baseURL, obj.ID).GetStudioImageUrl()
	return &imagePath, nil
}

func (r *studioResolver) SceneCount(ctx context.Context, obj *models.Studio) (*int, error) {
	qb := models.NewSceneQueryBuilder()
	res, err := qb.CountByStudioID(obj.ID)
	return &res, err
}
