package models

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/database"
)

type JoinsQueryBuilder struct{}

func NewJoinsQueryBuilder() JoinsQueryBuilder {
	return JoinsQueryBuilder{}
}

func (qb *JoinsQueryBuilder) GetScenePerformers(sceneID int, tx *sqlx.Tx) ([]PerformersScenes, error) {
	ensureTx(tx)

	// Delete the existing joins and then create new ones
	query := `SELECT * from performers_scenes WHERE scene_id = ?`

	var rows *sqlx.Rows
	var err error
	if tx != nil {
		rows, err = tx.Queryx(query, sceneID)
	} else {
		rows, err = database.DB.Queryx(query, sceneID)
	}

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	performerScenes := make([]PerformersScenes, 0)
	for rows.Next() {
		performerScene := PerformersScenes{}
		if err := rows.StructScan(&performerScene); err != nil {
			return nil, err
		}
		performerScenes = append(performerScenes, performerScene)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return performerScenes, nil
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

// AddPerformerScene adds a performer to a scene. It does not make any change
// if the performer already exists on the scene. It returns true if scene
// performer was added.
func (qb *JoinsQueryBuilder) AddPerformerScene(sceneID int, performerID int, tx *sqlx.Tx) (bool, error) {
	ensureTx(tx)

	existingPerformers, err := qb.GetScenePerformers(sceneID, tx)

	if err != nil {
		return false, err
	}

	// ensure not already present
	for _, p := range existingPerformers {
		if p.PerformerID == performerID && p.SceneID == sceneID {
			return false, nil
		}
	}

	performerJoin := PerformersScenes{
		PerformerID: performerID,
		SceneID:     sceneID,
	}
	performerJoins := append(existingPerformers, performerJoin)

	err = qb.UpdatePerformersScenes(sceneID, performerJoins, tx)

	return err == nil, err
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

func (qb *JoinsQueryBuilder) GetSceneTags(sceneID int, tx *sqlx.Tx) ([]ScenesTags, error) {
	ensureTx(tx)

	// Delete the existing joins and then create new ones
	query := `SELECT * from scenes_tags WHERE scene_id = ?`

	var rows *sqlx.Rows
	var err error
	if tx != nil {
		rows, err = tx.Queryx(query, sceneID)
	} else {
		rows, err = database.DB.Queryx(query, sceneID)
	}

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	sceneTags := make([]ScenesTags, 0)
	for rows.Next() {
		sceneTag := ScenesTags{}
		if err := rows.StructScan(&sceneTag); err != nil {
			return nil, err
		}
		sceneTags = append(sceneTags, sceneTag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sceneTags, nil
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

// AddSceneTag adds a tag to a scene. It does not make any change if the tag
// already exists on the scene. It returns true if scene tag was added.
func (qb *JoinsQueryBuilder) AddSceneTag(sceneID int, tagID int, tx *sqlx.Tx) (bool, error) {
	ensureTx(tx)

	existingTags, err := qb.GetSceneTags(sceneID, tx)

	if err != nil {
		return false, err
	}

	// ensure not already present
	for _, p := range existingTags {
		if p.TagID == tagID && p.SceneID == sceneID {
			return false, nil
		}
	}

	tagJoin := ScenesTags{
		TagID:   tagID,
		SceneID: sceneID,
	}
	tagJoins := append(existingTags, tagJoin)

	err = qb.UpdateScenesTags(sceneID, tagJoins, tx)

	return err == nil, err
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
