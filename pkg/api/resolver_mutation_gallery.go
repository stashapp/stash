package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *mutationResolver) GalleryCreate(ctx context.Context, input models.GalleryCreateInput) (*models.Gallery, error) {
	// name must be provided
	if input.Title == "" {
		return nil, errors.New("title must not be empty")
	}

	// for manually created galleries, generate checksum from title
	checksum := utils.MD5FromString(input.Title)

	// Populate a new performer from the input
	currentTime := time.Now()
	newGallery := models.Gallery{
		Title: sql.NullString{
			String: input.Title,
			Valid:  true,
		},
		Checksum:  checksum,
		CreatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
		UpdatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
	}
	if input.URL != nil {
		newGallery.URL = sql.NullString{String: *input.URL, Valid: true}
	}
	if input.Details != nil {
		newGallery.Details = sql.NullString{String: *input.Details, Valid: true}
	}
	if input.URL != nil {
		newGallery.URL = sql.NullString{String: *input.URL, Valid: true}
	}
	if input.Date != nil {
		newGallery.Date = models.SQLiteDate{String: *input.Date, Valid: true}
	}
	if input.Rating != nil {
		newGallery.Rating = sql.NullInt64{Int64: int64(*input.Rating), Valid: true}
	} else {
		// rating must be nullable
		newGallery.Rating = sql.NullInt64{Valid: false}
	}

	if input.StudioID != nil {
		studioID, _ := strconv.ParseInt(*input.StudioID, 10, 64)
		newGallery.StudioID = sql.NullInt64{Int64: studioID, Valid: true}
	} else {
		// studio must be nullable
		newGallery.StudioID = sql.NullInt64{Valid: false}
	}

	if input.SceneID != nil {
		sceneID, _ := strconv.ParseInt(*input.SceneID, 10, 64)
		newGallery.SceneID = sql.NullInt64{Int64: sceneID, Valid: true}
	} else {
		// studio must be nullable
		newGallery.SceneID = sql.NullInt64{Valid: false}
	}

	// Start the transaction and save the gallery
	var gallery *models.Gallery
	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Gallery()
		var err error
		gallery, err = qb.Create(newGallery)
		if err != nil {
			return err
		}

		// Save the performers
		if err := r.updateGalleryPerformers(qb, gallery.ID, input.PerformerIds); err != nil {
			return err
		}

		// Save the tags
		if err := r.updateGalleryTags(qb, gallery.ID, input.TagIds); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return gallery, nil
}

func (r *mutationResolver) updateGalleryPerformers(qb models.GalleryReaderWriter, galleryID int, performerIDs []string) error {
	ids, err := utils.StringSliceToIntSlice(performerIDs)
	if err != nil {
		return err
	}
	return qb.UpdatePerformers(galleryID, ids)
}

func (r *mutationResolver) updateGalleryTags(qb models.GalleryReaderWriter, galleryID int, tagIDs []string) error {
	ids, err := utils.StringSliceToIntSlice(tagIDs)
	if err != nil {
		return err
	}
	return qb.UpdateTags(galleryID, ids)
}

