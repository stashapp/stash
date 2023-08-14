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
	// Filter by created at
	CreatedAt *TimestampCriterionInput `json:"created_at"`
	// Filter by updated at
	UpdatedAt *TimestampCriterionInput `json:"updated_at"`
	// Filter by scenes date
	SceneDate *DateCriterionInput `json:"scene_date"`
	// Filter by scenes created at
	SceneCreatedAt *TimestampCriterionInput `json:"scene_created_at"`
	// Filter by scenes updated at
	SceneUpdatedAt *TimestampCriterionInput `json:"scene_updated_at"`
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
	Count(ctx context.Context) (int, error)
	All(ctx context.Context) ([]*SceneMarker, error)
	Query(ctx context.Context, sceneMarkerFilter *SceneMarkerFilterType, findFilter *FindFilterType) ([]*SceneMarker, int, error)
	QueryCount(ctx context.Context, sceneMarkerFilter *SceneMarkerFilterType, findFilter *FindFilterType) (int, error)
	GetTagIDs(ctx context.Context, imageID int) ([]int, error)
}

type SceneMarkerWriter interface {
	Create(ctx context.Context, newSceneMarker *SceneMarker) error
	Update(ctx context.Context, updatedSceneMarker *SceneMarker) error
	UpdatePartial(ctx context.Context, id int, updatedSceneMarker SceneMarkerPartial) (*SceneMarker, error)
	Destroy(ctx context.Context, id int) error
	UpdateTags(ctx context.Context, markerID int, tagIDs []int) error
}

type SceneMarkerReaderWriter interface {
	SceneMarkerReader
	SceneMarkerWriter
}
