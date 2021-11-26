package models

type GalleryReader interface {
	Find(id int) (*Gallery, error)
	FindMany(ids []int) ([]*Gallery, error)
	FindByChecksum(checksum string) (*Gallery, error)
	FindByChecksums(checksums []string) ([]*Gallery, error)
	FindByPath(path string) (*Gallery, error)
	FindByFileID(fileID int) ([]*Gallery, error)
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

	FileJoinReader
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

	FileJoinWriter
}

type GalleryReaderWriter interface {
	GalleryReader
	GalleryWriter
}
