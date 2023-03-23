package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *mutationResolver) getGallery(ctx context.Context, id int) (ret *models.Gallery, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Gallery.Find(ctx, id)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *mutationResolver) GalleryCreate(ctx context.Context, input GalleryCreateInput) (*models.Gallery, error) {
	// name must be provided
	if input.Title == "" {
		return nil, errors.New("title must not be empty")
	}

	// Populate a new performer from the input
	performerIDs, err := stringslice.StringSliceToIntSlice(input.PerformerIds)
	if err != nil {
		return nil, fmt.Errorf("converting performer ids: %w", err)
	}
	tagIDs, err := stringslice.StringSliceToIntSlice(input.TagIds)
	if err != nil {
		return nil, fmt.Errorf("converting tag ids: %w", err)
	}
	sceneIDs, err := stringslice.StringSliceToIntSlice(input.SceneIds)
	if err != nil {
		return nil, fmt.Errorf("converting scene ids: %w", err)
	}

	currentTime := time.Now()
	newGallery := models.Gallery{
		Title:        input.Title,
		PerformerIDs: models.NewRelatedIDs(performerIDs),
		TagIDs:       models.NewRelatedIDs(tagIDs),
		SceneIDs:     models.NewRelatedIDs(sceneIDs),
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
	}
	if input.URL != nil {
		newGallery.URL = *input.URL
	}
	if input.Details != nil {
		newGallery.Details = *input.Details
	}

	if input.Date != nil {
		d := models.NewDate(*input.Date)
		newGallery.Date = &d
	}

	if input.Rating100 != nil {
		newGallery.Rating = input.Rating100
	} else if input.Rating != nil {
		rating := models.Rating5To100(*input.Rating)
		newGallery.Rating = &rating
	}

	if input.StudioID != nil {
		studioID, _ := strconv.Atoi(*input.StudioID)
		newGallery.StudioID = &studioID
	}

	// Start the transaction and save the gallery
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Gallery
		if err := qb.Create(ctx, &newGallery, nil); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, newGallery.ID, plugin.GalleryCreatePost, input, nil)
	return r.getGallery(ctx, newGallery.ID)
}

type GallerySceneUpdater interface {
	UpdateScenes(ctx context.Context, galleryID int, sceneIDs []int) error
}

func (r *mutationResolver) GalleryUpdate(ctx context.Context, input models.GalleryUpdateInput) (ret *models.Gallery, err error) {
	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Start the transaction and save the gallery
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.galleryUpdate(ctx, input, translator)
		return err
	}); err != nil {
		return nil, err
	}

	// execute post hooks outside txn
	r.hookExecutor.ExecutePostHooks(ctx, ret.ID, plugin.GalleryUpdatePost, input, translator.getFields())
	return r.getGallery(ctx, ret.ID)
}

