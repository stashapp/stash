package models

import (
	"github.com/jmoiron/sqlx"
)

type SceneReader interface {
	// Find(id int) (*Scene, error)
	FindMany(ids []int) ([]*Scene, error)
	// FindByChecksum(checksum string) (*Scene, error)
	// FindByOSHash(oshash string) (*Scene, error)
	// FindByPath(path string) (*Scene, error)
	// FindByPerformerID(performerID int) ([]*Scene, error)
	// CountByPerformerID(performerID int) (int, error)
	// FindByStudioID(studioID int) ([]*Scene, error)
	// FindByMovieID(movieID int) ([]*Scene, error)
	// CountByMovieID(movieID int) (int, error)
	// Count() (int, error)
	// SizeCount() (string, error)
	// CountByStudioID(studioID int) (int, error)
	// CountByTagID(tagID int) (int, error)
	// CountMissingChecksum() (int, error)
	// CountMissingOSHash() (int, error)
	// Wall(q *string) ([]*Scene, error)
	All() ([]*Scene, error)
	// Query(sceneFilter *SceneFilterType, findFilter *FindFilterType) ([]*Scene, int)
	// QueryAllByPathRegex(regex string) ([]*Scene, error)
	// QueryByPathRegex(findFilter *FindFilterType) ([]*Scene, int)
	GetSceneCover(sceneID int) ([]byte, error)
}

type SceneWriter interface {
	// Create(newScene Scene) (*Scene, error)
	// Update(updatedScene ScenePartial) (*Scene, error)
	// IncrementOCounter(id int) (int, error)
	// DecrementOCounter(id int) (int, error)
	// ResetOCounter(id int) (int, error)
	// Destroy(id string) error
	// UpdateFormat(id int, format string) error
	// UpdateOSHash(id int, oshash string) error
	// UpdateChecksum(id int, checksum string) error
	// UpdateSceneCover(sceneID int, cover []byte) error
	// DestroySceneCover(sceneID int) error
}

type SceneReaderWriter interface {
	SceneReader
	SceneWriter
}

func NewSceneReaderWriter(tx *sqlx.Tx) SceneReaderWriter {
	return &sceneReaderWriter{
		tx: tx,
		qb: NewSceneQueryBuilder(),
	}
}

type sceneReaderWriter struct {
	tx *sqlx.Tx
	qb SceneQueryBuilder
}

func (t *sceneReaderWriter) FindMany(ids []int) ([]*Scene, error) {
	return t.qb.FindMany(ids)
}

func (t *sceneReaderWriter) All() ([]*Scene, error) {
	return t.qb.All()
}

func (t *sceneReaderWriter) GetSceneCover(sceneID int) ([]byte, error) {
	return t.qb.GetSceneCover(sceneID, t.tx)
}
