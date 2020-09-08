package models

import (
	"github.com/jmoiron/sqlx"
)

type JoinReader interface {
	// GetScenePerformers(sceneID int) ([]PerformersScenes, error)
	GetSceneMovies(sceneID int) ([]MoviesScenes, error)
	// GetSceneTags(sceneID int) ([]ScenesTags, error)
}

type JoinWriter interface {
	CreatePerformersScenes(newJoins []PerformersScenes) error
	// AddPerformerScene(sceneID int, performerID int) (bool, error)
	UpdatePerformersScenes(sceneID int, updatedJoins []PerformersScenes) error
	// DestroyPerformersScenes(sceneID int) error
	CreateMoviesScenes(newJoins []MoviesScenes) error
	// AddMoviesScene(sceneID int, movieID int, sceneIdx *int) (bool, error)
	UpdateMoviesScenes(sceneID int, updatedJoins []MoviesScenes) error
	// DestroyMoviesScenes(sceneID int) error
	// CreateScenesTags(newJoins []ScenesTags) error
	UpdateScenesTags(sceneID int, updatedJoins []ScenesTags) error
	// AddSceneTag(sceneID int, tagID int) (bool, error)
	// DestroyScenesTags(sceneID int) error
	// CreateSceneMarkersTags(newJoins []SceneMarkersTags) error
	UpdateSceneMarkersTags(sceneMarkerID int, updatedJoins []SceneMarkersTags) error
	// DestroySceneMarkersTags(sceneMarkerID int, updatedJoins []SceneMarkersTags) error
	// DestroyScenesGalleries(sceneID int) error
	// DestroyScenesMarkers(sceneID int) error
}

type JoinReaderWriter interface {
	JoinReader
	JoinWriter
}

func NewJoinReaderWriter(tx *sqlx.Tx) JoinReaderWriter {
	return &joinReaderWriter{
		tx: tx,
		qb: NewJoinsQueryBuilder(),
	}
}

type joinReaderWriter struct {
	tx *sqlx.Tx
	qb JoinsQueryBuilder
}

func (t *joinReaderWriter) GetSceneMovies(sceneID int) ([]MoviesScenes, error) {
	return t.qb.GetSceneMovies(sceneID, t.tx)
}

func (t *joinReaderWriter) CreatePerformersScenes(newJoins []PerformersScenes) error {
	return t.qb.CreatePerformersScenes(newJoins, t.tx)
}

func (t *joinReaderWriter) UpdatePerformersScenes(sceneID int, updatedJoins []PerformersScenes) error {
	return t.qb.UpdatePerformersScenes(sceneID, updatedJoins, t.tx)
}

func (t *joinReaderWriter) CreateMoviesScenes(newJoins []MoviesScenes) error {
	return t.qb.CreateMoviesScenes(newJoins, t.tx)
}

func (t *joinReaderWriter) UpdateMoviesScenes(sceneID int, updatedJoins []MoviesScenes) error {
	return t.qb.UpdateMoviesScenes(sceneID, updatedJoins, t.tx)
}

func (t *joinReaderWriter) UpdateScenesTags(sceneID int, updatedJoins []ScenesTags) error {
	return t.qb.UpdateScenesTags(sceneID, updatedJoins, t.tx)
}

func (t *joinReaderWriter) UpdateSceneMarkersTags(sceneMarkerID int, updatedJoins []SceneMarkersTags) error {
	return t.qb.UpdateSceneMarkersTags(sceneMarkerID, updatedJoins, t.tx)
}
