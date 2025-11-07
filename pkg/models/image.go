package models

import "context"

type ImageFilterType struct {
	OperatorFilter[ImageFilterType]
	ID           *IntCriterionInput    `json:"id"`
	Title        *StringCriterionInput `json:"title"`
	Code         *StringCriterionInput `json:"code"`
	Details      *StringCriterionInput `json:"details"`
	Photographer *StringCriterionInput `json:"photographer"`
	// Filter by file checksum
	Checksum *StringCriterionInput `json:"checksum"`
	// Filter by path
	Path *StringCriterionInput `json:"path"`
	// Filter by file count
	FileCount *IntCriterionInput `json:"file_count"`
	// Filter by rating expressed as 1-100
	Rating100 *IntCriterionInput `json:"rating100"`
	// Filter by date
	Date *DateCriterionInput `json:"date"`
	// Filter by url
	URL *StringCriterionInput `json:"url"`
	// Filter by organized
	Organized *bool `json:"organized"`
	// Filter by o-counter
	OCounter *IntCriterionInput `json:"o_counter"`
	// Filter by resolution
	Resolution *ResolutionCriterionInput `json:"resolution"`
	// Filter by landscape/portrait
	Orientation *OrientationCriterionInput `json:"orientation"`
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
	// Filter images by performer age at time of image
	PerformerAge *IntCriterionInput `json:"performer_age"`
	// Filter to only include images with these galleries
	Galleries *MultiCriterionInput `json:"galleries"`
	// Filter by related galleries that meet this criteria
	GalleriesFilter *GalleryFilterType `json:"galleries_filter"`
	// Filter by related performers that meet this criteria
	PerformersFilter *PerformerFilterType `json:"performers_filter"`
	// Filter by related studios that meet this criteria
	StudiosFilter *StudioFilterType `json:"studios_filter"`
	// Filter by related tags that meet this criteria
	TagsFilter *TagFilterType `json:"tags_filter"`
	// Filter by related files that meet this criteria
	FilesFilter *FileFilterType `json:"files_filter"`
	// Filter by created at
	CreatedAt *TimestampCriterionInput `json:"created_at"`
	// Filter by updated at
	UpdatedAt *TimestampCriterionInput `json:"updated_at"`
}

type ImageUpdateInput struct {
	ClientMutationID *string  `json:"clientMutationId"`
	ID               string   `json:"id"`
	Title            *string  `json:"title"`
	Code             *string  `json:"code"`
	Urls             []string `json:"urls"`
	Date             *string  `json:"date"`
	Details          *string  `json:"details"`
	Photographer     *string  `json:"photographer"`
	Rating100        *int     `json:"rating100"`
	Organized        *bool    `json:"organized"`
	SceneIds         []string `json:"scene_ids"`
	StudioID         *string  `json:"studio_id"`
	TagIds           []string `json:"tag_ids"`
	PerformerIds     []string `json:"performer_ids"`
	GalleryIds       []string `json:"gallery_ids"`
	PrimaryFileID    *string  `json:"primary_file_id"`

	// deprecated
	URL *string `json:"url"`
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
	QueryResult[int]
	Megapixels float64
	TotalSize  float64

	getter     ImageGetter
	images     []*Image
	resolveErr error
}

func NewImageQueryResult(getter ImageGetter) *ImageQueryResult {
	return &ImageQueryResult{
		getter: getter,
	}
}

func (r *ImageQueryResult) Resolve(ctx context.Context) ([]*Image, error) {
	// cache results
	if r.images == nil && r.resolveErr == nil {
		r.images, r.resolveErr = r.getter.FindMany(ctx, r.IDs)
	}
	return r.images, r.resolveErr
}
