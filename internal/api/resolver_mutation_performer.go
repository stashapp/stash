package api

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/performer"
	"github.com/stashapp/stash/pkg/plugin/hook"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/utils"
)

const (
	twitterURL   = "https://twitter.com"
	instagramURL = "https://instagram.com"
)

// used to refetch performer after hooks run
func (r *mutationResolver) getPerformer(ctx context.Context, id int) (ret *models.Performer, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Performer.Find(ctx, id)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *mutationResolver) PerformerCreate(ctx context.Context, input models.PerformerCreateInput) (*models.Performer, error) {
	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Populate a new performer from the input
	newPerformer := models.NewPerformer()

	newPerformer.Name = strings.TrimSpace(input.Name)
	newPerformer.Disambiguation = translator.string(input.Disambiguation)
	newPerformer.Aliases = models.NewRelatedStrings(stringslice.TrimSpace(input.AliasList))
	newPerformer.Gender = input.Gender
	newPerformer.Ethnicity = translator.string(input.Ethnicity)
	newPerformer.Country = translator.string(input.Country)
	newPerformer.EyeColor = translator.string(input.EyeColor)
	newPerformer.Measurements = translator.string(input.Measurements)
	newPerformer.FakeTits = translator.string(input.FakeTits)
	newPerformer.PenisLength = input.PenisLength
	newPerformer.Circumcised = input.Circumcised
	newPerformer.CareerLength = translator.string(input.CareerLength)
	newPerformer.Tattoos = translator.string(input.Tattoos)
	newPerformer.Piercings = translator.string(input.Piercings)
	newPerformer.Favorite = translator.bool(input.Favorite)
	newPerformer.Rating = input.Rating100
	newPerformer.Details = translator.string(input.Details)
	newPerformer.HairColor = translator.string(input.HairColor)
	newPerformer.Height = input.HeightCm
	newPerformer.Weight = input.Weight
	newPerformer.IgnoreAutoTag = translator.bool(input.IgnoreAutoTag)
	newPerformer.StashIDs = models.NewRelatedStashIDs(models.StashIDInputs(input.StashIds).ToStashIDs())

	newPerformer.URLs = models.NewRelatedStrings([]string{})
	if input.URL != nil {
		newPerformer.URLs.Add(strings.TrimSpace(*input.URL))
	}
	if input.Twitter != nil {
		newPerformer.URLs.Add(utils.URLFromHandle(strings.TrimSpace(*input.Twitter), twitterURL))
	}
	if input.Instagram != nil {
		newPerformer.URLs.Add(utils.URLFromHandle(strings.TrimSpace(*input.Instagram), instagramURL))
	}

	if input.Urls != nil {
		newPerformer.URLs.Add(stringslice.TrimSpace(input.Urls)...)
	}

	var err error

	newPerformer.Birthdate, err = translator.datePtr(input.Birthdate)
	if err != nil {
		return nil, fmt.Errorf("converting birthdate: %w", err)
	}
	newPerformer.DeathDate, err = translator.datePtr(input.DeathDate)
	if err != nil {
		return nil, fmt.Errorf("converting death date: %w", err)
	}

	newPerformer.TagIDs, err = translator.relatedIds(input.TagIds)
	if err != nil {
		return nil, fmt.Errorf("converting tag ids: %w", err)
	}

	// Process the base 64 encoded image string
	var imageData []byte
	if input.Image != nil {
		imageData, err = utils.ProcessImageInput(ctx, *input.Image)
		if err != nil {
			return nil, fmt.Errorf("processing image: %w", err)
		}
	}

	// Start the transaction and save the performer
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Performer

		if err := performer.ValidateCreate(ctx, newPerformer, qb); err != nil {
			return err
		}

		i := &models.CreatePerformerInput{
			Performer: &newPerformer,
			// convert json.Numbers to int/float
			CustomFields: convertMapJSONNumbers(input.CustomFields),
		}

		err = qb.Create(ctx, i)
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

	r.hookExecutor.ExecutePostHooks(ctx, newPerformer.ID, hook.PerformerCreatePost, input, nil)
	return r.getPerformer(ctx, newPerformer.ID)
}

