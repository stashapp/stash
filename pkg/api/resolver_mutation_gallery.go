package api

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/database"
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

	// Start the transaction and save the performer
	tx := database.DB.MustBeginTx(ctx, nil)
	qb := models.NewGalleryQueryBuilder()
	jqb := models.NewJoinsQueryBuilder()
	gallery, err := qb.Create(newGallery, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// Save the performers
	var performerJoins []models.PerformersGalleries
	for _, pid := range input.PerformerIds {
		performerID, _ := strconv.Atoi(pid)
		performerJoin := models.PerformersGalleries{
			PerformerID: performerID,
			GalleryID:   gallery.ID,
		}
		performerJoins = append(performerJoins, performerJoin)
	}
	if err := jqb.UpdatePerformersGalleries(gallery.ID, performerJoins, tx); err != nil {
		return nil, err
	}

	// Save the tags
	var tagJoins []models.GalleriesTags
	for _, tid := range input.TagIds {
		tagID, _ := strconv.Atoi(tid)
		tagJoin := models.GalleriesTags{
			GalleryID: gallery.ID,
			TagID:     tagID,
		}
		tagJoins = append(tagJoins, tagJoin)
	}
	if err := jqb.UpdateGalleriesTags(gallery.ID, tagJoins, tx); err != nil {
		return nil, err
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return gallery, nil
}

func (r *mutationResolver) GalleryUpdate(ctx context.Context, input models.GalleryUpdateInput) (*models.Gallery, error) {
	// Start the transaction and save the gallery
	tx := database.DB.MustBeginTx(ctx, nil)

	ret, err := r.galleryUpdate(input, tx)

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

func (r *mutationResolver) GalleriesUpdate(ctx context.Context, input []*models.GalleryUpdateInput) ([]*models.Gallery, error) {
	// Start the transaction and save the gallery
	tx := database.DB.MustBeginTx(ctx, nil)

	var ret []*models.Gallery

	for _, gallery := range input {
		thisGallery, err := r.galleryUpdate(*gallery, tx)
		ret = append(ret, thisGallery)

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

func (r *mutationResolver) galleryUpdate(input models.GalleryUpdateInput, tx *sqlx.Tx) (*models.Gallery, error) {
	qb := models.NewGalleryQueryBuilder()
	// Populate gallery from the input
	galleryID, _ := strconv.Atoi(input.ID)
	originalGallery, err := qb.Find(galleryID, nil)
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
	if input.Details != nil {
		updatedGallery.Details = &sql.NullString{String: *input.Details, Valid: true}
	}
	if input.URL != nil {
		updatedGallery.URL = &sql.NullString{String: *input.URL, Valid: true}
	}
	if input.Date != nil {
		updatedGallery.Date = &models.SQLiteDate{String: *input.Date, Valid: true}
	}

	if input.Rating != nil {
		updatedGallery.Rating = &sql.NullInt64{Int64: int64(*input.Rating), Valid: true}
	} else {
		// rating must be nullable
		updatedGallery.Rating = &sql.NullInt64{Valid: false}
	}

	if input.StudioID != nil {
		studioID, _ := strconv.ParseInt(*input.StudioID, 10, 64)
		updatedGallery.StudioID = &sql.NullInt64{Int64: studioID, Valid: true}
	} else {
		// studio must be nullable
		updatedGallery.StudioID = &sql.NullInt64{Valid: false}
	}

	// gallery scene is set from the scene only

	jqb := models.NewJoinsQueryBuilder()
	gallery, err := qb.UpdatePartial(updatedGallery, tx)
	if err != nil {
		return nil, err
	}

	// Save the performers
	var performerJoins []models.PerformersGalleries
	for _, pid := range input.PerformerIds {
		performerID, _ := strconv.Atoi(pid)
		performerJoin := models.PerformersGalleries{
			PerformerID: performerID,
			GalleryID:   galleryID,
		}
		performerJoins = append(performerJoins, performerJoin)
	}
	if err := jqb.UpdatePerformersGalleries(galleryID, performerJoins, tx); err != nil {
		return nil, err
	}

	// Save the tags
	var tagJoins []models.GalleriesTags
	for _, tid := range input.TagIds {
		tagID, _ := strconv.Atoi(tid)
		tagJoin := models.GalleriesTags{
			GalleryID: galleryID,
			TagID:     tagID,
		}
		tagJoins = append(tagJoins, tagJoin)
	}
	if err := jqb.UpdateGalleriesTags(galleryID, tagJoins, tx); err != nil {
		return nil, err
	}

	return gallery, nil
}

func (r *mutationResolver) BulkGalleryUpdate(ctx context.Context, input models.BulkGalleryUpdateInput) ([]*models.Gallery, error) {
	// Populate gallery from the input
	updatedTime := time.Now()

	// Start the transaction and save the gallery marker
	tx := database.DB.MustBeginTx(ctx, nil)
	qb := models.NewGalleryQueryBuilder()
	jqb := models.NewJoinsQueryBuilder()

	updatedGallery := models.GalleryPartial{
		UpdatedAt: &models.SQLiteTimestamp{Timestamp: updatedTime},
	}
	if input.Details != nil {
		updatedGallery.Details = &sql.NullString{String: *input.Details, Valid: true}
	}
	if input.URL != nil {
		updatedGallery.URL = &sql.NullString{String: *input.URL, Valid: true}
	}
	if input.Date != nil {
		updatedGallery.Date = &models.SQLiteDate{String: *input.Date, Valid: true}
	}
	if input.Rating != nil {
		// a rating of 0 means unset the rating
		if *input.Rating == 0 {
			updatedGallery.Rating = &sql.NullInt64{Int64: 0, Valid: false}
		} else {
			updatedGallery.Rating = &sql.NullInt64{Int64: int64(*input.Rating), Valid: true}
		}
	}
	if input.StudioID != nil {
		// empty string means unset the studio
		if *input.StudioID == "" {
			updatedGallery.StudioID = &sql.NullInt64{Int64: 0, Valid: false}
		} else {
			studioID, _ := strconv.ParseInt(*input.StudioID, 10, 64)
			updatedGallery.StudioID = &sql.NullInt64{Int64: studioID, Valid: true}
		}
	}
	if input.SceneID != nil {
		// empty string means unset the studio
		if *input.SceneID == "" {
			updatedGallery.SceneID = &sql.NullInt64{Int64: 0, Valid: false}
		} else {
			sceneID, _ := strconv.ParseInt(*input.SceneID, 10, 64)
			updatedGallery.SceneID = &sql.NullInt64{Int64: sceneID, Valid: true}
		}
	}

	ret := []*models.Gallery{}

	for _, galleryIDStr := range input.Ids {
		galleryID, _ := strconv.Atoi(galleryIDStr)
		updatedGallery.ID = galleryID

		gallery, err := qb.UpdatePartial(updatedGallery, tx)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}

		ret = append(ret, gallery)

		// Save the performers
		if wasFieldIncluded(ctx, "performer_ids") {
			performerIDs, err := adjustGalleryPerformerIDs(tx, galleryID, *input.PerformerIds)
			if err != nil {
				_ = tx.Rollback()
				return nil, err
			}

			var performerJoins []models.PerformersGalleries
			for _, performerID := range performerIDs {
				performerJoin := models.PerformersGalleries{
					PerformerID: performerID,
					GalleryID:   galleryID,
				}
				performerJoins = append(performerJoins, performerJoin)
			}
			if err := jqb.UpdatePerformersGalleries(galleryID, performerJoins, tx); err != nil {
				_ = tx.Rollback()
				return nil, err
			}
		}

		// Save the tags
		if wasFieldIncluded(ctx, "tag_ids") {
			tagIDs, err := adjustGalleryTagIDs(tx, galleryID, *input.TagIds)
			if err != nil {
				_ = tx.Rollback()
				return nil, err
			}

			var tagJoins []models.GalleriesTags
			for _, tagID := range tagIDs {
				tagJoin := models.GalleriesTags{
					GalleryID: galleryID,
					TagID:     tagID,
				}
				tagJoins = append(tagJoins, tagJoin)
			}
			if err := jqb.UpdateGalleriesTags(galleryID, tagJoins, tx); err != nil {
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

func adjustGalleryPerformerIDs(tx *sqlx.Tx, galleryID int, ids models.BulkUpdateIds) ([]int, error) {
	var ret []int

	jqb := models.NewJoinsQueryBuilder()
	if ids.Mode == models.BulkUpdateIDModeAdd || ids.Mode == models.BulkUpdateIDModeRemove {
		// adding to the joins
		performerJoins, err := jqb.GetGalleryPerformers(galleryID, tx)

		if err != nil {
			return nil, err
		}

		for _, join := range performerJoins {
			ret = append(ret, join.PerformerID)
		}
	}

	return adjustIDs(ret, ids), nil
}

func adjustGalleryTagIDs(tx *sqlx.Tx, galleryID int, ids models.BulkUpdateIds) ([]int, error) {
	var ret []int

	jqb := models.NewJoinsQueryBuilder()
	if ids.Mode == models.BulkUpdateIDModeAdd || ids.Mode == models.BulkUpdateIDModeRemove {
		// adding to the joins
		tagJoins, err := jqb.GetGalleryTags(galleryID, tx)

		if err != nil {
			return nil, err
		}

		for _, join := range tagJoins {
			ret = append(ret, join.TagID)
		}
	}

	return adjustIDs(ret, ids), nil
}

func (r *mutationResolver) GalleryDestroy(ctx context.Context, input models.GalleryDestroyInput) (bool, error) {
	qb := models.NewGalleryQueryBuilder()
	iqb := models.NewImageQueryBuilder()
	tx := database.DB.MustBeginTx(ctx, nil)

	var galleries []*models.Gallery
	var imgsToPostProcess []*models.Image
	var imgsToDelete []*models.Image

	for _, id := range input.Ids {
		galleryID, _ := strconv.Atoi(id)

		gallery, err := qb.Find(galleryID, tx)
		if gallery != nil {
			galleries = append(galleries, gallery)
		}
		err = qb.Destroy(galleryID, tx)

		if err != nil {
			tx.Rollback()
			return false, err
		}

		// if this is a zip-based gallery, delete the images as well
		if gallery.Zip {
			imgs, err := iqb.FindByGalleryID(galleryID)
			if err != nil {
				tx.Rollback()
				return false, err
			}

			for _, img := range imgs {
				err = iqb.Destroy(img.ID, tx)
				if err != nil {
					tx.Rollback()
					return false, err
				}

				imgsToPostProcess = append(imgsToPostProcess, img)
			}
		} else if input.DeleteFile != nil && *input.DeleteFile {
			// Delete image if it is only attached to this gallery
			imgs, err := iqb.FindByGalleryID(galleryID)
			if err != nil {
				tx.Rollback()
				return false, err
			}

			for _, img := range imgs {
				imgGalleries, err := qb.FindByImageID(img.ID, tx)
				if err != nil {
					tx.Rollback()
					return false, err
				}

				if len(imgGalleries) == 0 {
					err = iqb.Destroy(img.ID, tx)
					if err != nil {
						tx.Rollback()
						return false, err
					}

					imgsToDelete = append(imgsToDelete, img)
					imgsToPostProcess = append(imgsToPostProcess, img)
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
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
	galleryID, _ := strconv.Atoi(input.GalleryID)
	qb := models.NewGalleryQueryBuilder()
	gallery, err := qb.Find(galleryID, nil)
	if err != nil {
		return false, err
	}

	if gallery == nil {
		return false, errors.New("gallery not found")
	}

	if gallery.Zip {
		return false, errors.New("cannot modify zip gallery images")
	}

	jqb := models.NewJoinsQueryBuilder()
	tx := database.DB.MustBeginTx(ctx, nil)

	for _, id := range input.ImageIds {
		imageID, _ := strconv.Atoi(id)
		_, err := jqb.AddImageGallery(imageID, galleryID, tx)
		if err != nil {
			tx.Rollback()
			return false, err
		}
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) RemoveGalleryImages(ctx context.Context, input models.GalleryRemoveInput) (bool, error) {
	galleryID, _ := strconv.Atoi(input.GalleryID)
	qb := models.NewGalleryQueryBuilder()
	gallery, err := qb.Find(galleryID, nil)
	if err != nil {
		return false, err
	}

	if gallery == nil {
		return false, errors.New("gallery not found")
	}

	if gallery.Zip {
		return false, errors.New("cannot modify zip gallery images")
	}

	jqb := models.NewJoinsQueryBuilder()
	tx := database.DB.MustBeginTx(ctx, nil)

	for _, id := range input.ImageIds {
		imageID, _ := strconv.Atoi(id)
		_, err := jqb.RemoveImageGallery(imageID, galleryID, tx)
		if err != nil {
			tx.Rollback()
			return false, err
		}
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return true, nil
}
