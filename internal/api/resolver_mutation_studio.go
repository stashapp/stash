package api

import (
	"context"
	"database/sql"
	"github.com/stashapp/stash/internal/database"
	"github.com/stashapp/stash/internal/models"
	"github.com/stashapp/stash/internal/utils"
	"strconv"
	"time"
)

func (r *mutationResolver) StudioCreate(ctx context.Context, input models.StudioCreateInput) (*models.Studio, error) {
	// Process the base 64 encoded image string
	checksum, imageData, err := utils.ProcessBase64Image(input.Image)
	if err != nil {
		return nil, err
	}

	// Populate a new studio from the input
	currentTime := time.Now()
	newStudio := models.Studio{
		Image: imageData,
		Checksum: checksum,
		Name: sql.NullString{ String: input.Name, Valid: true },
		CreatedAt: models.SQLiteTimestamp{ Timestamp: currentTime },
		UpdatedAt: models.SQLiteTimestamp{ Timestamp: currentTime },
	}
	if input.URL != nil {
		newStudio.Url = sql.NullString{ String: *input.URL, Valid: true }
	}

	// Start the transaction and save the studio
	tx := database.DB.MustBeginTx(ctx, nil)
	qb := models.NewStudioQueryBuilder()
	studio, err := qb.Create(newStudio, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return studio, nil
}

func (r *mutationResolver) StudioUpdate(ctx context.Context, input models.StudioUpdateInput) (*models.Studio, error) {
	// Populate studio from the input
	studioID, _ := strconv.Atoi(input.ID)
	updatedStudio := models.Studio{
		ID: studioID,
		UpdatedAt: models.SQLiteTimestamp{ Timestamp: time.Now() },
	}
	if input.Image != nil {
		checksum, imageData, err := utils.ProcessBase64Image(*input.Image)
		if err != nil {
			return nil, err
		}
		updatedStudio.Image = imageData
		updatedStudio.Checksum = checksum
	}
	if input.Name != nil {
		updatedStudio.Name = sql.NullString{ String: *input.Name, Valid: true }
	}
	if input.URL != nil {
		updatedStudio.Url = sql.NullString{ String: *input.URL, Valid: true }
	}

	// Start the transaction and save the studio
	tx := database.DB.MustBeginTx(ctx, nil)
	qb := models.NewStudioQueryBuilder()
	studio, err := qb.Update(updatedStudio, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return studio, nil
}