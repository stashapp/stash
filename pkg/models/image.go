package models

type ImageQueryOptions struct {
	QueryOptions
	ImageFilter *ImageFilterType

	Megapixels bool
	TotalSize  bool
}

type ImageQueryResult struct {
	QueryResult
	Megapixels float64
	TotalSize  int

	finder     ImageFinder
	images     []*Image
	resolveErr error
}

func NewImageQueryResult(finder ImageFinder) *ImageQueryResult {
	return &ImageQueryResult{
		finder: finder,
	}
}

func (r *ImageQueryResult) Resolve() ([]*Image, error) {
	// cache results
	if r.images == nil && r.resolveErr == nil {
		r.images, r.resolveErr = r.finder.FindMany(r.IDs)
	}
	return r.images, r.resolveErr
}

type ImageFinder interface {
	// TODO - rename to Find and remove existing method
	FindMany(ids []int) ([]*Image, error)
}

type ImageReader interface {
	ImageFinder
	// TODO - remove this in another PR
	Find(id int) (*Image, error)
	FindByChecksum(checksum string) (*Image, error)
	FindByGalleryID(galleryID int) ([]*Image, error)
	CountByGalleryID(galleryID int) (int, error)
	FindByPath(path string) (*Image, error)
	// FindByPerformerID(performerID int) ([]*Image, error)
	// CountByPerformerID(performerID int) (int, error)
	// FindByStudioID(studioID int) ([]*Image, error)
	Count() (int, error)
	Size() (float64, error)
	// SizeCount() (string, error)
	// CountByStudioID(studioID int) (int, error)
	// CountByTagID(tagID int) (int, error)
	All() ([]*Image, error)
	Query(options ImageQueryOptions) (*ImageQueryResult, error)
	QueryCount(imageFilter *ImageFilterType, findFilter *FindFilterType) (int, error)
	GetGalleryIDs(imageID int) ([]int, error)
	GetTagIDs(imageID int) ([]int, error)
	GetPerformerIDs(imageID int) ([]int, error)
}

type ImageWriter interface {
	Create(newImage Image) (*Image, error)
	Update(updatedImage ImagePartial) (*Image, error)
	UpdateFull(updatedImage Image) (*Image, error)
	IncrementOCounter(id int) (int, error)
	DecrementOCounter(id int) (int, error)
	ResetOCounter(id int) (int, error)
	Destroy(id int) error
	UpdateGalleries(imageID int, galleryIDs []int) error
	UpdatePerformers(imageID int, performerIDs []int) error
	UpdateTags(imageID int, tagIDs []int) error
}

type ImageReaderWriter interface {
	ImageReader
	ImageWriter
}
