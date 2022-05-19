package models

import "context"

type GalleryFilterType struct {
	And     *GalleryFilterType    `json:"AND"`
	Or      *GalleryFilterType    `json:"OR"`
	Not     *GalleryFilterType    `json:"NOT"`
	Title   *StringCriterionInput `json:"title"`
	Details *StringCriterionInput `json:"details"`
	// Filter by file checksum
	Checksum *StringCriterionInput `json:"checksum"`
	// Filter by path
	Path *StringCriterionInput `json:"path"`
	// Filter to only include galleries missing this property
	IsMissing *string `json:"is_missing"`
	// Filter to include/exclude galleries that were created from zip
	IsZip *bool `json:"is_zip"`
	// Filter by rating
	Rating *IntCriterionInput `json:"rating"`
	// Filter by organized
	Organized *bool `json:"organized"`
	// Filter by average image resolution
	AverageResolution *ResolutionCriterionInput `json:"average_resolution"`
	// Filter to only include galleries with this studio
	Studios *HierarchicalMultiCriterionInput `json:"studios"`
	// Filter to only include galleries with these tags
	Tags *HierarchicalMultiCriterionInput `json:"tags"`
	// Filter by tag count
	TagCount *IntCriterionInput `json:"tag_count"`
	// Filter to only include galleries with performers with these tags
	PerformerTags *HierarchicalMultiCriterionInput `json:"performer_tags"`
	// Filter to only include galleries with these performers
	Performers *MultiCriterionInput `json:"performers"`
	// Filter by performer count
	PerformerCount *IntCriterionInput `json:"performer_count"`
	// Filter galleries that have performers that have been favorited
	PerformerFavorite *bool `json:"performer_favorite"`
	// Filter galleries by performer age at time of gallery
	PerformerAge *IntCriterionInput `json:"performer_age"`
	// Filter by number of images in this gallery
	ImageCount *IntCriterionInput `json:"image_count"`
	// Filter by url
	URL *StringCriterionInput `json:"url"`
}

type GalleryUpdateInput struct {
	ClientMutationID *string  `json:"clientMutationId"`
	ID               string   `json:"id"`
	Title            *string  `json:"title"`
	URL              *string  `json:"url"`
	Date             *string  `json:"date"`
	Details          *string  `json:"details"`
	Rating           *int     `json:"rating"`
	Organized        *bool    `json:"organized"`
	SceneIds         []string `json:"scene_ids"`
	StudioID         *string  `json:"studio_id"`
	TagIds           []string `json:"tag_ids"`
	PerformerIds     []string `json:"performer_ids"`
}

type GalleryDestroyInput struct {
	Ids []string `json:"ids"`
	// If true, then the zip file will be deleted if the gallery is zip-file-based.
	// If gallery is folder-based, then any files not associated with other
	// galleries will be deleted, along with the folder, if it is not empty.
	DeleteFile      *bool `json:"delete_file"`
	DeleteGenerated *bool `json:"delete_generated"`
}

type GalleryReader interface {
	Find(ctx context.Context, id int) (*Gallery, error)
	FindMany(ctx context.Context, ids []int) ([]*Gallery, error)
	FindByChecksum(ctx context.Context, checksum string) (*Gallery, error)
	FindByChecksums(ctx context.Context, checksums []string) ([]*Gallery, error)
	FindByPath(ctx context.Context, path string) (*Gallery, error)
	FindBySceneID(ctx context.Context, sceneID int) ([]*Gallery, error)
	FindByImageID(ctx context.Context, imageID int) ([]*Gallery, error)
	Count(ctx context.Context) (int, error)
	All(ctx context.Context) ([]*Gallery, error)
	Query(ctx context.Context, galleryFilter *GalleryFilterType, findFilter *FindFilterType) ([]*Gallery, int, error)
	QueryCount(ctx context.Context, galleryFilter *GalleryFilterType, findFilter *FindFilterType) (int, error)
	GetPerformerIDs(ctx context.Context, galleryID int) ([]int, error)
	GetTagIDs(ctx context.Context, galleryID int) ([]int, error)
	GetSceneIDs(ctx context.Context, galleryID int) ([]int, error)
	GetImageIDs(ctx context.Context, galleryID int) ([]int, error)
}

type GalleryWriter interface {
	Create(ctx context.Context, newGallery Gallery) (*Gallery, error)
	Update(ctx context.Context, updatedGallery Gallery) (*Gallery, error)
	UpdatePartial(ctx context.Context, updatedGallery GalleryPartial) (*Gallery, error)
	UpdateFileModTime(ctx context.Context, id int, modTime NullSQLiteTimestamp) error
	Destroy(ctx context.Context, id int) error
	UpdatePerformers(ctx context.Context, galleryID int, performerIDs []int) error
	UpdateTags(ctx context.Context, galleryID int, tagIDs []int) error
	UpdateScenes(ctx context.Context, galleryID int, sceneIDs []int) error
	UpdateImages(ctx context.Context, galleryID int, imageIDs []int) error
}

type GalleryReaderWriter interface {
	GalleryReader
	GalleryWriter
}
