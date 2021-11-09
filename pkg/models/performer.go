package models

type PerformerQueryOptions struct {
	QueryOptions
	PerformerFilter *PerformerFilterType

	TotalDuration bool
	TotalSize     bool
}

type PerformerQueryResult struct {
	QueryResult
	TotalDuration float64
	TotalSize     float64

	finder     PerformerFinder
	performers []*Performer
	resolveErr error
}

type PerformerFinder interface {
	// TODO - rename this to Find and remove existing method
	FindMany(ids []int) ([]*Performer, error)
}

func (r *PerformerQueryResult) Resolve() ([]*Performer, error) {
	// cache results
	if r.performers == nil && r.resolveErr == nil {
		r.performers, r.resolveErr = r.finder.FindMany(r.IDs)
	}
	return r.performers, r.resolveErr
}

type PerformerReader interface {
	Find(id int) (*Performer, error)
	FindMany(ids []int) ([]*Performer, error)
	FindBySceneID(sceneID int) ([]*Performer, error)
	FindNamesBySceneID(sceneID int) ([]*Performer, error)
	FindByImageID(imageID int) ([]*Performer, error)
	FindByGalleryID(galleryID int) ([]*Performer, error)
	FindByNames(names []string, nocase bool) ([]*Performer, error)
	FindByStashIDStatus(hasStashID bool, stashboxEndpoint string) ([]*Performer, error)
	CountByTagID(tagID int) (int, error)
	Count() (int, error)
	All() ([]*Performer, error)
	// TODO - this interface is temporary until the filter schema can fully
	// support the query needed
	QueryForAutoTag(words []string) ([]*Performer, error)
	Query(performerFilter *PerformerFilterType, findFilter *FindFilterType) ([]*Performer, int, error)
	GetImage(performerID int) ([]byte, error)
	GetStashIDs(performerID int) ([]*StashID, error)
	GetTagIDs(performerID int) ([]int, error)
}

type PerformerWriter interface {
	Create(newPerformer Performer) (*Performer, error)
	Update(updatedPerformer PerformerPartial) (*Performer, error)
	UpdateFull(updatedPerformer Performer) (*Performer, error)
	Destroy(id int) error
	UpdateImage(performerID int, image []byte) error
	DestroyImage(performerID int) error
	UpdateStashIDs(performerID int, stashIDs []StashID) error
	UpdateTags(performerID int, tagIDs []int) error
}

type PerformerReaderWriter interface {
	PerformerReader
	PerformerWriter
}
