package identify

import (
	"context"
	"fmt"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type StudioCreator interface {
	Create(ctx context.Context, newStudio *models.Studio) error
	UpdateStashIDs(ctx context.Context, studioID int, stashIDs []models.StashID) error
	UpdateImage(ctx context.Context, studioID int, image []byte) error
}

func createMissingStudio(ctx context.Context, endpoint string, w StudioCreator, studio *models.ScrapedStudio) (*int, error) {
	studioInput := scrapedToStudioInput(studio)
	err := w.Create(ctx, &studioInput)
	if err != nil {
		return nil, fmt.Errorf("error creating studio: %w", err)
	}

	// update image table
	if studio.Images != nil && len(studio.Images) > 0 {
		imageData, err := utils.ReadImageFromURL(ctx, *studio.Images[0])
		if err != nil {
			return nil, err
		}

		err = w.UpdateImage(ctx, studioInput.ID, imageData)
		if err != nil {
			return nil, err
		}
	}

	if endpoint != "" && studio.RemoteSiteID != nil {
		if err := w.UpdateStashIDs(ctx, studioInput.ID, []models.StashID{
			{
				Endpoint: endpoint,
				StashID:  *studio.RemoteSiteID,
			},
		}); err != nil {
			return nil, fmt.Errorf("error setting studio stash id: %w", err)
		}
	}

	return &studioInput.ID, nil
}

func scrapedToStudioInput(studio *models.ScrapedStudio) models.Studio {
	currentTime := time.Now()
	ret := models.Studio{
		Name:      studio.Name,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	if studio.URLS != nil {
		ret.URL = *studio.URLS[0]
	}

	return ret
}
