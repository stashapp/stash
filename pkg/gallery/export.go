package gallery

import (
	"github.com/stashapp/stash/pkg/manager/jsonschema"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

// ToBasicJSON converts a gallery object into its JSON object equivalent. It
// does not convert the relationships to other objects.
func ToBasicJSON(gallery *models.Gallery) (*jsonschema.Gallery, error) {
	newGalleryJSON := jsonschema.Gallery{
		Checksum:  gallery.Checksum,
		Zip:       gallery.Zip,
		CreatedAt: models.JSONTime{Time: gallery.CreatedAt.Timestamp},
		UpdatedAt: models.JSONTime{Time: gallery.UpdatedAt.Timestamp},
	}

	if gallery.Path.Valid {
		newGalleryJSON.Path = gallery.Path.String
	}

	if gallery.Title.Valid {
		newGalleryJSON.Title = gallery.Title.String
	}

	if gallery.URL.Valid {
		newGalleryJSON.URL = gallery.URL.String
	}

	if gallery.Date.Valid {
		newGalleryJSON.Date = utils.GetYMDFromDatabaseDate(gallery.Date.String)
	}

	if gallery.Rating.Valid {
		newGalleryJSON.Rating = int(gallery.Rating.Int64)
	}

	if gallery.Details.Valid {
		newGalleryJSON.Details = gallery.Details.String
	}

	return &newGalleryJSON, nil
}

// GetStudioName returns the name of the provided gallery's studio. It returns an
// empty string if there is no studio assigned to the gallery.
func GetStudioName(reader models.StudioReader, gallery *models.Gallery) (string, error) {
	if gallery.StudioID.Valid {
		studio, err := reader.Find(int(gallery.StudioID.Int64))
		if err != nil {
			return "", err
		}

		if studio != nil {
			return studio.Name.String, nil
		}
	}

	return "", nil
}
