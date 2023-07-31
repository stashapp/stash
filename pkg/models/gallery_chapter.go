package models

import "context"

type GalleryChapterReader interface {
	Find(ctx context.Context, id int) (*GalleryChapter, error)
	FindMany(ctx context.Context, ids []int) ([]*GalleryChapter, error)
	FindByGalleryID(ctx context.Context, galleryID int) ([]*GalleryChapter, error)
}

type GalleryChapterWriter interface {
	Create(ctx context.Context, newGalleryChapter *GalleryChapter) error
	Update(ctx context.Context, updatedGalleryChapter *GalleryChapter) error
	UpdatePartial(ctx context.Context, id int, updatedGalleryChapter GalleryChapterPartial) (*GalleryChapter, error)
	Destroy(ctx context.Context, id int) error
}

type GalleryChapterReaderWriter interface {
	GalleryChapterReader
	GalleryChapterWriter
}
