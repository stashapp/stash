package identify

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/models"
)

type PerformerCreator interface {
	Create(ctx context.Context, newPerformer *models.Performer) error
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
	performerInput := scrapedToPerformerInput(p)
	err := w.Create(ctx, &performerInput)
	if err != nil {
		return nil, fmt.Errorf("error creating performer: %w", err)
	}

	if endpoint != "" && p.RemoteSiteID != nil {
		if err := w.UpdateStashIDs(ctx, performerInput.ID, []models.StashID{
			{
				Endpoint: endpoint,
				StashID:  *p.RemoteSiteID,
			},
		}); err != nil {
			return nil, fmt.Errorf("error setting performer stash id: %w", err)
		}
	}

	return &performerInput.ID, nil
}

func scrapedToPerformerInput(performer *models.ScrapedPerformer) models.Performer {
	currentTime := time.Now()
	ret := models.Performer{
		Name:      *performer.Name,
		Checksum:  md5.FromString(*performer.Name),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
	if performer.Birthdate != nil {
		d := models.NewDate(*performer.Birthdate)
		ret.Birthdate = &d
	}
	if performer.DeathDate != nil {
		d := models.NewDate(*performer.DeathDate)
		ret.DeathDate = &d
	}
	if performer.Gender != nil {
		ret.Gender = models.GenderEnum(*performer.Gender)
	}
	if performer.Ethnicity != nil {
		ret.Ethnicity = *performer.Ethnicity
	}
	if performer.Country != nil {
		ret.Country = *performer.Country
	}
	if performer.EyeColor != nil {
		ret.EyeColor = *performer.EyeColor
	}
	if performer.HairColor != nil {
		ret.HairColor = *performer.HairColor
	}
	if performer.Height != nil {
		h, err := strconv.Atoi(*performer.Height) // height is stored as an int
		if err == nil {
			ret.Height = &h
		}
	}
	if performer.Weight != nil {
		h, err := strconv.Atoi(*performer.Weight)
		if err == nil {
			ret.Weight = &h
		}
	}
	if performer.Measurements != nil {
		ret.Measurements = *performer.Measurements
	}
	if performer.FakeTits != nil {
		ret.FakeTits = *performer.FakeTits
	}
	if performer.CareerLength != nil {
		ret.CareerLength = *performer.CareerLength
	}
	if performer.Tattoos != nil {
		ret.Tattoos = *performer.Tattoos
	}
	if performer.Piercings != nil {
		ret.Piercings = *performer.Piercings
	}
	if performer.Aliases != nil {
		ret.Aliases = *performer.Aliases
	}
	if performer.Twitter != nil {
		ret.Twitter = *performer.Twitter
	}
	if performer.Instagram != nil {
		ret.Instagram = *performer.Instagram
	}

	return ret
}
