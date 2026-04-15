package group

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/utils"
)

type GroupExportReader interface {
	GetFrontImage(ctx context.Context, groupID int) ([]byte, error)
	GetBackImage(ctx context.Context, groupID int) ([]byte, error)
	GetCustomFields(ctx context.Context, groupID int) (map[string]interface{}, error)
}

// ToJSON converts a Group into its JSON equivalent.
func ToJSON(ctx context.Context, reader GroupExportReader, studioReader models.StudioGetter, group *models.Group) (*jsonschema.Group, error) {
	newGroupJSON := jsonschema.Group{
		Name:      group.Name,
		Aliases:   group.Aliases,
		Director:  group.Director,
		Synopsis:  group.Synopsis,
		URLs:      group.URLs.List(),
		CreatedAt: json.JSONTime{Time: group.CreatedAt},
		UpdatedAt: json.JSONTime{Time: group.UpdatedAt},
	}

	if group.Date != nil {
		newGroupJSON.Date = group.Date.String()
	}
	if group.Rating != nil {
		newGroupJSON.Rating = *group.Rating
	}
	if group.Duration != nil {
		newGroupJSON.Duration = *group.Duration
	}

	if group.StudioID != nil {
		studio, err := studioReader.Find(ctx, *group.StudioID)
		if err != nil {
			return nil, fmt.Errorf("error getting movie studio: %v", err)
		}

		if studio != nil {
			newGroupJSON.Studio = studio.Name
		}
	}

	frontImage, err := reader.GetFrontImage(ctx, group.ID)
	if err != nil {
		logger.Errorf("Error getting movie front image: %v", err)
	}

	if len(frontImage) > 0 {
		newGroupJSON.FrontImage = utils.GetBase64StringFromData(frontImage)
	}

	backImage, err := reader.GetBackImage(ctx, group.ID)
	if err != nil {
		logger.Errorf("Error getting movie back image: %v", err)
	}

	if len(backImage) > 0 {
		newGroupJSON.BackImage = utils.GetBase64StringFromData(backImage)
	}

	newGroupJSON.CustomFields, err = reader.GetCustomFields(ctx, group.ID)
	if err != nil {
		return nil, fmt.Errorf("getting group custom fields: %v", err)
	}

	return &newGroupJSON, nil
}
