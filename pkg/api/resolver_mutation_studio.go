package api

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *mutationResolver) StudioCreate(ctx context.Context, input models.StudioCreateInput) (*models.Studio, error) {
	// generate checksum from studio name rather than image
	checksum := utils.MD5FromString(input.Name)

	var imageData []byte
	var err error

	// Process the base 64 encoded image string
	if input.Image != nil {
		_, imageData, err = utils.ProcessBase64Image(*input.Image)
		if err != nil {
			return nil, err
		}
	}

	// Populate a new studio from the input
	currentTime := time.Now()
	newStudio := models.Studio{
		Checksum:  checksum,
		Name:      sql.NullString{String: input.Name, Valid: true},
		CreatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
		UpdatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
	}
	if input.URL != nil {
		newStudio.URL = sql.NullString{String: *input.URL, Valid: true}
	}
	if input.ParentID != nil {
		parentID, _ := strconv.ParseInt(*input.ParentID, 10, 64)
		newStudio.ParentID = sql.NullInt64{Int64: parentID, Valid: true}
	}

	// Start the transaction and save the studio
	var studio *models.Studio
	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Studio()

		var err error
		studio, err = qb.Create(newStudio)
		if err != nil {
			return err
		}

		// update image table
		if len(imageData) > 0 {
			if err := qb.UpdateImage(studio.ID, imageData); err != nil {
				return err
			}
		}

		// Save the stash_ids
		if input.StashIds != nil {
			stashIDJoins := models.StashIDsFromInput(input.StashIds)
			if err := qb.UpdateStashIDs(studio.ID, stashIDJoins); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return studio, nil
}

func (r *mutationResolver) StudioUpdate(ctx context.Context, input models.StudioUpdateInput) (*models.Studio, error) {
	// Populate studio from the input
	studioID, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, err
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	updatedStudio := models.StudioPartial{
		ID:        studioID,
		UpdatedAt: &models.SQLiteTimestamp{Timestamp: time.Now()},
	}

	var imageData []byte
	imageIncluded := translator.hasField("image")
	if input.Image != nil {
		var err error
		_, imageData, err = utils.ProcessBase64Image(*input.Image)
		if err != nil {
			return nil, err
		}
	}
	if input.Name != nil {
		// generate checksum from studio name rather than image
		checksum := utils.MD5FromString(*input.Name)
		updatedStudio.Name = &sql.NullString{String: *input.Name, Valid: true}
		updatedStudio.Checksum = &checksum
	}

	updatedStudio.URL = translator.nullString(input.URL, "url")
	updatedStudio.ParentID = translator.nullInt64FromString(input.ParentID, "parent_id")

	// Start the transaction and save the studio
	var studio *models.Studio
	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Studio()

		if err := manager.ValidateModifyStudio(updatedStudio, qb); err != nil {
			return err
		}

		var err error
		studio, err = qb.Update(updatedStudio)
		if err != nil {
			return err
		}

		// update image table
		if len(imageData) > 0 {
			if err := qb.UpdateImage(studio.ID, imageData); err != nil {
				return err
			}
		} else if imageIncluded {
			// must be unsetting
			if err := qb.DestroyImage(studio.ID); err != nil {
				return err
			}
		}

		// Save the stash_ids
		if translator.hasField("stash_ids") {
			stashIDJoins := models.StashIDsFromInput(input.StashIds)
			if err := qb.UpdateStashIDs(studioID, stashIDJoins); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return studio, nil
}

func (r *mutationResolver) StudioDestroy(ctx context.Context, input models.StudioDestroyInput) (bool, error) {
	id, err := strconv.Atoi(input.ID)
	if err != nil {
		return false, err
	}

	if err := r.withTxn(ctx, func(repo models.Repository) error {
		return repo.Studio().Destroy(id)
	}); err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) StudiosDestroy(ctx context.Context, ids []string) (bool, error) {
	qb := models.NewStudioQueryBuilder()
	tx := database.DB.MustBeginTx(ctx, nil)
	for _, id := range ids {
		if err := qb.Destroy(id, tx); err != nil {
			_ = tx.Rollback()
			return false, err
		}
	}
	if err := tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}
