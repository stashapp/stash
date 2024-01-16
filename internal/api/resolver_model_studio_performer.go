package api

import (
	"context"

	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/movie"
	"github.com/stashapp/stash/pkg/scene"
)

func (r *studioPerformerResolver) SceneCount(ctx context.Context, obj *models.StudioPerformer) (ret int, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = scene.CountByPerformerIDStudioID(ctx, r.repository.Scene, obj.PerformerID, obj.StudioID, obj.Depth)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *studioPerformerResolver) GalleryCount(ctx context.Context, obj *models.StudioPerformer) (ret int, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = gallery.CountByPerformerIDStudioID(ctx, r.repository.Gallery, obj.PerformerID, obj.StudioID, obj.Depth)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *studioPerformerResolver) ImageCount(ctx context.Context, obj *models.StudioPerformer) (ret int, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = image.CountByPerformerIDStudioID(ctx, r.repository.Image, obj.PerformerID, obj.StudioID, obj.Depth)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *studioPerformerResolver) MovieCount(ctx context.Context, obj *models.StudioPerformer) (ret int, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = movie.CountByPerformerIDStudioID(ctx, r.repository.Movie, obj.PerformerID, obj.StudioID, obj.Depth)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}
