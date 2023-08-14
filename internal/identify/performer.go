package identify

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/utils"
)

type PerformerCreator interface {
	Create(ctx context.Context, newPerformer *models.Performer) error
	UpdateImage(ctx context.Context, performerID int, image []byte) error
}

func getPerformerID(ctx context.Context, endpoint string, w PerformerCreator, p *models.ScrapedPerformer, createMissing bool, skipSingleNamePerformers bool) (*int, error) {
	if p.StoredID != nil {
		// existing performer, just add it
		performerID, err := strconv.Atoi(*p.StoredID)
		if err != nil {
			return nil, fmt.Errorf("error converting performer ID %s: %w", *p.StoredID, err)
		}

		return &performerID, nil
	} else if createMissing && p.Name != nil { // name is mandatory
		// skip single name performers with no disambiguation
		if skipSingleNamePerformers && !strings.Contains(*p.Name, " ") && (p.Disambiguation == nil || len(*p.Disambiguation) == 0) {
			return nil, ErrSkipSingleNamePerformer
		}
		return createMissingPerformer(ctx, endpoint, w, p)
	}

	return nil, nil
}

func createMissingPerformer(ctx context.Context, endpoint string, w PerformerCreator, p *models.ScrapedPerformer) (*int, error) {
	performerInput := scrapedToPerformerInput(p)
	if endpoint != "" && p.RemoteSiteID != nil {
		performerInput.StashIDs = models.NewRelatedStashIDs([]models.StashID{
			{
				Endpoint: endpoint,
				StashID:  *p.RemoteSiteID,
			},
		})
	}

	err := w.Create(ctx, &performerInput)
	if err != nil {
		return nil, fmt.Errorf("error creating performer: %w", err)
	}

	// update image table
	if p.Image != nil && len(*p.Image) > 0 {
		imageData, err := utils.ReadImageFromURL(ctx, *p.Image)
		if err != nil {
			return nil, err
		}

		err = w.UpdateImage(ctx, performerInput.ID, imageData)
		if err != nil {
			return nil, err
		}
	}

	return &performerInput.ID, nil
}

func scrapedToPerformerInput(performer *models.ScrapedPerformer) models.Performer {
	currentTime := time.Now()
	ret := models.Performer{
		Name:      *performer.Name,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
	if performer.Disambiguation != nil {
		ret.Disambiguation = *performer.Disambiguation
	}
	if performer.Birthdate != nil {
		d, err := models.ParseDate(*performer.Birthdate)
		if err == nil {
			ret.Birthdate = &d
		}
	}
	if performer.DeathDate != nil {
		d, err := models.ParseDate(*performer.DeathDate)
		if err == nil {
			ret.DeathDate = &d
		}
	}
	if performer.Gender != nil {
		v := models.GenderEnum(*performer.Gender)
		ret.Gender = &v
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
	if performer.PenisLength != nil {
		h, err := strconv.ParseFloat(*performer.PenisLength, 64)
		if err == nil {
			ret.PenisLength = &h
		}
	}
	if performer.Circumcised != nil {
		v := models.CircumisedEnum(*performer.Circumcised)
		ret.Circumcised = &v
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
		ret.Aliases = models.NewRelatedStrings(stringslice.FromString(*performer.Aliases, ","))
	}
	if performer.Twitter != nil {
		ret.Twitter = *performer.Twitter
	}
	if performer.Instagram != nil {
		ret.Instagram = *performer.Instagram
	}
	if performer.URL != nil {
		ret.URL = *performer.URL
	}
	if performer.Details != nil {
		ret.Details = *performer.Details
	}

	return ret
}
