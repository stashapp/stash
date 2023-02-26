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
		Title:     gallery.Title,
		URL:       gallery.URL,
		Details:   gallery.Details,
		CreatedAt: json.JSONTime{Time: gallery.CreatedAt},
		UpdatedAt: json.JSONTime{Time: gallery.UpdatedAt},
	}

	if gallery.FolderID != nil {
		newGalleryJSON.FolderPath = gallery.Path
	}

	for _, f := range gallery.Files.List() {
		newGalleryJSON.ZipFiles = append(newGalleryJSON.ZipFiles, f.Base().Path)
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
			return studio.Name, nil
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

func GetRefs(galleries []*models.Gallery) []jsonschema.GalleryRef {
	var results []jsonschema.GalleryRef
	for _, gallery := range galleries {
		toAdd := jsonschema.GalleryRef{}
		switch {
		case gallery.FolderID != nil:
			toAdd.FolderPath = gallery.Path
		case len(gallery.Files.List()) > 0:
			for _, f := range gallery.Files.List() {
				toAdd.ZipFiles = append(toAdd.ZipFiles, f.Base().Path)
			}
		default:
			toAdd.Title = gallery.Title
		}

		results = append(results, toAdd)
	}

	return results
}
