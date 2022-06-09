package gallery

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/studio"
)

// ToBasicJSON converts a gallery object into its JSON object equivalent. It
// does not convert the relationships to other objects.
func ToBasicJSON(gallery *models.Gallery) (*jsonschema.Gallery, error) {
	newGalleryJSON := jsonschema.Gallery{
		Checksum:  gallery.Checksum,
		Title:     gallery.Title,
		URL:       gallery.URL,
		Details:   gallery.Details,
		Zip:       gallery.Zip,
		CreatedAt: json.JSONTime{Time: gallery.CreatedAt},
		UpdatedAt: json.JSONTime{Time: gallery.UpdatedAt},
	}

	if gallery.Path != nil {
		newGalleryJSON.Path = *gallery.Path
	}

	if gallery.FileModTime != nil {
		newGalleryJSON.FileModTime = json.JSONTime{Time: *gallery.FileModTime}
	}

	if gallery.Date != nil {
		newGalleryJSON.Date = gallery.Date.String()
	}

	if gallery.Rating != nil {
		newGalleryJSON.Rating = *gallery.Rating
	}

	newGalleryJSON.Organized = gallery.Organized

	return &newGalleryJSON, nil
}

// GetStudioName returns the name of the provided gallery's studio. It returns an
// empty string if there is no studio assigned to the gallery.
func GetStudioName(ctx context.Context, reader studio.Finder, gallery *models.Gallery) (string, error) {
	if gallery.StudioID != nil {
		studio, err := reader.Find(ctx, *gallery.StudioID)
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
