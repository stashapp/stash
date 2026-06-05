package api

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/internal/api/loaders"
	"github.com/stashapp/stash/internal/api/urlbuilders"
	"github.com/stashapp/stash/pkg/models"
)

func (r *mangaResolver) Studio(ctx context.Context, obj *models.Manga) (ret *models.Studio, err error) {
	if obj.StudioID == nil {
		return nil, nil
	}

	return loaders.From(ctx).StudioByID.Load(*obj.StudioID)
}

func (r *mangaResolver) Tags(ctx context.Context, obj *models.Manga) (ret []*models.Tag, err error) {
	if !obj.TagIDs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadTagIDs(ctx, r.repository.Manga)
		}); err != nil {
			return nil, err
		}
	}

	var errs []error
	ret, errs = loaders.From(ctx).TagByID.LoadAll(obj.TagIDs.List())
	return ret, firstError(errs)
}

func (r *mangaResolver) Performers(ctx context.Context, obj *models.Manga) (ret []*models.Performer, err error) {
	if !obj.PerformerIDs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadPerformerIDs(ctx, r.repository.Manga)
		}); err != nil {
			return nil, err
		}
	}

	var errs []error
	ret, errs = loaders.From(ctx).PerformerByID.LoadAll(obj.PerformerIDs.List())
	return ret, firstError(errs)
}

func (r *mangaResolver) Date(ctx context.Context, obj *models.Manga) (*string, error) {
	if obj.Date != nil {
		result := obj.Date.String()
		return &result, nil
	}
	return nil, nil
}

func (r *mangaResolver) Rating100(ctx context.Context, obj *models.Manga) (*int, error) {
	return obj.Rating, nil
}

func (r *mangaResolver) URL(ctx context.Context, obj *models.Manga) (*string, error) {
	return &obj.URL, nil
}

func (r *mangaResolver) Paths(ctx context.Context, obj *models.Manga) (*MangaPathsType, error) {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	builder := urlbuilders.NewMangaURLBuilder(baseURL, obj)

	return &MangaPathsType{
		Cover: builder.GetCoverURL(),
	}, nil
}
