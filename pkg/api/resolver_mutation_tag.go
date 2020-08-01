package api

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *mutationResolver) TagCreate(ctx context.Context, input models.TagCreateInput) (*models.Tag, error) {
	// Populate a new tag from the input
	currentTime := time.Now()
	newTag := models.Tag{
		Name:      input.Name,
		CreatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
		UpdatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
	}

	var imageData []byte
	var err error

	if input.Image != nil {
		_, imageData, err = utils.ProcessBase64Image(*input.Image)

		if err != nil {
			return nil, err
		}
	}

	// Start the transaction and save the tag
	tx := database.DB.MustBeginTx(ctx, nil)
	qb := models.NewTagQueryBuilder()

	// ensure name is unique
	if err := manager.EnsureTagNameUnique(newTag.Name, tx); err != nil {
		tx.Rollback()
		return nil, err
	}

	tag, err := qb.Create(newTag, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// update image table
	if len(imageData) > 0 {
		if err := qb.UpdateTagImage(tag.ID, imageData, tx); err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return tag, nil
}

func (r *mutationResolver) TagUpdate(ctx context.Context, input models.TagUpdateInput) (*models.Tag, error) {
	// Populate tag from the input
	tagID, _ := strconv.Atoi(input.ID)
	updatedTag := models.Tag{
		ID:        tagID,
		Name:      input.Name,
		UpdatedAt: models.SQLiteTimestamp{Timestamp: time.Now()},
	}

	var imageData []byte
	var err error

	if input.Image != nil {
		_, imageData, err = utils.ProcessBase64Image(*input.Image)

		if err != nil {
			return nil, err
		}
	}

	// Start the transaction and save the tag
	tx := database.DB.MustBeginTx(ctx, nil)
	qb := models.NewTagQueryBuilder()

	// ensure name is unique
	existing, err := qb.Find(tagID, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if existing == nil {
		tx.Rollback()
		return nil, fmt.Errorf("Tag with ID %d not found", tagID)
	}

	if existing.Name != updatedTag.Name {
		if err := manager.EnsureTagNameUnique(updatedTag.Name, tx); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tag, err := qb.Update(updatedTag, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// update image table
	if len(imageData) > 0 {
		if err := qb.UpdateTagImage(tag.ID, imageData, tx); err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return tag, nil
}

func (r *mutationResolver) TagDestroy(ctx context.Context, input models.TagDestroyInput) (bool, error) {
	qb := models.NewTagQueryBuilder()
	tx := database.DB.MustBeginTx(ctx, nil)
	if err := qb.Destroy(input.ID, tx); err != nil {
		_ = tx.Rollback()
		return false, err
	}
	if err := tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}
