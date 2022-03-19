package movie

import (
	"database/sql"
	"fmt"

	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/utils"
)

type Importer struct {
	ReaderWriter        models.MovieReaderWriter
	StudioWriter        models.StudioReaderWriter
	Input               jsonschema.Movie
	MissingRefBehaviour models.ImportMissingRefEnum

	movie          models.Movie
	frontImageData []byte
	backImageData  []byte
}

func (i *Importer) PreImport() error {
	i.movie = i.movieJSONToMovie(i.Input)

	if err := i.populateStudio(); err != nil {
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

func (i *Importer) populateStudio() error {
	if i.Input.Studio != "" {
		studio, err := i.StudioWriter.FindByName(i.Input.Studio, false)
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
				studioID, err := i.createStudio(i.Input.Studio)
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

func (i *Importer) createStudio(name string) (int, error) {
	newStudio := *models.NewStudio(name)

	created, err := i.StudioWriter.Create(newStudio)
	if err != nil {
		return 0, err
	}

	return created.ID, nil
}

func (i *Importer) PostImport(id int) error {
	if len(i.frontImageData) > 0 {
		if err := i.ReaderWriter.UpdateImages(id, i.frontImageData, i.backImageData); err != nil {
			return fmt.Errorf("error setting movie images: %v", err)
		}
	}

	return nil
}

func (i *Importer) Name() string {
	return i.Input.Name
}

func (i *Importer) FindExistingID() (*int, error) {
	const nocase = false
	existing, err := i.ReaderWriter.FindByName(i.Name(), nocase)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		id := existing.ID
		return &id, nil
	}

	return nil, nil
}

func (i *Importer) Create() (*int, error) {
	created, err := i.ReaderWriter.Create(i.movie)
	if err != nil {
		return nil, fmt.Errorf("error creating movie: %v", err)
	}

	id := created.ID
	return &id, nil
}

func (i *Importer) Update(id int) error {
	movie := i.movie
	movie.ID = id
	_, err := i.ReaderWriter.UpdateFull(movie)
	if err != nil {
		return fmt.Errorf("error updating existing movie: %v", err)
	}

	return nil
}
