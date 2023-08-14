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

func (r *mutationResolver) StudioCreate(ctx context.Context, input StudioCreateInput) (*models.Studio, error) {
	s, err := studioFromStudioCreateInput(ctx, input)
	if err != nil {
		return nil, err
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

		if s.Aliases.Loaded() && len(s.Aliases.List()) > 0 {
			if err := studio.EnsureAliasesUnique(ctx, 0, s.Aliases.List(), qb); err != nil {
				return err
			}
		}

		err = qb.Create(ctx, s)
		if err != nil {
			return err
		}

		if len(imageData) > 0 {
			if err := qb.UpdateImage(ctx, s.ID, imageData); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, s.ID, plugin.StudioCreatePost, input, nil)

	return s, nil
}

func studioFromStudioCreateInput(ctx context.Context, input StudioCreateInput) (*models.Studio, error) {
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

	if input.Aliases != nil {
		newStudio.Aliases = models.NewRelatedStrings(input.Aliases)
	}
	if input.StashIds != nil {
		newStudio.StashIDs = models.NewRelatedStashIDs(stashIDPtrSliceToSlice(input.StashIds))
	}

	return &newStudio, nil
}

func (r *mutationResolver) StudioUpdate(ctx context.Context, input StudioUpdateInput) (*models.Studio, error) {
	var updatedStudio *models.Studio
	var err error

	translator := changesetTranslator{
		inputMap: getNamedUpdateInputMap(ctx, updateInputField),
	}
	s := studioPartialFromStudioUpdateInput(input, &input.ID, translator)

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

		if err := studio.ValidateModify(ctx, *s, qb); err != nil {
			return err
		}

		updatedStudio, err = qb.UpdatePartial(ctx, *s)
		if err != nil {
			return err
		}

		if imageIncluded {
			if err := qb.UpdateImage(ctx, s.ID, imageData); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, updatedStudio.ID, plugin.StudioUpdatePost, input, translator.getFields())

	return updatedStudio, nil
}

// This is slightly different to studioPartialFromStudioCreateInput in that Name is handled differently
// and ImageIncluded is not hardcoded to true
func studioPartialFromStudioUpdateInput(input StudioUpdateInput, id *string, translator changesetTranslator) *models.StudioPartial {
	// Populate studio from the input
	updatedStudio := models.StudioPartial{
		Name:          translator.optionalString(input.Name, "name"),
		URL:           translator.optionalString(input.URL, "url"),
		Details:       translator.optionalString(input.Details, "details"),
		Rating:        translator.ratingConversionOptional(input.Rating, input.Rating100),
		IgnoreAutoTag: translator.optionalBool(input.IgnoreAutoTag, "ignore_auto_tag"),
		UpdatedAt:     models.NewOptionalTime(time.Now()),
	}

	updatedStudio.ID, _ = strconv.Atoi(*id)

	if input.ParentID != nil {
		parentID, _ := strconv.Atoi(*input.ParentID)
		if parentID > 0 {
			// This is to be set directly as we know it has a value and the translator won't have the field
			updatedStudio.ParentID = models.NewOptionalInt(parentID)
		}
	} else {
		updatedStudio.ParentID = translator.optionalInt(nil, "parent_id")
	}

	if translator.hasField("aliases") {
		updatedStudio.Aliases = &models.UpdateStrings{
			Values: input.Aliases,
			Mode:   models.RelationshipUpdateModeSet,
		}
	}

	if translator.hasField("stash_ids") {
		updatedStudio.StashIDs = &models.UpdateStashIDs{
			StashIDs: stashIDPtrSliceToSlice(input.StashIds),
			Mode:     models.RelationshipUpdateModeSet,
		}
	}

	return &updatedStudio
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
