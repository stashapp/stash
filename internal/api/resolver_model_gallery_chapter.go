package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

func (r *galleryChapterResolver) Gallery(ctx context.Context, obj *models.GalleryChapter) (ret *models.Gallery, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Gallery.Find(ctx, obj.GalleryID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