func validateNoLegacyURLs(translator changesetTranslator) error {
	// ensure url/twitter/instagram are not included in the input
	if translator.hasField("url") {
		return fmt.Errorf("url field must not be included if urls is included")
	}
	if translator.hasField("twitter") {
		return fmt.Errorf("twitter field must not be included if urls is included")
	}
	if translator.hasField("instagram") {
		return fmt.Errorf("instagram field must not be included if urls is included")
	}

	return nil
}

func (r *mutationResolver) handleLegacyURLs(ctx context.Context, performerID int, legacyURLs legacyPerformerURLs, updatedPerformer *models.PerformerPartial) error {
	qb := r.repository.Performer

	// we need to be careful with URL/Twitter/Instagram
	// treat URL as replacing the first non-Twitter/Instagram URL in the list
	// twitter should replace any existing twitter URL
	// instagram should replace any existing instagram URL
	p, err := qb.Find(ctx, performerID)
	if err != nil {
		return err
	}

	if err := p.LoadURLs(ctx, qb); err != nil {
		return fmt.Errorf("loading performer URLs: %w", err)
	}

	existingURLs := p.URLs.List()

	// performer partial URLs should be empty
	if legacyURLs.URL.Set {
		replaced := false
		for i, url := range existingURLs {
			if !performer.IsTwitterURL(url) && !performer.IsInstagramURL(url) {
				existingURLs[i] = legacyURLs.URL.Value
				replaced = true
				break
			}
		}

		if !replaced {
			existingURLs = append(existingURLs, legacyURLs.URL.Value)
		}
	}

	if legacyURLs.Twitter.Set {
		value := utils.URLFromHandle(legacyURLs.Twitter.Value, twitterURL)
		found := false
		// find and replace the first twitter URL
		for i, url := range existingURLs {
			if performer.IsTwitterURL(url) {
				existingURLs[i] = value
				found = true
				break
			}
		}

		if !found {
			existingURLs = append(existingURLs, value)
		}
	}
	if legacyURLs.Instagram.Set {
		found := false
		value := utils.URLFromHandle(legacyURLs.Instagram.Value, instagramURL)
		// find and replace the first instagram URL
		for i, url := range existingURLs {
			if performer.IsInstagramURL(url) {
				existingURLs[i] = value
				found = true
				break
			}
		}

		if !found {
			existingURLs = append(existingURLs, value)
		}
	}

	updatedPerformer.URLs = &models.UpdateStrings{
		Values: existingURLs,
		Mode:   models.RelationshipUpdateModeSet,
	}

	return nil
}

type legacyPerformerURLs struct {
	URL       models.OptionalString
	Twitter   models.OptionalString
	Instagram models.OptionalString
}

func (u *legacyPerformerURLs) AnySet() bool {
	return u.URL.Set || u.Twitter.Set || u.Instagram.Set
}

func legacyPerformerURLsFromInput(input models.PerformerUpdateInput, translator changesetTranslator) legacyPerformerURLs {
	return legacyPerformerURLs{
		URL:       translator.optionalString(input.URL, "url"),
		Twitter:   translator.optionalString(input.Twitter, "twitter"),
		Instagram: translator.optionalString(input.Instagram, "instagram"),
	}
}

