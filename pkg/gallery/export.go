package gallery

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/studio"
	"github.com/stashapp/stash/pkg/utils"
)

// ToBasicJSON converts a gallery object into its JSON object equivalent. It
// does not convert the relationships to other objects.
func ToBasicJSON(gallery *models.Gallery) (*jsonschema.Gallery, error) {
	newGalleryJSON := jsonschema.Gallery{
		Checksum:  gallery.Checksum,
		Zip:       gallery.Zip,
		CreatedAt: json.JSONTime{Time: gallery.CreatedAt.Timestamp},
		UpdatedAt: json.JSONTime{Time: gallery.UpdatedAt.Timestamp},
	}

	if gallery.Path.Valid {
		newGalleryJSON.Path = gallery.Path.String
	}

	if gallery.FileModTime.Valid {
		newGalleryJSON.FileModTime = json.JSONTime{Time: gallery.FileModTime.Timestamp}
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

	newGalleryJSON.Organized = gallery.Organized

	if gallery.Details.Valid {
		newGalleryJSON.Details = gallery.Details.String
	}

	return &newGalleryJSON, nil
}

// GetStudioName returns the name of the provided gallery's studio. It returns an
// empty string if there is no studio assigned to the gallery.
func GetStudioName(ctx context.Context, reader studio.Finder, gallery *models.Gallery) (string, error) {
	if gallery.StudioID.Valid {
		studio, err := reader.Find(ctx, int(gallery.StudioID.Int64))
		if err != nil {
			return "", err
		}

		if studio != nil {
			return studio.Name.String, nil
		}
	}

	return "", nil
}

func GetIDs(galleries []*models.Gallery) []int {
	var results []int
	for _, gallery := range galleries {
		results = append(results, gallery.ID)
	}

	return results
}

func GetChecksums(galleries []*models.Gallery) []string {
	var results []string
	for _, gallery := range galleries {
		if gallery.Checksum != "" {
			results = append(results, gallery.Checksum)
		}
	}

	return results
}
