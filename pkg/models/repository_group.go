package models

import "context"

// GroupGetter provides methods to get groups by ID.
type GroupGetter interface {
	// TODO - rename this to Find and remove existing method
	FindMany(ctx context.Context, ids []int) ([]*Group, error)
	Find(ctx context.Context, id int) (*Group, error)
}

// GroupFinder provides methods to find groups.
type GroupFinder interface {
	GroupGetter
	FindByPerformerID(ctx context.Context, performerID int) ([]*Group, error)
	FindByStudioID(ctx context.Context, studioID int) ([]*Group, error)
	FindByName(ctx context.Context, name string, nocase bool) (*Group, error)
	FindByNames(ctx context.Context, names []string, nocase bool) ([]*Group, error)
}

// GroupQueryer provides methods to query groups.
type GroupQueryer interface {
	Query(ctx context.Context, groupFilter *GroupFilterType, findFilter *FindFilterType) ([]*Group, int, error)
	QueryCount(ctx context.Context, groupFilter *GroupFilterType, findFilter *FindFilterType) (int, error)
}

// GroupCounter provides methods to count groups.
type GroupCounter interface {
	Count(ctx context.Context) (int, error)
	CountByPerformerID(ctx context.Context, performerID int) (int, error)
	CountByStudioID(ctx context.Context, studioID int) (int, error)
}

// GroupCreator provides methods to create groups.
type GroupCreator interface {
	Create(ctx context.Context, newGroup *Group) error
}

// GroupUpdater provides methods to update groups.
type GroupUpdater interface {
	Update(ctx context.Context, updatedGroup *Group) error
	UpdatePartial(ctx context.Context, id int, updatedGroup GroupPartial) (*Group, error)
	UpdateFrontImage(ctx context.Context, groupID int, frontImage []byte) error
	UpdateBackImage(ctx context.Context, groupID int, backImage []byte) error
}

// GroupDestroyer provides methods to destroy groups.
type GroupDestroyer interface {
	Destroy(ctx context.Context, id int) error
}

type GroupCreatorUpdater interface {
	GroupCreator
	GroupUpdater
}

type GroupFinderCreator interface {
	GroupFinder
	GroupCreator
}

// GroupReader provides all methods to read groups.
type GroupReader interface {
	GroupFinder
	GroupQueryer
	GroupCounter
	URLLoader
	TagIDLoader
	ContainingGroupLoader
	SubGroupLoader

	All(ctx context.Context) ([]*Group, error)
	GetFrontImage(ctx context.Context, groupID int) ([]byte, error)
	HasFrontImage(ctx context.Context, groupID int) (bool, error)
	GetBackImage(ctx context.Context, groupID int) ([]byte, error)
	HasBackImage(ctx context.Context, groupID int) (bool, error)
}

// GroupWriter provides all methods to modify groups.
type GroupWriter interface {
	GroupCreator
	GroupUpdater
	GroupDestroyer
}

// GroupReaderWriter provides all group methods.
type GroupReaderWriter interface {
	GroupReader
	GroupWriter
}
