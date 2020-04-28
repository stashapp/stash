package models

import "database/sql"

type PerformersScenes struct {
	PerformerID int `db:"performer_id" json:"performer_id"`
	SceneID     int `db:"scene_id" json:"scene_id"`
}

type MoviesScenes struct {
	MovieID    int           `db:"movie_id" json:"movie_id"`
	SceneID    int           `db:"scene_id" json:"scene_id"`
	SceneIndex sql.NullInt64 `db:"scene_index" json:"scene_index"`
}

type ScenesTags struct {
	SceneID int `db:"scene_id" json:"scene_id"`
	TagID   int `db:"tag_id" json:"tag_id"`
}

type SceneMarkersTags struct {
	SceneMarkerID int `db:"scene_marker_id" json:"scene_marker_id"`
	TagID         int `db:"tag_id" json:"tag_id"`
}

type StashID struct {
	StashID  string `db:"stash_id" json:"stash_id"`
	Endpoint string `db:"endpoint" json:"endpoint"`
}
