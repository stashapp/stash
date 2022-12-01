package api

import (
	"context"
	"database/sql"
	"errors"
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
	"github.com/stashapp/stash/pkg/txn"
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

func (r *mutationResolver) SceneCreate(ctx context.Context, input SceneCreateInput) (ret *models.Scene, err error) {
	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	performerIDs, err := stringslice.StringSliceToIntSlice(input.PerformerIds)
	if err != nil {
		return nil, fmt.Errorf("converting performer ids: %w", err)
	}
	tagIDs, err := stringslice.StringSliceToIntSlice(input.TagIds)
	if err != nil {
		return nil, fmt.Errorf("converting tag ids: %w", err)
	}
	galleryIDs, err := stringslice.StringSliceToIntSlice(input.GalleryIds)
	if err != nil {
		return nil, fmt.Errorf("converting gallery ids: %w", err)
	}

	moviesScenes, err := models.MoviesScenesFromInput(input.Movies)
	if err != nil {
		return nil, fmt.Errorf("converting movies scenes: %w", err)
	}

	fileIDsInt, err := stringslice.StringSliceToIntSlice(input.FileIds)
	if err != nil {
		return nil, fmt.Errorf("converting file ids: %w", err)
	}

	fileIDs := make([]file.ID, len(fileIDsInt))
	for i, v := range fileIDsInt {
		fileIDs[i] = file.ID(v)
	}

	newScene := models.Scene{
		Title:        translator.string(input.Title, "title"),
		Code:         translator.string(input.Code, "code"),
		Details:      translator.string(input.Details, "details"),
		Director:     translator.string(input.Director, "director"),
		URL:          translator.string(input.URL, "url"),
		Date:         translator.datePtr(input.Date, "date"),
		Rating:       translator.ratingConversionInt(input.Rating, input.Rating100),
		Organized:    translator.bool(input.Organized, "organized"),
		PerformerIDs: models.NewRelatedIDs(performerIDs),
		TagIDs:       models.NewRelatedIDs(tagIDs),
		GalleryIDs:   models.NewRelatedIDs(galleryIDs),
		Movies:       models.NewRelatedMovies(moviesScenes),
		StashIDs:     models.NewRelatedStashIDs(stashIDPtrSliceToSlice(input.StashIds)),
	}

	newScene.StudioID, err = translator.intPtrFromString(input.StudioID, "studio_id")
	if err != nil {
		return nil, fmt.Errorf("converting studio id: %w", err)
	}

	var coverImageData []byte
	if input.CoverImage != nil && *input.CoverImage != "" {
		var err error
		coverImageData, err = utils.ProcessImageInput(ctx, *input.CoverImage)
		if err != nil {
			return nil, err
		}
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.Resolver.sceneService.Create(ctx, &newScene, fileIDs, coverImageData)
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

func scenePartialFromInput(input models.SceneUpdateInput, translator changesetTranslator) (*models.ScenePartial, error) {
	updatedScene := models.NewScenePartial()
	updatedScene.Title = translator.optionalString(input.Title, "title")
	updatedScene.Code = translator.optionalString(input.Code, "code")
	updatedScene.Details = translator.optionalString(input.Details, "details")
	updatedScene.Director = translator.optionalString(input.Director, "director")
	updatedScene.URL = translator.optionalString(input.URL, "url")
	updatedScene.Date = translator.optionalDate(input.Date, "date")
	updatedScene.Rating = translator.ratingConversionOptional(input.Rating, input.Rating100)
	updatedScene.OCounter = translator.optionalInt(input.OCounter, "o_counter")
	updatedScene.PlayCount = translator.optionalInt(input.PlayCount, "play_count")
	updatedScene.PlayDuration = translator.optionalFloat64(input.PlayDuration, "play_duration")
	var err error
	updatedScene.StudioID, err = translator.optionalIntFromString(input.StudioID, "studio_id")
	if err != nil {
		return nil, fmt.Errorf("converting studio id: %w", err)
	}

	updatedScene.Organized = translator.optionalBool(input.Organized, "organized")

	if input.PrimaryFileID != nil {
		primaryFileID, err := strconv.Atoi(*input.PrimaryFileID)
		if err != nil {
			return nil, fmt.Errorf("converting primary file id: %w", err)
		}

		converted := file.ID(primaryFileID)
		updatedScene.PrimaryFileID = &converted
	}

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

	return &updatedScene, nil
}

func (r *mutationResolver) sceneUpdate(ctx context.Context, input models.SceneUpdateInput, translator changesetTranslator) (*models.Scene, error) {
	// Populate scene from the input
	sceneID, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, err
	}

	qb := r.repository.Scene

	s, err := qb.Find(ctx, sceneID)
	if err != nil {
		return nil, err
	}

	if s == nil {
		return nil, fmt.Errorf("scene with id %d not found", sceneID)
	}

	var coverImageData []byte

	updatedScene, err := scenePartialFromInput(input, translator)
	if err != nil {
		return nil, err
	}

	// ensure that title is set where scene has no file
	if updatedScene.Title.Set && updatedScene.Title.Value == "" {
		if err := s.LoadFiles(ctx, r.repository.Scene); err != nil {
			return nil, err
		}

		if len(s.Files.List()) == 0 {
			return nil, errors.New("title must be set if scene has no files")
		}
	}

	if updatedScene.PrimaryFileID != nil {
		newPrimaryFileID := *updatedScene.PrimaryFileID

		// if file hash has changed, we should migrate generated files
		// after commit
		if err := s.LoadFiles(ctx, r.repository.Scene); err != nil {
			return nil, err
		}

		// ensure that new primary file is associated with scene
		var f *file.VideoFile
		for _, ff := range s.Files.List() {
			if ff.ID == newPrimaryFileID {
				f = ff
			}
		}

		if f == nil {
			return nil, fmt.Errorf("file with id %d not associated with scene", newPrimaryFileID)
		}
	}

	if input.CoverImage != nil && *input.CoverImage != "" {
		var err error
		coverImageData, err = utils.ProcessImageInput(ctx, *input.CoverImage)
		if err != nil {
			return nil, err
		}
	}

	s, err = qb.UpdatePartial(ctx, sceneID, *updatedScene)
	if err != nil {
		return nil, err
	}

	if err := r.sceneUpdateCoverImage(ctx, s, coverImageData); err != nil {
		return nil, err
	}

	return s, nil
}

func (r *mutationResolver) sceneUpdateCoverImage(ctx context.Context, s *models.Scene, coverImageData []byte) error {
	if len(coverImageData) > 0 {
		qb := r.repository.Scene

		// update cover table
		if err := qb.UpdateCover(ctx, s.ID, coverImageData); err != nil {
			return err
		}

		if s.Path != "" {
			// update the file-based screenshot after commit
			txn.AddPostCommitHook(ctx, func(ctx context.Context) error {
				return scene.SetScreenshot(manager.GetInstance().Paths, s.GetHash(config.GetInstance().GetVideoFileNamingAlgorithm()), coverImageData)
			})
		}
	}

	return nil
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
	updatedScene.Code = translator.optionalString(input.Code, "code")
	updatedScene.Details = translator.optionalString(input.Details, "details")
	updatedScene.Director = translator.optionalString(input.Director, "director")
	updatedScene.URL = translator.optionalString(input.URL, "url")
	updatedScene.Date = translator.optionalDate(input.Date, "date")
	updatedScene.Rating = translator.ratingConversionOptional(input.Rating, input.Rating100)
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
	if translator.hasField("movie_ids") {
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

func (r *mutationResolver) SceneDestroy(ctx context.Context, input models.SceneDestroyInput) (bool, error) {
	sceneID, err := strconv.Atoi(input.ID)
	if err != nil {
		return false, err
	}

	fileNamingAlgo := manager.GetInstance().Config.GetVideoFileNamingAlgorithm()

	var s *models.Scene
	fileDeleter := &scene.FileDeleter{
		Deleter:        file.NewDeleter(),
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

		return r.sceneService.Destroy(ctx, s, fileDeleter, deleteGenerated, deleteFile)
	}); err != nil {
		fileDeleter.Rollback()
		return false, err
	}

	// perform the post-commit actions
	fileDeleter.Commit()

	// call post hook after performing the other actions
	r.hookExecutor.ExecutePostHooks(ctx, s.ID, plugin.SceneDestroyPost, plugin.SceneDestroyInput{
		SceneDestroyInput: input,
		Checksum:          s.Checksum,
		OSHash:            s.OSHash,
		Path:              s.Path,
	}, nil)

	return true, nil
}

func (r *mutationResolver) ScenesDestroy(ctx context.Context, input models.ScenesDestroyInput) (bool, error) {
	var scenes []*models.Scene
	fileNamingAlgo := manager.GetInstance().Config.GetVideoFileNamingAlgorithm()

	fileDeleter := &scene.FileDeleter{
		Deleter:        file.NewDeleter(),
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

			if err := r.sceneService.Destroy(ctx, s, fileDeleter, deleteGenerated, deleteFile); err != nil {
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
			Checksum:           scene.Checksum,
			OSHash:             scene.OSHash,
			Path:               scene.Path,
		}, nil)
	}

	return true, nil
}

func (r *mutationResolver) SceneAssignFile(ctx context.Context, input AssignSceneFileInput) (bool, error) {
	sceneID, err := strconv.Atoi(input.SceneID)
	if err != nil {
		return false, fmt.Errorf("converting scene ID: %w", err)
	}

	fileIDInt, err := strconv.Atoi(input.FileID)
	if err != nil {
		return false, fmt.Errorf("converting file ID: %w", err)
	}

	fileID := file.ID(fileIDInt)

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		return r.Resolver.sceneService.AssignFile(ctx, sceneID, fileID)
	}); err != nil {
		return false, fmt.Errorf("assigning file to scene: %w", err)
	}

	return true, nil
}

