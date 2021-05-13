package models

type ImageReader interface {
	Find(id int) (*Image, error)
	FindMany(ids []int) ([]*Image, error)
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
	Query(imageFilter *ImageFilterType, findFilter *FindFilterType) ([]*Image, int, error)
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
