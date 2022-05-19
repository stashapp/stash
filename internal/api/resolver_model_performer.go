package api

import (
	"context"
	"time"

	"github.com/stashapp/stash/internal/api/urlbuilders"
	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/image"
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
	imagePath := urlbuilders.NewPerformerURLBuilder(baseURL, obj).GetPerformerImageURL()
	return &imagePath, nil
}

func (r *performerResolver) Tags(ctx context.Context, obj *models.Performer) (ret []*models.Tag, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Tag.FindByPerformerID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *performerResolver) SceneCount(ctx context.Context, obj *models.Performer) (ret *int, err error) {
	var res int
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		res, err = r.repository.Scene.CountByPerformerID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *performerResolver) ImageCount(ctx context.Context, obj *models.Performer) (ret *int, err error) {
	var res int
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		res, err = image.CountByPerformerID(ctx, r.repository.Image, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *performerResolver) GalleryCount(ctx context.Context, obj *models.Performer) (ret *int, err error) {
	var res int
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		res, err = gallery.CountByPerformerID(ctx, r.repository.Gallery, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *performerResolver) Scenes(ctx context.Context, obj *models.Performer) (ret []*models.Scene, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Scene.FindByPerformerID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *performerResolver) StashIds(ctx context.Context, obj *models.Performer) (ret []*models.StashID, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Performer.GetStashIDs(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *performerResolver) Rating(ctx context.Context, obj *models.Performer) (*int, error) {
	if obj.Rating.Valid {
		rating := int(obj.Rating.Int64)
		return &rating, nil
	}
	return nil, nil
}

func (r *performerResolver) Details(ctx context.Context, obj *models.Performer) (*string, error) {
	if obj.Details.Valid {
		return &obj.Details.String, nil
	}
	return nil, nil
}

func (r *performerResolver) DeathDate(ctx context.Context, obj *models.Performer) (*string, error) {
	if obj.DeathDate.Valid {
		return &obj.DeathDate.String, nil
	}
	return nil, nil
}

func (r *performerResolver) HairColor(ctx context.Context, obj *models.Performer) (*string, error) {
	if obj.HairColor.Valid {
		return &obj.HairColor.String, nil
	}
	return nil, nil
}

func (r *performerResolver) Weight(ctx context.Context, obj *models.Performer) (*int, error) {
	if obj.Weight.Valid {
		weight := int(obj.Weight.Int64)
		return &weight, nil
	}
	return nil, nil
}

func (r *performerResolver) CreatedAt(ctx context.Context, obj *models.Performer) (*time.Time, error) {
	return &obj.CreatedAt.Timestamp, nil
}

func (r *performerResolver) UpdatedAt(ctx context.Context, obj *models.Performer) (*time.Time, error) {
	return &obj.UpdatedAt.Timestamp, nil
}

func (r *performerResolver) Movies(ctx context.Context, obj *models.Performer) (ret []*models.Movie, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Movie.FindByPerformerID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *performerResolver) MovieCount(ctx context.Context, obj *models.Performer) (ret *int, err error) {
	var res int
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		res, err = r.repository.Movie.CountByPerformerID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return &res, nil
}