func (r *mutationResolver) SceneMerge(ctx context.Context, input SceneMergeInput) (*models.Scene, error) {
	srcIDs, err := stringslice.StringSliceToIntSlice(input.Source)
	if err != nil {
		return nil, fmt.Errorf("converting source IDs: %w", err)
	}

	destID, err := strconv.Atoi(input.Destination)
	if err != nil {
		return nil, fmt.Errorf("converting destination ID %s: %w", input.Destination, err)
	}

	var values *models.ScenePartial
	if input.Values != nil {
		translator := changesetTranslator{
			inputMap: getNamedUpdateInputMap(ctx, "input.values"),
		}

		values, err = scenePartialFromInput(*input.Values, translator)
		if err != nil {
			return nil, err
		}
	} else {
		v := models.NewScenePartial()
		values = &v
	}

	var coverImageData []byte

	if input.Values.CoverImage != nil && *input.Values.CoverImage != "" {
		var err error
		coverImageData, err = utils.ProcessImageInput(ctx, *input.Values.CoverImage)
		if err != nil {
			return nil, err
		}
	}

	var ret *models.Scene
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		if err := r.Resolver.sceneService.Merge(ctx, srcIDs, destID, *values); err != nil {
			return err
		}

		ret, err = r.Resolver.repository.Scene.Find(ctx, destID)

		if err == nil && ret != nil {
			err = r.sceneUpdateCoverImage(ctx, ret, coverImageData)
		}

		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
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
		Deleter:        file.NewDeleter(),
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
		Deleter:        file.NewDeleter(),
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

func (r *mutationResolver) SceneSaveActivity(ctx context.Context, id string, resumeTime *float64, playDuration *float64) (ret bool, err error) {
	sceneID, err := strconv.Atoi(id)
	if err != nil {
		return false, err
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Scene

		ret, err = qb.SaveActivity(ctx, sceneID, resumeTime, playDuration)
		return err
	}); err != nil {
		return false, err
	}

	return ret, nil
}

func (r *mutationResolver) SceneIncrementPlayCount(ctx context.Context, id string) (ret int, err error) {
	sceneID, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Scene

		ret, err = qb.IncrementWatchCount(ctx, sceneID)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
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
