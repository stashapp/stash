package identify

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/models"
)

func createMissingStudio(endpoint string, repo models.Repository, studio *models.ScrapedStudio) (*int64, error) {
	created, err := repo.Studio().Create(scrapedToStudioInput(studio))
	if err != nil {
		return nil, fmt.Errorf("error creating studio: %w", err)
	}

	if endpoint != "" && studio.RemoteSiteID != nil {
		if err := repo.Studio().UpdateStashIDs(created.ID, []models.StashID{
			{
				Endpoint: endpoint,
				StashID:  *studio.RemoteSiteID,
			},
		}); err != nil {
			return nil, fmt.Errorf("error setting studio stash id: %w", err)
		}
	}

	createdID := int64(created.ID)
	return &createdID, nil
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
