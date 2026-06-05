package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/internal/api/loaders"
	"github.com/stashapp/stash/internal/api/urlbuilders"
	"github.com/stashapp/stash/pkg/models"
)

func (r *storyResolver) Studio(ctx context.Context, obj *models.Story) (ret *models.Studio, err error) {
	if obj.StudioID == nil {
		return nil, nil
	}
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Studio.Find(ctx, *obj.StudioID)
		return err
	}); err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *storyResolver) Tags(ctx context.Context, obj *models.Story) (ret []*models.Tag, err error) {
	if !obj.TagIDs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadTagIDs(ctx, r.repository.Story)
		}); err != nil {
			return nil, err
		}
	}
	var errs []error
	ret, errs = loaders.From(ctx).TagByID.LoadAll(obj.TagIDs.List())
	return ret, firstError(errs)
}

func (r *storyResolver) Performers(ctx context.Context, obj *models.Story) (ret []*models.Performer, err error) {
	if !obj.PerformerIDs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadPerformerIDs(ctx, r.repository.Story)
		}); err != nil {
			return nil, err
		}
	}
	var errs []error
	ret, errs = loaders.From(ctx).PerformerByID.LoadAll(obj.PerformerIDs.List())
	return ret, firstError(errs)
}

func (r *storyResolver) FrontImagePath(ctx context.Context, obj *models.Story) (*string, error) {
	hasImage, err := r.repository.Story.HasFrontImage(ctx, obj.ID)
	if err != nil {
		return nil, err
	}
	if !hasImage {
		return nil, nil
	}
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	builder := urlbuilders.NewStoryURLBuilder(baseURL, obj)
	imagePath := builder.GetFrontImageURL()
	return &imagePath, nil
}

func (r *storyResolver) BackImagePath(ctx context.Context, obj *models.Story) (*string, error) {
	hasImage, err := r.repository.Story.HasBackImage(ctx, obj.ID)
	if err != nil {
		return nil, err
	}
	if !hasImage {
		return nil, nil
	}
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	builder := urlbuilders.NewStoryURLBuilder(baseURL, obj)
	imagePath := builder.GetBackImageURL()
	return &imagePath, nil
}

func (r *storyResolver) Rating100(ctx context.Context, obj *models.Story) (*int, error) {
	return obj.Rating, nil
}

func (r *storyResolver) URL(ctx context.Context, obj *models.Story) (*string, error) {
	if !obj.URLs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadURLs(ctx, r.repository.Story)
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
