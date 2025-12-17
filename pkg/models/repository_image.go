package models

import "context"

// ImageGetter provides methods to get images by ID.
type ImageGetter interface {
	// TODO - rename this to Find and remove existing method
	FindMany(ctx context.Context, ids []int) ([]*Image, error)
	Find(ctx context.Context, id int) (*Image, error)
}

// ImageFinder provides methods to find images.
type ImageFinder interface {
	ImageGetter
	FindByFingerprints(ctx context.Context, fp []Fingerprint) ([]*Image, error)
	FindByChecksum(ctx context.Context, checksum string) ([]*Image, error)
	FindByFileID(ctx context.Context, fileID FileID) ([]*Image, error)
	FindByFolderID(ctx context.Context, fileID FolderID) ([]*Image, error)
	FindByZipFileID(ctx context.Context, zipFileID FileID) ([]*Image, error)
	FindByGalleryID(ctx context.Context, galleryID int) ([]*Image, error)
	FindByGalleryIDIndex(ctx context.Context, galleryID int, index uint) (*Image, error)
}

// ImageQueryer provides methods to query images.
type ImageQueryer interface {
	Query(ctx context.Context, options ImageQueryOptions) (*ImageQueryResult, error)
	QueryCount(ctx context.Context, imageFilter *ImageFilterType, findFilter *FindFilterType) (int, error)
}

type GalleryCoverFinder interface {
	CoverByGalleryID(ctx context.Context, galleryId int) (*Image, error)
}

// ImageCounter provides methods to count images.
type ImageCounter interface {
	Count(ctx context.Context) (int, error)
	CountByFileID(ctx context.Context, fileID FileID) (int, error)
	CountByGalleryID(ctx context.Context, galleryID int) (int, error)
	OCount(ctx context.Context) (int, error)
	OCountByPerformerID(ctx context.Context, performerID int) (int, error)
	OCountByStudioID(ctx context.Context, studioID int) (int, error)
}

// ImageCreator provides methods to create images.
type ImageCreator interface {
	Create(ctx context.Context, newImage *Image, fileIDs []FileID) error
}

// ImageUpdater provides methods to update images.
type ImageUpdater interface {
	Update(ctx context.Context, updatedImage *Image) error
	UpdatePartial(ctx context.Context, id int, partial ImagePartial) (*Image, error)
	UpdatePerformers(ctx context.Context, imageID int, performerIDs []int) error
	UpdateTags(ctx context.Context, imageID int, tagIDs []int) error
}

// ImageDestroyer provides methods to destroy images.
type ImageDestroyer interface {
	Destroy(ctx context.Context, id int) error
}

type ImageCreatorUpdater interface {
	ImageCreator
	ImageUpdater
}

// ImageReader provides all methods to read images.
type ImageReader interface {
	ImageFinder
	ImageQueryer
	ImageCounter

	URLLoader
	FileIDLoader
	GalleryIDLoader
	PerformerIDLoader
	TagIDLoader
	FileLoader

	GalleryCoverFinder

	All(ctx context.Context) ([]*Image, error)
	Size(ctx context.Context) (float64, error)
}

// ImageWriter provides all methods to modify images.
type ImageWriter interface {
	ImageCreator
	ImageUpdater
	ImageDestroyer

	AddFileID(ctx context.Context, id int, fileID FileID) error
	RemoveFileID(ctx context.Context, id int, fileID FileID) error
	IncrementOCounter(ctx context.Context, id int) (int, error)
	DecrementOCounter(ctx context.Context, id int) (int, error)
	ResetOCounter(ctx context.Context, id int) (int, error)
}

// ImageReaderWriter provides all image methods.
type ImageReaderWriter interface {
	ImageReader
	ImageWriter
}
