package api

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *mutationResolver) getScene(ctx context.Context, id int) (ret *models.Scene, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Scene.Find(ctx, id)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *mutationResolver) SceneUpdate(ctx context.Context, input models.SceneUpdateInput) (ret *models.Scene, err error) {
	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Start the transaction and save the scene
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.sceneUpdate(ctx, input, translator)
		return err
	}); err != nil {
		return nil, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, ret.ID, plugin.SceneUpdatePost, input, translator.getFields())
	return r.getScene(ctx, ret.ID)
}

func (r *mutationResolver) ScenesUpdate(ctx context.Context, input []*models.SceneUpdateInput) (ret []*models.Scene, err error) {
	inputMaps := getUpdateInputMaps(ctx)

	// Start the transaction and save the scene
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		for i, scene := range input {
			translator := changesetTranslator{
				inputMap: inputMaps[i],
			}

			thisScene, err := r.sceneUpdate(ctx, *scene, translator)
			ret = append(ret, thisScene)

			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// execute post hooks outside of txn
	var newRet []*models.Scene
	for i, scene := range ret {
		translator := changesetTranslator{
			inputMap: inputMaps[i],
		}

		r.hookExecutor.ExecutePostHooks(ctx, scene.ID, plugin.SceneUpdatePost, input, translator.getFields())

		scene, err = r.getScene(ctx, scene.ID)
		if err != nil {
			return nil, err
		}

		newRet = append(newRet, scene)
	}

	return newRet, nil
}

func (r *mutationResolver) sceneUpdate(ctx context.Context, input models.SceneUpdateInput, translator changesetTranslator) (*models.Scene, error) {
	// Populate scene from the input
	sceneID, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, err
	}

	var coverImageData []byte

	updatedScene := models.NewScenePartial()
	updatedScene.Title = translator.optionalString(input.Title, "title")
	updatedScene.Details = translator.optionalString(input.Details, "details")
	updatedScene.URL = translator.optionalString(input.URL, "url")
	updatedScene.Date = translator.optionalDate(input.Date, "date")
	updatedScene.Rating = translator.optionalInt(input.Rating, "rating")
	updatedScene.StudioID, err = translator.optionalIntFromString(input.StudioID, "studio_id")
	if err != nil {
		return nil, fmt.Errorf("converting studio id: %w", err)
	}

	updatedScene.Organized = translator.optionalBool(input.Organized, "organized")

	if translator.hasField("performer_ids") {
		updatedScene.PerformerIDs, err = translateUpdateIDs(input.PerformerIds, models.RelationshipUpdateModeSet)
		if err != nil {
			return nil, fmt.Errorf("converting performer ids: %w", err)
		}
	}

	if translator.hasField("tag_ids") {
		updatedScene.TagIDs, err = translateUpdateIDs(input.TagIds, models.RelationshipUpdateModeSet)
		if err != nil {
			return nil, fmt.Errorf("converting tag ids: %w", err)
		}
	}

	if translator.hasField("gallery_ids") {
		updatedScene.GalleryIDs, err = translateUpdateIDs(input.GalleryIds, models.RelationshipUpdateModeSet)
		if err != nil {
			return nil, fmt.Errorf("converting gallery ids: %w", err)
		}
	}

	// Save the movies
	if translator.hasField("movies") {
		updatedScene.MovieIDs, err = models.UpdateMovieIDsFromInput(input.Movies)
		if err != nil {
			return nil, fmt.Errorf("converting movie ids: %w", err)
		}
	}

	// Save the stash_ids
	if translator.hasField("stash_ids") {
		updatedScene.StashIDs = &models.UpdateStashIDs{
			StashIDs: input.StashIds,
			Mode:     models.RelationshipUpdateModeSet,
		}
	}

	if input.CoverImage != nil && *input.CoverImage != "" {
		var err error
		coverImageData, err = utils.ProcessImageInput(ctx, *input.CoverImage)
		if err != nil {
			return nil, err
		}

		// update the cover after updating the scene
	}

	qb := r.repository.Scene
	s, err := qb.UpdatePartial(ctx, sceneID, updatedScene)
	if err != nil {
		return nil, err
	}

	// update cover table
	if len(coverImageData) > 0 {
		if err := qb.UpdateCover(ctx, sceneID, coverImageData); err != nil {
			return nil, err
		}
	}

	// only update the cover image if provided and everything else was successful
	if coverImageData != nil {
		err = scene.SetScreenshot(manager.GetInstance().Paths, s.GetHash(config.GetInstance().GetVideoFileNamingAlgorithm()), coverImageData)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (r *mutationResolver) BulkSceneUpdate(ctx context.Context, input BulkSceneUpdateInput) ([]*models.Scene, error) {
	sceneIDs, err := stringslice.StringSliceToIntSlice(input.Ids)
	if err != nil {
		return nil, err
	}

	// Populate scene from the input
	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	updatedScene := models.NewScenePartial()
	updatedScene.Title = translator.optionalString(input.Title, "title")
	updatedScene.Details = translator.optionalString(input.Details, "details")
	updatedScene.URL = translator.optionalString(input.URL, "url")
	updatedScene.Date = translator.optionalDate(input.Date, "date")
	updatedScene.Rating = translator.optionalInt(input.Rating, "rating")
	updatedScene.StudioID, err = translator.optionalIntFromString(input.StudioID, "studio_id")
	if err != nil {
		return nil, fmt.Errorf("converting studio id: %w", err)
	}

	updatedScene.Organized = translator.optionalBool(input.Organized, "organized")

	if translator.hasField("performer_ids") {
		updatedScene.PerformerIDs, err = translateUpdateIDs(input.PerformerIds.Ids, input.PerformerIds.Mode)
		if err != nil {
			return nil, fmt.Errorf("converting performer ids: %w", err)
		}
	}

	if translator.hasField("tag_ids") {
		updatedScene.TagIDs, err = translateUpdateIDs(input.TagIds.Ids, input.TagIds.Mode)
		if err != nil {
			return nil, fmt.Errorf("converting tag ids: %w", err)
		}
	}

	if translator.hasField("gallery_ids") {
		updatedScene.GalleryIDs, err = translateUpdateIDs(input.GalleryIds.Ids, input.GalleryIds.Mode)
		if err != nil {
			return nil, fmt.Errorf("converting gallery ids: %w", err)
		}
	}

	// Save the movies
	if translator.hasField("movies") {
		updatedScene.MovieIDs, err = translateSceneMovieIDs(*input.MovieIds)
		if err != nil {
			return nil, fmt.Errorf("converting movie ids: %w", err)
		}
	}

	ret := []*models.Scene{}

	// Start the transaction and save the scene marker
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Scene

		for _, sceneID := range sceneIDs {
			scene, err := qb.UpdatePartial(ctx, sceneID, updatedScene)
			if err != nil {
				return err
			}

			ret = append(ret, scene)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// execute post hooks outside of txn
	var newRet []*models.Scene
	for _, scene := range ret {
		r.hookExecutor.ExecutePostHooks(ctx, scene.ID, plugin.SceneUpdatePost, input, translator.getFields())

		scene, err = r.getScene(ctx, scene.ID)
		if err != nil {
			return nil, err
		}

		newRet = append(newRet, scene)
	}

	return newRet, nil
}

func adjustIDs(existingIDs []int, updateIDs BulkUpdateIds) []int {
	// if we are setting the ids, just return the ids
	if updateIDs.Mode == models.RelationshipUpdateModeSet {
		existingIDs = []int{}
		for _, idStr := range updateIDs.Ids {
			id, _ := strconv.Atoi(idStr)
			existingIDs = append(existingIDs, id)
		}

		return existingIDs
	}

	for _, idStr := range updateIDs.Ids {
		id, _ := strconv.Atoi(idStr)

		// look for the id in the list
		foundExisting := false
		for idx, existingID := range existingIDs {
			if existingID == id {
				if updateIDs.Mode == models.RelationshipUpdateModeRemove {
					// remove from the list
					existingIDs = append(existingIDs[:idx], existingIDs[idx+1:]...)
				}

				foundExisting = true
				break
			}
		}

		if !foundExisting && updateIDs.Mode != models.RelationshipUpdateModeRemove {
			existingIDs = append(existingIDs, id)
		}
	}

	return existingIDs
}

type tagIDsGetter interface {
	GetTagIDs(ctx context.Context, id int) ([]int, error)
}

func adjustTagIDs(ctx context.Context, qb tagIDsGetter, sceneID int, ids BulkUpdateIds) (ret []int, err error) {
	ret, err = qb.GetTagIDs(ctx, sceneID)
	if err != nil {
		return nil, err
	}

	return adjustIDs(ret, ids), nil
}

func (r *mutationResolver) SceneDestroy(ctx context.Context, input models.SceneDestroyInput) (bool, error) {
	sceneID, err := strconv.Atoi(input.ID)
	if err != nil {
		return false, err
	}

	fileNamingAlgo := manager.GetInstance().Config.GetVideoFileNamingAlgorithm()

	var s *models.Scene
	fileDeleter := &scene.FileDeleter{
		Deleter:        *file.NewDeleter(),
		FileNamingAlgo: fileNamingAlgo,
		Paths:          manager.GetInstance().Paths,
	}

	deleteGenerated := utils.IsTrue(input.DeleteGenerated)
	deleteFile := utils.IsTrue(input.DeleteFile)

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Scene
		var err error
		s, err = qb.Find(ctx, sceneID)
		if err != nil {
			return err
		}

		if s == nil {
			return fmt.Errorf("scene with id %d not found", sceneID)
		}

		// kill any running encoders
		manager.KillRunningStreams(s, fileNamingAlgo)

		return scene.Destroy(ctx, s, r.repository.Scene, r.repository.SceneMarker, fileDeleter, deleteGenerated, deleteFile)
	}); err != nil {
		fileDeleter.Rollback()
		return false, err
	}

	// perform the post-commit actions
	fileDeleter.Commit()

	// call post hook after performing the other actions
	r.hookExecutor.ExecutePostHooks(ctx, s.ID, plugin.SceneDestroyPost, plugin.SceneDestroyInput{
		SceneDestroyInput: input,
		Checksum:          stringPtrToString(s.Checksum),
		OSHash:            stringPtrToString(s.OSHash),
		Path:              s.Path,
	}, nil)

	return true, nil
}

func (r *mutationResolver) ScenesDestroy(ctx context.Context, input models.ScenesDestroyInput) (bool, error) {
	var scenes []*models.Scene
	fileNamingAlgo := manager.GetInstance().Config.GetVideoFileNamingAlgorithm()

	fileDeleter := &scene.FileDeleter{
		Deleter:        *file.NewDeleter(),
		FileNamingAlgo: fileNamingAlgo,
		Paths:          manager.GetInstance().Paths,
	}

	deleteGenerated := utils.IsTrue(input.DeleteGenerated)
	deleteFile := utils.IsTrue(input.DeleteFile)

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Scene

		for _, id := range input.Ids {
			sceneID, _ := strconv.Atoi(id)

			s, err := qb.Find(ctx, sceneID)
			if err != nil {
				return err
			}
			if s != nil {
				scenes = append(scenes, s)
			}

			// kill any running encoders
			manager.KillRunningStreams(s, fileNamingAlgo)

			if err := scene.Destroy(ctx, s, r.repository.Scene, r.repository.SceneMarker, fileDeleter, deleteGenerated, deleteFile); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		fileDeleter.Rollback()
		return false, err
	}

	// perform the post-commit actions
	fileDeleter.Commit()

	for _, scene := range scenes {
		// call post hook after performing the other actions
		r.hookExecutor.ExecutePostHooks(ctx, scene.ID, plugin.SceneDestroyPost, plugin.ScenesDestroyInput{
			ScenesDestroyInput: input,
			Checksum:           stringPtrToString(scene.Checksum),
			OSHash:             stringPtrToString(scene.OSHash),
			Path:               scene.Path,
		}, nil)
	}

	return true, nil
}

func (r *mutationResolver) getSceneMarker(ctx context.Context, id int) (ret *models.SceneMarker, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.SceneMarker.Find(ctx, id)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *mutationResolver) SceneMarkerCreate(ctx context.Context, input SceneMarkerCreateInput) (*models.SceneMarker, error) {
	primaryTagID, err := strconv.Atoi(input.PrimaryTagID)
	if err != nil {
		return nil, err
	}

	sceneID, err := strconv.Atoi(input.SceneID)
	if err != nil {
		return nil, err
	}

	currentTime := time.Now()
	newSceneMarker := models.SceneMarker{
		Title:        input.Title,
		Seconds:      input.Seconds,
		PrimaryTagID: primaryTagID,
		SceneID:      sql.NullInt64{Int64: int64(sceneID), Valid: sceneID != 0},
		CreatedAt:    models.SQLiteTimestamp{Timestamp: currentTime},
		UpdatedAt:    models.SQLiteTimestamp{Timestamp: currentTime},
	}

	tagIDs, err := stringslice.StringSliceToIntSlice(input.TagIds)
	if err != nil {
		return nil, err
	}

	ret, err := r.changeMarker(ctx, create, newSceneMarker, tagIDs)
	if err != nil {
		return nil, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, ret.ID, plugin.SceneMarkerCreatePost, input, nil)
	return r.getSceneMarker(ctx, ret.ID)
}

func (r *mutationResolver) SceneMarkerUpdate(ctx context.Context, input SceneMarkerUpdateInput) (*models.SceneMarker, error) {
	// Populate scene marker from the input
	sceneMarkerID, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, err
	}

	primaryTagID, err := strconv.Atoi(input.PrimaryTagID)
	if err != nil {
		return nil, err
	}

	sceneID, err := strconv.Atoi(input.SceneID)
	if err != nil {
		return nil, err
	}

	updatedSceneMarker := models.SceneMarker{
		ID:           sceneMarkerID,
		Title:        input.Title,
		Seconds:      input.Seconds,
		SceneID:      sql.NullInt64{Int64: int64(sceneID), Valid: sceneID != 0},
		PrimaryTagID: primaryTagID,
		UpdatedAt:    models.SQLiteTimestamp{Timestamp: time.Now()},
	}

	tagIDs, err := stringslice.StringSliceToIntSlice(input.TagIds)
	if err != nil {
		return nil, err
	}

	ret, err := r.changeMarker(ctx, update, updatedSceneMarker, tagIDs)
	if err != nil {
		return nil, err
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}
	r.hookExecutor.ExecutePostHooks(ctx, ret.ID, plugin.SceneMarkerUpdatePost, input, translator.getFields())
	return r.getSceneMarker(ctx, ret.ID)
}

