package performer

import (
	"context"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/utils"
)

type ImageStashIDGetter interface {
	GetImage(ctx context.Context, performerID int) ([]byte, error)
	models.StashIDLoader
}

// ToJSON converts a Performer object into its JSON equivalent.
func ToJSON(ctx context.Context, reader ImageStashIDGetter, performer *models.Performer) (*jsonschema.Performer, error) {
	newPerformerJSON := jsonschema.Performer{
		Name:          performer.Name,
		Gender:        performer.Gender.String(),
		URL:           performer.URL,
		Ethnicity:     performer.Ethnicity,
		Country:       performer.Country,
		EyeColor:      performer.EyeColor,
		Measurements:  performer.Measurements,
		FakeTits:      performer.FakeTits,
		CareerLength:  performer.CareerLength,
		Tattoos:       performer.Tattoos,
		Piercings:     performer.Piercings,
		Aliases:       performer.Aliases,
		Twitter:       performer.Twitter,
		Instagram:     performer.Instagram,
		Favorite:      performer.Favorite,
		Details:       performer.Details,
		HairColor:     performer.HairColor,
		IgnoreAutoTag: performer.IgnoreAutoTag,
		CreatedAt:     json.JSONTime{Time: performer.CreatedAt},
		UpdatedAt:     json.JSONTime{Time: performer.UpdatedAt},
	}

	if performer.Birthdate != nil {
		newPerformerJSON.Birthdate = performer.Birthdate.String()
	}
	if performer.Rating != nil {
		newPerformerJSON.Rating = *performer.Rating
	}
	if performer.DeathDate != nil {
		newPerformerJSON.DeathDate = performer.DeathDate.String()
	}

	if performer.Height != nil {
		newPerformerJSON.Height = strconv.Itoa(*performer.Height)
	}

	if performer.Weight != nil {
		newPerformerJSON.Weight = *performer.Weight
	}

	image, err := reader.GetImage(ctx, performer.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting performers image: %v", err)
	}

	if len(image) > 0 {
		newPerformerJSON.Image = utils.GetBase64StringFromData(image)
	}

	stashIDs, _ := reader.GetStashIDs(ctx, performer.ID)
	var ret []models.StashID
	for _, stashID := range stashIDs {
		newJoin := models.StashID{
			StashID:  stashID.StashID,
			Endpoint: stashID.Endpoint,
		}
		ret = append(ret, newJoin)
	}

	newPerformerJSON.StashIDs = ret

	return &newPerformerJSON, nil
}

func GetIDs(performers []*models.Performer) []int {
	var results []int
	for _, performer := range performers {
		results = append(results, performer.ID)
	}

	return results
}

func GetNames(performers []*models.Performer) []string {
	var results []string
	for _, performer := range performers {
		if performer.Name != "" {
			results = append(results, performer.Name)
		}
	}

	return results
}