func performerPartialFromInput(input models.PerformerUpdateInput, translator changesetTranslator) (*models.PerformerPartial, error) {
	// Populate performer from the input
	updatedPerformer := models.NewPerformerPartial()

	updatedPerformer.Name = translator.optionalString(input.Name, "name")
	updatedPerformer.Disambiguation = translator.optionalString(input.Disambiguation, "disambiguation")
	updatedPerformer.Gender = translator.optionalString((*string)(input.Gender), "gender")
	updatedPerformer.Ethnicity = translator.optionalString(input.Ethnicity, "ethnicity")
	updatedPerformer.Country = translator.optionalString(input.Country, "country")
	updatedPerformer.EyeColor = translator.optionalString(input.EyeColor, "eye_color")
	updatedPerformer.Measurements = translator.optionalString(input.Measurements, "measurements")
	updatedPerformer.FakeTits = translator.optionalString(input.FakeTits, "fake_tits")
	updatedPerformer.PenisLength = translator.optionalFloat64(input.PenisLength, "penis_length")
	updatedPerformer.Circumcised = translator.optionalString((*string)(input.Circumcised), "circumcised")
	updatedPerformer.CareerLength = translator.optionalString(input.CareerLength, "career_length")
	updatedPerformer.Tattoos = translator.optionalString(input.Tattoos, "tattoos")
	updatedPerformer.Piercings = translator.optionalString(input.Piercings, "piercings")
	updatedPerformer.Favorite = translator.optionalBool(input.Favorite, "favorite")
	updatedPerformer.Rating = translator.optionalInt(input.Rating100, "rating100")
	updatedPerformer.Details = translator.optionalString(input.Details, "details")
	updatedPerformer.HairColor = translator.optionalString(input.HairColor, "hair_color")
	updatedPerformer.Weight = translator.optionalInt(input.Weight, "weight")
	updatedPerformer.IgnoreAutoTag = translator.optionalBool(input.IgnoreAutoTag, "ignore_auto_tag")
	updatedPerformer.StashIDs = translator.updateStashIDs(input.StashIds, "stash_ids")

	var err error

	if translator.hasField("urls") {
		// ensure url/twitter/instagram are not included in the input
		if err := validateNoLegacyURLs(translator); err != nil {
			return nil, err
		}

		updatedPerformer.URLs = translator.updateStrings(input.Urls, "urls")
	}

	updatedPerformer.Birthdate, err = translator.optionalDate(input.Birthdate, "birthdate")
	if err != nil {
		return nil, fmt.Errorf("converting birthdate: %w", err)
	}
	updatedPerformer.DeathDate, err = translator.optionalDate(input.DeathDate, "death_date")
	if err != nil {
		return nil, fmt.Errorf("converting death date: %w", err)
	}

	// prefer height_cm over height
	if translator.hasField("height_cm") {
		updatedPerformer.Height = translator.optionalInt(input.HeightCm, "height_cm")
	}

	// prefer alias_list over aliases
	if translator.hasField("alias_list") {
		updatedPerformer.Aliases = translator.updateStrings(input.AliasList, "alias_list")
	}

	updatedPerformer.TagIDs, err = translator.updateIds(input.TagIds, "tag_ids")
	if err != nil {
		return nil, fmt.Errorf("converting tag ids: %w", err)
	}

	updatedPerformer.CustomFields = handleUpdateCustomFields(input.CustomFields)

	return &updatedPerformer, nil
}

func (r *mutationResolver) PerformerUpdate(ctx context.Context, input models.PerformerUpdateInput) (*models.Performer, error) {
	performerID, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, fmt.Errorf("converting id: %w", err)
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	updatedPerformer, err := performerPartialFromInput(input, translator)
	if err != nil {
		return nil, err
	}

	legacyURLs := legacyPerformerURLsFromInput(input, translator)

	var imageData []byte
	imageIncluded := translator.hasField("image")
	if input.Image != nil {
		imageData, err = utils.ProcessImageInput(ctx, *input.Image)
		if err != nil {
			return nil, fmt.Errorf("processing image: %w", err)
		}
	}

	// Start the transaction and save the performer
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Performer

		if legacyURLs.AnySet() {
			if err := r.handleLegacyURLs(ctx, performerID, legacyURLs, updatedPerformer); err != nil {
				return err
			}
		}

		if err := performer.ValidateUpdate(ctx, performerID, *updatedPerformer, qb); err != nil {
			return err
		}

		_, err = qb.UpdatePartial(ctx, performerID, *updatedPerformer)
		if err != nil {
			return err
		}

		// update image table
		if imageIncluded {
			if err := qb.UpdateImage(ctx, performerID, imageData); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, performerID, hook.PerformerUpdatePost, input, translator.getFields())
	return r.getPerformer(ctx, performerID)
}

