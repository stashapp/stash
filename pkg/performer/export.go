package performer

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/utils"
)

type ImageStashIDGetter interface {
	GetImage(ctx context.Context, performerID int) ([]byte, error)
	GetStashIDs(ctx context.Context, performerID int) ([]*models.StashID, error)
}

// ToJSON converts a Performer object into its JSON equivalent.
func ToJSON(ctx context.Context, reader ImageStashIDGetter, performer *models.Performer) (*jsonschema.Performer, error) {
	newPerformerJSON := jsonschema.Performer{
		IgnoreAutoTag: performer.IgnoreAutoTag,
		CreatedAt:     json.JSONTime{Time: performer.CreatedAt.Timestamp},
		UpdatedAt:     json.JSONTime{Time: performer.UpdatedAt.Timestamp},
	}

	if performer.Name.Valid {
		newPerformerJSON.Name = performer.Name.String
	}
	if performer.Gender.Valid {
		newPerformerJSON.Gender = performer.Gender.String
	}
	if performer.URL.Valid {
		newPerformerJSON.URL = performer.URL.String
	}
	if performer.Birthdate.Valid {
		newPerformerJSON.Birthdate = utils.GetYMDFromDatabaseDate(performer.Birthdate.String)
	}
	if performer.Ethnicity.Valid {
		newPerformerJSON.Ethnicity = performer.Ethnicity.String
	}
	if performer.Country.Valid {
		newPerformerJSON.Country = performer.Country.String
	}
	if performer.EyeColor.Valid {
		newPerformerJSON.EyeColor = performer.EyeColor.String
	}
	if performer.Height.Valid {
		newPerformerJSON.Height = performer.Height.String
	}
	if performer.Measurements.Valid {
		newPerformerJSON.Measurements = performer.Measurements.String
	}
	if performer.FakeTits.Valid {
		newPerformerJSON.FakeTits = performer.FakeTits.String
	}
	if performer.CareerLength.Valid {
		newPerformerJSON.CareerLength = performer.CareerLength.String
	}
	if performer.Tattoos.Valid {
		newPerformerJSON.Tattoos = performer.Tattoos.String
	}
	if performer.Piercings.Valid {
		newPerformerJSON.Piercings = performer.Piercings.String
	}
	if performer.Aliases.Valid {
		newPerformerJSON.Aliases = performer.Aliases.String
	}
	if performer.Twitter.Valid {
		newPerformerJSON.Twitter = performer.Twitter.String
	}
	if performer.Instagram.Valid {
		newPerformerJSON.Instagram = performer.Instagram.String
	}
	if performer.Favorite.Valid {
		newPerformerJSON.Favorite = performer.Favorite.Bool
	}
	if performer.Rating.Valid {
		newPerformerJSON.Rating = int(performer.Rating.Int64)
	}
	if performer.Details.Valid {
		newPerformerJSON.Details = performer.Details.String
	}
	if performer.DeathDate.Valid {
		newPerformerJSON.DeathDate = utils.GetYMDFromDatabaseDate(performer.DeathDate.String)
	}
	if performer.HairColor.Valid {
		newPerformerJSON.HairColor = performer.HairColor.String
	}
	if performer.Weight.Valid {
		newPerformerJSON.Weight = int(performer.Weight.Int64)
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
		if performer.Name.Valid {
			results = append(results, performer.Name.String)
		}
	}

	return results
}
