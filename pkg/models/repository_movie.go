package models

import "context"

// MovieGetter provides methods to get movies by ID.
type MovieGetter interface {
	// TODO - rename this to Find and remove existing method
	FindMany(ctx context.Context, ids []int) ([]*Group, error)
	Find(ctx context.Context, id int) (*Group, error)
}

// MovieFinder provides methods to find movies.
type MovieFinder interface {
	MovieGetter
	FindByPerformerID(ctx context.Context, performerID int) ([]*Group, error)
	FindByStudioID(ctx context.Context, studioID int) ([]*Group, error)
	FindByName(ctx context.Context, name string, nocase bool) (*Group, error)
	FindByNames(ctx context.Context, names []string, nocase bool) ([]*Group, error)
}

// MovieQueryer provides methods to query movies.
type MovieQueryer interface {
	Query(ctx context.Context, movieFilter *MovieFilterType, findFilter *FindFilterType) ([]*Group, int, error)
	QueryCount(ctx context.Context, movieFilter *MovieFilterType, findFilter *FindFilterType) (int, error)
}

// MovieCounter provides methods to count movies.
type MovieCounter interface {
	Count(ctx context.Context) (int, error)
	CountByPerformerID(ctx context.Context, performerID int) (int, error)
	CountByStudioID(ctx context.Context, studioID int) (int, error)
}

// MovieCreator provides methods to create movies.
type MovieCreator interface {
	Create(ctx context.Context, newMovie *Group) error
}

// MovieUpdater provides methods to update movies.
type MovieUpdater interface {
	Update(ctx context.Context, updatedMovie *Group) error
	UpdatePartial(ctx context.Context, id int, updatedMovie GroupPartial) (*Group, error)
	UpdateFrontImage(ctx context.Context, movieID int, frontImage []byte) error
	UpdateBackImage(ctx context.Context, movieID int, backImage []byte) error
}

// MovieDestroyer provides methods to destroy movies.
type MovieDestroyer interface {
	Destroy(ctx context.Context, id int) error
}

type MovieCreatorUpdater interface {
	MovieCreator
	MovieUpdater
}

type MovieFinderCreator interface {
	MovieFinder
	MovieCreator
}

// MovieReader provides all methods to read movies.
type MovieReader interface {
	MovieFinder
	MovieQueryer
	MovieCounter
	URLLoader
	TagIDLoader

	All(ctx context.Context) ([]*Group, error)
	GetFrontImage(ctx context.Context, movieID int) ([]byte, error)
	HasFrontImage(ctx context.Context, movieID int) (bool, error)
	GetBackImage(ctx context.Context, movieID int) ([]byte, error)
	HasBackImage(ctx context.Context, movieID int) (bool, error)
}

// MovieWriter provides all methods to modify movies.
type MovieWriter interface {
	MovieCreator
	MovieUpdater
	MovieDestroyer
}

// MovieReaderWriter provides all movie methods.
type MovieReaderWriter interface {
	MovieReader
	MovieWriter
}
