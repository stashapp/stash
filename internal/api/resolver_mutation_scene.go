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
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		ret, err = repo.Scene().Find(id)
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
	if err := r.withTxn(ctx, func(repo models.Repository) error {
		ret, err = r.sceneUpdate(ctx, input, translator, repo)
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
	if err := r.withTxn(ctx, func(repo models.Repository) error {
		for i, scene := range input {
			translator := changesetTranslator{
				inputMap: inputMaps[i],
			}

			thisScene, err := r.sceneUpdate(ctx, *scene, translator, repo)
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

func (r *mutationResolver) sceneUpdate(ctx context.Context, input models.SceneUpdateInput, translator changesetTranslator, repo models.Repository) (*models.Scene, error) {
	// Populate scene from the input
	sceneID, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, err
	}

	var coverImageData []byte

	updatedTime := time.Now()
	updatedScene := models.ScenePartial{
		ID:        sceneID,
		UpdatedAt: &models.SQLiteTimestamp{Timestamp: updatedTime},
	}

	updatedScene.Title = translator.nullString(input.Title, "title")
	updatedScene.Details = translator.nullString(input.Details, "details")
	updatedScene.URL = translator.nullString(input.URL, "url")
	updatedScene.Date = translator.sqliteDate(input.Date, "date")
	updatedScene.Rating = translator.nullInt64(input.Rating, "rating")
	updatedScene.StudioID = translator.nullInt64FromString(input.StudioID, "studio_id")
	updatedScene.Organized = input.Organized

	if input.CoverImage != nil && *input.CoverImage != "" {
		var err error
		coverImageData, err = utils.ProcessImageInput(ctx, *input.CoverImage)
		if err != nil {
			return nil, err
		}

		// update the cover after updating the scene
	}

	qb := repo.Scene()
	s, err := qb.Update(updatedScene)
	if err != nil {
		return nil, err
	}

	// update cover table
	if len(coverImageData) > 0 {
		if err := qb.UpdateCover(sceneID, coverImageData); err != nil {
			return nil, err
		}
	}

	// Save the performers
	if translator.hasField("performer_ids") {
		if err := r.updateScenePerformers(qb, sceneID, input.PerformerIds); err != nil {
			return nil, err
		}
	}

	// Save the movies
	if translator.hasField("movies") {
		if err := r.updateSceneMovies(qb, sceneID, input.Movies); err != nil {
			return nil, err
		}
	}

	// Save the tags
	if translator.hasField("tag_ids") {
		if err := r.updateSceneTags(qb, sceneID, input.TagIds); err != nil {
			return nil, err
		}
	}

	// Save the galleries
	if translator.hasField("gallery_ids") {
		if err := r.updateSceneGalleries(qb, sceneID, input.GalleryIds); err != nil {
			return nil, err
		}
	}

	// Save the stash_ids
	if translator.hasField("stash_ids") {
		stashIDJoins := models.StashIDsFromInput(input.StashIds)
		if err := qb.UpdateStashIDs(sceneID, stashIDJoins); err != nil {
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

func (r *mutationResolver) updateScenePerformers(qb models.SceneReaderWriter, sceneID int, performerIDs []string) error {
	ids, err := stringslice.StringSliceToIntSlice(performerIDs)
	if err != nil {
		return err
	}
	return qb.UpdatePerformers(sceneID, ids)
}

func (r *mutationResolver) updateSceneMovies(qb models.SceneReaderWriter, sceneID int, movies []*models.SceneMovieInput) error {
	var movieJoins []models.MoviesScenes

	for _, movie := range movies {
		movieID, err := strconv.Atoi(movie.MovieID)
		if err != nil {
			return err
		}

		movieJoin := models.MoviesScenes{
			MovieID: movieID,
		}

		if movie.SceneIndex != nil {
			movieJoin.SceneIndex = sql.NullInt64{
				Int64: int64(*movie.SceneIndex),
				Valid: true,
			}
		}

		movieJoins = append(movieJoins, movieJoin)
	}

	return qb.UpdateMovies(sceneID, movieJoins)
}

func (r *mutationResolver) updateSceneTags(qb models.SceneReaderWriter, sceneID int, tagsIDs []string) error {
	ids, err := stringslice.StringSliceToIntSlice(tagsIDs)
	if err != nil {
		return err
	}
	return qb.UpdateTags(sceneID, ids)
}

func (r *mutationResolver) updateSceneGalleries(qb models.SceneReaderWriter, sceneID int, galleryIDs []string) error {
	ids, err := stringslice.StringSliceToIntSlice(galleryIDs)
	if err != nil {
		return err
	}
	return qb.UpdateGalleries(sceneID, ids)
}

func (r *mutationResolver) BulkSceneUpdate(ctx context.Context, input BulkSceneUpdateInput) ([]*models.Scene, error) {
	sceneIDs, err := stringslice.StringSliceToIntSlice(input.Ids)
	if err != nil {
		return nil, err
	}

	// Populate scene from the input
	updatedTime := time.Now()

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	updatedScene := models.ScenePartial{
		UpdatedAt: &models.SQLiteTimestamp{Timestamp: updatedTime},
	}

	updatedScene.Title = translator.nullString(input.Title, "title")
	updatedScene.Details = translator.nullString(input.Details, "details")
	updatedScene.URL = translator.nullString(input.URL, "url")
	updatedScene.Date = translator.sqliteDate(input.Date, "date")
	updatedScene.Rating = translator.nullInt64(input.Rating, "rating")
	updatedScene.StudioID = translator.nullInt64FromString(input.StudioID, "studio_id")
	updatedScene.Organized = input.Organized

	ret := []*models.Scene{}

	// Start the transaction and save the scene marker
	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Scene()

		for _, sceneID := range sceneIDs {
			updatedScene.ID = sceneID

			scene, err := qb.Update(updatedScene)
			if err != nil {
				return err
			}

			ret = append(ret, scene)

			// Save the performers
			if translator.hasField("performer_ids") {
				performerIDs, err := adjustScenePerformerIDs(qb, sceneID, *input.PerformerIds)
				if err != nil {
					return err
				}

				if err := qb.UpdatePerformers(sceneID, performerIDs); err != nil {
					return err
				}
			}

			// Save the tags
			if translator.hasField("tag_ids") {
				tagIDs, err := adjustTagIDs(qb, sceneID, *input.TagIds)
				if err != nil {
					return err
				}

				if err := qb.UpdateTags(sceneID, tagIDs); err != nil {
					return err
				}
			}

			// Save the galleries
			if translator.hasField("gallery_ids") {
				galleryIDs, err := adjustSceneGalleryIDs(qb, sceneID, *input.GalleryIds)
				if err != nil {
					return err
				}

				if err := qb.UpdateGalleries(sceneID, galleryIDs); err != nil {
					return err
				}
			}

			// Save the movies
			if translator.hasField("movie_ids") {
				movies, err := adjustSceneMovieIDs(qb, sceneID, *input.MovieIds)
				if err != nil {
					return err
				}

				if err := qb.UpdateMovies(sceneID, movies); err != nil {
					return err
				}
			}
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
	if updateIDs.Mode == BulkUpdateIDModeSet {
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
				if updateIDs.Mode == BulkUpdateIDModeRemove {
					// remove from the list
					existingIDs = append(existingIDs[:idx], existingIDs[idx+1:]...)
				}

				foundExisting = true
				break
			}
		}

		if !foundExisting && updateIDs.Mode != BulkUpdateIDModeRemove {
			existingIDs = append(existingIDs, id)
		}
	}

	return existingIDs
}

func adjustScenePerformerIDs(qb models.SceneReader, sceneID int, ids BulkUpdateIds) (ret []int, err error) {
	ret, err = qb.GetPerformerIDs(sceneID)
	if err != nil {
		return nil, err
	}

	return adjustIDs(ret, ids), nil
}

type tagIDsGetter interface {
	GetTagIDs(id int) ([]int, error)
}

func adjustTagIDs(qb tagIDsGetter, sceneID int, ids BulkUpdateIds) (ret []int, err error) {
	ret, err = qb.GetTagIDs(sceneID)
	if err != nil {
		return nil, err
	}

	return adjustIDs(ret, ids), nil
}

func adjustSceneGalleryIDs(qb models.SceneReader, sceneID int, ids BulkUpdateIds) (ret []int, err error) {
	ret, err = qb.GetGalleryIDs(sceneID)
	if err != nil {
		return nil, err
	}

	return adjustIDs(ret, ids), nil
}

func adjustSceneMovieIDs(qb models.SceneReader, sceneID int, updateIDs BulkUpdateIds) ([]models.MoviesScenes, error) {
	existingMovies, err := qb.GetMovies(sceneID)
	if err != nil {
		return nil, err
	}

	// if we are setting the ids, just return the ids
	if updateIDs.Mode == BulkUpdateIDModeSet {
		existingMovies = []models.MoviesScenes{}
		for _, idStr := range updateIDs.Ids {
			id, _ := strconv.Atoi(idStr)
			existingMovies = append(existingMovies, models.MoviesScenes{MovieID: id})
		}

		return existingMovies, nil
	}

	for _, idStr := range updateIDs.Ids {
		id, _ := strconv.Atoi(idStr)

		// look for the id in the list
		foundExisting := false
		for idx, existingMovie := range existingMovies {
			if existingMovie.MovieID == id {
				if updateIDs.Mode == BulkUpdateIDModeRemove {
					// remove from the list
					existingMovies = append(existingMovies[:idx], existingMovies[idx+1:]...)
				}

				foundExisting = true
				break
			}
		}

		if !foundExisting && updateIDs.Mode != BulkUpdateIDModeRemove {
			existingMovies = append(existingMovies, models.MoviesScenes{MovieID: id})
		}
	}

	return existingMovies, err
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

	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Scene()
		var err error
		s, err = qb.Find(sceneID)
		if err != nil {
			return err
		}

		if s == nil {
			return fmt.Errorf("scene with id %d not found", sceneID)
		}

		// kill any running encoders
		manager.KillRunningStreams(s, fileNamingAlgo)

		return scene.Destroy(s, repo, fileDeleter, deleteGenerated, deleteFile)
	}); err != nil {
		fileDeleter.Rollback()
		return false, err
	}

	// perform the post-commit actions
	fileDeleter.Commit()

	// call post hook after performing the other actions
	r.hookExecutor.ExecutePostHooks(ctx, s.ID, plugin.SceneDestroyPost, plugin.SceneDestroyInput{
		SceneDestroyInput: input,
		Checksum:          s.Checksum.String,
		OSHash:            s.OSHash.String,
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

	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Scene()

		for _, id := range input.Ids {
			sceneID, _ := strconv.Atoi(id)

			s, err := qb.Find(sceneID)
			if err != nil {
				return err
			}
			if s != nil {
				scenes = append(scenes, s)
			}

			// kill any running encoders
			manager.KillRunningStreams(s, fileNamingAlgo)

			if err := scene.Destroy(s, repo, fileDeleter, deleteGenerated, deleteFile); err != nil {
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
			Checksum:           scene.Checksum.String,
			OSHash:             scene.OSHash.String,
			Path:               scene.Path,
		}, nil)
	}

	return true, nil
}

func (r *mutationResolver) getSceneMarker(ctx context.Context, id int) (ret *models.SceneMarker, err error) {
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		ret, err = repo.SceneMarker().Find(id)
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

	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.SceneMarker()
		sqb := repo.Scene()

		marker, err := qb.Find(markerID)

		if err != nil {
			return err
		}

		if marker == nil {
			return fmt.Errorf("scene marker with id %d not found", markerID)
		}

		s, err := sqb.Find(int(marker.SceneID.Int64))
		if err != nil {
			return err
		}

		return scene.DestroyMarker(s, marker, qb, fileDeleter)
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
	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.SceneMarker()
		sqb := repo.Scene()

		var err error
		switch changeType {
		case create:
			sceneMarker, err = qb.Create(changedMarker)
		case update:
			// check to see if timestamp was changed
			existingMarker, err = qb.Find(changedMarker.ID)
			if err != nil {
				return err
			}
			sceneMarker, err = qb.Update(changedMarker)
			if err != nil {
				return err
			}

			s, err = sqb.Find(int(existingMarker.SceneID.Int64))
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
		return qb.UpdateTags(sceneMarker.ID, tagIDs)
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

	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Scene()

		ret, err = qb.IncrementOCounter(sceneID)
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

	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Scene()

		ret, err = qb.DecrementOCounter(sceneID)
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

	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Scene()

		ret, err = qb.ResetOCounter(sceneID)
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
