package api

import (
	"context"

	"github.com/stashapp/stash/pkg/api/urlbuilders"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *movieResolver) Name(ctx context.Context, obj *models.Movie) (string, error) {
	if obj.Name.Valid {
		return obj.Name.String, nil
	}
	panic("null name") // TODO make name required
}

func (r *movieResolver) URL(ctx context.Context, obj *models.Movie) (*string, error) {
	if obj.URL.Valid {
		return &obj.URL.String, nil
	}
	return nil, nil
}

func (r *movieResolver) Aliases(ctx context.Context, obj *models.Movie) (*string, error) {
	if obj.Aliases.Valid {
		return &obj.Aliases.String, nil
	}
	return nil, nil
}

func (r *movieResolver) DurationMovie(ctx context.Context, obj *models.Movie) (*string, error) {
	if obj.Duration_movie.Valid {
		return &obj.Duration_movie.String, nil
	}
	return nil, nil
}

func (r *movieResolver) DateMovie(ctx context.Context, obj *models.Movie) (*string, error) {
	if obj.Date_movie.Valid {
		result := utils.GetYMDFromDatabaseDate(obj.Date_movie.String)
		return &result, nil
	}
	return nil, nil
}

func (r *movieResolver) RatingMovie(ctx context.Context, obj *models.Movie) (*string, error) {
	if obj.Rating_movie.Valid {
		return &obj.Rating_movie.String, nil
	}
	return nil, nil
}

func (r *movieResolver) Director(ctx context.Context, obj *models.Movie) (*string, error) {
	if obj.Director.Valid {
		return &obj.Director.String, nil
	}
	return nil, nil
}

func (r *movieResolver) Synopsis(ctx context.Context, obj *models.Movie) (*string, error) {
	if obj.Synopsis.Valid {
		return &obj.Synopsis.String, nil
	}
	return nil, nil
}

func (r *movieResolver) FrontImagePath(ctx context.Context, obj *models.Movie) (*string, error) {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	frontimagePath := urlbuilders.NewMovieURLBuilder(baseURL, obj.ID).GetMovieFrontImageURL()
	return &frontimagePath, nil
}

func (r *movieResolver) BackImagePath(ctx context.Context, obj *models.Movie) (*string, error) {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	backimagePath := urlbuilders.NewMovieURLBuilder(baseURL, obj.ID).GetMovieBackImageURL()
	return &backimagePath, nil
}

func (r *movieResolver) SceneCount(ctx context.Context, obj *models.Movie) (*int, error) {
	qb := models.NewSceneQueryBuilder()
	res, err := qb.CountByMovieID(obj.ID)
	return &res, err
}

func (r *movieResolver) Scenes(ctx context.Context, obj *models.Movie) ([]*models.Scene, error) {
	qb := models.NewSceneQueryBuilder()
	return qb.FindByPerformerID(obj.ID)
}