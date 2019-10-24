package models

import "github.com/jmoiron/sqlx"

type JoinsQueryBuilder struct{}

func NewJoinsQueryBuilder() JoinsQueryBuilder {
	return JoinsQueryBuilder{}
}

func (qb *JoinsQueryBuilder) CreatePerformersScenes(newJoins []PerformersScenes, tx *sqlx.Tx) error {
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

func (qb *JoinsQueryBuilder) UpdatePerformersScenes(sceneID int, updatedJoins []PerformersScenes, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins and then create new ones
	_, err := tx.Exec("DELETE FROM performers_scenes WHERE scene_id = ?", sceneID)
	if err != nil {
		return err
	}
	return qb.CreatePerformersScenes(updatedJoins, tx)
}

func (qb *JoinsQueryBuilder) DestroyPerformersScenes(sceneID int, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins
	_, err := tx.Exec("DELETE FROM performers_scenes WHERE scene_id = ?", sceneID)
	return err
}

func (qb *JoinsQueryBuilder) CreateScenesTags(newJoins []ScenesTags, tx *sqlx.Tx) error {
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

func (qb *JoinsQueryBuilder) UpdateScenesTags(sceneID int, updatedJoins []ScenesTags, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins and then create new ones
	_, err := tx.Exec("DELETE FROM scenes_tags WHERE scene_id = ?", sceneID)
	if err != nil {
		return err
	}
	return qb.CreateScenesTags(updatedJoins, tx)
}

func (qb *JoinsQueryBuilder) DestroyScenesTags(sceneID int, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins
	_, err := tx.Exec("DELETE FROM scenes_tags WHERE scene_id = ?", sceneID)

	return err
}

func (qb *JoinsQueryBuilder) CreateSceneMarkersTags(newJoins []SceneMarkersTags, tx *sqlx.Tx) error {
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

func (qb *JoinsQueryBuilder) UpdateSceneMarkersTags(sceneMarkerID int, updatedJoins []SceneMarkersTags, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins and then create new ones
	_, err := tx.Exec("DELETE FROM scene_markers_tags WHERE scene_marker_id = ?", sceneMarkerID)
	if err != nil {
		return err
	}
	return qb.CreateSceneMarkersTags(updatedJoins, tx)
}

func (qb *JoinsQueryBuilder) DestroySceneMarkersTags(sceneMarkerID int, updatedJoins []SceneMarkersTags, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins
	_, err := tx.Exec("DELETE FROM scene_markers_tags WHERE scene_marker_id = ?", sceneMarkerID)
	return err
}

func (qb *JoinsQueryBuilder) DestroyScenesGalleries(sceneID int, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Unset the existing scene id from galleries
	_, err := tx.Exec("UPDATE galleries SET scene_id = null WHERE scene_id = ?", sceneID)

	return err
}

func (qb *JoinsQueryBuilder) DestroyScenesMarkers(sceneID int, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the scene marker tags
	_, err := tx.Exec("DELETE t FROM scene_markers_tags t join scene_markers m on t.scene_marker_id = m.id WHERE m.scene_id = ?", sceneID)

	// Delete the existing joins
	_, err = tx.Exec("DELETE FROM scene_markers WHERE scene_id = ?", sceneID)

	return err
}
