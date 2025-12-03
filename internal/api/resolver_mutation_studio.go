package api

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin/hook"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/studio"
	"github.com/stashapp/stash/pkg/utils"
)

// used to refetch studio after hooks run
func (r *mutationResolver) getStudio(ctx context.Context, id int) (ret *models.Studio, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Studio.Find(ctx, id)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *mutationResolver) StudioCreate(ctx context.Context, input models.StudioCreateInput) (*models.Studio, error) {
	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Populate a new studio from the input
	newStudio := models.NewStudio()

	newStudio.Name = strings.TrimSpace(input.Name)
	newStudio.Rating = input.Rating100
	newStudio.Favorite = translator.bool(input.Favorite)
	newStudio.Details = translator.string(input.Details)
	newStudio.IgnoreAutoTag = translator.bool(input.IgnoreAutoTag)
	newStudio.Aliases = models.NewRelatedStrings(stringslice.TrimSpace(input.Aliases))
	newStudio.StashIDs = models.NewRelatedStashIDs(models.StashIDInputs(input.StashIds).ToStashIDs())

	var err error

	newStudio.URLs = models.NewRelatedStrings([]string{})
	if input.URL != nil {
		newStudio.URLs.Add(strings.TrimSpace(*input.URL))
	}

	if input.Urls != nil {
		newStudio.URLs.Add(stringslice.TrimSpace(input.Urls)...)
	}

	newStudio.ParentID, err = translator.intPtrFromString(input.ParentID)
	if err != nil {
		return nil, fmt.Errorf("converting parent id: %w", err)
	}

	newStudio.TagIDs, err = translator.relatedIds(input.TagIds)
	if err != nil {
		return nil, fmt.Errorf("converting tag ids: %w", err)
	}

	// Process the base 64 encoded image string
	var imageData []byte
	if input.Image != nil {
		var err error
		imageData, err = utils.ProcessImageInput(ctx, *input.Image)
		if err != nil {
			return nil, fmt.Errorf("processing image: %w", err)
		}
	}

	// Start the transaction and save the studio
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Studio

		if err := studio.ValidateCreate(ctx, newStudio, qb); err != nil {
			return err
		}

		err = qb.Create(ctx, &newStudio)
		if err != nil {
			return err
		}

		if len(imageData) > 0 {
			if err := qb.UpdateImage(ctx, newStudio.ID, imageData); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, newStudio.ID, hook.StudioCreatePost, input, nil)
	return r.getStudio(ctx, newStudio.ID)
}

func (r *mutationResolver) StudioUpdate(ctx context.Context, input models.StudioUpdateInput) (*models.Studio, error) {
	studioID, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, fmt.Errorf("converting id: %w", err)
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Populate studio from the input
	updatedStudio := models.NewStudioPartial()

	updatedStudio.ID = studioID
	updatedStudio.Name = translator.optionalString(input.Name, "name")
	updatedStudio.Details = translator.optionalString(input.Details, "details")
	updatedStudio.Rating = translator.optionalInt(input.Rating100, "rating100")
	updatedStudio.Favorite = translator.optionalBool(input.Favorite, "favorite")
	updatedStudio.IgnoreAutoTag = translator.optionalBool(input.IgnoreAutoTag, "ignore_auto_tag")
	updatedStudio.Aliases = translator.updateStrings(input.Aliases, "aliases")
	updatedStudio.StashIDs = translator.updateStashIDs(input.StashIds, "stash_ids")

	updatedStudio.ParentID, err = translator.optionalIntFromString(input.ParentID, "parent_id")
	if err != nil {
		return nil, fmt.Errorf("converting parent id: %w", err)
	}

	updatedStudio.TagIDs, err = translator.updateIds(input.TagIds, "tag_ids")
	if err != nil {
		return nil, fmt.Errorf("converting tag ids: %w", err)
	}

	if translator.hasField("urls") {
		// ensure url not included in the input
		if err := validateNoLegacyURLs(translator); err != nil {
			return nil, err
		}

		updatedStudio.URLs = translator.updateStrings(input.Urls, "urls")
	} else if translator.hasField("url") {
		// handle legacy url field
		legacyURLs := []string{}
		if input.URL != nil {
			legacyURLs = append(legacyURLs, *input.URL)
		}

		updatedStudio.URLs = &models.UpdateStrings{
			Mode:   models.RelationshipUpdateModeSet,
			Values: legacyURLs,
		}
	}

	// Process the base 64 encoded image string
	var imageData []byte
	imageIncluded := translator.hasField("image")
	if input.Image != nil {
		var err error
		imageData, err = utils.ProcessImageInput(ctx, *input.Image)
		if err != nil {
			return nil, fmt.Errorf("processing image: %w", err)
		}
	}

	// Start the transaction and update the studio
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Studio

		if err := studio.ValidateModify(ctx, updatedStudio, qb); err != nil {
			return err
		}

		_, err = qb.UpdatePartial(ctx, updatedStudio)
		if err != nil {
			return err
		}

		if imageIncluded {
			if err := qb.UpdateImage(ctx, studioID, imageData); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, studioID, hook.StudioUpdatePost, input, translator.getFields())
	return r.getStudio(ctx, studioID)
}

