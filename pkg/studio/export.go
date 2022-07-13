package studio

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/utils"
)

type FinderImageStashIDGetter interface {
	Finder
	GetAliases(ctx context.Context, studioID int) ([]string, error)
	GetImage(ctx context.Context, studioID int) ([]byte, error)
	GetStashIDs(ctx context.Context, studioID int) ([]*models.StashID, error)
}

// ToJSON converts a Studio object into its JSON equivalent.
func ToJSON(ctx context.Context, reader FinderImageStashIDGetter, studio *models.Studio) (*jsonschema.Studio, error) {
	newStudioJSON := jsonschema.Studio{
		IgnoreAutoTag: studio.IgnoreAutoTag,
		CreatedAt:     json.JSONTime{Time: studio.CreatedAt.Timestamp},
		UpdatedAt:     json.JSONTime{Time: studio.UpdatedAt.Timestamp},
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
		parent, err := reader.Find(ctx, int(studio.ParentID.Int64))
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

	aliases, err := reader.GetAliases(ctx, studio.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting studio aliases: %v", err)
	}

	newStudioJSON.Aliases = aliases

	image, err := reader.GetImage(ctx, studio.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting studio image: %v", err)
	}

	if len(image) > 0 {
		newStudioJSON.Image = utils.GetBase64StringFromData(image)
	}

	stashIDs, _ := reader.GetStashIDs(ctx, studio.ID)
	var ret []*models.StashID
	for _, stashID := range stashIDs {
		newJoin := &models.StashID{
			StashID:  stashID.StashID,
			Endpoint: stashID.Endpoint,
		}
		ret = append(ret, newJoin)
	}

	newStudioJSON.StashIDs = ret

	return &newStudioJSON, nil
}
