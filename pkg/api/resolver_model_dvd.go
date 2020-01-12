package api

import (
	"context"

	"github.com/stashapp/stash/pkg/api/urlbuilders"
	"github.com/stashapp/stash/pkg/models"
)

func (r *dvdResolver) Name(ctx context.Context, obj *models.Dvd) (string, error) {
	if obj.Name.Valid {
		return obj.Name.String, nil
	}
	panic("null name") // TODO make name required
}

func (r *dvdResolver) URL(ctx context.Context, obj *models.Dvd) (*string, error) {
	if obj.URL.Valid {
		return &obj.URL.String, nil
	}
	return nil, nil
}

func (r *dvdResolver) Aliases(ctx context.Context, obj *models.Dvd) (*string, error) {
	if obj.Aliases.Valid {
		return &obj.Aliases.String, nil
	}
	return nil, nil
}

func (r *dvdResolver) Durationdvd(ctx context.Context, obj *models.Dvd) (*string, error) {
	if obj.Durationdvd.Valid {
		return &obj.Durationdvd.String, nil
	}
	return nil, nil
}

func (r *dvdResolver) Year(ctx context.Context, obj *models.Dvd) (*string, error) {
	if obj.Year.Valid {
		return &obj.Year.String, nil
	}
	return nil, nil
}

func (r *dvdResolver) Director(ctx context.Context, obj *models.Dvd) (*string, error) {
	if obj.Director.Valid {
		return &obj.Director.String, nil
	}
	return nil, nil
}

func (r *dvdResolver) Synopsis(ctx context.Context, obj *models.Dvd) (*string, error) {
	if obj.Synopsis.Valid {
		return &obj.Synopsis.String, nil
	}
	return nil, nil
}

func (r *dvdResolver) FrontimagePath(ctx context.Context, obj *models.Dvd) (*string, error) {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	frontimagePath := urlbuilders.NewDvdURLBuilder(baseURL, obj.ID).GetDvdFrontImageURL()
	return &frontimagePath, nil
}

func (r *dvdResolver) BackimagePath(ctx context.Context, obj *models.Dvd) (*string, error) {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	backimagePath := urlbuilders.NewDvdURLBuilder(baseURL, obj.ID).GetDvdBackImageURL()
	return &backimagePath, nil
}

func (r *dvdResolver) SceneCount(ctx context.Context, obj *models.Dvd) (*int, error) {
	qb := models.NewSceneQueryBuilder()
	res, err := qb.CountByDvdID(obj.ID)
	return &res, err
}
