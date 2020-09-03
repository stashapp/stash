package performer

import (
	"database/sql"
	"fmt"

	"github.com/stashapp/stash/pkg/manager/jsonschema"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type Importer struct {
	ReaderWriter       models.PerformerReaderWriter
	DuplicateBehaviour models.ImportDuplicateEnum
}

func (i Importer) Import(performerJSON *jsonschema.Performer) error {
	p := performerJSONToPerformer(performerJSON)

	var imageData []byte
	var err error
	if len(performerJSON.Image) > 0 {
		_, imageData, err = utils.ProcessBase64Image(performerJSON.Image)
		if err != nil {
			return fmt.Errorf("invalid image: %s", err.Error())
		}
	}

	// try to find an existing performer with the same name
	const nocase = false
	existing, err := i.ReaderWriter.FindByNames([]string{p.Name.String}, nocase)
	if err != nil {
		return fmt.Errorf("error finding existing performers: %s", err.Error())
	}

	var id int

	if len(existing) > 0 {
		if i.DuplicateBehaviour == models.ImportDuplicateEnumFail {
			return fmt.Errorf("existing performer with name '%s'", p.Name.String)
		} else if i.DuplicateBehaviour == models.ImportDuplicateEnumIgnore {
			return nil
		}

		// must be overwriting
		p.ID = existing[0].ID
		id = p.ID
		_, err := i.ReaderWriter.Update(*p)
		if err != nil {
			return fmt.Errorf("error updating existing performer: %s", err.Error())
		}
	} else {
		// creating
		created, err := i.ReaderWriter.Create(*p)
		if err != nil {
			return fmt.Errorf("error creating performer: %s", err.Error())
		}

		id = created.ID
	}

	// regardless of whether we created or updated, update the image if provided
	if len(imageData) > 0 {
		if err := i.ReaderWriter.UpdatePerformerImage(id, imageData); err != nil {
			return fmt.Errorf("error setting performer image: %s", err.Error())
		}
	}

	return nil
}

func performerJSONToPerformer(performerJSON *jsonschema.Performer) *models.Performer {
	checksum := utils.MD5FromString(performerJSON.Name)

	newPerformer := models.Performer{
		Checksum:  checksum,
		Favorite:  sql.NullBool{Bool: performerJSON.Favorite, Valid: true},
		CreatedAt: models.SQLiteTimestamp{Timestamp: performerJSON.CreatedAt.GetTime()},
		UpdatedAt: models.SQLiteTimestamp{Timestamp: performerJSON.UpdatedAt.GetTime()},
	}

	if performerJSON.Name != "" {
		newPerformer.Name = sql.NullString{String: performerJSON.Name, Valid: true}
	}
	if performerJSON.Gender != "" {
		newPerformer.Gender = sql.NullString{String: performerJSON.Gender, Valid: true}
	}
	if performerJSON.URL != "" {
		newPerformer.URL = sql.NullString{String: performerJSON.URL, Valid: true}
	}
	if performerJSON.Birthdate != "" {
		newPerformer.Birthdate = models.SQLiteDate{String: performerJSON.Birthdate, Valid: true}
	}
	if performerJSON.Ethnicity != "" {
		newPerformer.Ethnicity = sql.NullString{String: performerJSON.Ethnicity, Valid: true}
	}
	if performerJSON.Country != "" {
		newPerformer.Country = sql.NullString{String: performerJSON.Country, Valid: true}
	}
	if performerJSON.EyeColor != "" {
		newPerformer.EyeColor = sql.NullString{String: performerJSON.EyeColor, Valid: true}
	}
	if performerJSON.Height != "" {
		newPerformer.Height = sql.NullString{String: performerJSON.Height, Valid: true}
	}
	if performerJSON.Measurements != "" {
		newPerformer.Measurements = sql.NullString{String: performerJSON.Measurements, Valid: true}
	}
	if performerJSON.FakeTits != "" {
		newPerformer.FakeTits = sql.NullString{String: performerJSON.FakeTits, Valid: true}
	}
	if performerJSON.CareerLength != "" {
		newPerformer.CareerLength = sql.NullString{String: performerJSON.CareerLength, Valid: true}
	}
	if performerJSON.Tattoos != "" {
		newPerformer.Tattoos = sql.NullString{String: performerJSON.Tattoos, Valid: true}
	}
	if performerJSON.Piercings != "" {
		newPerformer.Piercings = sql.NullString{String: performerJSON.Piercings, Valid: true}
	}
	if performerJSON.Aliases != "" {
		newPerformer.Aliases = sql.NullString{String: performerJSON.Aliases, Valid: true}
	}
	if performerJSON.Twitter != "" {
		newPerformer.Twitter = sql.NullString{String: performerJSON.Twitter, Valid: true}
	}
	if performerJSON.Instagram != "" {
		newPerformer.Instagram = sql.NullString{String: performerJSON.Instagram, Valid: true}
	}

	return &newPerformer
}
