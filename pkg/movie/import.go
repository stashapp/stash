package movie

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/studio"
	"github.com/stashapp/stash/pkg/utils"
)

type NameFinderCreatorUpdater interface {
	NameFinderCreator
	UpdateFull(ctx context.Context, updatedMovie models.Movie) (*models.Movie, error)
	UpdateImages(ctx context.Context, movieID int, frontImage []byte, backImage []byte) error
}

type Importer struct {
	ReaderWriter        NameFinderCreatorUpdater
	StudioWriter        studio.NameFinderCreator
	Input               jsonschema.Movie
	MissingRefBehaviour models.ImportMissingRefEnum

	movie          models.Movie
	frontImageData []byte
	backImageData  []byte
}

func (i *Importer) PreImport(ctx context.Context) error {
	i.movie = i.movieJSONToMovie(i.Input)

	if err := i.populateStudio(ctx); err != nil {
		return err
	}

	var err error
	if len(i.Input.FrontImage) > 0 {
		i.frontImageData, err = utils.ProcessBase64Image(i.Input.FrontImage)
		if err != nil {
			return fmt.Errorf("invalid front_image: %v", err)
		}
	}
	if len(i.Input.BackImage) > 0 {
		i.backImageData, err = utils.ProcessBase64Image(i.Input.BackImage)
		if err != nil {
			return fmt.Errorf("invalid back_image: %v", err)
		}
	}

	return nil
}

func (i *Importer) movieJSONToMovie(movieJSON jsonschema.Movie) models.Movie {
	checksum := md5.FromString(movieJSON.Name)

	newMovie := models.Movie{
		Checksum:  checksum,
		Name:      sql.NullString{String: movieJSON.Name, Valid: true},
		Aliases:   sql.NullString{String: movieJSON.Aliases, Valid: true},
		Date:      models.SQLiteDate{String: movieJSON.Date, Valid: true},
		Director:  sql.NullString{String: movieJSON.Director, Valid: true},
		Synopsis:  sql.NullString{String: movieJSON.Synopsis, Valid: true},
		URL:       sql.NullString{String: movieJSON.URL, Valid: true},
		CreatedAt: models.SQLiteTimestamp{Timestamp: movieJSON.CreatedAt.GetTime()},
		UpdatedAt: models.SQLiteTimestamp{Timestamp: movieJSON.UpdatedAt.GetTime()},
	}

	if movieJSON.Rating != 0 {
		newMovie.Rating = sql.NullInt64{Int64: int64(movieJSON.Rating), Valid: true}
	}

	if movieJSON.Duration != 0 {
		newMovie.Duration = sql.NullInt64{Int64: int64(movieJSON.Duration), Valid: true}
	}

	return newMovie
}

func (i *Importer) populateStudio(ctx context.Context) error {
	if i.Input.Studio != "" {
		studio, err := i.StudioWriter.FindByName(ctx, i.Input.Studio, false)
		if err != nil {
			return fmt.Errorf("error finding studio by name: %v", err)
		}

		if studio == nil {
			if i.MissingRefBehaviour == models.ImportMissingRefEnumFail {
				return fmt.Errorf("movie studio '%s' not found", i.Input.Studio)
			}

			if i.MissingRefBehaviour == models.ImportMissingRefEnumIgnore {
				return nil
			}

			if i.MissingRefBehaviour == models.ImportMissingRefEnumCreate {
				studioID, err := i.createStudio(ctx, i.Input.Studio)
				if err != nil {
					return err
				}
				i.movie.StudioID = sql.NullInt64{
					Int64: int64(studioID),
					Valid: true,
				}
			}
		} else {
			i.movie.StudioID = sql.NullInt64{Int64: int64(studio.ID), Valid: true}
		}
	}

	return nil
}

func (i *Importer) createStudio(ctx context.Context, name string) (int, error) {
	newStudio := *models.NewStudio(name)

	created, err := i.StudioWriter.Create(ctx, newStudio)
	if err != nil {
		return 0, err
	}

	return created.ID, nil
}

func (i *Importer) PostImport(ctx context.Context, id int) error {
	if len(i.frontImageData) > 0 {
		if err := i.ReaderWriter.UpdateImages(ctx, id, i.frontImageData, i.backImageData); err != nil {
			return fmt.Errorf("error setting movie images: %v", err)
		}
	}

	return nil
}

func (i *Importer) Name() string {
	return i.Input.Name
}

func (i *Importer) FindExistingID(ctx context.Context) (*int, error) {
	const nocase = false
	existing, err := i.ReaderWriter.FindByName(ctx, i.Name(), nocase)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		id := existing.ID
		return &id, nil
	}

	return nil, nil
}

func (i *Importer) Create(ctx context.Context) (*int, error) {
	created, err := i.ReaderWriter.Create(ctx, i.movie)
	if err != nil {
		return nil, fmt.Errorf("error creating movie: %v", err)
	}

	id := created.ID
	return &id, nil
}

func (i *Importer) Update(ctx context.Context, id int) error {
	movie := i.movie
	movie.ID = id
	_, err := i.ReaderWriter.UpdateFull(ctx, movie)
	if err != nil {
		return fmt.Errorf("error updating existing movie: %v", err)
	}

	return nil
}
