package database

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

type GalleryChapterStore interface {
	Create(ctx context.Context, newObject *models.GalleryChapter) error
	Destroy(ctx context.Context, id int) error
	Find(ctx context.Context, id int) (*models.GalleryChapter, error)
	FindByGalleryID(ctx context.Context, galleryID int) ([]*models.GalleryChapter, error)
	FindMany(ctx context.Context, ids []int) ([]*models.GalleryChapter, error)
	Update(ctx context.Context, updatedObject *models.GalleryChapter) error
	UpdatePartial(ctx context.Context, id int, partial models.GalleryChapterPartial) (*models.GalleryChapter, error)
}
