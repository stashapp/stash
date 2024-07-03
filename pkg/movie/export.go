package movie

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/utils"
)

type ImageGetter interface {
	GetFrontImage(ctx context.Context, movieID int) ([]byte, error)
	GetBackImage(ctx context.Context, movieID int) ([]byte, error)
}

// ToJSON converts a Movie into its JSON equivalent.
func ToJSON(ctx context.Context, reader ImageGetter, studioReader models.StudioGetter, movie *models.Group) (*jsonschema.Group, error) {
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

	if movie.StudioID != nil {
		studio, err := studioReader.Find(ctx, *movie.StudioID)
		if err != nil {
			return nil, fmt.Errorf("error getting movie studio: %v", err)
		}

		if studio != nil {
			newMovieJSON.Studio = studio.Name
		}
	}

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
