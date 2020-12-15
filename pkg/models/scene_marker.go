package models

type SceneMarkerReader interface {
	Find(id int) (*SceneMarker, error)
	// FindMany(ids []int) ([]*SceneMarker, error)
	FindBySceneID(sceneID int) ([]*SceneMarker, error)
	CountByTagID(tagID int) (int, error)
	// GetMarkerStrings(q *string, sort *string) ([]*MarkerStringsResultType, error)
	// Wall(q *string) ([]*SceneMarker, error)
	// Query(sceneMarkerFilter *SceneMarkerFilterType, findFilter *FindFilterType) ([]*SceneMarker, int)
	GetTagIDs(imageID int) ([]int, error)
}

type SceneMarkerWriter interface {
	Create(newSceneMarker SceneMarker) (*SceneMarker, error)
	Update(updatedSceneMarker SceneMarker) (*SceneMarker, error)
	Destroy(id int) error
	UpdateTags(markerID int, tagIDs []int) error
}

type SceneMarkerReaderWriter interface {
	SceneMarkerReader
	SceneMarkerWriter
}
