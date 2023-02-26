package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/match"
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
	// Parent studio data is being passed in, but there is no local ID, so it needs to be created or matched first
	if input.ParentID == nil && input.Parent != nil {
		// If the parent studio matched an existing one, the user has chosen in the UI to link and update it
		if err := r.withTxn(ctx, func(ctx context.Context) error {
			st := &models.ScrapedStudio{
				Name:         input.Parent.Name,
				RemoteSiteID: &input.Parent.StashIds[0].StashID,
			}

			err := match.ScrapedStudio(ctx, r.repository.Studio, st, &input.Parent.StashIds[0].Endpoint)
			if err != nil {
				return err
			}

			// Found the local ID for the studio to link and update
			input.ParentID = st.StoredID
			return nil
		}); err != nil {
			return nil, err
		}

		if input.ParentID == nil {
			// Create
			studioID, err := r.studioCreate(ctx, *input.Parent)
			if err != nil {
				return nil, err
			}

			// Assign the new parent studio ID so the child studio will be linked
			id_as_string := strconv.Itoa(studioID)
			input.ParentID = &id_as_string
		} else {
			// Update
			su := &StudioUpdateInput{
				ID:       *input.ParentID,
				Name:     &input.Parent.Name,
				URL:      input.Parent.URL,
				Image:    input.Parent.Image,
				StashIds: input.Parent.StashIds,
			}

			_, err := r.StudioUpdate(ctx, *su)
			if err != nil {
				return nil, err
			}
		}
	}

	// Now create the main studio
	studioID, err := r.studioCreate(ctx, input)
	if err != nil {
		return nil, err
	}

	return r.getStudio(ctx, studioID)
}

func (r *mutationResolver) studioCreate(ctx context.Context, input StudioCreateInput) (int, error) {
	var imageData []byte
	var err error

	// Process the base 64 encoded image string
	if input.Image != nil {
		imageData, err = utils.ProcessImageInput(ctx, *input.Image)
		if err != nil {
			return 0, err
		}
	}

	// Populate a new studio from the input
	newStudio := models.Studio{
		Name: input.Name,
	}
	if input.URL != nil {
		newStudio.URL = *input.URL
	}
	if input.ParentID != nil {
		parentID, _ := strconv.ParseInt(*input.ParentID, 10, 32)
		var parentID_32 = int(parentID)
		newStudio.ParentID = &parentID_32
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

	if input.Aliases != nil {
		newStudio.Aliases = models.NewRelatedStrings(input.Aliases)
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

		return nil
	}); err != nil {
		return 0, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, newStudio.ID, plugin.StudioCreatePost, input, nil)

	return newStudio.ID, nil
}

func (r *mutationResolver) StudioUpdate(ctx context.Context, input StudioUpdateInput) (*models.Studio, error) {
	// Populate studio from the input
	studioID, _ := strconv.Atoi(input.ID)
	updatedStudio := models.NewStudioPartial()
	updatedStudio.ID = studioID

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	var imageData []byte
	imageIncluded := translator.hasField("image")
	if input.Image != nil {
		var err error
		imageData, err = utils.ProcessImageInput(ctx, *input.Image)
		if err != nil {
			return nil, err
		}
	}

	if input.Name != nil {
		// generate checksum from studio name rather than image
		checksum := md5.FromString(*input.Name)
		updatedStudio.Name = translator.optionalString(input.Name, "name")
		updatedStudio.Checksum = translator.optionalString(&checksum, "checksum")
	}

	updatedStudio.URL = translator.optionalString(input.URL, "url")

	if input.ParentID != nil {
		var parentIDTemp, _ = strconv.Atoi(*input.ParentID)
		updatedStudio.ParentID = translator.optionalInt(&parentIDTemp, "parent_id")
	} else {
		updatedStudio.ParentID = translator.optionalInt(nil, "parent_id")
	}

	updatedStudio.Details = translator.optionalString(input.Details, "details")
	updatedStudio.Rating = translator.ratingConversionOptional(input.Rating, input.Rating100)
	updatedStudio.IgnoreAutoTag = translator.optionalBool(input.IgnoreAutoTag, "ignore_auto_tag")

	if translator.hasField("aliases") {
		updatedStudio.Aliases = &models.UpdateStrings{
			Values: input.Aliases,
			Mode:   models.RelationshipUpdateModeSet,
		}
	}

	// Save the stash_ids
	if translator.hasField("stash_ids") {
		updatedStudio.StashIDs = &models.UpdateStashIDs{
			StashIDs: stashIDPtrSliceToSlice(input.StashIds),
			Mode:     models.RelationshipUpdateModeSet,
		}
	}

	// Start the transaction and save the studio
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Studio

		if err := manager.ValidateModifyStudio(ctx, updatedStudio, qb); err != nil {
			return err
		}

		var err error
		_, err = qb.UpdatePartial(ctx, studioID, updatedStudio)
		if err != nil {
			return err
		}

		// update image table
		if len(imageData) > 0 {
			if err := qb.UpdateImage(ctx, studioID, imageData); err != nil {
				return err
			}
		} else if imageIncluded {
			// must be unsetting
			if err := qb.DestroyImage(ctx, studioID); err != nil {
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
