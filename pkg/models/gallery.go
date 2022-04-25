package models

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
	Find(id int) (*Gallery, error)
	FindMany(ids []int) ([]*Gallery, error)
	FindByChecksum(checksum string) (*Gallery, error)
	FindByChecksums(checksums []string) ([]*Gallery, error)
	FindByPath(path string) (*Gallery, error)
	FindBySceneID(sceneID int) ([]*Gallery, error)
	FindByImageID(imageID int) ([]*Gallery, error)
	Count() (int, error)
	All() ([]*Gallery, error)
	Query(galleryFilter *GalleryFilterType, findFilter *FindFilterType) ([]*Gallery, int, error)
	QueryCount(galleryFilter *GalleryFilterType, findFilter *FindFilterType) (int, error)
	GetPerformerIDs(galleryID int) ([]int, error)
	GetTagIDs(galleryID int) ([]int, error)
	GetSceneIDs(galleryID int) ([]int, error)
	GetImageIDs(galleryID int) ([]int, error)
}

type GalleryWriter interface {
	Create(newGallery Gallery) (*Gallery, error)
	Update(updatedGallery Gallery) (*Gallery, error)
	UpdatePartial(updatedGallery GalleryPartial) (*Gallery, error)
	UpdateFileModTime(id int, modTime NullSQLiteTimestamp) error
	Destroy(id int) error
	UpdatePerformers(galleryID int, performerIDs []int) error
	UpdateTags(galleryID int, tagIDs []int) error
	UpdateScenes(galleryID int, sceneIDs []int) error
	UpdateImages(galleryID int, imageIDs []int) error
}

type GalleryReaderWriter interface {
	GalleryReader
	GalleryWriter
}