func (r *mutationResolver) BulkStudioUpdate(ctx context.Context, input BulkStudioUpdateInput) ([]*models.Studio, error) {
	ids, err := stringslice.StringSliceToIntSlice(input.Ids)
	if err != nil {
		return nil, fmt.Errorf("converting ids: %w", err)
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Populate performer from the input
	partial := models.NewStudioPartial()

	partial.ParentID, err = translator.optionalIntFromString(input.ParentID, "parent_id")
	if err != nil {
		return nil, fmt.Errorf("converting parent id: %w", err)
	}

	if translator.hasField("urls") {
		// ensure url/twitter/instagram are not included in the input
		if err := validateNoLegacyURLs(translator); err != nil {
			return nil, err
		}

		partial.URLs = translator.updateStringsBulk(input.Urls, "urls")
	} else if translator.hasField("url") {
		// handle legacy url field
		legacyURLs := []string{}
		if input.URL != nil {
			legacyURLs = append(legacyURLs, *input.URL)
		}

		partial.URLs = &models.UpdateStrings{
			Mode:   models.RelationshipUpdateModeSet,
			Values: legacyURLs,
		}
	}

	partial.Favorite = translator.optionalBool(input.Favorite, "favorite")
	partial.Rating = translator.optionalInt(input.Rating100, "rating100")
	partial.Details = translator.optionalString(input.Details, "details")
	partial.IgnoreAutoTag = translator.optionalBool(input.IgnoreAutoTag, "ignore_auto_tag")

	partial.TagIDs, err = translator.updateIdsBulk(input.TagIds, "tag_ids")
	if err != nil {
		return nil, fmt.Errorf("converting tag ids: %w", err)
	}

	ret := []*models.Studio{}

	// Start the transaction and save the performers
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Studio

		for _, id := range ids {
			local := partial
			local.ID = id
			if err := studio.ValidateModify(ctx, local, qb); err != nil {
				return err
			}

			updated, err := qb.UpdatePartial(ctx, local)
			if err != nil {
				return err
			}

			ret = append(ret, updated)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// execute post hooks outside of txn
	var newRet []*models.Studio
	for _, studio := range ret {
		r.hookExecutor.ExecutePostHooks(ctx, studio.ID, hook.StudioUpdatePost, input, translator.getFields())

		studio, err = r.getStudio(ctx, studio.ID)
		if err != nil {
			return nil, err
		}

		newRet = append(newRet, studio)
	}

	return newRet, nil
}

func (r *mutationResolver) StudioDestroy(ctx context.Context, input StudioDestroyInput) (bool, error) {
	id, err := strconv.Atoi(input.ID)
	if err != nil {
		return false, fmt.Errorf("converting id: %w", err)
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		return r.repository.Studio.Destroy(ctx, id)
	}); err != nil {
		return false, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, id, hook.StudioDestroyPost, input, nil)

	return true, nil
}

func (r *mutationResolver) StudiosDestroy(ctx context.Context, studioIDs []string) (bool, error) {
	ids, err := stringslice.StringSliceToIntSlice(studioIDs)
	if err != nil {
		return false, fmt.Errorf("converting ids: %w", err)
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Studio
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
		r.hookExecutor.ExecutePostHooks(ctx, id, hook.StudioDestroyPost, studioIDs, nil)
	}

	return true, nil
}
