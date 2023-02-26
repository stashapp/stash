package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
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
	var studioID *int
	var dbInput models.StudioDBInput
	var err error
	var parentTranslator changesetTranslator
	runParentCreateHook := false
	runParentUpdateHook := false

	if input.Parent != nil {
		if input.ParentID == nil {
			// The parent needs to be created
			dbInput.ParentCreate, err = studioFromStudioCreateInput(ctx, *input.Parent)
			runParentCreateHook = true
		} else {
			parentTranslator = changesetTranslator{
				inputMap: getNamedUpdateInputMap(ctx, updateInputField+".parent"),
			}

			// The parent studio matched an existing one and the user has chosen in the UI to link and/or update it
			dbInput.ParentUpdate, err = studioPartialFromStudioCreateInput(ctx, *input.Parent, input.ParentID, parentTranslator)
			runParentUpdateHook = true
		}
		if err != nil {
			return nil, err
		}
	}

	dbInput.StudioCreate, err = studioFromStudioCreateInput(ctx, input)
	if err != nil {
		return nil, err
	}

	// Start the transaction and save the studio
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Studio
		studioID, err = qb.Create(ctx, dbInput)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	newStudio, err := r.getStudio(ctx, *studioID)
	if err != nil {
		return nil, fmt.Errorf("finding after create: %w", err)
	}

	if runParentCreateHook {
		r.hookExecutor.ExecutePostHooks(ctx, *newStudio.ParentID, plugin.StudioCreatePost, input, nil)
	} else if runParentUpdateHook {
		r.hookExecutor.ExecutePostHooks(ctx, *newStudio.ParentID, plugin.StudioUpdatePost, input, parentTranslator.getFields())
	}
	r.hookExecutor.ExecutePostHooks(ctx, *studioID, plugin.StudioCreatePost, input, nil)

	return newStudio, nil
}

func studioFromStudioCreateInput(ctx context.Context, input StudioCreateInput) (*models.Studio, error) {
	// Populate a new studio from the input
	newStudio := models.Studio{
		Name: input.Name,
	}
	if input.URL != nil {
		newStudio.URL = *input.URL
	}
	if input.ParentID != nil {
		parentID, _ := strconv.Atoi(*input.ParentID)
		newStudio.ParentID = &parentID
	}
	if input.Details != nil {
		newStudio.Details = *input.Details
	}
	if input.Rating100 != nil {
		newStudio.Rating = input.Rating100
	} else if input.Rating != nil {
		rating := models.Rating5To100(*input.Rating)
		newStudio.Rating = &rating
	}
	if input.IgnoreAutoTag != nil {
		newStudio.IgnoreAutoTag = *input.IgnoreAutoTag
	}

	// Process the base 64 encoded image string
	if input.Image != nil {
		var err error
		newStudio.ImageBytes, err = utils.ProcessImageInput(ctx, *input.Image)
		if err != nil {
			return nil, err
		}
	}

	if input.Aliases != nil {
		newStudio.Aliases = models.NewRelatedStrings(input.Aliases)
	}
	if input.StashIds != nil {
		newStudio.StashIDs = models.NewRelatedStashIDs(stashIDPtrSliceToSlice(input.StashIds))
	}

	return &newStudio, nil
}

