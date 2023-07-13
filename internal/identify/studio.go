package identify

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/studio"
)

func createMissingStudio(ctx context.Context, endpoint string, w models.StudioReaderWriter, s *models.ScrapedStudio) (*int, error) {
	var err error

	if s.Parent != nil {
		if s.Parent.StoredID == nil {
			// The parent needs to be created
			newParentStudio, err := s.Parent.ToStudio(ctx, endpoint, nil)
			if err != nil {
				logger.Errorf("Failed to make parent studio from scraped studio %s: %s", s.Parent.Name, err.Error())
				return nil, err
			}

			// Create the studio
			err = w.Create(ctx, newParentStudio)
			if err != nil {
				return nil, err
			}
			storedId := strconv.Itoa(newParentStudio.ID)
			s.Parent.StoredID = &storedId
		} else {
			// The parent studio matched an existing one and the user has chosen in the UI to link and/or update it
			existingStashIDs := getStashIDsForStudio(ctx, *s.Parent.StoredID, w)
			studioPartial, err := s.Parent.ToPartial(ctx, s.Parent.StoredID, endpoint, nil, existingStashIDs)
			if err != nil {
				return nil, err
			}

			if err := studio.ValidateModify(ctx, *studioPartial, w); err != nil {
				return nil, err
			}

			_, err = w.UpdatePartial(ctx, *studioPartial)
			if err != nil {
				return nil, err
			}
		}
	}

	newStudio, err := s.ToStudio(ctx, endpoint, nil)
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
