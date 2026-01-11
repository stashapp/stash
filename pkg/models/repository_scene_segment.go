package models

import "context"

type SceneSegmentReader interface {
	Find(ctx context.Context, id int) (*SceneSegment, error)
	FindMany(ctx context.Context, ids []int) ([]*SceneSegment, error)
	FindBySceneID(ctx context.Context, sceneID int) ([]*SceneSegment, error)
	All(ctx context.Context) ([]*SceneSegment, error)
}

type SceneSegmentWriter interface {
	Create(ctx context.Context, newSegment *SceneSegment) error
	UpdatePartial(ctx context.Context, id int, updatedSegment SceneSegmentPartial) (*SceneSegment, error)
	Destroy(ctx context.Context, id int) error
}

type SceneSegmentCreator interface {
	Create(ctx context.Context, newSegment *SceneSegment) error
}

type SceneSegmentUpdater interface {
	UpdatePartial(ctx context.Context, id int, updatedSegment SceneSegmentPartial) (*SceneSegment, error)
}

type SceneSegmentDestroyer interface {
	Destroy(ctx context.Context, id int) error
}

type SceneSegmentFinder interface {
	SceneSegmentReader
}

type SceneSegmentReaderWriter interface {
	SceneSegmentReader
	SceneSegmentWriter
}
