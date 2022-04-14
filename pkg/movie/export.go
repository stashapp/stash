package movie

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/studio"
	"github.com/stashapp/stash/pkg/utils"
)

type ImageGetter interface {
	GetFrontImage(ctx context.Context, movieID int) ([]byte, error)
	GetBackImage(ctx context.Context, movieID int) ([]byte, error)
}

// ToJSON converts a Movie into its JSON equivalent.
func ToJSON(ctx context.Context, reader ImageGetter, studioReader studio.Finder, movie *models.Movie) (*jsonschema.Movie, error) {
	newMovieJSON := jsonschema.Movie{
		CreatedAt: json.JSONTime{Time: movie.CreatedAt.Timestamp},
		UpdatedAt: json.JSONTime{Time: movie.UpdatedAt.Timestamp},
	}

	if movie.Name.Valid {
		newMovieJSON.Name = movie.Name.String
	}
	if movie.Aliases.Valid {
		newMovieJSON.Aliases = movie.Aliases.String
	}
	if movie.Date.Valid {
		newMovieJSON.Date = utils.GetYMDFromDatabaseDate(movie.Date.String)
	}
	if movie.Rating.Valid {
		newMovieJSON.Rating = int(movie.Rating.Int64)
	}
	if movie.Duration.Valid {
		newMovieJSON.Duration = int(movie.Duration.Int64)
	}

	if movie.Director.Valid {
		newMovieJSON.Director = movie.Director.String
	}

	if movie.Synopsis.Valid {
		newMovieJSON.Synopsis = movie.Synopsis.String
	}

	if movie.URL.Valid {
		newMovieJSON.URL = movie.URL.String
	}

	if movie.StudioID.Valid {
		studio, err := studioReader.Find(ctx, int(movie.StudioID.Int64))
		if err != nil {
			return nil, fmt.Errorf("error getting movie studio: %v", err)
		}

		if studio != nil {
			newMovieJSON.Studio = studio.Name.String
		}
	}

	frontImage, err := reader.GetFrontImage(ctx, movie.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting movie front image: %v", err)
	}

	if len(frontImage) > 0 {
		newMovieJSON.FrontImage = utils.GetBase64StringFromData(frontImage)
	}

	backImage, err := reader.GetBackImage(ctx, movie.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting movie back image: %v", err)
	}

	if len(backImage) > 0 {
		newMovieJSON.BackImage = utils.GetBase64StringFromData(backImage)
	}

	return &newMovieJSON, nil
}
