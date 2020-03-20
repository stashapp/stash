package api

import (
	"context"

	"github.com/stashapp/stash/pkg/api/urlbuilders"
	"github.com/stashapp/stash/pkg/models"
)

func (r *performerResolver) Name(ctx context.Context, obj *models.Performer) (*string, error) {
	if obj.Name.Valid {
		return &obj.Name.String, nil
	}
	return nil, nil
}

func (r *performerResolver) URL(ctx context.Context, obj *models.Performer) (*string, error) {
	if obj.URL.Valid {
		return &obj.URL.String, nil
	}
	return nil, nil
}

func (r *performerResolver) Gender(ctx context.Context, obj *models.Performer) (*models.GenderEnum, error) {
	var ret models.GenderEnum

	if obj.Gender.Valid {
		ret = models.GenderEnum(obj.Gender.String)
		if ret.IsValid() {
			return &ret, nil
		}
	}

	return nil, nil
}

func (r *performerResolver) Twitter(ctx context.Context, obj *models.Performer) (*string, error) {
	if obj.Twitter.Valid {
		return &obj.Twitter.String, nil
	}
	return nil, nil
}

func (r *performerResolver) Instagram(ctx context.Context, obj *models.Performer) (*string, error) {
	if obj.Instagram.Valid {
		return &obj.Instagram.String, nil
	}
	return nil, nil
}

func (r *performerResolver) Birthdate(ctx context.Context, obj *models.Performer) (*string, error) {
	if obj.Birthdate.Valid {
		return &obj.Birthdate.String, nil
	}
	return nil, nil
}

func (r *performerResolver) Ethnicity(ctx context.Context, obj *models.Performer) (*string, error) {
	if obj.Ethnicity.Valid {
		return &obj.Ethnicity.String, nil
	}
	return nil, nil
}

func (r *performerResolver) Country(ctx context.Context, obj *models.Performer) (*string, error) {
	if obj.Country.Valid {
		return &obj.Country.String, nil
	}
	return nil, nil
}

func (r *performerResolver) EyeColor(ctx context.Context, obj *models.Performer) (*string, error) {
	if obj.EyeColor.Valid {
		return &obj.EyeColor.String, nil
	}
	return nil, nil
}

func (r *performerResolver) Height(ctx context.Context, obj *models.Performer) (*string, error) {
	if obj.Height.Valid {
		return &obj.Height.String, nil
	}
	return nil, nil
}

func (r *performerResolver) Measurements(ctx context.Context, obj *models.Performer) (*string, error) {
	if obj.Measurements.Valid {
		return &obj.Measurements.String, nil
	}
	return nil, nil
}

func (r *performerResolver) FakeTits(ctx context.Context, obj *models.Performer) (*string, error) {
	if obj.FakeTits.Valid {
		return &obj.FakeTits.String, nil
	}
	return nil, nil
}

func (r *performerResolver) CareerLength(ctx context.Context, obj *models.Performer) (*string, error) {
	if obj.CareerLength.Valid {
		return &obj.CareerLength.String, nil
	}
	return nil, nil
}

func (r *performerResolver) Tattoos(ctx context.Context, obj *models.Performer) (*string, error) {
	if obj.Tattoos.Valid {
		return &obj.Tattoos.String, nil
	}
	return nil, nil
}

func (r *performerResolver) Piercings(ctx context.Context, obj *models.Performer) (*string, error) {
	if obj.Piercings.Valid {
		return &obj.Piercings.String, nil
	}
	return nil, nil
}

func (r *performerResolver) Aliases(ctx context.Context, obj *models.Performer) (*string, error) {
	if obj.Aliases.Valid {
		return &obj.Aliases.String, nil
	}
	return nil, nil
}

func (r *performerResolver) Favorite(ctx context.Context, obj *models.Performer) (bool, error) {
	if obj.Favorite.Valid {
		return obj.Favorite.Bool, nil
	}
	return false, nil
}

func (r *performerResolver) ImagePath(ctx context.Context, obj *models.Performer) (*string, error) {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	imagePath := urlbuilders.NewPerformerURLBuilder(baseURL, obj.ID).GetPerformerImageURL()
	return &imagePath, nil
}

func (r *performerResolver) SceneCount(ctx context.Context, obj *models.Performer) (*int, error) {
	qb := models.NewSceneQueryBuilder()
	res, err := qb.CountByPerformerID(obj.ID)
	return &res, err
}

func (r *performerResolver) Scenes(ctx context.Context, obj *models.Performer) ([]*models.Scene, error) {
	qb := models.NewSceneQueryBuilder()
	return qb.FindByPerformerID(obj.ID)
}

func (r *performerResolver) StashID(ctx context.Context, obj *models.Performer) (*string, error) {
	if obj.StashID.Valid {
		return &obj.StashID.String, nil
	}
	return nil, nil
}
