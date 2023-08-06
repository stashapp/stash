package api

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
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
	currentTime := time.Now()
	newStudio := models.Studio{
		Name:          input.Name,
		CreatedAt:     currentTime,
		UpdatedAt:     currentTime,
		URL:           translator.string(input.URL, "url"),
		Rating:        translator.ratingConversion(input.Rating, input.Rating100),
		Details:       translator.string(input.Details, "details"),
		IgnoreAutoTag: translator.bool(input.IgnoreAutoTag, "ignore_auto_tag"),
		Aliases:       models.NewRelatedStrings(input.Aliases),
		StashIDs:      models.NewRelatedStashIDs(input.StashIds),
	}

	var err error

	newStudio.ParentID, err = translator.intPtrFromString(input.ParentID, "parent_id")
	if err != nil {
		return nil, fmt.Errorf("converting parent id: %w", err)
	}

	// Process the base 64 encoded image string
	var imageData []byte
	if input.Image != nil {
		var err error
		imageData, err = utils.ProcessImageInput(ctx, *input.Image)
		if err != nil {
			return nil, err
		}
	}

	// Start the transaction and save the studio
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Studio

		if len(input.Aliases) > 0 {
			if err := studio.EnsureAliasesUnique(ctx, 0, input.Aliases, qb); err != nil {
				return err
			}
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

	r.hookExecutor.ExecutePostHooks(ctx, newStudio.ID, plugin.StudioCreatePost, input, nil)
	return r.getStudio(ctx, newStudio.ID)
}

func (r *mutationResolver) StudioUpdate(ctx context.Context, input models.StudioUpdateInput) (*models.Studio, error) {
	studioID, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, err
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Populate studio from the input
	updatedStudio := models.StudioPartial{
		ID:            studioID,
		Name:          translator.optionalString(input.Name, "name"),
		URL:           translator.optionalString(input.URL, "url"),
		Details:       translator.optionalString(input.Details, "details"),
		Rating:        translator.optionalRatingConversion(input.Rating, input.Rating100),
		IgnoreAutoTag: translator.optionalBool(input.IgnoreAutoTag, "ignore_auto_tag"),
		Aliases:       translator.updateStrings(input.Aliases, "aliases"),
		StashIDs:      translator.updateStashIDs(input.StashIds, "stash_ids"),
		UpdatedAt:     models.NewOptionalTime(time.Now()),
	}

	updatedStudio.ParentID, err = translator.optionalIntFromString(input.ParentID, "parent_id")
	if err != nil {
		return nil, fmt.Errorf("converting parent id: %w", err)
	}

	// Process the base 64 encoded image string
	var imageData []byte
	imageIncluded := translator.hasField("image")
	if input.Image != nil {
		var err error
		imageData, err = utils.ProcessImageInput(ctx, *input.Image)
		if err != nil {
			return nil, err
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

	r.hookExecutor.ExecutePostHooks(ctx, studioID, plugin.StudioUpdatePost, input, translator.getFields())
	return r.getStudio(ctx, studioID)
}

func (r *mutationResolver) StudioDestroy(ctx context.Context, input StudioDestroyInput) (bool, error) {
	id, err := strconv.Atoi(input.ID)
	if err != nil {
		return false, err
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		return r.repository.Studio.Destroy(ctx, id)
	}); err != nil {
		return false, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, id, plugin.StudioDestroyPost, input, nil)

	return true, nil
}

func (r *mutationResolver) StudiosDestroy(ctx context.Context, studioIDs []string) (bool, error) {
	ids, err := stringslice.StringSliceToIntSlice(studioIDs)
	if err != nil {
		return false, err
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
		r.hookExecutor.ExecutePostHooks(ctx, id, plugin.StudioDestroyPost, studioIDs, nil)
	}

	return true, nil
}