func (r *mutationResolver) BulkPerformerUpdate(ctx context.Context, input BulkPerformerUpdateInput) ([]*models.Performer, error) {
	performerIDs, err := stringslice.StringSliceToIntSlice(input.Ids)
	if err != nil {
		return nil, fmt.Errorf("converting ids: %w", err)
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Populate performer from the input
	updatedPerformer := models.NewPerformerPartial()

	updatedPerformer.Disambiguation = translator.optionalString(input.Disambiguation, "disambiguation")

	updatedPerformer.Gender = translator.optionalString((*string)(input.Gender), "gender")
	updatedPerformer.Ethnicity = translator.optionalString(input.Ethnicity, "ethnicity")
	updatedPerformer.Country = translator.optionalString(input.Country, "country")
	updatedPerformer.EyeColor = translator.optionalString(input.EyeColor, "eye_color")
	updatedPerformer.Measurements = translator.optionalString(input.Measurements, "measurements")
	updatedPerformer.FakeTits = translator.optionalString(input.FakeTits, "fake_tits")
	updatedPerformer.PenisLength = translator.optionalFloat64(input.PenisLength, "penis_length")
	updatedPerformer.Circumcised = translator.optionalString((*string)(input.Circumcised), "circumcised")
	updatedPerformer.CareerLength = translator.optionalString(input.CareerLength, "career_length")
	updatedPerformer.Tattoos = translator.optionalString(input.Tattoos, "tattoos")
	updatedPerformer.Piercings = translator.optionalString(input.Piercings, "piercings")

	updatedPerformer.Favorite = translator.optionalBool(input.Favorite, "favorite")
	updatedPerformer.Rating = translator.optionalInt(input.Rating100, "rating100")
	updatedPerformer.Details = translator.optionalString(input.Details, "details")
	updatedPerformer.HairColor = translator.optionalString(input.HairColor, "hair_color")
	updatedPerformer.Weight = translator.optionalInt(input.Weight, "weight")
	updatedPerformer.IgnoreAutoTag = translator.optionalBool(input.IgnoreAutoTag, "ignore_auto_tag")

	if translator.hasField("urls") {
		// ensure url/twitter/instagram are not included in the input
		if err := validateNoLegacyURLs(translator); err != nil {
			return nil, err
		}

		updatedPerformer.URLs = translator.updateStringsBulk(input.Urls, "urls")
	}

	legacyURLs := legacyPerformerURLs{
		URL:       translator.optionalString(input.URL, "url"),
		Twitter:   translator.optionalString(input.Twitter, "twitter"),
		Instagram: translator.optionalString(input.Instagram, "instagram"),
	}

	updatedPerformer.Birthdate, err = translator.optionalDate(input.Birthdate, "birthdate")
	if err != nil {
		return nil, fmt.Errorf("converting birthdate: %w", err)
	}
	updatedPerformer.DeathDate, err = translator.optionalDate(input.DeathDate, "death_date")
	if err != nil {
		return nil, fmt.Errorf("converting death date: %w", err)
	}

	// prefer height_cm over height
	if translator.hasField("height_cm") {
		updatedPerformer.Height = translator.optionalInt(input.HeightCm, "height_cm")
	}

	// prefer alias_list over aliases
	if translator.hasField("alias_list") {
		updatedPerformer.Aliases = translator.updateStringsBulk(input.AliasList, "alias_list")
	}

	updatedPerformer.TagIDs, err = translator.updateIdsBulk(input.TagIds, "tag_ids")
	if err != nil {
		return nil, fmt.Errorf("converting tag ids: %w", err)
	}

	if input.CustomFields != nil {
		updatedPerformer.CustomFields = handleUpdateCustomFields(*input.CustomFields)
	}

	ret := []*models.Performer{}

	// Start the transaction and save the performers
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Performer

		for _, performerID := range performerIDs {
			if legacyURLs.AnySet() {
				if err := r.handleLegacyURLs(ctx, performerID, legacyURLs, &updatedPerformer); err != nil {
					return err
				}
			}

			if err := performer.ValidateUpdate(ctx, performerID, updatedPerformer, qb); err != nil {
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
		r.hookExecutor.ExecutePostHooks(ctx, performer.ID, hook.PerformerUpdatePost, input, translator.getFields())

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
		return false, fmt.Errorf("converting id: %w", err)
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		return r.repository.Performer.Destroy(ctx, id)
	}); err != nil {
		return false, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, id, hook.PerformerDestroyPost, input, nil)

	return true, nil
}

func (r *mutationResolver) PerformersDestroy(ctx context.Context, performerIDs []string) (bool, error) {
	ids, err := stringslice.StringSliceToIntSlice(performerIDs)
	if err != nil {
		return false, fmt.Errorf("converting ids: %w", err)
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
		r.hookExecutor.ExecutePostHooks(ctx, id, hook.PerformerDestroyPost, performerIDs, nil)
	}

	return true, nil
}

func (r *mutationResolver) PerformerMerge(ctx context.Context, input PerformerMergeInput) (*models.Performer, error) {
	srcIDs, err := stringslice.StringSliceToIntSlice(input.Source)
	if err != nil {
		return nil, fmt.Errorf("converting source ids: %w", err)
	}

	// ensure source ids are unique
	srcIDs = sliceutil.AppendUniques(nil, srcIDs)

	destID, err := strconv.Atoi(input.Destination)
	if err != nil {
		return nil, fmt.Errorf("converting destination id: %w", err)
	}

	// ensure destination is not in source list
	if slices.Contains(srcIDs, destID) {
		return nil, errors.New("destination performer cannot be in source list")
	}

	var values *models.PerformerPartial
	var imageData []byte

	if input.Values != nil {
		translator := changesetTranslator{
			inputMap: getNamedUpdateInputMap(ctx, "input.values"),
		}

		values, err = performerPartialFromInput(*input.Values, translator)
		if err != nil {
			return nil, err
		}
		legacyURLs := legacyPerformerURLsFromInput(*input.Values, translator)
		if legacyURLs.AnySet() {
			return nil, errors.New("Merging legacy performer URLs is not supported")
		}

		if input.Values.Image != nil {
			var err error
			imageData, err = utils.ProcessImageInput(ctx, *input.Values.Image)
			if err != nil {
				return nil, fmt.Errorf("processing cover image: %w", err)
			}
		}
	} else {
		v := models.NewPerformerPartial()
		values = &v
	}

	var dest *models.Performer
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Performer

		dest, err = qb.Find(ctx, destID)
		if err != nil {
			return fmt.Errorf("finding destination performer ID %d: %w", destID, err)
		}

		// ensure source performers exist
		if _, err := qb.FindMany(ctx, srcIDs); err != nil {
			return fmt.Errorf("finding source performers: %w", err)
		}

		if _, err := qb.UpdatePartial(ctx, destID, *values); err != nil {
			return fmt.Errorf("updating performer: %w", err)
		}

		if err := qb.Merge(ctx, srcIDs, destID); err != nil {
			return fmt.Errorf("merging performers: %w", err)
		}

		if len(imageData) > 0 {
			if err := qb.UpdateImage(ctx, destID, imageData); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return dest, nil
}
