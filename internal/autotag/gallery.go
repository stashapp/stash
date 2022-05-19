package autotag

import (
	"context"

	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
)

func getGalleryFileTagger(s *models.Gallery, cache *match.Cache) tagger {
	// only trim the extension if gallery is file-based
	trimExt := s.Zip

	return tagger{
		ID:      s.ID,
		Type:    "gallery",
		Name:    s.GetTitle(),
		Path:    s.Path.String,
		trimExt: trimExt,
		cache:   cache,
	}
}

// GalleryPerformers tags the provided gallery with performers whose name matches the gallery's path.
func GalleryPerformers(ctx context.Context, s *models.Gallery, rw gallery.PerformerUpdater, performerReader match.PerformerAutoTagQueryer, cache *match.Cache) error {
	t := getGalleryFileTagger(s, cache)

	return t.tagPerformers(ctx, performerReader, func(subjectID, otherID int) (bool, error) {
		return gallery.AddPerformer(ctx, rw, subjectID, otherID)
	})
}

// GalleryStudios tags the provided gallery with the first studio whose name matches the gallery's path.
//
// Gallerys will not be tagged if studio is already set.
func GalleryStudios(ctx context.Context, s *models.Gallery, rw GalleryFinderUpdater, studioReader match.StudioAutoTagQueryer, cache *match.Cache) error {
	if s.StudioID.Valid {
		// don't modify
		return nil
	}

	t := getGalleryFileTagger(s, cache)

	return t.tagStudios(ctx, studioReader, func(subjectID, otherID int) (bool, error) {
		return addGalleryStudio(ctx, rw, subjectID, otherID)
	})
}

// GalleryTags tags the provided gallery with tags whose name matches the gallery's path.
func GalleryTags(ctx context.Context, s *models.Gallery, rw gallery.TagUpdater, tagReader match.TagAutoTagQueryer, cache *match.Cache) error {
	t := getGalleryFileTagger(s, cache)

	return t.tagTags(ctx, tagReader, func(subjectID, otherID int) (bool, error) {
		return gallery.AddTag(ctx, rw, subjectID, otherID)
	})
}
