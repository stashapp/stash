package api

import (
	"context"
	"time"

	"github.com/stashapp/stash/pkg/models"
)

func (r *galleryChapterResolver) Gallery(ctx context.Context, obj *models.GalleryChapter) (ret *models.Gallery, err error) {
	if !obj.GalleryID.Valid {
		panic("Invalid gallery id")
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		galleryID := int(obj.GalleryID.Int64)
		ret, err = r.repository.Gallery.Find(ctx, galleryID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *galleryChapterResolver) CreatedAt(ctx context.Context, obj *models.GalleryChapter) (*time.Time, error) {
	return &obj.CreatedAt.Timestamp, nil
}

func (r *galleryChapterResolver) UpdatedAt(ctx context.Context, obj *models.GalleryChapter) (*time.Time, error) {
	return &obj.UpdatedAt.Timestamp, nil
}
