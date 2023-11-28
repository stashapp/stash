package image

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
	"github.com/stashapp/stash/pkg/models/jsonschema"
)

// ToBasicJSON converts a image object into its JSON object equivalent. It
// does not convert the relationships to other objects, with the exception
// of cover image.
func ToBasicJSON(image *models.Image) *jsonschema.Image {
	newImageJSON := jsonschema.Image{
		Title:        image.Title,
		Code:         image.Code,
		URLs:         image.URLs.List(),
		Details:      image.Details,
		Photographer: image.Photographer,
		CreatedAt:    json.JSONTime{Time: image.CreatedAt},
		UpdatedAt:    json.JSONTime{Time: image.UpdatedAt},
	}

	if image.Rating != nil {
		newImageJSON.Rating = *image.Rating
	}

	if image.Date != nil {
		newImageJSON.Date = image.Date.String()
	}

	newImageJSON.Organized = image.Organized
	newImageJSON.OCounter = image.OCounter

	for _, f := range image.Files.List() {
		newImageJSON.Files = append(newImageJSON.Files, f.Base().Path)
	}

	return &newImageJSON
}

// GetStudioName returns the name of the provided image's studio. It returns an
// empty string if there is no studio assigned to the image.
func GetStudioName(ctx context.Context, reader models.StudioGetter, image *models.Image) (string, error) {
	if image.StudioID != nil {
		studio, err := reader.Find(ctx, *image.StudioID)
		if err != nil {
			return "", err
		}

		if studio != nil {
			return studio.Name, nil
		}
	}

	return "", nil
}

// GetGalleryChecksum returns the checksum of the provided image. It returns an
// empty string if there is no gallery assigned to the image.
// func GetGalleryChecksum(reader models.GalleryReader, image *models.Image) (string, error) {
// 	gallery, err := reader.FindByImageID(image.ID)
// 	if err != nil {
// 		return "", fmt.Errorf("error getting image gallery: %v", err)
// 	}

// 	if gallery != nil {
// 		return gallery.Checksum, nil
// 	}

// 	return "", nil
// }
