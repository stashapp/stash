package models

type SceneSegmentReader interface {
	Find(id int) (*SceneSegment, error)
	FindMany(ids []int) ([]*SceneSegment, error)
	FindBySceneID(sceneID int) ([]*SceneSegment, error)
	All() ([]*SceneSegment, error)
}

type SceneSegmentWriter interface {
	Create(newSegment *SceneSegment) error
	Update(id int, updatedSegment SceneSegmentPartial) error
	Destroy(id int) error
}

type SceneSegmentCreator interface {
	Create(newSegment *SceneSegment) error
}

type SceneSegmentUpdater interface {
	Update(id int, updatedSegment SceneSegmentPartial) error
}

type SceneSegmentDestroyer interface {
	Destroy(id int) error
}

type SceneSegmentFinder interface {
	SceneSegmentReader
}

type SceneSegmentReaderWriter interface {
	SceneSegmentReader
	SceneSegmentWriter
}
