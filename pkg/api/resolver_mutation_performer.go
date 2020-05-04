package api

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *mutationResolver) PerformerCreate(ctx context.Context, input models.PerformerCreateInput) (*models.Performer, error) {
	// generate checksum from performer name rather than image
	checksum := utils.MD5FromString(*input.Name)

	var imageData []byte
	var err error

	if input.Image == nil {
		gender := ""
		if input.Gender != nil {
			gender = input.Gender.String()
		}
		imageData, err = getRandomPerformerImage(gender)
	} else {
		_, imageData, err = utils.ProcessBase64Image(*input.Image)
	}

	if err != nil {
		return nil, err
	}

	// Populate a new performer from the input
	currentTime := time.Now()
	newPerformer := models.Performer{
		Checksum:  checksum,
		CreatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
		UpdatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
	}
	if input.Name != nil {
		newPerformer.Name = sql.NullString{String: *input.Name, Valid: true}
	}
	if input.URL != nil {
		newPerformer.URL = sql.NullString{String: *input.URL, Valid: true}
	}
	if input.Gender != nil {
		newPerformer.Gender = sql.NullString{String: input.Gender.String(), Valid: true}
	}
	if input.Birthdate != nil {
		newPerformer.Birthdate = models.SQLiteDate{String: *input.Birthdate, Valid: true}
	}
	if input.Ethnicity != nil {
		newPerformer.Ethnicity = sql.NullString{String: *input.Ethnicity, Valid: true}
	}
	if input.Country != nil {
		newPerformer.Country = sql.NullString{String: *input.Country, Valid: true}
	}
	if input.EyeColor != nil {
		newPerformer.EyeColor = sql.NullString{String: *input.EyeColor, Valid: true}
	}
	if input.Height != nil {
		newPerformer.Height = sql.NullString{String: *input.Height, Valid: true}
	}
	if input.Measurements != nil {
		newPerformer.Measurements = sql.NullString{String: *input.Measurements, Valid: true}
	}
	if input.FakeTits != nil {
		newPerformer.FakeTits = sql.NullString{String: *input.FakeTits, Valid: true}
	}
	if input.CareerLength != nil {
		newPerformer.CareerLength = sql.NullString{String: *input.CareerLength, Valid: true}
	}
	if input.Tattoos != nil {
		newPerformer.Tattoos = sql.NullString{String: *input.Tattoos, Valid: true}
	}
	if input.Piercings != nil {
		newPerformer.Piercings = sql.NullString{String: *input.Piercings, Valid: true}
	}
	if input.Aliases != nil {
		newPerformer.Aliases = sql.NullString{String: *input.Aliases, Valid: true}
	}
	if input.Twitter != nil {
		newPerformer.Twitter = sql.NullString{String: *input.Twitter, Valid: true}
	}
	if input.Instagram != nil {
		newPerformer.Instagram = sql.NullString{String: *input.Instagram, Valid: true}
	}
	if input.Favorite != nil {
		newPerformer.Favorite = sql.NullBool{Bool: *input.Favorite, Valid: true}
	} else {
		newPerformer.Favorite = sql.NullBool{Bool: false, Valid: true}
	}

	// Start the transaction and save the performer
	tx := database.DB.MustBeginTx(ctx, nil)
	qb := models.NewPerformerQueryBuilder()
	performer, err := qb.Create(newPerformer, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// update image table
	if len(imageData) > 0 {
		if err := qb.UpdatePerformerImage(performer.ID, imageData, tx); err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return performer, nil
}

func (r *mutationResolver) PerformerUpdate(ctx context.Context, input models.PerformerUpdateInput) (*models.Performer, error) {
	// Populate performer from the input
	performerID, _ := strconv.Atoi(input.ID)
	updatedPerformer := models.Performer{
		ID:        performerID,
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
	if input.Name != nil {
		// generate checksum from performer name rather than image
		checksum := utils.MD5FromString(*input.Name)

		updatedPerformer.Name = sql.NullString{String: *input.Name, Valid: true}
		updatedPerformer.Checksum = checksum
	}
	if input.URL != nil {
		updatedPerformer.URL = sql.NullString{String: *input.URL, Valid: true}
	}
	if input.Gender != nil {
		updatedPerformer.Gender = sql.NullString{String: input.Gender.String(), Valid: true}
	}
	if input.Birthdate != nil {
		updatedPerformer.Birthdate = models.SQLiteDate{String: *input.Birthdate, Valid: true}
	}
	if input.Ethnicity != nil {
		updatedPerformer.Ethnicity = sql.NullString{String: *input.Ethnicity, Valid: true}
	}
	if input.Country != nil {
		updatedPerformer.Country = sql.NullString{String: *input.Country, Valid: true}
	}
	if input.EyeColor != nil {
		updatedPerformer.EyeColor = sql.NullString{String: *input.EyeColor, Valid: true}
	}
	if input.Height != nil {
		updatedPerformer.Height = sql.NullString{String: *input.Height, Valid: true}
	}
	if input.Measurements != nil {
		updatedPerformer.Measurements = sql.NullString{String: *input.Measurements, Valid: true}
	}
	if input.FakeTits != nil {
		updatedPerformer.FakeTits = sql.NullString{String: *input.FakeTits, Valid: true}
	}
	if input.CareerLength != nil {
		updatedPerformer.CareerLength = sql.NullString{String: *input.CareerLength, Valid: true}
	}
	if input.Tattoos != nil {
		updatedPerformer.Tattoos = sql.NullString{String: *input.Tattoos, Valid: true}
	}
	if input.Piercings != nil {
		updatedPerformer.Piercings = sql.NullString{String: *input.Piercings, Valid: true}
	}
	if input.Aliases != nil {
		updatedPerformer.Aliases = sql.NullString{String: *input.Aliases, Valid: true}
	}
	if input.Twitter != nil {
		updatedPerformer.Twitter = sql.NullString{String: *input.Twitter, Valid: true}
	}
	if input.Instagram != nil {
		updatedPerformer.Instagram = sql.NullString{String: *input.Instagram, Valid: true}
	}
	if input.Favorite != nil {
		updatedPerformer.Favorite = sql.NullBool{Bool: *input.Favorite, Valid: true}
	} else {
		updatedPerformer.Favorite = sql.NullBool{Bool: false, Valid: true}
	}

	// Start the transaction and save the performer
	tx := database.DB.MustBeginTx(ctx, nil)
	qb := models.NewPerformerQueryBuilder()
	performer, err := qb.Update(updatedPerformer, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// update image table
	if len(imageData) > 0 {
		if err := qb.UpdatePerformerImage(performer.ID, imageData, tx); err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return performer, nil
}

func (r *mutationResolver) PerformerDestroy(ctx context.Context, input models.PerformerDestroyInput) (bool, error) {
	qb := models.NewPerformerQueryBuilder()
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
