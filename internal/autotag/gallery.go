package autotag

import (
	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
)

func getGalleryFileTagger(s *models.Gallery, cache *match.Cache) tagger {
	return tagger{
		ID:    s.ID,
		Type:  "gallery",
		Name:  s.GetTitle(),
		Path:  s.Path.String,
		cache: cache,
	}
}

// GalleryPerformers tags the provided gallery with performers whose name matches the gallery's path.
func GalleryPerformers(s *models.Gallery, rw models.GalleryReaderWriter, performerReader models.PerformerReader, cache *match.Cache) error {
	t := getGalleryFileTagger(s, cache)

	return t.tagPerformers(performerReader, func(subjectID, otherID int) (bool, error) {
		return gallery.AddPerformer(rw, subjectID, otherID)
	})
}

// GalleryStudios tags the provided gallery with the first studio whose name matches the gallery's path.
//
// Gallerys will not be tagged if studio is already set.
func GalleryStudios(s *models.Gallery, rw models.GalleryReaderWriter, studioReader models.StudioReader, cache *match.Cache) error {
	if s.StudioID.Valid {
		// don't modify
		return nil
	}

	t := getGalleryFileTagger(s, cache)

	return t.tagStudios(studioReader, func(subjectID, otherID int) (bool, error) {
		return addGalleryStudio(rw, subjectID, otherID)
	})
}

// GalleryTags tags the provided gallery with tags whose name matches the gallery's path.
func GalleryTags(s *models.Gallery, rw models.GalleryReaderWriter, tagReader models.TagReader, cache *match.Cache) error {
	t := getGalleryFileTagger(s, cache)

	return t.tagTags(tagReader, func(subjectID, otherID int) (bool, error) {
		return gallery.AddTag(rw, subjectID, otherID)
	})
}