func (r *mutationResolver) GalleriesUpdate(ctx context.Context, input []*models.GalleryUpdateInput) (ret []*models.Gallery, err error) {
	inputMaps := getUpdateInputMaps(ctx)

	// Start the transaction and save the gallery
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		for i, gallery := range input {
			translator := changesetTranslator{
				inputMap: inputMaps[i],
			}

			thisGallery, err := r.galleryUpdate(ctx, *gallery, translator)
			if err != nil {
				return err
			}

			ret = append(ret, thisGallery)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// execute post hooks outside txn
	var newRet []*models.Gallery
	for i, gallery := range ret {
		translator := changesetTranslator{
			inputMap: inputMaps[i],
		}

		r.hookExecutor.ExecutePostHooks(ctx, gallery.ID, plugin.GalleryUpdatePost, input, translator.getFields())
		gallery, err = r.getGallery(ctx, gallery.ID)
		if err != nil {
			return nil, err
		}

		newRet = append(newRet, gallery)
	}

	return newRet, nil
}

func (r *mutationResolver) galleryUpdate(ctx context.Context, input models.GalleryUpdateInput, translator changesetTranslator) (*models.Gallery, error) {
	qb := r.repository.Gallery

	// Populate gallery from the input
	galleryID, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, err
	}

	originalGallery, err := qb.Find(ctx, galleryID)
	if err != nil {
		return nil, err
	}

	if originalGallery == nil {
		return nil, errors.New("not found")
	}

	updatedGallery := models.NewGalleryPartial()

	if input.Title != nil {
		// ensure title is not empty
		if *input.Title == "" && originalGallery.IsUserCreated() {
			return nil, errors.New("title must not be empty for user-created galleries")
		}

		updatedGallery.Title = models.NewOptionalString(*input.Title)
	}

	updatedGallery.Details = translator.optionalString(input.Details, "details")
	updatedGallery.URL = translator.optionalString(input.URL, "url")
	updatedGallery.Date = translator.optionalDate(input.Date, "date")
	updatedGallery.Rating = translator.ratingConversionOptional(input.Rating, input.Rating100)
	updatedGallery.StudioID, err = translator.optionalIntFromString(input.StudioID, "studio_id")
	if err != nil {
		return nil, fmt.Errorf("converting studio id: %w", err)
	}
	updatedGallery.Organized = translator.optionalBool(input.Organized, "organized")

	if input.PrimaryFileID != nil {
		primaryFileID, err := strconv.Atoi(*input.PrimaryFileID)
		if err != nil {
			return nil, fmt.Errorf("converting primary file id: %w", err)
		}

		converted := file.ID(primaryFileID)
		updatedGallery.PrimaryFileID = &converted

		if err := originalGallery.LoadFiles(ctx, r.repository.Gallery); err != nil {
			return nil, err
		}

		// ensure that new primary file is associated with scene
		var f file.File
		for _, ff := range originalGallery.Files.List() {
			if ff.Base().ID == converted {
				f = ff
			}
		}

		if f == nil {
			return nil, fmt.Errorf("file with id %d not associated with gallery", converted)
		}
	}

	if translator.hasField("performer_ids") {
		updatedGallery.PerformerIDs, err = translateUpdateIDs(input.PerformerIds, models.RelationshipUpdateModeSet)
		if err != nil {
			return nil, fmt.Errorf("converting performer ids: %w", err)
		}
	}

	if translator.hasField("tag_ids") {
		updatedGallery.TagIDs, err = translateUpdateIDs(input.TagIds, models.RelationshipUpdateModeSet)
		if err != nil {
			return nil, fmt.Errorf("converting tag ids: %w", err)
		}
	}

	if translator.hasField("scene_ids") {
		updatedGallery.SceneIDs, err = translateUpdateIDs(input.SceneIds, models.RelationshipUpdateModeSet)
		if err != nil {
			return nil, fmt.Errorf("converting scene ids: %w", err)
		}
	}

	// gallery scene is set from the scene only

	gallery, err := qb.UpdatePartial(ctx, galleryID, updatedGallery)
	if err != nil {
		return nil, err
	}

	return gallery, nil
}

