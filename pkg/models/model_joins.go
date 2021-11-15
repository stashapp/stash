package models

import "database/sql"

type MoviesScenes struct {
	MovieID    int           `db:"movie_id" json:"movie_id"`
	SceneID    int           `db:"scene_id" json:"scene_id"`
	SceneIndex sql.NullInt64 `db:"scene_index" json:"scene_index"`
}

type StashID struct {
	StashID  string `db:"stash_id" json:"stash_id"`
	Endpoint string `db:"endpoint" json:"endpoint"`
}

func (s StashID) StashIDInput() StashIDInput {
	return StashIDInput{
		Endpoint: s.Endpoint,
		StashID:  s.StashID,
	}
}
