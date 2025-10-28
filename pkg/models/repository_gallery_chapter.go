package models

import "context"

// GalleryChapterGetter provides methods to get gallery chapters by ID.
type GalleryChapterGetter interface {
	// TODO - rename this to Find and remove existing method
	FindMany(ctx context.Context, ids []int) ([]*GalleryChapter, error)
	Find(ctx context.Context, id int) (*GalleryChapter, error)
}

// GalleryChapterFinder provides methods to find gallery chapters.
type GalleryChapterFinder interface {
	GalleryChapterGetter
	FindByGalleryID(ctx context.Context, galleryID int) ([]*GalleryChapter, error)
}

// GalleryChapterCreator provides methods to create gallery chapters.
type GalleryChapterCreator interface {
	Create(ctx context.Context, newGalleryChapter *GalleryChapter) error
}

// GalleryChapterUpdater provides methods to update gallery chapters.
type GalleryChapterUpdater interface {
	Update(ctx context.Context, updatedGalleryChapter *GalleryChapter) error
	UpdatePartial(ctx context.Context, id int, updatedGalleryChapter GalleryChapterPartial) (*GalleryChapter, error)
}

// GalleryChapterDestroyer provides methods to destroy gallery chapters.
type GalleryChapterDestroyer interface {
	Destroy(ctx context.Context, id int) error
}

type GalleryChapterCreatorUpdater interface {
	GalleryChapterCreator
	GalleryChapterUpdater
}

// GalleryChapterReader provides all methods to read gallery chapters.
type GalleryChapterReader interface {
	GalleryChapterFinder
}

// GalleryChapterWriter provides all methods to modify gallery chapters.
type GalleryChapterWriter interface {
	GalleryChapterCreator
	GalleryChapterUpdater
	GalleryChapterDestroyer
}

// GalleryChapterReaderWriter provides all gallery chapter methods.
type GalleryChapterReaderWriter interface {
	GalleryChapterReader
	GalleryChapterWriter
}
