package models

import "context"

// StudioGetter provides methods to get studios by ID.
type StudioGetter interface {
	// TODO - rename this to Find and remove existing method
	FindMany(ctx context.Context, ids []int) ([]*Studio, error)
	Find(ctx context.Context, id int) (*Studio, error)
}

// StudioFinder provides methods to find studios.
type StudioFinder interface {
	StudioGetter
	FindChildren(ctx context.Context, id int) ([]*Studio, error)
	FindBySceneID(ctx context.Context, sceneID int) (*Studio, error)
	FindByStashID(ctx context.Context, stashID StashID) ([]*Studio, error)
	FindByStashIDStatus(ctx context.Context, hasStashID bool, stashboxEndpoint string) ([]*Studio, error)
	FindByName(ctx context.Context, name string, nocase bool) (*Studio, error)
}

// StudioQueryer provides methods to query studios.
type StudioQueryer interface {
	Query(ctx context.Context, studioFilter *StudioFilterType, findFilter *FindFilterType) ([]*Studio, int, error)
	QueryCount(ctx context.Context, studioFilter *StudioFilterType, findFilter *FindFilterType) (int, error)
}

type StudioAutoTagQueryer interface {
	StudioQueryer
	AliasLoader

	// TODO - this interface is temporary until the filter schema can fully
	// support the query needed
	QueryForAutoTag(ctx context.Context, words []string) ([]*Studio, error)
}

// StudioCounter provides methods to count studios.
type StudioCounter interface {
	Count(ctx context.Context) (int, error)
	CountByTagID(ctx context.Context, tagID int) (int, error)
}

// StudioCreator provides methods to create studios.
type StudioCreator interface {
	Create(ctx context.Context, newStudio *Studio) error
}

// StudioUpdater provides methods to update studios.
type StudioUpdater interface {
	Update(ctx context.Context, updatedStudio *Studio) error
	UpdatePartial(ctx context.Context, updatedStudio StudioPartial) (*Studio, error)
	UpdateImage(ctx context.Context, studioID int, image []byte) error
}

// StudioDestroyer provides methods to destroy studios.
type StudioDestroyer interface {
	Destroy(ctx context.Context, id int) error
}

type StudioFinderCreator interface {
	StudioFinder
	StudioCreator
}

type StudioCreatorUpdater interface {
	StudioCreator
	StudioUpdater
}

// StudioReader provides all methods to read studios.
type StudioReader interface {
	StudioFinder
	StudioQueryer
	StudioAutoTagQueryer
	StudioCounter

	AliasLoader
	StashIDLoader
	TagIDLoader

	All(ctx context.Context) ([]*Studio, error)
	GetImage(ctx context.Context, studioID int) ([]byte, error)
	HasImage(ctx context.Context, studioID int) (bool, error)
}

// StudioWriter provides all methods to modify studios.
type StudioWriter interface {
	StudioCreator
	StudioUpdater
	StudioDestroyer
}

// StudioReaderWriter provides all studio methods.
type StudioReaderWriter interface {
	StudioReader
	StudioWriter
}
