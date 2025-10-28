package models

import "context"

// SceneMarkerGetter provides methods to get scene markers by ID.
type SceneMarkerGetter interface {
	// TODO - rename this to Find and remove existing method
	FindMany(ctx context.Context, ids []int) ([]*SceneMarker, error)
	Find(ctx context.Context, id int) (*SceneMarker, error)
}

// SceneMarkerFinder provides methods to find scene markers.
type SceneMarkerFinder interface {
	SceneMarkerGetter
	FindBySceneID(ctx context.Context, sceneID int) ([]*SceneMarker, error)
}

// SceneMarkerQueryer provides methods to query scene markers.
type SceneMarkerQueryer interface {
	Query(ctx context.Context, sceneMarkerFilter *SceneMarkerFilterType, findFilter *FindFilterType) ([]*SceneMarker, int, error)
	QueryCount(ctx context.Context, sceneMarkerFilter *SceneMarkerFilterType, findFilter *FindFilterType) (int, error)
}

// SceneMarkerCounter provides methods to count scene markers.
type SceneMarkerCounter interface {
	Count(ctx context.Context) (int, error)
	CountByTagID(ctx context.Context, tagID int) (int, error)
}

// SceneMarkerCreator provides methods to create scene markers.
type SceneMarkerCreator interface {
	Create(ctx context.Context, newSceneMarker *SceneMarker) error
}

// SceneMarkerUpdater provides methods to update scene markers.
type SceneMarkerUpdater interface {
	Update(ctx context.Context, updatedSceneMarker *SceneMarker) error
	UpdatePartial(ctx context.Context, id int, updatedSceneMarker SceneMarkerPartial) (*SceneMarker, error)
	UpdateTags(ctx context.Context, markerID int, tagIDs []int) error
}

// SceneMarkerDestroyer provides methods to destroy scene markers.
type SceneMarkerDestroyer interface {
	Destroy(ctx context.Context, id int) error
}

type SceneMarkerCreatorUpdater interface {
	SceneMarkerCreator
	SceneMarkerUpdater
}

// SceneMarkerReader provides all methods to read scene markers.
type SceneMarkerReader interface {
	SceneMarkerFinder
	SceneMarkerQueryer
	SceneMarkerCounter

	TagIDLoader

	All(ctx context.Context) ([]*SceneMarker, error)
	Wall(ctx context.Context, q *string) ([]*SceneMarker, error)
	GetMarkerStrings(ctx context.Context, q *string, sort *string) ([]*MarkerStringsResultType, error)
}

// SceneMarkerWriter provides all methods to modify scene markers.
type SceneMarkerWriter interface {
	SceneMarkerCreator
	SceneMarkerUpdater
	SceneMarkerDestroyer
}

// SceneMarkerReaderWriter provides all scene marker methods.
type SceneMarkerReaderWriter interface {
	SceneMarkerReader
	SceneMarkerWriter
}
