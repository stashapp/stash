package models

type ImageFilterType struct {
	And   *ImageFilterType      `json:"AND"`
	Or    *ImageFilterType      `json:"OR"`
	Not   *ImageFilterType      `json:"NOT"`
	Title *StringCriterionInput `json:"title"`
	// Filter by file checksum
	Checksum *StringCriterionInput `json:"checksum"`
	// Filter by path
	Path *StringCriterionInput `json:"path"`
	// Filter by rating
	Rating *IntCriterionInput `json:"rating"`
	// Filter by organized
	Organized *bool `json:"organized"`
	// Filter by o-counter
	OCounter *IntCriterionInput `json:"o_counter"`
	// Filter by resolution
	Resolution *ResolutionCriterionInput `json:"resolution"`
	// Filter to only include images missing this property
	IsMissing *string `json:"is_missing"`
	// Filter to only include images with this studio
	Studios *HierarchicalMultiCriterionInput `json:"studios"`
	// Filter to only include images with these tags
	Tags *HierarchicalMultiCriterionInput `json:"tags"`
	// Filter by tag count
	TagCount *IntCriterionInput `json:"tag_count"`
	// Filter to only include images with performers with these tags
	PerformerTags *HierarchicalMultiCriterionInput `json:"performer_tags"`
	// Filter to only include images with these performers
	Performers *MultiCriterionInput `json:"performers"`
	// Filter by performer count
	PerformerCount *IntCriterionInput `json:"performer_count"`
	// Filter images that have performers that have been favorited
	PerformerFavorite *bool `json:"performer_favorite"`
	// Filter to only include images with these galleries
	Galleries *MultiCriterionInput `json:"galleries"`
}

type ImageDestroyInput struct {
	ID              string `json:"id"`
	DeleteFile      *bool  `json:"delete_file"`
	DeleteGenerated *bool  `json:"delete_generated"`
}

type ImagesDestroyInput struct {
	Ids             []string `json:"ids"`
	DeleteFile      *bool    `json:"delete_file"`
	DeleteGenerated *bool    `json:"delete_generated"`
}

type ImageQueryOptions struct {
	QueryOptions
	ImageFilter *ImageFilterType

	Megapixels bool
	TotalSize  bool
}

type ImageQueryResult struct {
	QueryResult
	Megapixels float64
	TotalSize  float64

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
