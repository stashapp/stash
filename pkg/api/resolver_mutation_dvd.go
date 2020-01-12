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

func (r *mutationResolver) DvdCreate(ctx context.Context, input models.DvdCreateInput) (*models.Dvd, error) {
	// generate checksum from dvd name rather than image
	checksum := utils.MD5FromString(input.Name)

	var frontimageData []byte
	var backimageData []byte
	var err error

	if input.Frontimage == nil {
		input.Frontimage = &models.DefaultDvdImage
	}
	if input.Backimage == nil {
		input.Backimage = &models.DefaultDvdImage
	}
	// Process the base 64 encoded image string
	_, frontimageData, err = utils.ProcessBase64Image(*input.Frontimage)
	if err != nil {
		return nil, err
	}
	// Process the base 64 encoded image string
	_, backimageData, err = utils.ProcessBase64Image(*input.Backimage)
	if err != nil {
		return nil, err
	}

	// Populate a new dvd from the input
	currentTime := time.Now()
	newDvd := models.Dvd{
		BackImage:  backimageData,
		FrontImage: frontimageData,
		Checksum:   checksum,
		Name:       sql.NullString{String: input.Name, Valid: true},
		CreatedAt:  models.SQLiteTimestamp{Timestamp: currentTime},
		UpdatedAt:  models.SQLiteTimestamp{Timestamp: currentTime},
	}

	if input.Aliases != nil {
		newDvd.Aliases = sql.NullString{String: *input.Aliases, Valid: true}
	}
	if input.Durationdvd != nil {
		newDvd.Durationdvd = sql.NullString{String: *input.Durationdvd, Valid: true}
	}

	if input.Year != nil {
		newDvd.Year = sql.NullString{String: *input.Year, Valid: true}
	}

	if input.Director != nil {
		newDvd.Director = sql.NullString{String: *input.Director, Valid: true}
	}

	if input.Synopsis != nil {
		newDvd.Synopsis = sql.NullString{String: *input.Synopsis, Valid: true}
	}

	if input.URL != nil {
		newDvd.URL = sql.NullString{String: *input.URL, Valid: true}
	}

	// Start the transaction and save the dvd
	tx := database.DB.MustBeginTx(ctx, nil)
	qb := models.NewDvdQueryBuilder()
	dvd, err := qb.Create(newDvd, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return dvd, nil
}

func (r *mutationResolver) DvdUpdate(ctx context.Context, input models.DvdUpdateInput) (*models.Dvd, error) {
	// Populate dvd from the input
	dvdID, _ := strconv.Atoi(input.ID)
	updatedDvd := models.Dvd{
		ID:        dvdID,
		UpdatedAt: models.SQLiteTimestamp{Timestamp: time.Now()},
	}
	if input.Frontimage != nil {
		_, frontimageData, err := utils.ProcessBase64Image(*input.Frontimage)
		if err != nil {
			return nil, err
		}
		updatedDvd.FrontImage = frontimageData
	}
	if input.Backimage != nil {
		_, backimageData, err := utils.ProcessBase64Image(*input.Backimage)
		if err != nil {
			return nil, err
		}
		updatedDvd.BackImage = backimageData
	}

	if input.Name != nil {
		// generate checksum from dvd name rather than image
		checksum := utils.MD5FromString(*input.Name)
		updatedDvd.Name = sql.NullString{String: *input.Name, Valid: true}
		updatedDvd.Checksum = checksum
	}

	if input.Aliases != nil {
		updatedDvd.Aliases = sql.NullString{String: *input.Aliases, Valid: true}
	}
	if input.Durationdvd != nil {
		updatedDvd.Durationdvd = sql.NullString{String: *input.Durationdvd, Valid: true}
	}

	if input.Year != nil {
		updatedDvd.Year = sql.NullString{String: *input.Year, Valid: true}
	}

	if input.Director != nil {
		updatedDvd.Director = sql.NullString{String: *input.Director, Valid: true}
	}

	if input.Synopsis != nil {
		updatedDvd.Synopsis = sql.NullString{String: *input.Synopsis, Valid: true}
	}

	if input.URL != nil {
		updatedDvd.URL = sql.NullString{String: *input.URL, Valid: true}
	}

	// Start the transaction and save the dvd
	tx := database.DB.MustBeginTx(ctx, nil)
	qb := models.NewDvdQueryBuilder()
	dvd, err := qb.Update(updatedDvd, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return dvd, nil
}

func (r *mutationResolver) DvdDestroy(ctx context.Context, input models.DvdDestroyInput) (bool, error) {
	qb := models.NewDvdQueryBuilder()
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
