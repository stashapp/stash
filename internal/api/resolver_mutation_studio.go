package api

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/studio"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *mutationResolver) getStudio(ctx context.Context, id int) (ret *models.Studio, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Studio.Find(ctx, id)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *mutationResolver) StudioCreate(ctx context.Context, input StudioCreateInput) (*models.Studio, error) {
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
		Rating:        translator.ratingConversionInt(input.Rating, input.Rating100),
		Details:       translator.string(input.Details, "details"),
		IgnoreAutoTag: translator.bool(input.IgnoreAutoTag, "ignore_auto_tag"),
	}

	var err error

	newStudio.ParentID, err = translator.intPtrFromString(input.ParentID, "parent_id")
	if err != nil {
		return nil, fmt.Errorf("converting parent id: %w", err)
	}

	// Process the base 64 encoded image string
	var imageData []byte
	if input.Image != nil {
		imageData, err = utils.ProcessImageInput(ctx, *input.Image)
		if err != nil {
			return nil, err
		}
	}

	// Start the transaction and save the studio
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Studio

		err = qb.Create(ctx, &newStudio)
		if err != nil {
			return err
		}

		// update image table
		if len(imageData) > 0 {
			if err := qb.UpdateImage(ctx, newStudio.ID, imageData); err != nil {
				return err
			}
		}

		// Save the stash_ids
		if input.StashIds != nil {
			stashIDJoins := stashIDPtrSliceToSlice(input.StashIds)
			if err := qb.UpdateStashIDs(ctx, newStudio.ID, stashIDJoins); err != nil {
				return err
			}
		}

		if len(input.Aliases) > 0 {
			if err := studio.EnsureAliasesUnique(ctx, newStudio.ID, input.Aliases, qb); err != nil {
				return err
			}

			if err := qb.UpdateAliases(ctx, newStudio.ID, input.Aliases); err != nil {
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

func (r *mutationResolver) StudioUpdate(ctx context.Context, input StudioUpdateInput) (*models.Studio, error) {
	studioID, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, err
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Populate studio from the input
	updatedStudio := models.NewStudioPartial()

	updatedStudio.Name = translator.optionalString(input.Name, "name")
	updatedStudio.URL = translator.optionalString(input.URL, "url")
	updatedStudio.Details = translator.optionalString(input.Details, "details")
	updatedStudio.Rating = translator.ratingConversionOptional(input.Rating, input.Rating100)
	updatedStudio.IgnoreAutoTag = translator.optionalBool(input.IgnoreAutoTag, "ignore_auto_tag")
	updatedStudio.ParentID, err = translator.optionalIntFromString(input.ParentID, "parent_id")
	if err != nil {
		return nil, fmt.Errorf("converting parent id: %w", err)
	}

	var imageData []byte
	imageIncluded := translator.hasField("image")
	if input.Image != nil {
		imageData, err = utils.ProcessImageInput(ctx, *input.Image)
		if err != nil {
			return nil, err
		}
	}

	// Start the transaction and save the studio
	var s *models.Studio
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Studio

		if err := manager.ValidateModifyStudio(ctx, studioID, updatedStudio, qb); err != nil {
			return err
		}

		var err error
		s, err = qb.UpdatePartial(ctx, studioID, updatedStudio)
		if err != nil {
			return err
		}

		// update image table
		if imageIncluded {
			if err := qb.UpdateImage(ctx, s.ID, imageData); err != nil {
				return err
			}
		}

		// Save the stash_ids
		if translator.hasField("stash_ids") {
			stashIDJoins := stashIDPtrSliceToSlice(input.StashIds)
			if err := qb.UpdateStashIDs(ctx, studioID, stashIDJoins); err != nil {
				return err
			}
		}

		if translator.hasField("aliases") {
			if err := studio.EnsureAliasesUnique(ctx, studioID, input.Aliases, qb); err != nil {
				return err
			}

			if err := qb.UpdateAliases(ctx, studioID, input.Aliases); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, s.ID, plugin.StudioUpdatePost, input, translator.getFields())
	return r.getStudio(ctx, s.ID)
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
