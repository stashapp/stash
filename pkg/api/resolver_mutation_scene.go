package api

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/models"
)

func (r *mutationResolver) SceneUpdate(ctx context.Context, input models.SceneUpdateInput) (*models.Scene, error) {
	// Populate scene from the input
	sceneID, _ := strconv.Atoi(input.ID)
	updatedTime := time.Now()
	updatedScene := models.Scene{
		ID:        sceneID,
		UpdatedAt: models.SQLiteTimestamp{Timestamp: updatedTime},
	}
	if input.Title != nil {
		updatedScene.Title = sql.NullString{String: *input.Title, Valid: true}
	}
	if input.Details != nil {
		updatedScene.Details = sql.NullString{String: *input.Details, Valid: true}
	}
	if input.URL != nil {
		updatedScene.URL = sql.NullString{String: *input.URL, Valid: true}
	}
	if input.Date != nil {
		updatedScene.Date = models.SQLiteDate{String: *input.Date, Valid: true}
	}

	if input.Rating != nil {
		updatedScene.Rating = sql.NullInt64{Int64: int64(*input.Rating), Valid: true}
	} else {
		// rating must be nullable
		updatedScene.Rating = sql.NullInt64{Valid: false}
	}

	if input.StudioID != nil {
		studioID, _ := strconv.ParseInt(*input.StudioID, 10, 64)
		updatedScene.StudioID = sql.NullInt64{Int64: studioID, Valid: true}
	} else {
		// studio must be nullable
		updatedScene.StudioID = sql.NullInt64{Valid: false}
	}

	// Start the transaction and save the scene marker
	tx := database.DB.MustBeginTx(ctx, nil)
	qb := models.NewSceneQueryBuilder()
	jqb := models.NewJoinsQueryBuilder()
	scene, err := qb.Update(updatedScene, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// Clear the existing gallery value
	gqb := models.NewGalleryQueryBuilder()
	err = gqb.ClearGalleryId(sceneID, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if input.GalleryID != nil {
		// Save the gallery
		galleryID, _ := strconv.Atoi(*input.GalleryID)
		updatedGallery := models.Gallery{
			ID:        galleryID,
			SceneID:   sql.NullInt64{Int64: int64(sceneID), Valid: true},
			UpdatedAt: models.SQLiteTimestamp{Timestamp: updatedTime},
		}
		gqb := models.NewGalleryQueryBuilder()
		_, err := gqb.Update(updatedGallery, tx)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}

	// Save the performers
	var performerJoins []models.PerformersScenes
	for _, pid := range input.PerformerIds {
		performerID, _ := strconv.Atoi(pid)
		performerJoin := models.PerformersScenes{
			PerformerID: performerID,
			SceneID:     sceneID,
		}
		performerJoins = append(performerJoins, performerJoin)
	}
	if err := jqb.UpdatePerformersScenes(sceneID, performerJoins, tx); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// Save the tags
	var tagJoins []models.ScenesTags
	for _, tid := range input.TagIds {
		tagID, _ := strconv.Atoi(tid)
		tagJoin := models.ScenesTags{
			SceneID: sceneID,
			TagID:   tagID,
		}
		tagJoins = append(tagJoins, tagJoin)
	}
	if err := jqb.UpdateScenesTags(sceneID, tagJoins, tx); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return scene, nil
}

func (r *mutationResolver) SceneMarkerCreate(ctx context.Context, input models.SceneMarkerCreateInput) (*models.SceneMarker, error) {
	primaryTagID, _ := strconv.Atoi(input.PrimaryTagID)
	sceneID, _ := strconv.Atoi(input.SceneID)
	currentTime := time.Now()
	newSceneMarker := models.SceneMarker{
		Title:        input.Title,
		Seconds:      input.Seconds,
		PrimaryTagID: primaryTagID,
		SceneID:      sql.NullInt64{Int64: int64(sceneID), Valid: sceneID != 0},
		CreatedAt:    models.SQLiteTimestamp{Timestamp: currentTime},
		UpdatedAt:    models.SQLiteTimestamp{Timestamp: currentTime},
	}

	return changeMarker(ctx, create, newSceneMarker, input.TagIds)
}

func (r *mutationResolver) SceneMarkerUpdate(ctx context.Context, input models.SceneMarkerUpdateInput) (*models.SceneMarker, error) {
	// Populate scene marker from the input
	sceneMarkerID, _ := strconv.Atoi(input.ID)
	sceneID, _ := strconv.Atoi(input.SceneID)
	primaryTagID, _ := strconv.Atoi(input.PrimaryTagID)
	updatedSceneMarker := models.SceneMarker{
		ID:           sceneMarkerID,
		Title:        input.Title,
		Seconds:      input.Seconds,
		SceneID:      sql.NullInt64{Int64: int64(sceneID), Valid: sceneID != 0},
		PrimaryTagID: primaryTagID,
		UpdatedAt:    models.SQLiteTimestamp{Timestamp: time.Now()},
	}

	return changeMarker(ctx, update, updatedSceneMarker, input.TagIds)
}

func (r *mutationResolver) SceneMarkerDestroy(ctx context.Context, id string) (bool, error) {
	qb := models.NewSceneMarkerQueryBuilder()
	tx := database.DB.MustBeginTx(ctx, nil)
	if err := qb.Destroy(id, tx); err != nil {
		_ = tx.Rollback()
		return false, err
	}
	if err := tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}

func changeMarker(ctx context.Context, changeType int, changedMarker models.SceneMarker, tagIds []string) (*models.SceneMarker, error) {
	// Start the transaction and save the scene marker
	tx := database.DB.MustBeginTx(ctx, nil)
	qb := models.NewSceneMarkerQueryBuilder()
	jqb := models.NewJoinsQueryBuilder()

	var sceneMarker *models.SceneMarker
	var err error
	switch changeType {
	case create:
		sceneMarker, err = qb.Create(changedMarker, tx)
	case update:
		sceneMarker, err = qb.Update(changedMarker, tx)
	}
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// Save the marker tags
	var markerTagJoins []models.SceneMarkersTags
	for _, tid := range tagIds {
		tagID, _ := strconv.Atoi(tid)
		if tagID == changedMarker.PrimaryTagID {
			continue // If this tag is the primary tag, then let's not add it.
		}
		markerTag := models.SceneMarkersTags{
			SceneMarkerID: sceneMarker.ID,
			TagID:         tagID,
		}
		markerTagJoins = append(markerTagJoins, markerTag)
	}
	switch changeType {
	case create:
		if err := jqb.CreateSceneMarkersTags(markerTagJoins, tx); err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	case update:
		if err := jqb.UpdateSceneMarkersTags(changedMarker.ID, markerTagJoins, tx); err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return sceneMarker, nil
}
