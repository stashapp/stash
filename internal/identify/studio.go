package identify

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type StudioCreator interface {
	Create(ctx context.Context, input models.StudioDBInput) (*int, error)
}

func createMissingStudio(ctx context.Context, endpoint string, w StudioCreator, studio *models.ScrapedStudio) (*int, error) {
	var dbInput models.StudioDBInput
	var err error

	if studio.Parent != nil {
		if studio.Parent.StoredID == nil {
			// The parent needs to be created
			dbInput.ParentCreate, err = studioFromScrapedStudio(ctx, studio.Parent, endpoint)
			if err != nil {
				return nil, err
			}
		} else {
			// The parent studio matched an existing one and the user has chosen in the UI to link and/or update it
			dbInput.ParentUpdate, err = studioPartialFromScrapedStudio(ctx, studio.Parent, studio.Parent.StoredID, endpoint)
			if err != nil {
				return nil, err
			}
		}
	}

	dbInput.StudioCreate, err = studioFromScrapedStudio(ctx, studio, endpoint)
	if err != nil {
		return nil, err
	}

	studioID, err := w.Create(ctx, dbInput)
	if err != nil {
		return nil, err
	}

	return studioID, nil
}

// Duplicated in task_stash_box_tag.go
func studioFromScrapedStudio(ctx context.Context, input *models.ScrapedStudio, endpoint string) (*models.Studio, error) {
	// Populate a new studio from the input
	newStudio := models.Studio{
		Name: input.Name,
		StashIDs: models.NewRelatedStashIDs([]models.StashID{
			{
				Endpoint: endpoint,
				StashID:  *input.RemoteSiteID,
			},
		}),
	}

	if input.URL != nil {
		newStudio.URL = *input.URL
	}

	if input.Parent != nil && input.Parent.StoredID != nil {
		parentId, _ := strconv.Atoi(*input.Parent.StoredID)
		newStudio.ParentID = &parentId
	}

	// Process the base 64 encoded image string
	if input.Image != nil {
		var err error
		newStudio.ImageBytes, err = utils.ProcessImageInput(ctx, *input.Image)
		if err != nil {
			return nil, err
		}
	}

	return &newStudio, nil
}

// Duplicated in task_stash_box_tag.go
func studioPartialFromScrapedStudio(ctx context.Context, input *models.ScrapedStudio, id *string, endpoint string) (*models.StudioPartial, error) {
	partial := models.NewStudioPartial()
	partial.ID, _ = strconv.Atoi(*id)

	if input.Name != "" {
		partial.Name = models.NewOptionalString(input.Name)

	}

	if input.URL != nil {
		partial.URL = models.NewOptionalString(*input.URL)
	}

	if input.Parent != nil {
		if input.Parent.StoredID != nil {
			parentID, _ := strconv.Atoi(*input.Parent.StoredID)
			if parentID > 0 {
				// This is to be set directly as we know it has a value and the translator won't have the field
				partial.ParentID = models.NewOptionalInt(parentID)
			}
		}
	} else {
		partial.ParentID = models.NewOptionalIntPtr(nil)
	}

	// Process the base 64 encoded image string
	if len(input.Images) > 0 {
		partial.ImageIncluded = true
		var err error
		partial.ImageBytes, err = utils.ProcessImageInput(ctx, input.Images[0])
		if err != nil {
			return nil, err
		}
	}

	partial.StashIDs = &models.UpdateStashIDs{
		StashIDs: []models.StashID{
			{
				Endpoint: endpoint,
				StashID:  *input.RemoteSiteID,
			},
		},
		Mode: models.RelationshipUpdateModeSet,
	}

	return &partial, nil
}
