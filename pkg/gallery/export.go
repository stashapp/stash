package gallery

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
	"github.com/stashapp/stash/pkg/models/jsonschema"
)

// ToBasicJSON converts a gallery object into its JSON object equivalent. It
// does not convert the relationships to other objects.
func ToBasicJSON(gallery *models.Gallery) (*jsonschema.Gallery, error) {
	newGalleryJSON := jsonschema.Gallery{
		Title:        gallery.Title,
		Code:         gallery.Code,
		URLs:         gallery.URLs.List(),
		Details:      gallery.Details,
		Photographer: gallery.Photographer,
		CreatedAt:    json.JSONTime{Time: gallery.CreatedAt},
		UpdatedAt:    json.JSONTime{Time: gallery.UpdatedAt},
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

// GetStudioNames returns the names of the provided gallery's studios.
func GetStudioNames(ctx context.Context, reader models.StudioGetter, gallery *models.Gallery) ([]string, error) {
	// Assume studio IDs are already loaded
	studioIDs := gallery.StudioIDs.List()
	if len(studioIDs) == 0 {
		return nil, nil
	}

	var studioNames []string
	for _, studioID := range studioIDs {
		studio, err := reader.Find(ctx, studioID)
		if err != nil {
			return nil, err
		}

		if studio != nil && studio.Name != "" {
			studioNames = append(studioNames, studio.Name)
		}
	}

	return studioNames, nil
}

// GetGalleryChaptersJSON returns a slice of GalleryChapter JSON representation
// objects corresponding to the provided gallery's chapters.
func GetGalleryChaptersJSON(ctx context.Context, chapterReader models.GalleryChapterFinder, gallery *models.Gallery) ([]jsonschema.GalleryChapter, error) {
	galleryChapters, err := chapterReader.FindByGalleryID(ctx, gallery.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting gallery chapters: %v", err)
	}

	var results []jsonschema.GalleryChapter

	for _, galleryChapter := range galleryChapters {
		galleryChapterJSON := jsonschema.GalleryChapter{
			Title:      galleryChapter.Title,
			ImageIndex: galleryChapter.ImageIndex,
			CreatedAt:  json.JSONTime{Time: galleryChapter.CreatedAt},
			UpdatedAt:  json.JSONTime{Time: galleryChapter.UpdatedAt},
		}

		results = append(results, galleryChapterJSON)
	}

	return results, nil
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
