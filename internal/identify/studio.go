package identify

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

func createMissingStudio(ctx context.Context, endpoint string, w models.StudioReaderWriter, studio *models.ScrapedStudio) (*int, error) {
	var err error

	if studio.Parent != nil {
		if studio.Parent.StoredID == nil {
			// The parent needs to be created
			newParentStudio, err := studio.Parent.ToStudio(ctx, endpoint, nil)
			if err != nil {
				logger.Errorf("Failed to make parent studio from scraped studio %s: %s", studio.Parent.Name, err.Error())
				return nil, err
			}

			// Create the studio
			err = w.Create(ctx, newParentStudio)
			if err != nil {
				return nil, err
			}
			storedId := strconv.Itoa(newParentStudio.ID)
			studio.Parent.StoredID = &storedId
		} else {
			// The parent studio matched an existing one and the user has chosen in the UI to link and/or update it
			existingStashIDs := getStashIDsForStudio(ctx, *studio.Parent.StoredID, w)
			studioPartial, err := studio.Parent.ToPartial(ctx, studio.Parent.StoredID, endpoint, nil, existingStashIDs)
			if err != nil {
				return nil, err
			}

			if err := studioPartial.ValidateModifyStudio(ctx, w); err != nil {
				return nil, err
			}

			_, err = w.UpdatePartial(ctx, *studioPartial)
			if err != nil {
				return nil, err
			}
		}
	}

	newStudio, err := studio.ToStudio(ctx, endpoint, nil)
	if err != nil {
		return nil, err
	}

	err = w.Create(ctx, newStudio)
	if err != nil {
		return nil, err
	}

	return &newStudio.ID, nil
}

func getStashIDsForStudio(ctx context.Context, studioID string, w models.StudioReaderWriter) []models.StashID {
	id, _ := strconv.Atoi(studioID)
	tempStudio := &models.Studio{ID: id}

	err := tempStudio.LoadStashIDs(ctx, w)
	if err != nil {
		return nil
	}
	return tempStudio.StashIDs.List()
}
