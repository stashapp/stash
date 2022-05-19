package api

import (
	"context"
	"time"

	"github.com/stashapp/stash/internal/api/urlbuilders"
	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/image"
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
	imagePath := urlbuilders.NewStudioURLBuilder(baseURL, obj).GetStudioImageURL()

	var hasImage bool
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		var err error
		hasImage, err = r.repository.Studio.HasImage(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	// indicate that image is missing by setting default query param to true
	if !hasImage {
		imagePath += "?default=true"
	}

	return &imagePath, nil
}

func (r *studioResolver) Aliases(ctx context.Context, obj *models.Studio) (ret []string, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Studio.GetAliases(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, err
}

func (r *studioResolver) SceneCount(ctx context.Context, obj *models.Studio) (ret *int, err error) {
	var res int
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		res, err = r.repository.Scene.CountByStudioID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return &res, err
}

func (r *studioResolver) ImageCount(ctx context.Context, obj *models.Studio) (ret *int, err error) {
	var res int
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		res, err = image.CountByStudioID(ctx, r.repository.Image, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *studioResolver) GalleryCount(ctx context.Context, obj *models.Studio) (ret *int, err error) {
	var res int
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		res, err = gallery.CountByStudioID(ctx, r.repository.Gallery, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *studioResolver) ParentStudio(ctx context.Context, obj *models.Studio) (ret *models.Studio, err error) {
	if !obj.ParentID.Valid {
		return nil, nil
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Studio.Find(ctx, int(obj.ParentID.Int64))
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *studioResolver) ChildStudios(ctx context.Context, obj *models.Studio) (ret []*models.Studio, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Studio.FindChildren(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *studioResolver) StashIds(ctx context.Context, obj *models.Studio) (ret []*models.StashID, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Studio.GetStashIDs(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *studioResolver) Rating(ctx context.Context, obj *models.Studio) (*int, error) {
	if obj.Rating.Valid {
		rating := int(obj.Rating.Int64)
		return &rating, nil
	}
	return nil, nil
}

func (r *studioResolver) Details(ctx context.Context, obj *models.Studio) (*string, error) {
	if obj.Details.Valid {
		return &obj.Details.String, nil
	}
	return nil, nil
}

func (r *studioResolver) CreatedAt(ctx context.Context, obj *models.Studio) (*time.Time, error) {
	return &obj.CreatedAt.Timestamp, nil
}

func (r *studioResolver) UpdatedAt(ctx context.Context, obj *models.Studio) (*time.Time, error) {
	return &obj.UpdatedAt.Timestamp, nil
}

func (r *studioResolver) Movies(ctx context.Context, obj *models.Studio) (ret []*models.Movie, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Movie.FindByStudioID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *studioResolver) MovieCount(ctx context.Context, obj *models.Studio) (ret *int, err error) {
	var res int
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		res, err = r.repository.Movie.CountByStudioID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return &res, nil
}
