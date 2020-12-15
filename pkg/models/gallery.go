package models

type GalleryReader interface {
	Find(id int) (*Gallery, error)
	FindMany(ids []int) ([]*Gallery, error)
	FindByChecksum(checksum string) (*Gallery, error)
	FindByPath(path string) (*Gallery, error)
	FindBySceneID(sceneID int) (*Gallery, error)
	FindByImageID(imageID int) ([]*Gallery, error)
	// ValidGalleriesForScenePath(scenePath string) ([]*Gallery, error)
	// Count() (int, error)
	All() ([]*Gallery, error)
	// Query(galleryFilter *GalleryFilterType, findFilter *FindFilterType) ([]*Gallery, int)
	GetPerformerIDs(galleryID int) ([]int, error)
	GetTagIDs(galleryID int) ([]int, error)
	GetImageIDs(galleryID int) ([]int, error)
}

type GalleryWriter interface {
	Create(newGallery Gallery) (*Gallery, error)
	Update(updatedGallery Gallery) (*Gallery, error)
	UpdatePartial(updatedGallery GalleryPartial) (*Gallery, error)
	Destroy(id int) error
	ClearGalleryId(sceneID int) error
	UpdatePerformers(galleryID int, performerIDs []int) error
	UpdateTags(galleryID int, tagIDs []int) error
	UpdateImages(galleryID int, imageIDs []int) error
}

type GalleryReaderWriter interface {
	GalleryReader
	GalleryWriter
}
