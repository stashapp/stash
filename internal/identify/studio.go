package identify

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type StudioCreator interface {
	Create(ctx context.Context, newStudio models.Studio) (*models.Studio, error)
	UpdateStashIDs(ctx context.Context, studioID int, stashIDs []models.StashID) error
	UpdateImage(ctx context.Context, studioID int, image []byte) error
}

func createMissingStudio(ctx context.Context, endpoint string, w StudioCreator, studio *models.ScrapedStudio) (*int, error) {
	created, err := w.Create(ctx, scrapedToStudioInput(studio))
	if err != nil {
		return nil, fmt.Errorf("error creating studio: %w", err)
	}

	// update image table
	if studio.Image != nil && len(*studio.Image) > 0 {
		imageData, err := utils.ReadImageFromURL(ctx, *studio.Image)
		if err != nil {
			return nil, err
		}

		err = w.UpdateImage(ctx, created.ID, imageData)
		if err != nil {
			return nil, err
		}
	}

	if endpoint != "" && studio.RemoteSiteID != nil {
		if err := w.UpdateStashIDs(ctx, created.ID, []models.StashID{
			{
				Endpoint: endpoint,
				StashID:  *studio.RemoteSiteID,
			},
		}); err != nil {
			return nil, fmt.Errorf("error setting studio stash id: %w", err)
		}
	}

	return &created.ID, nil
}

func scrapedToStudioInput(studio *models.ScrapedStudio) models.Studio {
	currentTime := time.Now()
	ret := models.Studio{
		Name:      sql.NullString{String: studio.Name, Valid: true},
		Checksum:  md5.FromString(studio.Name),
		CreatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
		UpdatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
	}

	if studio.URL != nil {
		ret.URL = sql.NullString{String: *studio.URL, Valid: true}
	}

	return ret
}