func (r *mutationResolver) BulkGalleryUpdate(ctx context.Context, input BulkGalleryUpdateInput) ([]*models.Gallery, error) {
	// Populate gallery from the input
	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	updatedGallery := models.NewGalleryPartial()

	updatedGallery.Details = translator.optionalString(input.Details, "details")
	updatedGallery.URL = translator.optionalString(input.URL, "url")
	updatedGallery.Date = translator.optionalDate(input.Date, "date")
	updatedGallery.Rating = translator.ratingConversionOptional(input.Rating, input.Rating100)
	var err error
	updatedGallery.StudioID, err = translator.optionalIntFromString(input.StudioID, "studio_id")
	if err != nil {
		return nil, fmt.Errorf("converting studio id: %w", err)
	}
	updatedGallery.Organized = translator.optionalBool(input.Organized, "organized")

	if translator.hasField("performer_ids") {
		updatedGallery.PerformerIDs, err = translateUpdateIDs(input.PerformerIds.Ids, input.PerformerIds.Mode)
		if err != nil {
			return nil, fmt.Errorf("converting performer ids: %w", err)
		}
	}

	if translator.hasField("tag_ids") {
		updatedGallery.TagIDs, err = translateUpdateIDs(input.TagIds.Ids, input.TagIds.Mode)
		if err != nil {
			return nil, fmt.Errorf("converting tag ids: %w", err)
		}
	}

	if translator.hasField("scene_ids") {
		updatedGallery.SceneIDs, err = translateUpdateIDs(input.SceneIds.Ids, input.SceneIds.Mode)
		if err != nil {
			return nil, fmt.Errorf("converting scene ids: %w", err)
		}
	}

	ret := []*models.Gallery{}

	// Start the transaction and save the galleries
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Gallery

		for _, galleryIDStr := range input.Ids {
			galleryID, _ := strconv.Atoi(galleryIDStr)

			gallery, err := qb.UpdatePartial(ctx, galleryID, updatedGallery)
			if err != nil {
				return err
			}

			ret = append(ret, gallery)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// execute post hooks outside of txn
	var newRet []*models.Gallery
	for _, gallery := range ret {
		r.hookExecutor.ExecutePostHooks(ctx, gallery.ID, plugin.GalleryUpdatePost, input, translator.getFields())

		gallery, err := r.getGallery(ctx, gallery.ID)
		if err != nil {
			return nil, err
		}

		newRet = append(newRet, gallery)
	}

	return newRet, nil
}

type GallerySceneGetter interface {
	GetSceneIDs(ctx context.Context, galleryID int) ([]int, error)
}

func (r *mutationResolver) GalleryDestroy(ctx context.Context, input models.GalleryDestroyInput) (bool, error) {
	galleryIDs, err := stringslice.StringSliceToIntSlice(input.Ids)
	if err != nil {
		return false, err
	}

	var galleries []*models.Gallery
	var imgsDestroyed []*models.Image
	fileDeleter := &image.FileDeleter{
		Deleter: file.NewDeleter(),
		Paths:   manager.GetInstance().Paths,
	}

	deleteGenerated := utils.IsTrue(input.DeleteGenerated)
	deleteFile := utils.IsTrue(input.DeleteFile)

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Gallery

		for _, id := range galleryIDs {
			gallery, err := qb.Find(ctx, id)
			if err != nil {
				return err
			}

			if gallery == nil {
				return fmt.Errorf("gallery with id %d not found", id)
			}

			if err := gallery.LoadFiles(ctx, qb); err != nil {
				return fmt.Errorf("loading files for gallery %d", id)
			}

			galleries = append(galleries, gallery)

			imgsDestroyed, err = r.galleryService.Destroy(ctx, gallery, fileDeleter, deleteGenerated, deleteFile)
			if err != nil {
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

	for _, gallery := range galleries {
		// don't delete stash library paths
		path := gallery.Path
		if deleteFile && path != "" && !isStashPath(path) {
			// try to remove the folder - it is possible that it is not empty
			// so swallow the error if present
			_ = os.Remove(path)
		}
	}

	// call post hook after performing the other actions
	for _, gallery := range galleries {
		r.hookExecutor.ExecutePostHooks(ctx, gallery.ID, plugin.GalleryDestroyPost, plugin.GalleryDestroyInput{
			GalleryDestroyInput: input,
			Checksum:            gallery.PrimaryChecksum(),
			Path:                gallery.Path,
		}, nil)
	}

	// call image destroy post hook as well
	for _, img := range imgsDestroyed {
		r.hookExecutor.ExecutePostHooks(ctx, img.ID, plugin.ImageDestroyPost, plugin.ImageDestroyInput{
			Checksum: img.Checksum,
			Path:     img.Path,
		}, nil)
	}

	return true, nil
}

func isStashPath(path string) bool {
	stashConfigs := manager.GetInstance().Config.GetStashPaths()
	for _, config := range stashConfigs {
		if path == config.Path {
			return true
		}
	}

	return false
}

func (r *mutationResolver) AddGalleryImages(ctx context.Context, input GalleryAddInput) (bool, error) {
	galleryID, err := strconv.Atoi(input.GalleryID)
	if err != nil {
		return false, err
	}

	imageIDs, err := stringslice.StringSliceToIntSlice(input.ImageIds)
	if err != nil {
		return false, err
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Gallery
		gallery, err := qb.Find(ctx, galleryID)
		if err != nil {
			return err
		}

		if gallery == nil {
			return errors.New("gallery not found")
		}

		return r.galleryService.AddImages(ctx, gallery, imageIDs...)
	}); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) RemoveGalleryImages(ctx context.Context, input GalleryRemoveInput) (bool, error) {
	galleryID, err := strconv.Atoi(input.GalleryID)
	if err != nil {
		return false, err
	}

	imageIDs, err := stringslice.StringSliceToIntSlice(input.ImageIds)
	if err != nil {
		return false, err
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Gallery
		gallery, err := qb.Find(ctx, galleryID)
		if err != nil {
			return err
		}

		if gallery == nil {
			return errors.New("gallery not found")
		}

		return r.galleryService.RemoveImages(ctx, gallery, imageIDs...)
	}); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) getGalleryChapter(ctx context.Context, id int) (ret *models.GalleryChapter, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.GalleryChapter.Find(ctx, id)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *mutationResolver) GalleryChapterCreate(ctx context.Context, input GalleryChapterCreateInput) (*models.GalleryChapter, error) {
	galleryID, err := strconv.Atoi(input.GalleryID)
	if err != nil {
		return nil, err
	}

	var imageCount int
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		imageCount, err = r.repository.Image.CountByGalleryID(ctx, galleryID)
		return err
	}); err != nil {
		return nil, err
	}
	// Sanity Check of Index
	if input.ImageIndex > imageCount || input.ImageIndex < 1 {
		return nil, errors.New("Image # must greater than zero and in range of the gallery images")
	}

	currentTime := time.Now()
	newGalleryChapter := models.GalleryChapter{
		Title:      input.Title,
		ImageIndex: input.ImageIndex,
		GalleryID:  sql.NullInt64{Int64: int64(galleryID), Valid: galleryID != 0},
		CreatedAt:  models.SQLiteTimestamp{Timestamp: currentTime},
		UpdatedAt:  models.SQLiteTimestamp{Timestamp: currentTime},
	}

	if err != nil {
		return nil, err
	}

	ret, err := r.changeChapter(ctx, create, newGalleryChapter)
	if err != nil {
		return nil, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, ret.ID, plugin.GalleryChapterCreatePost, input, nil)
	return r.getGalleryChapter(ctx, ret.ID)
}

