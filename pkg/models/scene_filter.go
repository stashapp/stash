package models

import "context"

type SceneFilterFilterType struct {
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

type SceneFilterReader interface {
	Find(ctx context.Context, id int) (*SceneFilter, error)
	FindMany(ctx context.Context, ids []int) ([]*SceneFilter, error)
	FindBySceneID(ctx context.Context, sceneID int) ([]*SceneFilter, error)
	Count(ctx context.Context) (int, error)
	All(ctx context.Context) ([]*SceneFilter, error)
	Query(ctx context.Context, sceneFilterFilter *SceneFilterFilterType, findFilter *FindFilterType) ([]*SceneFilter, int, error)
	QueryCount(ctx context.Context, sceneFilterFilter *SceneFilterFilterType, findFilter *FindFilterType) (int, error)
}

type SceneFilterWriter interface {
	Create(ctx context.Context, newSceneFilter *SceneFilter) error
	Update(ctx context.Context, updatedSceneFilter *SceneFilter) error
	Destroy(ctx context.Context, id int) error
}

type SceneFilterReaderWriter interface {
	SceneFilterReader
	SceneFilterWriter
}
