package models

import "context"

// PerformerGetter provides methods to get performers by ID.
type PerformerGetter interface {
	// TODO - rename this to Find and remove existing method
	FindMany(ctx context.Context, ids []int) ([]*Performer, error)
	Find(ctx context.Context, id int) (*Performer, error)
}

// PerformerFinder provides methods to find performers.
type PerformerFinder interface {
	PerformerGetter
	FindBySceneID(ctx context.Context, sceneID int) ([]*Performer, error)
	FindByImageID(ctx context.Context, imageID int) ([]*Performer, error)
	FindByGalleryID(ctx context.Context, galleryID int) ([]*Performer, error)
	FindByStashID(ctx context.Context, stashID StashID) ([]*Performer, error)
	FindByStashIDStatus(ctx context.Context, hasStashID bool, stashboxEndpoint string) ([]*Performer, error)
	FindByNames(ctx context.Context, names []string, nocase bool) ([]*Performer, error)
}

// PerformerQueryer provides methods to query performers.
type PerformerQueryer interface {
	Query(ctx context.Context, performerFilter *PerformerFilterType, findFilter *FindFilterType) ([]*Performer, int, error)
	QueryCount(ctx context.Context, performerFilter *PerformerFilterType, findFilter *FindFilterType) (int, error)
}

type PerformerAutoTagQueryer interface {
	PerformerQueryer
	AliasLoader

	// TODO - this interface is temporary until the filter schema can fully
	// support the query needed
	QueryForAutoTag(ctx context.Context, words []string) ([]*Performer, error)
}

// PerformerCounter provides methods to count performers.
type PerformerCounter interface {
	Count(ctx context.Context) (int, error)
	CountByTagID(ctx context.Context, tagID int) (int, error)
}

// PerformerCreator provides methods to create performers.
type PerformerCreator interface {
	Create(ctx context.Context, newPerformer *Performer) error
}

// PerformerUpdater provides methods to update performers.
type PerformerUpdater interface {
	Update(ctx context.Context, updatedPerformer *Performer) error
	UpdatePartial(ctx context.Context, id int, updatedPerformer PerformerPartial) (*Performer, error)
	UpdateImage(ctx context.Context, performerID int, image []byte) error
}

// PerformerDestroyer provides methods to destroy performers.
type PerformerDestroyer interface {
	Destroy(ctx context.Context, id int) error
}

type PerformerFinderCreator interface {
	PerformerFinder
	PerformerCreator
}

type PerformerCreatorUpdater interface {
	PerformerCreator
	PerformerUpdater
}

// PerformerReader provides all methods to read performers.
type PerformerReader interface {
	PerformerFinder
	PerformerQueryer
	PerformerAutoTagQueryer
	PerformerCounter

	AliasLoader
	StashIDLoader
	TagIDLoader
	URLLoader

	All(ctx context.Context) ([]*Performer, error)
	GetImage(ctx context.Context, performerID int) ([]byte, error)
	HasImage(ctx context.Context, performerID int) (bool, error)
}

// PerformerWriter provides all methods to modify performers.
type PerformerWriter interface {
	PerformerCreator
	PerformerUpdater
	PerformerDestroyer
}

// PerformerReaderWriter provides all performer methods.
type PerformerReaderWriter interface {
	PerformerReader
	PerformerWriter
}