func (r *mutationResolver) SceneMarkerDestroy(ctx context.Context, id string) (bool, error) {
	markerID, err := strconv.Atoi(id)
	if err != nil {
		return false, err
	}

	fileNamingAlgo := manager.GetInstance().Config.GetVideoFileNamingAlgorithm()

	fileDeleter := &scene.FileDeleter{
		Deleter:        *file.NewDeleter(),
		FileNamingAlgo: fileNamingAlgo,
		Paths:          manager.GetInstance().Paths,
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.SceneMarker
		sqb := r.repository.Scene

		marker, err := qb.Find(ctx, markerID)

		if err != nil {
			return err
		}

		if marker == nil {
			return fmt.Errorf("scene marker with id %d not found", markerID)
		}

		s, err := sqb.Find(ctx, int(marker.SceneID.Int64))
		if err != nil {
			return err
		}

		return scene.DestroyMarker(ctx, s, marker, qb, fileDeleter)
	}); err != nil {
		fileDeleter.Rollback()
		return false, err
	}

	// perform the post-commit actions
	fileDeleter.Commit()

	r.hookExecutor.ExecutePostHooks(ctx, markerID, plugin.SceneMarkerDestroyPost, id, nil)

	return true, nil
}

func (r *mutationResolver) changeMarker(ctx context.Context, changeType int, changedMarker models.SceneMarker, tagIDs []int) (*models.SceneMarker, error) {
	var existingMarker *models.SceneMarker
	var sceneMarker *models.SceneMarker
	var s *models.Scene

	fileNamingAlgo := manager.GetInstance().Config.GetVideoFileNamingAlgorithm()

	fileDeleter := &scene.FileDeleter{
		Deleter:        *file.NewDeleter(),
		FileNamingAlgo: fileNamingAlgo,
		Paths:          manager.GetInstance().Paths,
	}

	// Start the transaction and save the scene marker
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.SceneMarker
		sqb := r.repository.Scene

		var err error
		switch changeType {
		case create:
			sceneMarker, err = qb.Create(ctx, changedMarker)
		case update:
			// check to see if timestamp was changed
			existingMarker, err = qb.Find(ctx, changedMarker.ID)
			if err != nil {
				return err
			}
			sceneMarker, err = qb.Update(ctx, changedMarker)
			if err != nil {
				return err
			}

			s, err = sqb.Find(ctx, int(existingMarker.SceneID.Int64))
		}
		if err != nil {
			return err
		}

		// remove the marker preview if the timestamp was changed
		if s != nil && existingMarker != nil && existingMarker.Seconds != changedMarker.Seconds {
			seconds := int(existingMarker.Seconds)
			if err := fileDeleter.MarkMarkerFiles(s, seconds); err != nil {
				return err
			}
		}

		// Save the marker tags
		// If this tag is the primary tag, then let's not add it.
		tagIDs = intslice.IntExclude(tagIDs, []int{changedMarker.PrimaryTagID})
		return qb.UpdateTags(ctx, sceneMarker.ID, tagIDs)
	}); err != nil {
		fileDeleter.Rollback()
		return nil, err
	}

	// perform the post-commit actions
	fileDeleter.Commit()
	return sceneMarker, nil
}

func (r *mutationResolver) SceneIncrementO(ctx context.Context, id string) (ret int, err error) {
	sceneID, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Scene

		ret, err = qb.IncrementOCounter(ctx, sceneID)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *mutationResolver) SceneDecrementO(ctx context.Context, id string) (ret int, err error) {
	sceneID, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Scene

		ret, err = qb.DecrementOCounter(ctx, sceneID)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *mutationResolver) SceneResetO(ctx context.Context, id string) (ret int, err error) {
	sceneID, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Scene

		ret, err = qb.ResetOCounter(ctx, sceneID)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *mutationResolver) SceneGenerateScreenshot(ctx context.Context, id string, at *float64) (string, error) {
	if at != nil {
		manager.GetInstance().GenerateScreenshot(ctx, id, *at)
	} else {
		manager.GetInstance().GenerateDefaultScreenshot(ctx, id)
	}

	return "todo", nil
}
