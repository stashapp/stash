package models

import "context"

// MovieGetter provides methods to get movies by ID.
type MovieGetter interface {
	// TODO - rename this to Find and remove existing method
	FindMany(ctx context.Context, ids []int) ([]*Movie, error)
	Find(ctx context.Context, id int) (*Movie, error)
}

// MovieFinder provides methods to find movies.
type MovieFinder interface {
	MovieGetter
	FindByPerformerID(ctx context.Context, performerID int) ([]*Movie, error)
	FindByStudioID(ctx context.Context, studioID int) ([]*Movie, error)
	FindByName(ctx context.Context, name string, nocase bool) (*Movie, error)
	FindByNames(ctx context.Context, names []string, nocase bool) ([]*Movie, error)
}

// MovieQueryer provides methods to query movies.
type MovieQueryer interface {
	Query(ctx context.Context, movieFilter *MovieFilterType, findFilter *FindFilterType) ([]*Movie, int, error)
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
	Create(ctx context.Context, newMovie *Movie) error
}

// MovieUpdater provides methods to update movies.
type MovieUpdater interface {
	Update(ctx context.Context, updatedMovie *Movie) error
	UpdatePartial(ctx context.Context, id int, updatedMovie MoviePartial) (*Movie, error)
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

	All(ctx context.Context) ([]*Movie, error)
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
