package api

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/performer"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *mutationResolver) getPerformer(ctx context.Context, id int) (ret *models.Performer, err error) {
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		ret, err = repo.Performer().Find(id)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *mutationResolver) PerformerCreate(ctx context.Context, input PerformerCreateInput) (*models.Performer, error) {
	// generate checksum from performer name rather than image
	checksum := md5.FromString(input.Name)

	var imageData []byte
	var err error

	if input.Image != nil {
		imageData, err = utils.ProcessImageInput(ctx, *input.Image)
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
	newPerformer.Name = sql.NullString{String: input.Name, Valid: true}
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
	if input.Rating != nil {
		newPerformer.Rating = sql.NullInt64{Int64: int64(*input.Rating), Valid: true}
	} else {
		newPerformer.Rating = sql.NullInt64{Valid: false}
	}
	if input.Details != nil {
		newPerformer.Details = sql.NullString{String: *input.Details, Valid: true}
	}
	if input.DeathDate != nil {
		newPerformer.DeathDate = models.SQLiteDate{String: *input.DeathDate, Valid: true}
	}
	if input.HairColor != nil {
		newPerformer.HairColor = sql.NullString{String: *input.HairColor, Valid: true}
	}
	if input.Weight != nil {
		weight := int64(*input.Weight)
		newPerformer.Weight = sql.NullInt64{Int64: weight, Valid: true}
	}
	if input.IgnoreAutoTag != nil {
		newPerformer.IgnoreAutoTag = *input.IgnoreAutoTag
	}

	if err := performer.ValidateDeathDate(nil, input.Birthdate, input.DeathDate); err != nil {
		if err != nil {
			return nil, err
		}
	}

	// Start the transaction and save the performer
	var performer *models.Performer
	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Performer()

		performer, err = qb.Create(newPerformer)
		if err != nil {
			return err
		}

		if len(input.TagIds) > 0 {
			if err := r.updatePerformerTags(qb, performer.ID, input.TagIds); err != nil {
				return err
			}
		}

		// update image table
		if len(imageData) > 0 {
			if err := qb.UpdateImage(performer.ID, imageData); err != nil {
				return err
			}
		}

		// Save the stash_ids
		if input.StashIds != nil {
			stashIDJoins := models.StashIDsFromInput(input.StashIds)
			if err := qb.UpdateStashIDs(performer.ID, stashIDJoins); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, performer.ID, plugin.PerformerCreatePost, input, nil)
	return r.getPerformer(ctx, performer.ID)
}

func (r *mutationResolver) PerformerUpdate(ctx context.Context, input PerformerUpdateInput) (*models.Performer, error) {
	// Populate performer from the input
	performerID, _ := strconv.Atoi(input.ID)
	updatedPerformer := models.PerformerPartial{
		ID:        performerID,
		UpdatedAt: &models.SQLiteTimestamp{Timestamp: time.Now()},
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	var imageData []byte
	var err error
	imageIncluded := translator.hasField("image")
	if input.Image != nil {
		imageData, err = utils.ProcessImageInput(ctx, *input.Image)
		if err != nil {
			return nil, err
		}
	}

	if input.Name != nil {
		// generate checksum from performer name rather than image
		checksum := md5.FromString(*input.Name)

		updatedPerformer.Name = &sql.NullString{String: *input.Name, Valid: true}
		updatedPerformer.Checksum = &checksum
	}

	updatedPerformer.URL = translator.nullString(input.URL, "url")

	if translator.hasField("gender") {
		if input.Gender != nil {
			updatedPerformer.Gender = &sql.NullString{String: input.Gender.String(), Valid: true}
		} else {
			updatedPerformer.Gender = &sql.NullString{String: "", Valid: false}
		}
	}

	updatedPerformer.Birthdate = translator.sqliteDate(input.Birthdate, "birthdate")
	updatedPerformer.Country = translator.nullString(input.Country, "country")
	updatedPerformer.EyeColor = translator.nullString(input.EyeColor, "eye_color")
	updatedPerformer.Measurements = translator.nullString(input.Measurements, "measurements")
	updatedPerformer.Height = translator.nullString(input.Height, "height")
	updatedPerformer.Ethnicity = translator.nullString(input.Ethnicity, "ethnicity")
	updatedPerformer.FakeTits = translator.nullString(input.FakeTits, "fake_tits")
	updatedPerformer.CareerLength = translator.nullString(input.CareerLength, "career_length")
	updatedPerformer.Tattoos = translator.nullString(input.Tattoos, "tattoos")
	updatedPerformer.Piercings = translator.nullString(input.Piercings, "piercings")
	updatedPerformer.Aliases = translator.nullString(input.Aliases, "aliases")
	updatedPerformer.Twitter = translator.nullString(input.Twitter, "twitter")
	updatedPerformer.Instagram = translator.nullString(input.Instagram, "instagram")
	updatedPerformer.Favorite = translator.nullBool(input.Favorite, "favorite")
	updatedPerformer.Rating = translator.nullInt64(input.Rating, "rating")
	updatedPerformer.Details = translator.nullString(input.Details, "details")
	updatedPerformer.DeathDate = translator.sqliteDate(input.DeathDate, "death_date")
	updatedPerformer.HairColor = translator.nullString(input.HairColor, "hair_color")
	updatedPerformer.Weight = translator.nullInt64(input.Weight, "weight")
	updatedPerformer.IgnoreAutoTag = input.IgnoreAutoTag

	// Start the transaction and save the p
	var p *models.Performer
	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Performer()

		// need to get existing performer
		existing, err := qb.Find(updatedPerformer.ID)
		if err != nil {
			return err
		}

		if existing == nil {
			return fmt.Errorf("performer with id %d not found", updatedPerformer.ID)
		}

		if err := performer.ValidateDeathDate(existing, input.Birthdate, input.DeathDate); err != nil {
			if err != nil {
				return err
			}
		}

		p, err = qb.Update(updatedPerformer)
		if err != nil {
			return err
		}

		// Save the tags
		if translator.hasField("tag_ids") {
			if err := r.updatePerformerTags(qb, p.ID, input.TagIds); err != nil {
				return err
			}
		}

		// update image table
		if len(imageData) > 0 {
			if err := qb.UpdateImage(p.ID, imageData); err != nil {
				return err
			}
		} else if imageIncluded {
			// must be unsetting
			if err := qb.DestroyImage(p.ID); err != nil {
				return err
			}
		}

		// Save the stash_ids
		if translator.hasField("stash_ids") {
			stashIDJoins := models.StashIDsFromInput(input.StashIds)
			if err := qb.UpdateStashIDs(performerID, stashIDJoins); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, p.ID, plugin.PerformerUpdatePost, input, translator.getFields())
	return r.getPerformer(ctx, p.ID)
}

func (r *mutationResolver) updatePerformerTags(qb models.PerformerReaderWriter, performerID int, tagsIDs []string) error {
	ids, err := stringslice.StringSliceToIntSlice(tagsIDs)
	if err != nil {
		return err
	}
	return qb.UpdateTags(performerID, ids)
}

func (r *mutationResolver) BulkPerformerUpdate(ctx context.Context, input BulkPerformerUpdateInput) ([]*models.Performer, error) {
	performerIDs, err := stringslice.StringSliceToIntSlice(input.Ids)
	if err != nil {
		return nil, err
	}

	// Populate performer from the input
	updatedTime := time.Now()

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	updatedPerformer := models.PerformerPartial{
		UpdatedAt: &models.SQLiteTimestamp{Timestamp: updatedTime},
	}

	updatedPerformer.URL = translator.nullString(input.URL, "url")
	updatedPerformer.Birthdate = translator.sqliteDate(input.Birthdate, "birthdate")
	updatedPerformer.Ethnicity = translator.nullString(input.Ethnicity, "ethnicity")
	updatedPerformer.Country = translator.nullString(input.Country, "country")
	updatedPerformer.EyeColor = translator.nullString(input.EyeColor, "eye_color")
	updatedPerformer.Height = translator.nullString(input.Height, "height")
	updatedPerformer.Measurements = translator.nullString(input.Measurements, "measurements")
	updatedPerformer.FakeTits = translator.nullString(input.FakeTits, "fake_tits")
	updatedPerformer.CareerLength = translator.nullString(input.CareerLength, "career_length")
	updatedPerformer.Tattoos = translator.nullString(input.Tattoos, "tattoos")
	updatedPerformer.Piercings = translator.nullString(input.Piercings, "piercings")
	updatedPerformer.Aliases = translator.nullString(input.Aliases, "aliases")
	updatedPerformer.Twitter = translator.nullString(input.Twitter, "twitter")
	updatedPerformer.Instagram = translator.nullString(input.Instagram, "instagram")
	updatedPerformer.Favorite = translator.nullBool(input.Favorite, "favorite")
	updatedPerformer.Rating = translator.nullInt64(input.Rating, "rating")
	updatedPerformer.Details = translator.nullString(input.Details, "details")
	updatedPerformer.DeathDate = translator.sqliteDate(input.DeathDate, "death_date")
	updatedPerformer.HairColor = translator.nullString(input.HairColor, "hair_color")
	updatedPerformer.Weight = translator.nullInt64(input.Weight, "weight")
	updatedPerformer.IgnoreAutoTag = input.IgnoreAutoTag

	if translator.hasField("gender") {
		if input.Gender != nil {
			updatedPerformer.Gender = &sql.NullString{String: input.Gender.String(), Valid: true}
		} else {
			updatedPerformer.Gender = &sql.NullString{String: "", Valid: false}
		}
	}

	ret := []*models.Performer{}

	// Start the transaction and save the scene marker
	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Performer()

		for _, performerID := range performerIDs {
			updatedPerformer.ID = performerID

			// need to get existing performer
			existing, err := qb.Find(performerID)
			if err != nil {
				return err
			}

			if existing == nil {
				return fmt.Errorf("performer with id %d not found", performerID)
			}

			if err := performer.ValidateDeathDate(existing, input.Birthdate, input.DeathDate); err != nil {
				return err
			}

			performer, err := qb.Update(updatedPerformer)
			if err != nil {
				return err
			}

			ret = append(ret, performer)

			// Save the tags
			if translator.hasField("tag_ids") {
				tagIDs, err := adjustTagIDs(qb, performerID, *input.TagIds)
				if err != nil {
					return err
				}

				if err := qb.UpdateTags(performerID, tagIDs); err != nil {
					return err
				}
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// execute post hooks outside of txn
	var newRet []*models.Performer
	for _, performer := range ret {
		r.hookExecutor.ExecutePostHooks(ctx, performer.ID, plugin.ImageUpdatePost, input, translator.getFields())

		performer, err = r.getPerformer(ctx, performer.ID)
		if err != nil {
			return nil, err
		}

		newRet = append(newRet, performer)
	}

	return newRet, nil
}

func (r *mutationResolver) PerformerDestroy(ctx context.Context, input PerformerDestroyInput) (bool, error) {
	id, err := strconv.Atoi(input.ID)
	if err != nil {
		return false, err
	}

	if err := r.withTxn(ctx, func(repo models.Repository) error {
		return repo.Performer().Destroy(id)
	}); err != nil {
		return false, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, id, plugin.PerformerDestroyPost, input, nil)

	return true, nil
}

func (r *mutationResolver) PerformersDestroy(ctx context.Context, performerIDs []string) (bool, error) {
	ids, err := stringslice.StringSliceToIntSlice(performerIDs)
	if err != nil {
		return false, err
	}

	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Performer()
		for _, id := range ids {
			if err := qb.Destroy(id); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return false, err
	}

	for _, id := range ids {
		r.hookExecutor.ExecutePostHooks(ctx, id, plugin.PerformerDestroyPost, performerIDs, nil)
	}

	return true, nil
}
