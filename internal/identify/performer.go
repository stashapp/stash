package identify

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/models"
)

type PerformerCreator interface {
	Create(ctx context.Context, newPerformer models.Performer) (*models.Performer, error)
	UpdateStashIDs(ctx context.Context, performerID int, stashIDs []models.StashID) error
}

func getPerformerID(ctx context.Context, endpoint string, w PerformerCreator, p *models.ScrapedPerformer, createMissing bool) (*int, error) {
	if p.StoredID != nil {
		// existing performer, just add it
		performerID, err := strconv.Atoi(*p.StoredID)
		if err != nil {
			return nil, fmt.Errorf("error converting performer ID %s: %w", *p.StoredID, err)
		}

		return &performerID, nil
	} else if createMissing && p.Name != nil { // name is mandatory
		return createMissingPerformer(ctx, endpoint, w, p)
	}

	return nil, nil
}

func createMissingPerformer(ctx context.Context, endpoint string, w PerformerCreator, p *models.ScrapedPerformer) (*int, error) {
	created, err := w.Create(ctx, scrapedToPerformerInput(p))
	if err != nil {
		return nil, fmt.Errorf("error creating performer: %w", err)
	}

	if endpoint != "" && p.RemoteSiteID != nil {
		if err := w.UpdateStashIDs(ctx, created.ID, []models.StashID{
			{
				Endpoint: endpoint,
				StashID:  *p.RemoteSiteID,
			},
		}); err != nil {
			return nil, fmt.Errorf("error setting performer stash id: %w", err)
		}
	}

	return &created.ID, nil
}

func scrapedToPerformerInput(performer *models.ScrapedPerformer) models.Performer {
	currentTime := time.Now()
	ret := models.Performer{
		Name:      sql.NullString{String: *performer.Name, Valid: true},
		Checksum:  md5.FromString(*performer.Name),
		CreatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
		UpdatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
		Favorite:  sql.NullBool{Bool: false, Valid: true},
	}
	if performer.Birthdate != nil {
		ret.Birthdate = models.SQLiteDate{String: *performer.Birthdate, Valid: true}
	}
	if performer.DeathDate != nil {
		ret.DeathDate = models.SQLiteDate{String: *performer.DeathDate, Valid: true}
	}
	if performer.Gender != nil {
		ret.Gender = sql.NullString{String: *performer.Gender, Valid: true}
	}
	if performer.Ethnicity != nil {
		ret.Ethnicity = sql.NullString{String: *performer.Ethnicity, Valid: true}
	}
	if performer.Country != nil {
		ret.Country = sql.NullString{String: *performer.Country, Valid: true}
	}
	if performer.EyeColor != nil {
		ret.EyeColor = sql.NullString{String: *performer.EyeColor, Valid: true}
	}
	if performer.HairColor != nil {
		ret.HairColor = sql.NullString{String: *performer.HairColor, Valid: true}
	}
	if performer.Height != nil {
		ret.Height = sql.NullString{String: *performer.Height, Valid: true}
	}
	if performer.Measurements != nil {
		ret.Measurements = sql.NullString{String: *performer.Measurements, Valid: true}
	}
	if performer.FakeTits != nil {
		ret.FakeTits = sql.NullString{String: *performer.FakeTits, Valid: true}
	}
	if performer.CareerLength != nil {
		ret.CareerLength = sql.NullString{String: *performer.CareerLength, Valid: true}
	}
	if performer.Tattoos != nil {
		ret.Tattoos = sql.NullString{String: *performer.Tattoos, Valid: true}
	}
	if performer.Piercings != nil {
		ret.Piercings = sql.NullString{String: *performer.Piercings, Valid: true}
	}
	if performer.Aliases != nil {
		ret.Aliases = sql.NullString{String: *performer.Aliases, Valid: true}
	}
	if performer.Twitter != nil {
		ret.Twitter = sql.NullString{String: *performer.Twitter, Valid: true}
	}
	if performer.Instagram != nil {
		ret.Instagram = sql.NullString{String: *performer.Instagram, Valid: true}
	}

	return ret
}
