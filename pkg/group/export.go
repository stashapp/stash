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

// GetStudioNames returns the names of studios associated with the group.
func GetStudioNames(ctx context.Context, studioReader models.StudioGetter, studioIDLoader models.StudioIDLoader, group *models.Group) ([]string, error) {
	if err := group.LoadStudioIDs(ctx, studioIDLoader); err != nil {
		return nil, fmt.Errorf("error loading studio IDs: %v", err)
	}

	var names []string
	for _, studioID := range group.StudioIDs.List() {
		studio, err := studioReader.Find(ctx, studioID)
		if err != nil {
			return nil, fmt.Errorf("error getting studio: %v", err)
		}

		if studio != nil {
			names = append(names, studio.Name)
		}
	}

	return names, nil
}

type ImageGetter interface {
	GetFrontImage(ctx context.Context, movieID int) ([]byte, error)
	GetBackImage(ctx context.Context, movieID int) ([]byte, error)
}

// ToJSON converts a Movie into its JSON equivalent.
func ToJSON(ctx context.Context, reader ImageGetter, studioReader models.StudioGetter, studioIDLoader models.StudioIDLoader, movie *models.Group) (*jsonschema.Group, error) {
	newMovieJSON := jsonschema.Group{
		Name:      movie.Name,
		Aliases:   movie.Aliases,
		Director:  movie.Director,
		Synopsis:  movie.Synopsis,
		URLs:      movie.URLs.List(),
		CreatedAt: json.JSONTime{Time: movie.CreatedAt},
		UpdatedAt: json.JSONTime{Time: movie.UpdatedAt},
	}

	if movie.Date != nil {
		newMovieJSON.Date = movie.Date.String()
	}
	if movie.Rating != nil {
		newMovieJSON.Rating = *movie.Rating
	}
	if movie.Duration != nil {
		newMovieJSON.Duration = *movie.Duration
	}

	studioNames, err := GetStudioNames(ctx, studioReader, studioIDLoader, movie)
	if err != nil {
		return nil, fmt.Errorf("error getting movie studio names: %v", err)
	}
	newMovieJSON.Studios = studioNames

	frontImage, err := reader.GetFrontImage(ctx, movie.ID)
	if err != nil {
		logger.Errorf("Error getting movie front image: %v", err)
	}

	if len(frontImage) > 0 {
		newMovieJSON.FrontImage = utils.GetBase64StringFromData(frontImage)
	}

	backImage, err := reader.GetBackImage(ctx, movie.ID)
	if err != nil {
		logger.Errorf("Error getting movie back image: %v", err)
	}

	if len(backImage) > 0 {
		newMovieJSON.BackImage = utils.GetBase64StringFromData(backImage)
	}

	return &newMovieJSON, nil
}
