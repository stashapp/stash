package image

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/studio"
)

// ToBasicJSON converts a image object into its JSON object equivalent. It
// does not convert the relationships to other objects, with the exception
// of cover image.
func ToBasicJSON(image *models.Image) *jsonschema.Image {
	newImageJSON := jsonschema.Image{
		Checksum:  image.Checksum,
		CreatedAt: json.JSONTime{Time: image.CreatedAt.Timestamp},
		UpdatedAt: json.JSONTime{Time: image.UpdatedAt.Timestamp},
	}

	if image.Title.Valid {
		newImageJSON.Title = image.Title.String
	}

	if image.Rating.Valid {
		newImageJSON.Rating = int(image.Rating.Int64)
	}

	newImageJSON.Organized = image.Organized
	newImageJSON.OCounter = image.OCounter

	newImageJSON.File = getImageFileJSON(image)

	return &newImageJSON
}

func getImageFileJSON(image *models.Image) *jsonschema.ImageFile {
	ret := &jsonschema.ImageFile{}

	if image.FileModTime.Valid {
		ret.ModTime = json.JSONTime{Time: image.FileModTime.Timestamp}
	}

	if image.Size.Valid {
		ret.Size = int(image.Size.Int64)
	}

	if image.Width.Valid {
		ret.Width = int(image.Width.Int64)
	}

	if image.Height.Valid {
		ret.Height = int(image.Height.Int64)
	}

	return ret
}

// GetStudioName returns the name of the provided image's studio. It returns an
// empty string if there is no studio assigned to the image.
func GetStudioName(ctx context.Context, reader studio.Finder, image *models.Image) (string, error) {
	if image.StudioID.Valid {
		studio, err := reader.Find(ctx, int(image.StudioID.Int64))
		if err != nil {
			return "", err
		}

		if studio != nil {
			return studio.Name.String, nil
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
