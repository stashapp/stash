package models

import "context"

// SceneFilterGetter provides methods to get scene filters by ID.
type SceneFilterGetter interface {
	// TODO - rename this to Find and remove existing method
	FindMany(ctx context.Context, ids []int) ([]*SceneFilter, error)
	Find(ctx context.Context, id int) (*SceneFilter, error)
}

// SceneFilterFinder provides methods to find scene filters.
type SceneFilterFinder interface {
	SceneFilterGetter
	FindBySceneID(ctx context.Context, sceneID int) ([]*SceneFilter, error)
}

// SceneFilterQueryer provides methods to query scene filters.
type SceneFilterQueryer interface {
	Query(ctx context.Context, sceneFilterFilter *SceneFilterFilterType, findFilter *FindFilterType) ([]*SceneFilter, int, error)
	QueryCount(ctx context.Context, sceneFilterFilter *SceneFilterFilterType, findFilter *FindFilterType) (int, error)
}

// SceneFilterCounter provides methods to count scene filters.
type SceneFilterCounter interface {
	Count(ctx context.Context) (int, error)
}

// SceneFilterCreator provides methods to create scene filters.
type SceneFilterCreator interface {
	Create(ctx context.Context, newSceneFilter *SceneFilter) error
}

// SceneFilterUpdater provides methods to update scene filters.
type SceneFilterUpdater interface {
	Update(ctx context.Context, updatedSceneFilter *SceneFilter) error
}

// SceneFilterDestroyer provides methods to destroy scene filters.
type SceneFilterDestroyer interface {
	Destroy(ctx context.Context, id int) error
}

type SceneFilterCreatorUpdater interface {
	SceneFilterCreator
	SceneFilterUpdater
}

// SceneFilterReader provides all methods to read scene filters.
type SceneFilterReader interface {
	SceneFilterFinder
	SceneFilterQueryer
	SceneFilterCounter

	All(ctx context.Context) ([]*SceneFilter, error)
}

// SceneFilterWriter provides all methods to modify scene filters.
type SceneFilterWriter interface {
	SceneFilterCreator
	SceneFilterUpdater
	SceneFilterDestroyer
}

// SceneFilterReaderWriter provides all scene filter methods.
type SceneFilterReaderWriter interface {
	SceneFilterReader
	SceneFilterWriter
}
