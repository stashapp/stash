package api

import (
	"context"

	"github.com/stashapp/stash/pkg/api/urlbuilders"
	"github.com/stashapp/stash/pkg/models"
)

func (r *studioResolver) Name(ctx context.Context, obj *models.Studio) (string, error) {
	if obj.Name.Valid {
		return obj.Name.String, nil
	}
	panic("null name") // TODO make name required
}

func (r *studioResolver) URL(ctx context.Context, obj *models.Studio) (*string, error) {
	if obj.URL.Valid {
		return &obj.URL.String, nil
	}
	return nil, nil
}

func (r *studioResolver) ImagePath(ctx context.Context, obj *models.Studio) (*string, error) {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	imagePath := urlbuilders.NewStudioURLBuilder(baseURL, obj.ID).GetStudioImageURL()
	return &imagePath, nil
}

func (r *studioResolver) SceneCount(ctx context.Context, obj *models.Studio) (*int, error) {
	qb := models.NewSceneQueryBuilder()
	res, err := qb.CountByStudioID(obj.ID)
	return &res, err
}

func (r *studioResolver) ParentStudio(ctx context.Context, obj *models.Studio) (*models.Studio, error) {
	if !obj.ParentID.Valid {
		return nil, nil
	}

	qb := models.NewStudioQueryBuilder()
	return qb.Find(int(obj.ParentID.Int64), nil)
}

func (r *studioResolver) ChildStudios(ctx context.Context, obj *models.Studio) ([]*models.Studio, error) {
	qb := models.NewStudioQueryBuilder()
	return qb.FindChildren(obj.ID, nil)
}
