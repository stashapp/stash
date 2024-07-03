package api

import (
	"context"

	"github.com/stashapp/stash/internal/api/loaders"
	"github.com/stashapp/stash/internal/api/urlbuilders"
	"github.com/stashapp/stash/pkg/models"
)

func (r *groupResolver) Date(ctx context.Context, obj *models.Group) (*string, error) {
	if obj.Date != nil {
		result := obj.Date.String()
		return &result, nil
	}
	return nil, nil
}

func (r *groupResolver) Rating100(ctx context.Context, obj *models.Group) (*int, error) {
	return obj.Rating, nil
}

func (r *groupResolver) URL(ctx context.Context, obj *models.Group) (*string, error) {
	if !obj.URLs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadURLs(ctx, r.repository.Group)
		}); err != nil {
			return nil, err
		}
	}

	urls := obj.URLs.List()
	if len(urls) == 0 {
		return nil, nil
	}

	return &urls[0], nil
}

func (r *groupResolver) Urls(ctx context.Context, obj *models.Group) ([]string, error) {
	if !obj.URLs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadURLs(ctx, r.repository.Group)
		}); err != nil {
			return nil, err
		}
	}

	return obj.URLs.List(), nil
}

func (r *groupResolver) Studio(ctx context.Context, obj *models.Group) (ret *models.Studio, err error) {
	if obj.StudioID == nil {
		return nil, nil
	}

	return loaders.From(ctx).StudioByID.Load(*obj.StudioID)
}

func (r groupResolver) Tags(ctx context.Context, obj *models.Group) (ret []*models.Tag, err error) {
	if !obj.TagIDs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadTagIDs(ctx, r.repository.Group)
		}); err != nil {
			return nil, err
		}
	}

	var errs []error
	ret, errs = loaders.From(ctx).TagByID.LoadAll(obj.TagIDs.List())
	return ret, firstError(errs)
}

func (r *groupResolver) FrontImagePath(ctx context.Context, obj *models.Group) (*string, error) {
	var hasImage bool
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var err error
		hasImage, err = r.repository.Group.HasFrontImage(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	imagePath := urlbuilders.NewGroupURLBuilder(baseURL, obj).GetGroupFrontImageURL(hasImage)
	return &imagePath, nil
}

func (r *groupResolver) BackImagePath(ctx context.Context, obj *models.Group) (*string, error) {
	var hasImage bool
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var err error
		hasImage, err = r.repository.Group.HasBackImage(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	// don't return anything if there is no back image
	if !hasImage {
		return nil, nil
	}

	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	imagePath := urlbuilders.NewGroupURLBuilder(baseURL, obj).GetGroupBackImageURL()
	return &imagePath, nil
}

func (r *groupResolver) SceneCount(ctx context.Context, obj *models.Group) (ret int, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Scene.CountByGroupID(ctx, obj.ID)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *groupResolver) Scenes(ctx context.Context, obj *models.Group) (ret []*models.Scene, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var err error
		ret, err = r.repository.Scene.FindByGroupID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
