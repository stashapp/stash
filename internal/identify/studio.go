package identify

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

type StudioCreator interface {
	Create(ctx context.Context, newStudio *models.Studio) error
}

func createMissingStudio(ctx context.Context, endpoint string, w StudioCreator, studio *models.ScrapedStudio) (*int, error) {
	studioInput := scrapedToStudioInput(studio)
	if endpoint != "" && studio.RemoteSiteID != nil {
		studioInput.StashIDs = models.NewRelatedStashIDs([]models.StashID{
			{
				Endpoint: endpoint,
				StashID:  *studio.RemoteSiteID,
			},
		})
	}

	err := w.Create(ctx, &studioInput)
	if err != nil {
		return nil, fmt.Errorf("error creating studio: %w", err)
	}

	return &studioInput.ID, nil
}

func scrapedToStudioInput(studio *models.ScrapedStudio) models.Studio {
	ret := models.Studio{
		Name: studio.Name,
	}

	if studio.URL != nil {
		ret.URL = *studio.URL
	}

	return ret
}
