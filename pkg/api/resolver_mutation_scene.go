package api

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *mutationResolver) SceneUpdate(ctx context.Context, input models.SceneUpdateInput) (*models.Scene, error) {
	// Start the transaction and save the scene
	tx := database.DB.MustBeginTx(ctx, nil)

	ret, err := r.sceneUpdate(input, tx)

	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *mutationResolver) ScenesUpdate(ctx context.Context, input []*models.SceneUpdateInput) ([]*models.Scene, error) {
	// Start the transaction and save the scene
	tx := database.DB.MustBeginTx(ctx, nil)

	var ret []*models.Scene

	for _, scene := range input {
		thisScene, err := r.sceneUpdate(*scene, tx)
		ret = append(ret, thisScene)

		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *mutationResolver) sceneUpdate(input models.SceneUpdateInput, tx *sqlx.Tx) (*models.Scene, error) {
	// Populate scene from the input
	sceneID, _ := strconv.Atoi(input.ID)

	var coverImageData []byte

	updatedTime := time.Now()
	updatedScene := models.ScenePartial{
		ID:        sceneID,
		UpdatedAt: &models.SQLiteTimestamp{Timestamp: updatedTime},
	}
	if input.Title != nil {
		updatedScene.Title = &sql.NullString{String: *input.Title, Valid: true}
	}
	if input.Details != nil {
		updatedScene.Details = &sql.NullString{String: *input.Details, Valid: true}
	}
	if input.URL != nil {
		updatedScene.URL = &sql.NullString{String: *input.URL, Valid: true}
	}
	if input.Date != nil {
		updatedScene.Date = &models.SQLiteDate{String: *input.Date, Valid: true}
	}

	if input.CoverImage != nil && *input.CoverImage != "" {
		var err error
		_, coverImageData, err = utils.ProcessBase64Image(*input.CoverImage)
		if err != nil {
			return nil, err
		}

		// update the cover after updating the scene
	}

	if input.Rating != nil {
		updatedScene.Rating = &sql.NullInt64{Int64: int64(*input.Rating), Valid: true}
	} else {
		// rating must be nullable
		updatedScene.Rating = &sql.NullInt64{Valid: false}
	}

	if input.StudioID != nil {
		studioID, _ := strconv.ParseInt(*input.StudioID, 10, 64)
		updatedScene.StudioID = &sql.NullInt64{Int64: studioID, Valid: true}
	} else {
		// studio must be nullable
		updatedScene.StudioID = &sql.NullInt64{Valid: false}
	}

	qb := models.NewSceneQueryBuilder()
	jqb := models.NewJoinsQueryBuilder()
	scene, err := qb.Update(updatedScene, tx)
	if err != nil {
		return nil, err
	}

	// update cover table
	if len(coverImageData) > 0 {
		if err := qb.UpdateSceneCover(sceneID, coverImageData, tx); err != nil {
			return nil, err
		}
	}

	// Clear the existing gallery value
	gqb := models.NewGalleryQueryBuilder()
	err = gqb.ClearGalleryId(sceneID, tx)
	if err != nil {
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
		return nil, err
	}

	// Save the movies
	var movieJoins []models.MoviesScenes

	for _, movie := range input.Movies {

		movieID, _ := strconv.Atoi(movie.MovieID)

		movieJoin := models.MoviesScenes{
			MovieID: movieID,
			SceneID: sceneID,
		}

		if movie.SceneIndex != nil {
			movieJoin.SceneIndex = sql.NullInt64{
				Int64: int64(*movie.SceneIndex),
				Valid: true,
			}
		}

		movieJoins = append(movieJoins, movieJoin)
	}
	if err := jqb.UpdateMoviesScenes(sceneID, movieJoins, tx); err != nil {
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
		return nil, err
	}

	// only update the cover image if provided and everything else was successful
	if coverImageData != nil {
		err = manager.SetSceneScreenshot(scene.Checksum, coverImageData)
		if err != nil {
			return nil, err
		}
	}

	return scene, nil
}

func (r *mutationResolver) BulkSceneUpdate(ctx context.Context, input models.BulkSceneUpdateInput) ([]*models.Scene, error) {
	// Populate scene from the input
	updatedTime := time.Now()

	// Start the transaction and save the scene marker
	tx := database.DB.MustBeginTx(ctx, nil)
	qb := models.NewSceneQueryBuilder()
	jqb := models.NewJoinsQueryBuilder()

	updatedScene := models.ScenePartial{
		UpdatedAt: &models.SQLiteTimestamp{Timestamp: updatedTime},
	}
	if input.Title != nil {
		updatedScene.Title = &sql.NullString{String: *input.Title, Valid: true}
	}
	if input.Details != nil {
		updatedScene.Details = &sql.NullString{String: *input.Details, Valid: true}
	}
	if input.URL != nil {
		updatedScene.URL = &sql.NullString{String: *input.URL, Valid: true}
	}
	if input.Date != nil {
		updatedScene.Date = &models.SQLiteDate{String: *input.Date, Valid: true}
	}
	if input.Rating != nil {
		// a rating of 0 means unset the rating
		if *input.Rating == 0 {
			updatedScene.Rating = &sql.NullInt64{Int64: 0, Valid: false}
		} else {
			updatedScene.Rating = &sql.NullInt64{Int64: int64(*input.Rating), Valid: true}
		}
	}
	if input.StudioID != nil {
		// empty string means unset the studio
		if *input.StudioID == "" {
			updatedScene.StudioID = &sql.NullInt64{Int64: 0, Valid: false}
		} else {
			studioID, _ := strconv.ParseInt(*input.StudioID, 10, 64)
			updatedScene.StudioID = &sql.NullInt64{Int64: studioID, Valid: true}
		}
	}

	ret := []*models.Scene{}

	for _, sceneIDStr := range input.Ids {
		sceneID, _ := strconv.Atoi(sceneIDStr)
		updatedScene.ID = sceneID

		scene, err := qb.Update(updatedScene, tx)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}

		ret = append(ret, scene)

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
		if wasFieldIncluded(ctx, "performer_ids") {
			performerIDs, err := adjustScenePerformerIDs(tx, sceneID, *input.PerformerIds)
			if err != nil {
				_ = tx.Rollback()
				return nil, err
			}

			var performerJoins []models.PerformersScenes
			for _, performerID := range performerIDs {
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
		}

		// Save the tags
		if wasFieldIncluded(ctx, "tag_ids") {
			tagIDs, err := adjustSceneTagIDs(tx, sceneID, *input.TagIds)
			if err != nil {
				_ = tx.Rollback()
				return nil, err
			}

			var tagJoins []models.ScenesTags
			for _, tagID := range tagIDs {
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
		}
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return ret, nil
}

func adjustIDs(existingIDs []int, updateIDs models.BulkUpdateIds) []int {
	for _, idStr := range updateIDs.Ids {
		id, _ := strconv.Atoi(idStr)

		// look for the id in the list
		foundExisting := false
		for idx, existingID := range existingIDs {
			if existingID == id {
				if updateIDs.Mode == models.BulkUpdateIDModeRemove {
					// remove from the list
					existingIDs = append(existingIDs[:idx], existingIDs[idx+1:]...)
				}

				foundExisting = true
				break
			}
		}

		if !foundExisting && updateIDs.Mode != models.BulkUpdateIDModeRemove {
			existingIDs = append(existingIDs, id)
		}
	}

	return existingIDs
}

func adjustScenePerformerIDs(tx *sqlx.Tx, sceneID int, ids models.BulkUpdateIds) ([]int, error) {
	var ret []int

	jqb := models.NewJoinsQueryBuilder()
	if ids.Mode == models.BulkUpdateIDModeAdd || ids.Mode == models.BulkUpdateIDModeRemove {
		// adding to the joins
		performerJoins, err := jqb.GetScenePerformers(sceneID, tx)

		if err != nil {
			return nil, err
		}

		for _, join := range performerJoins {
			ret = append(ret, join.PerformerID)
		}
	}

	return adjustIDs(ret, ids), nil
}

func adjustSceneTagIDs(tx *sqlx.Tx, sceneID int, ids models.BulkUpdateIds) ([]int, error) {
	var ret []int

	jqb := models.NewJoinsQueryBuilder()
	if ids.Mode == models.BulkUpdateIDModeAdd || ids.Mode == models.BulkUpdateIDModeRemove {
		// adding to the joins
		tagJoins, err := jqb.GetSceneTags(sceneID, tx)

		if err != nil {
			return nil, err
		}

		for _, join := range tagJoins {
			ret = append(ret, join.TagID)
		}
	}

	return adjustIDs(ret, ids), nil
}

func (r *mutationResolver) SceneDestroy(ctx context.Context, input models.SceneDestroyInput) (bool, error) {
	qb := models.NewSceneQueryBuilder()
	tx := database.DB.MustBeginTx(ctx, nil)

	sceneID, _ := strconv.Atoi(input.ID)
	scene, err := qb.Find(sceneID)
	err = manager.DestroyScene(sceneID, tx)

	if err != nil {
		tx.Rollback()
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	// if delete generated is true, then delete the generated files
	// for the scene
	if input.DeleteGenerated != nil && *input.DeleteGenerated {
		manager.DeleteGeneratedSceneFiles(scene)
	}

	// if delete file is true, then delete the file as well
	// if it fails, just log a message
	if input.DeleteFile != nil && *input.DeleteFile {
		manager.DeleteSceneFile(scene)
	}

	return true, nil
}

func (r *mutationResolver) ScenesDestroy(ctx context.Context, input models.ScenesDestroyInput) (bool, error) {
	qb := models.NewSceneQueryBuilder()
	tx := database.DB.MustBeginTx(ctx, nil)

	var scenes []*models.Scene
	for _, id := range input.Ids {
		sceneID, _ := strconv.Atoi(id)

		scene, err := qb.Find(sceneID)
		if scene != nil {
			scenes = append(scenes, scene)
		}
		err = manager.DestroyScene(sceneID, tx)

		if err != nil {
			tx.Rollback()
			return false, err
		}
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	for _, scene := range scenes {
		// if delete generated is true, then delete the generated files
		// for the scene
		if input.DeleteGenerated != nil && *input.DeleteGenerated {
			manager.DeleteGeneratedSceneFiles(scene)
		}

		// if delete file is true, then delete the file as well
		// if it fails, just log a message
		if input.DeleteFile != nil && *input.DeleteFile {
			manager.DeleteSceneFile(scene)
		}
	}

	return true, nil
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

	markerID, _ := strconv.Atoi(id)
	marker, err := qb.Find(markerID)

	if err != nil {
		return false, err
	}

	if err := qb.Destroy(id, tx); err != nil {
		_ = tx.Rollback()
		return false, err
	}
	if err := tx.Commit(); err != nil {
		return false, err
	}

	// delete the preview for the marker
	sqb := models.NewSceneQueryBuilder()
	scene, _ := sqb.Find(int(marker.SceneID.Int64))

	if scene != nil {
		seconds := int(marker.Seconds)
		manager.DeleteSceneMarkerFiles(scene, seconds)
	}

	return true, nil
}

func changeMarker(ctx context.Context, changeType int, changedMarker models.SceneMarker, tagIds []string) (*models.SceneMarker, error) {
	// Start the transaction and save the scene marker
	tx := database.DB.MustBeginTx(ctx, nil)
	qb := models.NewSceneMarkerQueryBuilder()
	jqb := models.NewJoinsQueryBuilder()

	var existingMarker *models.SceneMarker
	var sceneMarker *models.SceneMarker
	var err error
	switch changeType {
	case create:
		sceneMarker, err = qb.Create(changedMarker, tx)
	case update:
		// check to see if timestamp was changed
		existingMarker, err = qb.Find(changedMarker.ID)
		if err == nil {
			sceneMarker, err = qb.Update(changedMarker, tx)
		}
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

	// remove the marker preview if the timestamp was changed
	if existingMarker != nil && existingMarker.Seconds != changedMarker.Seconds {
		sqb := models.NewSceneQueryBuilder()

		scene, _ := sqb.Find(int(existingMarker.SceneID.Int64))

		if scene != nil {
			seconds := int(existingMarker.Seconds)
			manager.DeleteSceneMarkerFiles(scene, seconds)
		}
	}

	return sceneMarker, nil
}

func (r *mutationResolver) SceneIncrementO(ctx context.Context, id string) (int, error) {
	sceneID, _ := strconv.Atoi(id)

	tx := database.DB.MustBeginTx(ctx, nil)
	qb := models.NewSceneQueryBuilder()

	newVal, err := qb.IncrementOCounter(sceneID, tx)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return newVal, nil
}

func (r *mutationResolver) SceneDecrementO(ctx context.Context, id string) (int, error) {
	sceneID, _ := strconv.Atoi(id)

	tx := database.DB.MustBeginTx(ctx, nil)
	qb := models.NewSceneQueryBuilder()

	newVal, err := qb.DecrementOCounter(sceneID, tx)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return newVal, nil
}

func (r *mutationResolver) SceneResetO(ctx context.Context, id string) (int, error) {
	sceneID, _ := strconv.Atoi(id)

	tx := database.DB.MustBeginTx(ctx, nil)
	qb := models.NewSceneQueryBuilder()

	newVal, err := qb.ResetOCounter(sceneID, tx)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return newVal, nil
}

func (r *mutationResolver) SceneGenerateScreenshot(ctx context.Context, id string, at *float64) (string, error) {
	if at != nil {
		manager.GetInstance().GenerateScreenshot(id, *at)
	} else {
		manager.GetInstance().GenerateDefaultScreenshot(id)
	}

	return "todo", nil
}