func (r *mutationResolver) GalleryChapterUpdate(ctx context.Context, input GalleryChapterUpdateInput) (*models.GalleryChapter, error) {
	// Populate gallery chapter from the input
	galleryChapterID, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, err
	}

	galleryID, err := strconv.Atoi(input.GalleryID)
	if err != nil {
		return nil, err
	}

	var imageCount int
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		imageCount, err = r.repository.Image.CountByGalleryID(ctx, galleryID)
		return err
	}); err != nil {
		return nil, err
	}
	// Sanity Check of Index
	if input.ImageIndex > imageCount || input.ImageIndex < 1 {
		return nil, errors.New("Image # must greater than zero and in range of the gallery images")
	}

	updatedGalleryChapter := models.GalleryChapter{
		ID:         galleryChapterID,
		Title:      input.Title,
		ImageIndex: input.ImageIndex,
		GalleryID:  sql.NullInt64{Int64: int64(galleryID), Valid: galleryID != 0},
		UpdatedAt:  models.SQLiteTimestamp{Timestamp: time.Now()},
	}

	ret, err := r.changeChapter(ctx, update, updatedGalleryChapter)
	if err != nil {
		return nil, err
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}
	r.hookExecutor.ExecutePostHooks(ctx, ret.ID, plugin.GalleryChapterUpdatePost, input, translator.getFields())
	return r.getGalleryChapter(ctx, ret.ID)
}

func (r *mutationResolver) GalleryChapterDestroy(ctx context.Context, id string) (bool, error) {
	chapterID, err := strconv.Atoi(id)
	if err != nil {
		return false, err
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.GalleryChapter

		chapter, err := qb.Find(ctx, chapterID)

		if err != nil {
			return err
		}

		if chapter == nil {
			return fmt.Errorf("Chapter with id %d not found", chapterID)
		}

		return gallery.DestroyChapter(ctx, chapter, qb)
	}); err != nil {
		return false, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, chapterID, plugin.GalleryChapterDestroyPost, id, nil)

	return true, nil
}

func (r *mutationResolver) changeChapter(ctx context.Context, changeType int, changedChapter models.GalleryChapter) (*models.GalleryChapter, error) {
	var galleryChapter *models.GalleryChapter

	// Start the transaction and save the gallery chapter
	var err = r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.GalleryChapter
		var err error

		switch changeType {
		case create:
			galleryChapter, err = qb.Create(ctx, changedChapter)
		case update:
			galleryChapter, err = qb.Update(ctx, changedChapter)
			if err != nil {
				return err
			}
		}
		return err
	})

	return galleryChapter, err
}
