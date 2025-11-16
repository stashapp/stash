package models

import "context"

// TagGetter provides methods to get tags by ID.
type TagGetter interface {
	// TODO - rename this to Find and remove existing method
	FindMany(ctx context.Context, ids []int) ([]*Tag, error)
	Find(ctx context.Context, id int) (*Tag, error)
}

// TagFinder provides methods to find tags.
type TagFinder interface {
	TagGetter
	FindAllAncestors(ctx context.Context, tagID int, excludeIDs []int) ([]*TagPath, error)
	FindAllDescendants(ctx context.Context, tagID int, excludeIDs []int) ([]*TagPath, error)
	FindByParentTagID(ctx context.Context, parentID int) ([]*Tag, error)
	FindByChildTagID(ctx context.Context, childID int) ([]*Tag, error)
	FindBySceneID(ctx context.Context, sceneID int) ([]*Tag, error)
	FindByImageID(ctx context.Context, imageID int) ([]*Tag, error)
	FindByGalleryID(ctx context.Context, galleryID int) ([]*Tag, error)
	FindByPerformerID(ctx context.Context, performerID int) ([]*Tag, error)
	FindByGroupID(ctx context.Context, groupID int) ([]*Tag, error)
	FindBySceneMarkerID(ctx context.Context, sceneMarkerID int) ([]*Tag, error)
	FindByStudioID(ctx context.Context, studioID int) ([]*Tag, error)
	FindByName(ctx context.Context, name string, nocase bool) (*Tag, error)
	FindByNames(ctx context.Context, names []string, nocase bool) ([]*Tag, error)
	FindByStashID(ctx context.Context, stashID StashID) ([]*Tag, error)
}

// TagQueryer provides methods to query tags.
type TagQueryer interface {
	Query(ctx context.Context, tagFilter *TagFilterType, findFilter *FindFilterType) ([]*Tag, int, error)
}

type TagAutoTagQueryer interface {
	TagQueryer
	AliasLoader

	// TODO - this interface is temporary until the filter schema can fully
	// support the query needed
	QueryForAutoTag(ctx context.Context, words []string) ([]*Tag, error)
}

// TagCounter provides methods to count tags.
type TagCounter interface {
	Count(ctx context.Context) (int, error)
	CountByParentTagID(ctx context.Context, parentID int) (int, error)
	CountByChildTagID(ctx context.Context, childID int) (int, error)
}

// TagCreator provides methods to create tags.
type TagCreator interface {
	Create(ctx context.Context, newTag *Tag) error
}

// TagUpdater provides methods to update tags.
type TagUpdater interface {
	Update(ctx context.Context, updatedTag *Tag) error
	UpdatePartial(ctx context.Context, id int, updateTag TagPartial) (*Tag, error)
	UpdateAliases(ctx context.Context, tagID int, aliases []string) error
	UpdateImage(ctx context.Context, tagID int, image []byte) error
	UpdateParentTags(ctx context.Context, tagID int, parentIDs []int) error
	UpdateChildTags(ctx context.Context, tagID int, parentIDs []int) error
}

// TagDestroyer provides methods to destroy tags.
type TagDestroyer interface {
	Destroy(ctx context.Context, id int) error
}

type TagFinderCreator interface {
	TagFinder
	TagCreator
}

type TagCreatorUpdater interface {
	TagCreator
	TagUpdater
}

// TagReader provides all methods to read tags.
type TagReader interface {
	TagFinder
	TagQueryer
	TagAutoTagQueryer
	TagCounter

	AliasLoader
	TagRelationLoader
	StashIDLoader

	All(ctx context.Context) ([]*Tag, error)
	GetImage(ctx context.Context, tagID int) ([]byte, error)
	HasImage(ctx context.Context, tagID int) (bool, error)
}

// TagWriter provides all methods to modify tags.
type TagWriter interface {
	TagCreator
	TagUpdater
	TagDestroyer

	Merge(ctx context.Context, source []int, destination int) error
}

// TagReaderWriter provides all tags methods.
type TagReaderWriter interface {
	TagReader
	TagWriter
}
