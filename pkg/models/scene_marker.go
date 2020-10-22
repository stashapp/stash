package models

import (
	"github.com/jmoiron/sqlx"
)

type SceneMarkerReader interface {
	// Find(id int) (*SceneMarker, error)
	// FindMany(ids []int) ([]*SceneMarker, error)
	FindBySceneID(sceneID int) ([]*SceneMarker, error)
	// CountByTagID(tagID int) (int, error)
	// GetMarkerStrings(q *string, sort *string) ([]*MarkerStringsResultType, error)
	// Wall(q *string) ([]*SceneMarker, error)
	// Query(sceneMarkerFilter *SceneMarkerFilterType, findFilter *FindFilterType) ([]*SceneMarker, int)
}

type SceneMarkerWriter interface {
	Create(newSceneMarker SceneMarker) (*SceneMarker, error)
	Update(updatedSceneMarker SceneMarker) (*SceneMarker, error)
	// Destroy(id string) error
}

type SceneMarkerReaderWriter interface {
	SceneMarkerReader
	SceneMarkerWriter
}

func NewSceneMarkerReaderWriter(tx *sqlx.Tx) SceneMarkerReaderWriter {
	return &sceneMarkerReaderWriter{
		tx: tx,
		qb: NewSceneMarkerQueryBuilder(),
	}
}

type sceneMarkerReaderWriter struct {
	tx *sqlx.Tx
	qb SceneMarkerQueryBuilder
}

func (t *sceneMarkerReaderWriter) FindBySceneID(sceneID int) ([]*SceneMarker, error) {
	return t.qb.FindBySceneID(sceneID, t.tx)
}

func (t *sceneMarkerReaderWriter) Create(newSceneMarker SceneMarker) (*SceneMarker, error) {
	return t.qb.Create(newSceneMarker, t.tx)
}

func (t *sceneMarkerReaderWriter) Update(updatedSceneMarker SceneMarker) (*SceneMarker, error) {
	return t.qb.Update(updatedSceneMarker, t.tx)
}
