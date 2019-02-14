package models

type PerformersScenes struct {
	PerformerID int `db:"performer_id" json:"performer_id"`
	SceneID     int `db:"scene_id" json:"scene_id"`
}

type ScenesTags struct {
	SceneID int `db:"scene_id" json:"scene_id"`
	TagID   int `db:"tag_id" json:"tag_id"`
}

type SceneMarkersTags struct {
	SceneMarkerID int `db:"scene_marker_id" json:"scene_marker_id"`
	TagID         int `db:"tag_id" json:"tag_id"`
}
