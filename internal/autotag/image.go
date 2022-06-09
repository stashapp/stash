package autotag

import (
	"context"

	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
)

func getImageFileTagger(s *models.Image, cache *match.Cache) tagger {
	return tagger{
		ID:    s.ID,
		Type:  "image",
		Name:  s.GetTitle(),
		Path:  s.Path,
		cache: cache,
	}
}

// ImagePerformers tags the provided image with performers whose name matches the image's path.
func ImagePerformers(ctx context.Context, s *models.Image, rw image.PartialUpdater, performerReader match.PerformerAutoTagQueryer, cache *match.Cache) error {
	t := getImageFileTagger(s, cache)

	return t.tagPerformers(ctx, performerReader, func(subjectID, otherID int) (bool, error) {
		return image.AddPerformer(ctx, rw, s, otherID)
	})
}

// ImageStudios tags the provided image with the first studio whose name matches the image's path.
//
// Images will not be tagged if studio is already set.
func ImageStudios(ctx context.Context, s *models.Image, rw ImageFinderUpdater, studioReader match.StudioAutoTagQueryer, cache *match.Cache) error {
	if s.StudioID != nil {
		// don't modify
		return nil
	}

	t := getImageFileTagger(s, cache)

	return t.tagStudios(ctx, studioReader, func(subjectID, otherID int) (bool, error) {
		return addImageStudio(ctx, rw, s, otherID)
	})
}

// ImageTags tags the provided image with tags whose name matches the image's path.
func ImageTags(ctx context.Context, s *models.Image, rw image.PartialUpdater, tagReader match.TagAutoTagQueryer, cache *match.Cache) error {
	t := getImageFileTagger(s, cache)

	return t.tagTags(ctx, tagReader, func(subjectID, otherID int) (bool, error) {
		return image.AddTag(ctx, rw, s, otherID)
	})
}