func studioPartialFromStudioCreateInput(ctx context.Context, input StudioCreateInput, id *string, translator changesetTranslator) (*models.StudioPartial, error) {
	// Populate studio from the input
	updatedStudio := models.NewStudioPartial()
	updatedStudio.ID, _ = strconv.Atoi(*id)

	if input.Name != "" {
		updatedStudio.Name = translator.optionalString(&input.Name, "name")
	}

	updatedStudio.URL = translator.optionalString(input.URL, "url")

	if input.ParentID != nil {
		parentID, _ := strconv.Atoi(*input.ParentID)
		if parentID > 0 {
			// This is to be set directly as we know it has a value and the translator won't have the field
			updatedStudio.ParentID = models.NewOptionalInt(parentID)
		}
	} else {
		updatedStudio.ParentID = translator.optionalInt(nil, "parent_id")
	}

	updatedStudio.Details = translator.optionalString(input.Details, "details")
	updatedStudio.Rating = translator.ratingConversionOptional(input.Rating, input.Rating100)
	updatedStudio.IgnoreAutoTag = translator.optionalBool(input.IgnoreAutoTag, "ignore_auto_tag")

	// Process the base 64 encoded image string
	if input.Parent.Image != nil {
		updatedStudio.ImageIncluded = true
		var err error
		updatedStudio.ImageBytes, err = utils.ProcessImageInput(ctx, *input.Image)
		if err != nil {
			return nil, err
		}
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

	// Process the base 64 encoded image string
	updatedStudio.ImageIncluded = translator.hasField("image")
	if input.Image != nil {
		var err error
		updatedStudio.ImageBytes, err = utils.ProcessImageInput(ctx, *input.Image)
		if err != nil {
			return nil, err
		}
	}

	return &updatedStudio, nil
}

func (r *mutationResolver) StudioUpdate(ctx context.Context, input StudioUpdateInput) (*models.Studio, error) {
	var updatedStudio *models.Studio
	var dbInput models.StudioDBInput
	var err error
	var parentTranslator changesetTranslator
	runParentCreateHook := false
	runParentUpdateHook := false

	if input.Parent != nil {
		if input.ParentID == nil {
			// The parent needs to be created
			dbInput.ParentCreate, err = studioFromStudioCreateInput(ctx, *input.Parent)
			runParentCreateHook = true
		} else {
			parentTranslator = changesetTranslator{
				inputMap: getNamedUpdateInputMap(ctx, updateInputField+".parent"),
			}

			// The parent studio matched an existing one and the user has chosen in the UI to link and/or update it
			dbInput.ParentUpdate, err = studioPartialFromStudioCreateInput(ctx, *input.Parent, input.ParentID, parentTranslator)
			runParentUpdateHook = true
		}
		if err != nil {
			return nil, err
		}
	}

	translator := changesetTranslator{
		inputMap: getNamedUpdateInputMap(ctx, updateInputField),
	}
	dbInput.StudioUpdate, err = studioPartialFromStudioUpdateInput(ctx, input, &input.ID, translator)
	if err != nil {
		return nil, err
	}

	// Start the transaction and update the studio
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Studio

		if err := manager.ValidateModifyStudio(ctx, *dbInput.StudioUpdate, qb); err != nil {
			return err
		}

		updatedStudio, err = qb.UpdatePartial(ctx, dbInput)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	if runParentCreateHook {
		r.hookExecutor.ExecutePostHooks(ctx, *updatedStudio.ParentID, plugin.StudioCreatePost, input, nil)
	} else if runParentUpdateHook {
		r.hookExecutor.ExecutePostHooks(ctx, *updatedStudio.ParentID, plugin.StudioUpdatePost, input, parentTranslator.getFields())
	}
	r.hookExecutor.ExecutePostHooks(ctx, updatedStudio.ID, plugin.StudioUpdatePost, input, translator.getFields())

	return updatedStudio, nil
}

// This is slightly different to studioPartialFromStudioCreateInput in that Name is handled differently
// and ImageIncluded is not hardcoded to true
func studioPartialFromStudioUpdateInput(ctx context.Context, input StudioUpdateInput, id *string, translator changesetTranslator) (*models.StudioPartial, error) {
	// Populate studio from the input
	updatedStudio := models.NewStudioPartial()
	updatedStudio.ID, _ = strconv.Atoi(*id)

	if input.Name != nil && *input.Name != "" {
		updatedStudio.Name = translator.optionalString(input.Name, "name")
	}

	updatedStudio.URL = translator.optionalString(input.URL, "url")

	if input.ParentID != nil {
		parentID, _ := strconv.Atoi(*input.ParentID)
		if parentID > 0 {
			// This is to be set directly as we know it has a value and the translator won't have the field
			updatedStudio.ParentID = models.NewOptionalInt(parentID)
		}
	} else {
		updatedStudio.ParentID = translator.optionalInt(nil, "parent_id")
	}

	updatedStudio.Details = translator.optionalString(input.Details, "details")
	updatedStudio.Rating = translator.ratingConversionOptional(input.Rating, input.Rating100)
	updatedStudio.IgnoreAutoTag = translator.optionalBool(input.IgnoreAutoTag, "ignore_auto_tag")

	// Process the base 64 encoded image string
	updatedStudio.ImageIncluded = translator.hasField("image")
	if input.Image != nil {
		var err error
		updatedStudio.ImageBytes, err = utils.ProcessImageInput(ctx, *input.Image)
		if err != nil {
			return nil, err
		}
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

	return &updatedStudio, nil
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
