package models

type PerformerReader interface {
	// Find(id int) (*Performer, error)
	FindMany(ids []int) ([]*Performer, error)
	FindBySceneID(sceneID int) ([]*Performer, error)
	FindNamesBySceneID(sceneID int) ([]*Performer, error)
	FindByImageID(imageID int) ([]*Performer, error)
	FindByGalleryID(galleryID int) ([]*Performer, error)
	FindByNames(names []string, nocase bool) ([]*Performer, error)
	// Count() (int, error)
	All() ([]*Performer, error)
	// AllSlim() ([]*Performer, error)
	// Query(performerFilter *PerformerFilterType, findFilter *FindFilterType) ([]*Performer, int)
	GetPerformerImage(performerID int) ([]byte, error)
}

type PerformerWriter interface {
	Create(newPerformer Performer) (*Performer, error)
	Update(updatedPerformer PerformerPartial) (*Performer, error)
	UpdateFull(updatedPerformer Performer) (*Performer, error)
	// Destroy(id string) error
	UpdatePerformerImage(performerID int, image []byte) error
	// DestroyPerformerImage(performerID int) error
}

type PerformerReaderWriter interface {
	PerformerReader
	PerformerWriter
}
