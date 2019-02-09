package models

import "github.com/jmoiron/sqlx"

type joinsQueryBuilder struct {}

func NewJoinsQueryBuilder() joinsQueryBuilder {
	return joinsQueryBuilder{}
}

func (qb *joinsQueryBuilder) CreatePerformersScenes(newJoins []PerformersScenes, tx *sqlx.Tx) error {
	ensureTx(tx)
	for _, join := range newJoins {
		_, err := tx.NamedExec(
			`INSERT INTO performers_scenes (performer_id, scene_id) VALUES (:performer_id, :scene_id)`,
			join,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (qb *joinsQueryBuilder) UpdatePerformersScenes(sceneID int, updatedJoins []PerformersScenes, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins and then create new ones
	_, err := tx.Exec("DELETE FROM performers_scenes WHERE scene_id = ?", sceneID)
	if err != nil {
		return err
	}
	return qb.CreatePerformersScenes(updatedJoins, tx)
}

func (qb *joinsQueryBuilder) CreateScenesTags(newJoins []ScenesTags, tx *sqlx.Tx) error {
	ensureTx(tx)
	for _, join := range newJoins {
		_, err := tx.NamedExec(
			`INSERT INTO scenes_tags (scene_id, tag_id) VALUES (:scene_id, :tag_id)`,
			join,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (qb *joinsQueryBuilder) UpdateScenesTags(sceneID int, updatedJoins []ScenesTags, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins and then create new ones
	_, err := tx.Exec("DELETE FROM scenes_tags WHERE scene_id = ?", sceneID)
	if err != nil {
		return err
	}
	return qb.CreateScenesTags(updatedJoins, tx)
}

func (qb *joinsQueryBuilder) CreateSceneMarkersTags(newJoins []SceneMarkersTags, tx *sqlx.Tx) error {
	ensureTx(tx)
	for _, join := range newJoins {
		_, err := tx.NamedExec(
			`INSERT INTO scene_markers_tags (scene_marker_id, tag_id) VALUES (:scene_marker_id, :tag_id)`,
			join,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (qb *joinsQueryBuilder) UpdateSceneMarkersTags(sceneMarkerID int, updatedJoins []SceneMarkersTags, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins and then create new ones
	_, err := tx.Exec("DELETE FROM scene_markers_tags WHERE scene_marker_id = ?", sceneMarkerID)
	if err != nil {
		return err
	}
	return qb.CreateSceneMarkersTags(updatedJoins, tx)
}
