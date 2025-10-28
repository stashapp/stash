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
			newParentStudio := s.Parent.ToStudio(endpoint, nil)
			parentImage, err := s.Parent.GetImage(ctx, nil)
			if err != nil {
				logger.Errorf("Failed to make parent studio from scraped studio %s: %s", s.Parent.Name, err.Error())
				return nil, err
			}

			// Create the studio
			err = w.Create(ctx, newParentStudio)
			if err != nil {
				return nil, err
			}

			// Update image table
			if len(parentImage) > 0 {
				if err := w.UpdateImage(ctx, newParentStudio.ID, parentImage); err != nil {
					return nil, err
				}
			}

			storedId := strconv.Itoa(newParentStudio.ID)
			s.Parent.StoredID = &storedId
		} else {
			// The parent studio matched an existing one and the user has chosen in the UI to link and/or update it
			storedID, _ := strconv.Atoi(*s.Parent.StoredID)

			existingStashIDs, err := w.GetStashIDs(ctx, storedID)
			if err != nil {
				return nil, err
			}

			studioPartial := s.Parent.ToPartial(*s.Parent.StoredID, endpoint, nil, existingStashIDs)
			parentImage, err := s.Parent.GetImage(ctx, nil)
			if err != nil {
				return nil, err
			}

			if err := studio.ValidateModify(ctx, studioPartial, w); err != nil {
				return nil, err
			}

			_, err = w.UpdatePartial(ctx, studioPartial)
			if err != nil {
				return nil, err
			}

			if len(parentImage) > 0 {
				if err := w.UpdateImage(ctx, studioPartial.ID, parentImage); err != nil {
					return nil, err
				}
			}
		}
	}

	newStudio := s.ToStudio(endpoint, nil)
	studioImage, err := s.GetImage(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = w.Create(ctx, newStudio)
	if err != nil {
		return nil, err
	}

	// Update image table
	if len(studioImage) > 0 {
		if err := w.UpdateImage(ctx, newStudio.ID, studioImage); err != nil {
			return nil, err
		}
	}

	return &newStudio.ID, nil
}
