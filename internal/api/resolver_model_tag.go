package api

import (
	"context"
	"time"

	"github.com/stashapp/stash/internal/api/urlbuilders"
	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
)

func (r *tagResolver) Parents(ctx context.Context, obj *models.Tag) (ret []*models.Tag, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Tag.FindByChildTagID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *tagResolver) Children(ctx context.Context, obj *models.Tag) (ret []*models.Tag, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Tag.FindByParentTagID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *tagResolver) Aliases(ctx context.Context, obj *models.Tag) (ret []string, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Tag.GetAliases(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, err
}

func (r *tagResolver) SceneCount(ctx context.Context, obj *models.Tag) (ret *int, err error) {
	var count int
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		count, err = r.repository.Scene.CountByTagID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return &count, err
}

func (r *tagResolver) SceneMarkerCount(ctx context.Context, obj *models.Tag) (ret *int, err error) {
	var count int
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		count, err = r.repository.SceneMarker.CountByTagID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return &count, err
}

func (r *tagResolver) ImageCount(ctx context.Context, obj *models.Tag) (ret *int, err error) {
	var res int
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		res, err = image.CountByTagID(ctx, r.repository.Image, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *tagResolver) GalleryCount(ctx context.Context, obj *models.Tag) (ret *int, err error) {
	var res int
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		res, err = gallery.CountByTagID(ctx, r.repository.Gallery, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *tagResolver) PerformerCount(ctx context.Context, obj *models.Tag) (ret *int, err error) {
	var count int
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		count, err = r.repository.Performer.CountByTagID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return &count, err
}

func (r *tagResolver) ImagePath(ctx context.Context, obj *models.Tag) (*string, error) {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	imagePath := urlbuilders.NewTagURLBuilder(baseURL, obj).GetTagImageURL()
	return &imagePath, nil
}

func (r *tagResolver) CreatedAt(ctx context.Context, obj *models.Tag) (*time.Time, error) {
	return &obj.CreatedAt.Timestamp, nil
}

func (r *tagResolver) UpdatedAt(ctx context.Context, obj *models.Tag) (*time.Time, error) {
	return &obj.UpdatedAt.Timestamp, nil
}
