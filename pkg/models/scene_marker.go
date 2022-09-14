package models

import "context"

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
	Find(ctx context.Context, id int) (*SceneMarker, error)
	FindMany(ctx context.Context, ids []int) ([]*SceneMarker, error)
	FindBySceneID(ctx context.Context, sceneID int) ([]*SceneMarker, error)
	CountByTagID(ctx context.Context, tagID int) (int, error)
	GetMarkerStrings(ctx context.Context, q *string, sort *string) ([]*MarkerStringsResultType, error)
	Wall(ctx context.Context, q *string) ([]*SceneMarker, error)
	Query(ctx context.Context, sceneMarkerFilter *SceneMarkerFilterType, findFilter *FindFilterType) ([]*SceneMarker, int, error)
	GetTagIDs(ctx context.Context, imageID int) ([]int, error)
}

type SceneMarkerWriter interface {
	Create(ctx context.Context, newSceneMarker SceneMarker) (*SceneMarker, error)
	Update(ctx context.Context, updatedSceneMarker SceneMarker) (*SceneMarker, error)
	Destroy(ctx context.Context, id int) error
	UpdateTags(ctx context.Context, markerID int, tagIDs []int) error
}

type SceneMarkerReaderWriter interface {
	SceneMarkerReader
	SceneMarkerWriter
}