func (r *mutationResolver) GalleryUpdate(ctx context.Context, input models.GalleryUpdateInput) (ret *models.Gallery, err error) {
	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Start the transaction and save the gallery
	if err := r.withTxn(ctx, func(repo models.Repository) error {
		ret, err = r.galleryUpdate(input, translator, repo)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *mutationResolver) GalleriesUpdate(ctx context.Context, input []*models.GalleryUpdateInput) (ret []*models.Gallery, err error) {
	inputMaps := getUpdateInputMaps(ctx)

	// Start the transaction and save the gallery
	if err := r.withTxn(ctx, func(repo models.Repository) error {
		for i, gallery := range input {
			translator := changesetTranslator{
				inputMap: inputMaps[i],
			}

			thisGallery, err := r.galleryUpdate(*gallery, translator, repo)
			if err != nil {
				return err
			}

			ret = append(ret, thisGallery)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *mutationResolver) galleryUpdate(input models.GalleryUpdateInput, translator changesetTranslator, repo models.Repository) (*models.Gallery, error) {
	qb := repo.Gallery()

	// Populate gallery from the input
	galleryID, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, err
	}

	originalGallery, err := qb.Find(galleryID)
	if err != nil {
		return nil, err
	}

	if originalGallery == nil {
		return nil, errors.New("not found")
	}

	updatedTime := time.Now()
	updatedGallery := models.GalleryPartial{
		ID:        galleryID,
		UpdatedAt: &models.SQLiteTimestamp{Timestamp: updatedTime},
	}

	if input.Title != nil {
		// ensure title is not empty
		if *input.Title == "" {
			return nil, errors.New("title must not be empty")
		}

		// if gallery is not zip-based, then generate the checksum from the title
		if !originalGallery.Path.Valid {
			checksum := utils.MD5FromString(*input.Title)
			updatedGallery.Checksum = &checksum
		}

		updatedGallery.Title = &sql.NullString{String: *input.Title, Valid: true}
	}

	updatedGallery.Details = translator.nullString(input.Details, "details")
	updatedGallery.URL = translator.nullString(input.URL, "url")
	updatedGallery.Date = translator.sqliteDate(input.Date, "date")
	updatedGallery.Rating = translator.nullInt64(input.Rating, "rating")
	updatedGallery.StudioID = translator.nullInt64FromString(input.StudioID, "studio_id")
	updatedGallery.Organized = input.Organized

	// gallery scene is set from the scene only

	gallery, err := qb.UpdatePartial(updatedGallery)
	if err != nil {
		return nil, err
	}

	// Save the performers
	if translator.hasField("performer_ids") {
		if err := r.updateGalleryPerformers(qb, galleryID, input.PerformerIds); err != nil {
			return nil, err
		}
	}

	// Save the tags
	if translator.hasField("tag_ids") {
		if err := r.updateGalleryTags(qb, galleryID, input.TagIds); err != nil {
			return nil, err
		}
	}

	return gallery, nil
}

func (r *mutationResolver) BulkGalleryUpdate(ctx context.Context, input models.BulkGalleryUpdateInput) ([]*models.Gallery, error) {
	// Populate gallery from the input
	updatedTime := time.Now()

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	updatedGallery := models.GalleryPartial{
		UpdatedAt: &models.SQLiteTimestamp{Timestamp: updatedTime},
	}

	updatedGallery.Details = translator.nullString(input.Details, "details")
	updatedGallery.URL = translator.nullString(input.URL, "url")
	updatedGallery.Date = translator.sqliteDate(input.Date, "date")
	updatedGallery.Rating = translator.nullInt64(input.Rating, "rating")
	updatedGallery.StudioID = translator.nullInt64FromString(input.StudioID, "studio_id")
	updatedGallery.SceneID = translator.nullInt64FromString(input.SceneID, "scene_id")
	updatedGallery.Organized = input.Organized

	ret := []*models.Gallery{}

	// Start the transaction and save the galleries
	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Gallery()

		for _, galleryIDStr := range input.Ids {
			galleryID, _ := strconv.Atoi(galleryIDStr)
			updatedGallery.ID = galleryID

			gallery, err := qb.UpdatePartial(updatedGallery)
			if err != nil {
				return err
			}

			ret = append(ret, gallery)

			// Save the performers
			if translator.hasField("performer_ids") {
				performerIDs, err := adjustGalleryPerformerIDs(qb, galleryID, *input.PerformerIds)
				if err != nil {
					return err
				}

				if err := qb.UpdatePerformers(galleryID, performerIDs); err != nil {
					return err
				}
			}

			// Save the tags
			if translator.hasField("tag_ids") {
				tagIDs, err := adjustGalleryTagIDs(qb, galleryID, *input.TagIds)
				if err != nil {
					return err
				}

				if err := qb.UpdateTags(galleryID, tagIDs); err != nil {
					return err
				}
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func adjustGalleryPerformerIDs(qb models.GalleryReader, galleryID int, ids models.BulkUpdateIds) (ret []int, err error) {
	ret, err = qb.GetPerformerIDs(galleryID)
	if err != nil {
		return nil, err
	}

	return adjustIDs(ret, ids), nil
}

func adjustGalleryTagIDs(qb models.GalleryReader, galleryID int, ids models.BulkUpdateIds) (ret []int, err error) {
	ret, err = qb.GetTagIDs(galleryID)
	if err != nil {
		return nil, err
	}

	return adjustIDs(ret, ids), nil
}

func (r *mutationResolver) GalleryDestroy(ctx context.Context, input models.GalleryDestroyInput) (bool, error) {
	galleryIDs, err := utils.StringSliceToIntSlice(input.Ids)
	if err != nil {
		return false, err
	}

	var galleries []*models.Gallery
	var imgsToPostProcess []*models.Image
	var imgsToDelete []*models.Image

	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Gallery()
		iqb := repo.Image()

		for _, id := range galleryIDs {
			gallery, err := qb.Find(id)
			if err != nil {
				return err
			}

			if gallery == nil {
				return fmt.Errorf("gallery with id %d not found", id)
			}

			galleries = append(galleries, gallery)

			// if this is a zip-based gallery, delete the images as well first
			if gallery.Zip {
				imgs, err := iqb.FindByGalleryID(id)
				if err != nil {
					return err
				}

				for _, img := range imgs {
					if err := iqb.Destroy(img.ID); err != nil {
						return err
					}

					imgsToPostProcess = append(imgsToPostProcess, img)
				}
			} else if input.DeleteFile != nil && *input.DeleteFile {
				// Delete image if it is only attached to this gallery
				imgs, err := iqb.FindByGalleryID(id)
				if err != nil {
					return err
				}

				for _, img := range imgs {
					imgGalleries, err := qb.FindByImageID(img.ID)
					if err != nil {
						return err
					}

					if len(imgGalleries) == 0 {
						if err := iqb.Destroy(img.ID); err != nil {
							return err
						}

						imgsToDelete = append(imgsToDelete, img)
						imgsToPostProcess = append(imgsToPostProcess, img)
					}
				}
			}

			if err := qb.Destroy(id); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return false, err
	}

	// if delete file is true, then delete the file as well
	// if it fails, just log a message
	if input.DeleteFile != nil && *input.DeleteFile {
		for _, gallery := range galleries {
			manager.DeleteGalleryFile(gallery)
		}

		for _, img := range imgsToDelete {
			manager.DeleteImageFile(img)
		}
	}

	// if delete generated is true, then delete the generated files
	// for the gallery
	if input.DeleteGenerated != nil && *input.DeleteGenerated {
		for _, img := range imgsToPostProcess {
			manager.DeleteGeneratedImageFiles(img)
		}
	}

	return true, nil
}

func (r *mutationResolver) AddGalleryImages(ctx context.Context, input models.GalleryAddInput) (bool, error) {
	galleryID, err := strconv.Atoi(input.GalleryID)
	if err != nil {
		return false, err
	}

	imageIDs, err := utils.StringSliceToIntSlice(input.ImageIds)
	if err != nil {
		return false, err
	}

	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Gallery()
		gallery, err := qb.Find(galleryID)
		if err != nil {
			return err
		}

		if gallery == nil {
			return errors.New("gallery not found")
		}

		if gallery.Zip {
			return errors.New("cannot modify zip gallery images")
		}

		newIDs, err := qb.GetImageIDs(galleryID)
		if err != nil {
			return err
		}

		newIDs = utils.IntAppendUniques(newIDs, imageIDs)
		return qb.UpdateImages(galleryID, newIDs)
	}); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) RemoveGalleryImages(ctx context.Context, input models.GalleryRemoveInput) (bool, error) {
	galleryID, err := strconv.Atoi(input.GalleryID)
	if err != nil {
		return false, err
	}

	imageIDs, err := utils.StringSliceToIntSlice(input.ImageIds)
	if err != nil {
		return false, err
	}

	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Gallery()
		gallery, err := qb.Find(galleryID)
		if err != nil {
			return err
		}

		if gallery == nil {
			return errors.New("gallery not found")
		}

		if gallery.Zip {
			return errors.New("cannot modify zip gallery images")
		}

		newIDs, err := qb.GetImageIDs(galleryID)
		if err != nil {
			return err
		}

		newIDs = utils.IntExclude(newIDs, imageIDs)
		return qb.UpdateImages(galleryID, newIDs)
	}); err != nil {
		return false, err
	}

	return true, nil
}
