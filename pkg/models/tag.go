package models

type TagReader interface {
	Find(id int) (*Tag, error)
	FindMany(ids []int) ([]*Tag, error)
	FindBySceneID(sceneID int) ([]*Tag, error)
	FindByPerformerID(performerID int) ([]*Tag, error)
	FindBySceneMarkerID(sceneMarkerID int) ([]*Tag, error)
	FindByImageID(imageID int) ([]*Tag, error)
	FindByGalleryID(galleryID int) ([]*Tag, error)
	FindByName(name string, nocase bool) (*Tag, error)
	FindByNames(names []string, nocase bool) ([]*Tag, error)
	Count() (int, error)
	All() ([]*Tag, error)
	AllSlim() ([]*Tag, error)
	Query(tagFilter *TagFilterType, findFilter *FindFilterType) ([]*Tag, int, error)
	GetImage(tagID int) ([]byte, error)
}

type TagWriter interface {
	Create(newTag Tag) (*Tag, error)
	Update(updatedTag Tag) (*Tag, error)
	Destroy(id int) error
	UpdateImage(tagID int, image []byte) error
	DestroyImage(tagID int) error
}

type TagReaderWriter interface {
	TagReader
	TagWriter
}
