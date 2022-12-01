package api

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/performer"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *mutationResolver) getPerformer(ctx context.Context, id int) (ret *models.Performer, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Performer.Find(ctx, id)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func stashIDPtrSliceToSlice(v []*models.StashID) []models.StashID {
	ret := make([]models.StashID, len(v))
	for i, vv := range v {
		c := vv
		ret[i] = *c
	}

	return ret
}

func (r *mutationResolver) PerformerCreate(ctx context.Context, input PerformerCreateInput) (*models.Performer, error) {
	var imageData []byte
	var err error

	if input.Image != nil {
		imageData, err = utils.ProcessImageInput(ctx, *input.Image)
	}

	if err != nil {
		return nil, err
	}

	tagIDs, err := stringslice.StringSliceToIntSlice(input.TagIds)
	if err != nil {
		return nil, fmt.Errorf("converting tag ids: %w", err)
	}

	// Populate a new performer from the input
	currentTime := time.Now()
	newPerformer := models.Performer{
		Name:      input.Name,
		TagIDs:    models.NewRelatedIDs(tagIDs),
		StashIDs:  models.NewRelatedStashIDs(stashIDPtrSliceToSlice(input.StashIds)),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
	if input.Disambiguation != nil {
		newPerformer.Disambiguation = *input.Disambiguation
	}
	if input.URL != nil {
		newPerformer.URL = *input.URL
	}
	if input.Gender != nil {
		newPerformer.Gender = *input.Gender
	}
	if input.Birthdate != nil {
		d := models.NewDate(*input.Birthdate)
		newPerformer.Birthdate = &d
	}
	if input.Ethnicity != nil {
		newPerformer.Ethnicity = *input.Ethnicity
	}
	if input.Country != nil {
		newPerformer.Country = *input.Country
	}
	if input.EyeColor != nil {
		newPerformer.EyeColor = *input.EyeColor
	}
	// prefer height_cm over height
	if input.HeightCm != nil {
		newPerformer.Height = input.HeightCm
	} else if input.Height != nil {
		h, err := strconv.Atoi(*input.Height)
		if err != nil {
			return nil, fmt.Errorf("invalid height: %s", *input.Height)
		}
		newPerformer.Height = &h
	}
	if input.Measurements != nil {
		newPerformer.Measurements = *input.Measurements
	}
	if input.FakeTits != nil {
		newPerformer.FakeTits = *input.FakeTits
	}
	if input.CareerLength != nil {
		newPerformer.CareerLength = *input.CareerLength
	}
	if input.Tattoos != nil {
		newPerformer.Tattoos = *input.Tattoos
	}
	if input.Piercings != nil {
		newPerformer.Piercings = *input.Piercings
	}
	if input.AliasList != nil {
		newPerformer.Aliases = models.NewRelatedStrings(input.AliasList)
	} else if input.Aliases != nil {
		newPerformer.Aliases = models.NewRelatedStrings(stringslice.FromString(*input.Aliases, ","))
	}
	if input.Twitter != nil {
		newPerformer.Twitter = *input.Twitter
	}
	if input.Instagram != nil {
		newPerformer.Instagram = *input.Instagram
	}
	if input.Favorite != nil {
		newPerformer.Favorite = *input.Favorite
	}
	if input.Rating100 != nil {
		newPerformer.Rating = input.Rating100
	} else if input.Rating != nil {
		rating := models.Rating5To100(*input.Rating)
		newPerformer.Rating = &rating
	}
	if input.Details != nil {
		newPerformer.Details = *input.Details
	}
	if input.DeathDate != nil {
		d := models.NewDate(*input.DeathDate)
		newPerformer.DeathDate = &d
	}
	if input.HairColor != nil {
		newPerformer.HairColor = *input.HairColor
	}
	if input.Weight != nil {
		newPerformer.Weight = input.Weight
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
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Performer

		err = qb.Create(ctx, &newPerformer)
		if err != nil {
			return err
		}

		// update image table
		if len(imageData) > 0 {
			if err := qb.UpdateImage(ctx, newPerformer.ID, imageData); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, newPerformer.ID, plugin.PerformerCreatePost, input, nil)
	return r.getPerformer(ctx, newPerformer.ID)
}

func (r *mutationResolver) PerformerUpdate(ctx context.Context, input PerformerUpdateInput) (*models.Performer, error) {
	// Populate performer from the input
	performerID, _ := strconv.Atoi(input.ID)
	updatedPerformer := models.NewPerformerPartial()

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

	updatedPerformer.Name = translator.optionalString(input.Name, "name")
	updatedPerformer.Disambiguation = translator.optionalString(input.Disambiguation, "disambiguation")
	updatedPerformer.URL = translator.optionalString(input.URL, "url")

	if translator.hasField("gender") {
		if input.Gender != nil {
			updatedPerformer.Gender = models.NewOptionalString(input.Gender.String())
		} else {
			updatedPerformer.Gender = models.NewOptionalStringPtr(nil)
		}
	}

	updatedPerformer.Birthdate = translator.optionalDate(input.Birthdate, "birthdate")
	updatedPerformer.Country = translator.optionalString(input.Country, "country")
	updatedPerformer.EyeColor = translator.optionalString(input.EyeColor, "eye_color")
	updatedPerformer.Measurements = translator.optionalString(input.Measurements, "measurements")
	// prefer height_cm over height
	if translator.hasField("height_cm") {
		updatedPerformer.Height = translator.optionalInt(input.HeightCm, "height_cm")
	} else if translator.hasField("height") {
		updatedPerformer.Height, err = translator.optionalIntFromString(input.Height, "height")
		if err != nil {
			return nil, err
		}
	}

	updatedPerformer.Ethnicity = translator.optionalString(input.Ethnicity, "ethnicity")
	updatedPerformer.FakeTits = translator.optionalString(input.FakeTits, "fake_tits")
	updatedPerformer.CareerLength = translator.optionalString(input.CareerLength, "career_length")
	updatedPerformer.Tattoos = translator.optionalString(input.Tattoos, "tattoos")
	updatedPerformer.Piercings = translator.optionalString(input.Piercings, "piercings")
	updatedPerformer.Twitter = translator.optionalString(input.Twitter, "twitter")
	updatedPerformer.Instagram = translator.optionalString(input.Instagram, "instagram")
	updatedPerformer.Favorite = translator.optionalBool(input.Favorite, "favorite")
	updatedPerformer.Rating = translator.ratingConversionOptional(input.Rating, input.Rating100)
	updatedPerformer.Details = translator.optionalString(input.Details, "details")
	updatedPerformer.DeathDate = translator.optionalDate(input.DeathDate, "death_date")
	updatedPerformer.HairColor = translator.optionalString(input.HairColor, "hair_color")
	updatedPerformer.Weight = translator.optionalInt(input.Weight, "weight")
	updatedPerformer.IgnoreAutoTag = translator.optionalBool(input.IgnoreAutoTag, "ignore_auto_tag")

	if translator.hasField("alias_list") {
		updatedPerformer.Aliases = &models.UpdateStrings{
			Values: input.AliasList,
			Mode:   models.RelationshipUpdateModeSet,
		}
	} else if translator.hasField("aliases") {
		updatedPerformer.Aliases = &models.UpdateStrings{
			Values: stringslice.FromString(*input.Aliases, ","),
			Mode:   models.RelationshipUpdateModeSet,
		}
	}

	if translator.hasField("tag_ids") {
		updatedPerformer.TagIDs, err = translateUpdateIDs(input.TagIds, models.RelationshipUpdateModeSet)
		if err != nil {
			return nil, fmt.Errorf("converting tag ids: %w", err)
		}
	}

	// Save the stash_ids
	if translator.hasField("stash_ids") {
		updatedPerformer.StashIDs = &models.UpdateStashIDs{
			StashIDs: stashIDPtrSliceToSlice(input.StashIds),
			Mode:     models.RelationshipUpdateModeSet,
		}
	}

	// Start the transaction and save the p
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Performer

		// need to get existing performer
		existing, err := qb.Find(ctx, performerID)
		if err != nil {
			return err
		}

		if existing == nil {
			return fmt.Errorf("performer with id %d not found", performerID)
		}

		if err := performer.ValidateDeathDate(existing, input.Birthdate, input.DeathDate); err != nil {
			if err != nil {
				return err
			}
		}

		_, err = qb.UpdatePartial(ctx, performerID, updatedPerformer)
		if err != nil {
			return err
		}

		// update image table
		if len(imageData) > 0 {
			if err := qb.UpdateImage(ctx, performerID, imageData); err != nil {
				return err
			}
		} else if imageIncluded {
			// must be unsetting
			if err := qb.DestroyImage(ctx, performerID); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, performerID, plugin.PerformerUpdatePost, input, translator.getFields())
	return r.getPerformer(ctx, performerID)
}

func (r *mutationResolver) BulkPerformerUpdate(ctx context.Context, input BulkPerformerUpdateInput) ([]*models.Performer, error) {
	performerIDs, err := stringslice.StringSliceToIntSlice(input.Ids)
	if err != nil {
		return nil, err
	}

	// Populate performer from the input
	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	updatedPerformer := models.NewPerformerPartial()

	updatedPerformer.Disambiguation = translator.optionalString(input.Disambiguation, "disambiguation")
	updatedPerformer.URL = translator.optionalString(input.URL, "url")
	updatedPerformer.Birthdate = translator.optionalDate(input.Birthdate, "birthdate")
	updatedPerformer.Ethnicity = translator.optionalString(input.Ethnicity, "ethnicity")
	updatedPerformer.Country = translator.optionalString(input.Country, "country")
	updatedPerformer.EyeColor = translator.optionalString(input.EyeColor, "eye_color")
	// prefer height_cm over height
	if translator.hasField("height_cm") {
		updatedPerformer.Height = translator.optionalInt(input.HeightCm, "height_cm")
	} else if translator.hasField("height") {
		updatedPerformer.Height, err = translator.optionalIntFromString(input.Height, "height")
		if err != nil {
			return nil, err
		}
	}

	updatedPerformer.Measurements = translator.optionalString(input.Measurements, "measurements")
	updatedPerformer.FakeTits = translator.optionalString(input.FakeTits, "fake_tits")
	updatedPerformer.CareerLength = translator.optionalString(input.CareerLength, "career_length")
	updatedPerformer.Tattoos = translator.optionalString(input.Tattoos, "tattoos")
	updatedPerformer.Piercings = translator.optionalString(input.Piercings, "piercings")
	updatedPerformer.Twitter = translator.optionalString(input.Twitter, "twitter")
	updatedPerformer.Instagram = translator.optionalString(input.Instagram, "instagram")
	updatedPerformer.Favorite = translator.optionalBool(input.Favorite, "favorite")
	updatedPerformer.Rating = translator.ratingConversionOptional(input.Rating, input.Rating100)
	updatedPerformer.Details = translator.optionalString(input.Details, "details")
	updatedPerformer.DeathDate = translator.optionalDate(input.DeathDate, "death_date")
	updatedPerformer.HairColor = translator.optionalString(input.HairColor, "hair_color")
	updatedPerformer.Weight = translator.optionalInt(input.Weight, "weight")
	updatedPerformer.IgnoreAutoTag = translator.optionalBool(input.IgnoreAutoTag, "ignore_auto_tag")

	if translator.hasField("alias_list") {
		updatedPerformer.Aliases = &models.UpdateStrings{
			Values: input.AliasList.Values,
			Mode:   input.AliasList.Mode,
		}
	} else if translator.hasField("aliases") {
		updatedPerformer.Aliases = &models.UpdateStrings{
			Values: stringslice.FromString(*input.Aliases, ","),
			Mode:   models.RelationshipUpdateModeSet,
		}
	}

	if translator.hasField("gender") {
		if input.Gender != nil {
			updatedPerformer.Gender = models.NewOptionalString(input.Gender.String())
		} else {
			updatedPerformer.Gender = models.NewOptionalStringPtr(nil)
		}
	}

	if translator.hasField("tag_ids") {
		updatedPerformer.TagIDs, err = translateUpdateIDs(input.TagIds.Ids, input.TagIds.Mode)
		if err != nil {
			return nil, fmt.Errorf("converting tag ids: %w", err)
		}
	}

	ret := []*models.Performer{}

	// Start the transaction and save the scene marker
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Performer

		for _, performerID := range performerIDs {
			updatedPerformer.ID = performerID

			// need to get existing performer
			existing, err := qb.Find(ctx, performerID)
			if err != nil {
				return err
			}

			if existing == nil {
				return fmt.Errorf("performer with id %d not found", performerID)
			}

			if err := performer.ValidateDeathDate(existing, input.Birthdate, input.DeathDate); err != nil {
				return err
			}

			performer, err := qb.UpdatePartial(ctx, performerID, updatedPerformer)
			if err != nil {
				return err
			}

			ret = append(ret, performer)
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

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		return r.repository.Performer.Destroy(ctx, id)
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

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Performer
		for _, id := range ids {
			if err := qb.Destroy(ctx, id); err != nil {
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
