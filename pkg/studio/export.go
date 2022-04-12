package studio

import (
	"fmt"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/utils"
)

// ToJSON converts a Studio object into its JSON equivalent.
func ToJSON(reader models.StudioReader, studio *models.Studio) (*jsonschema.Studio, error) {
	newStudioJSON := jsonschema.Studio{
		IgnoreAutoTag: studio.IgnoreAutoTag,
		CreatedAt:     models.JSONTime{Time: studio.CreatedAt.Timestamp},
		UpdatedAt:     models.JSONTime{Time: studio.UpdatedAt.Timestamp},
	}

	if studio.Name.Valid {
		newStudioJSON.Name = studio.Name.String
	}

	if studio.URL.Valid {
		newStudioJSON.URL = studio.URL.String
	}

	if studio.Details.Valid {
		newStudioJSON.Details = studio.Details.String
	}

	if studio.ParentID.Valid {
		parent, err := reader.Find(int(studio.ParentID.Int64))
		if err != nil {
			return nil, fmt.Errorf("error getting parent studio: %v", err)
		}

		if parent != nil {
			newStudioJSON.ParentStudio = parent.Name.String
		}
	}

	if studio.Rating.Valid {
		newStudioJSON.Rating = int(studio.Rating.Int64)
	}

	aliases, err := reader.GetAliases(studio.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting studio aliases: %v", err)
	}

	newStudioJSON.Aliases = aliases

	image, err := reader.GetImage(studio.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting studio image: %v", err)
	}

	if len(image) > 0 {
		newStudioJSON.Image = utils.GetBase64StringFromData(image)
	}

	stashIDs, _ := reader.GetStashIDs(studio.ID)
	var ret []models.StashID
	for _, stashID := range stashIDs {
		newJoin := models.StashID{
			StashID:  stashID.StashID,
			Endpoint: stashID.Endpoint,
		}
		ret = append(ret, newJoin)
	}

	newStudioJSON.StashIDs = ret

	return &newStudioJSON, nil
}
