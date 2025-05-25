package autotag

import (
	"context"
	"slices"

	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
)

type ImageFinderUpdater interface {
	models.ImageQueryer
	models.ImageUpdater
}

type ImagePerformerUpdater interface {
	models.PerformerIDLoader
	models.ImageUpdater
}

type ImageTagUpdater interface {
	models.TagIDLoader
	models.ImageUpdater
}

type ImageStudioUpdater interface {
	models.ImageQueryer
	models.StudioIDLoader
	models.ImageUpdater
}

func getImageFileTagger(s *models.Image, cache *match.Cache) tagger {
	return tagger{
		ID:      s.ID,
		Type:    "image",
		Name:    s.DisplayName(),
		Path:    s.Path,
		trimExt: true, // trim extension for image files
		cache:   cache,
	}
}

// ImagePerformers tags the provided image with performers whose name matches the image's path.
func ImagePerformers(ctx context.Context, s *models.Image, rw ImagePerformerUpdater, performerReader models.PerformerAutoTagQueryer, cache *match.Cache) error {
	t := getImageFileTagger(s, cache)

	return t.tagPerformers(ctx, performerReader, func(subjectID, otherID int) (bool, error) {
		if err := s.LoadPerformerIDs(ctx, rw); err != nil {
			return false, err
		}
		existing := s.PerformerIDs.List()

		if slices.Contains(existing, otherID) {
			return false, nil
		}

		if err := image.AddPerformer(ctx, rw, s, otherID); err != nil {
			return false, err
		}

		return true, nil
	})
}

// ImageStudios tags the provided image with studios whose names match the image's path.
func ImageStudios(ctx context.Context, s *models.Image, rw ImageStudioUpdater, studioReader models.StudioAutoTagQueryer, cache *match.Cache) error {
	t := getImageFileTagger(s, cache)

	return t.tagStudios(ctx, studioReader, func(subjectID, otherID int) (bool, error) {
		return addImageStudio(ctx, rw, s, otherID)
	})
}

// ImageTags tags the provided image with tags whose name matches the image's path.
func ImageTags(ctx context.Context, s *models.Image, rw ImageTagUpdater, tagReader models.TagAutoTagQueryer, cache *match.Cache) error {
	t := getImageFileTagger(s, cache)

	return t.tagTags(ctx, tagReader, func(subjectID, otherID int) (bool, error) {
		if err := s.LoadTagIDs(ctx, rw); err != nil {
			return false, err
		}
		existing := s.TagIDs.List()

		if slices.Contains(existing, otherID) {
			return false, nil
		}

		if err := image.AddTag(ctx, rw, s, otherID); err != nil {
			return false, err
		}

		return true, nil
	})
}
