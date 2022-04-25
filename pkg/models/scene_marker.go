package models

type SceneMarkerFilterType struct {
	// Filter to only include scene markers with this tag
	TagID *string `json:"tag_id"`
	// Filter to only include scene markers with these tags
	Tags *HierarchicalMultiCriterionInput `json:"tags"`
	// Filter to only include scene markers attached to a scene with these tags
	SceneTags *HierarchicalMultiCriterionInput `json:"scene_tags"`
	// Filter to only include scene markers with these performers
	Performers *MultiCriterionInput `json:"performers"`
}

type MarkerStringsResultType struct {
	Count int    `json:"count"`
	ID    string `json:"id"`
	Title string `json:"title"`
}

type SceneMarkerReader interface {
	Find(id int) (*SceneMarker, error)
	FindMany(ids []int) ([]*SceneMarker, error)
	FindBySceneID(sceneID int) ([]*SceneMarker, error)
	CountByTagID(tagID int) (int, error)
	GetMarkerStrings(q *string, sort *string) ([]*MarkerStringsResultType, error)
	Wall(q *string) ([]*SceneMarker, error)
	Query(sceneMarkerFilter *SceneMarkerFilterType, findFilter *FindFilterType) ([]*SceneMarker, int, error)
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
