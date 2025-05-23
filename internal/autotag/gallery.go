package autotag

import (
	"context"
	"slices"
	"time"

	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
)

type GalleryFinderUpdater interface {
	models.GalleryQueryer
	models.GalleryUpdater
}

type GalleryPerformerUpdater interface {
	models.PerformerIDLoader
	models.GalleryUpdater
}

type GalleryTagUpdater interface {
	models.TagIDLoader
	models.GalleryUpdater
}

func getGalleryFileTagger(s *models.Gallery, cache *match.Cache) tagger {
	var path string
	if s.Path != "" {
		path = s.Path
	}

	// only trim the extension if gallery is file-based
	trimExt := s.PrimaryFileID != nil

	return tagger{
		ID:      s.ID,
		Type:    "gallery",
		Name:    s.DisplayName(),
		Path:    path,
		trimExt: trimExt,
		cache:   cache,
	}
}

// GalleryPerformers tags the provided gallery with performers whose name matches the gallery's path.
func GalleryPerformers(ctx context.Context, s *models.Gallery, rw GalleryPerformerUpdater, performerReader models.PerformerAutoTagQueryer, cache *match.Cache) error {
	t := getGalleryFileTagger(s, cache)

	return t.tagPerformers(ctx, performerReader, func(subjectID, otherID int) (bool, error) {
		if err := s.LoadPerformerIDs(ctx, rw); err != nil {
			return false, err
		}
		existing := s.PerformerIDs.List()

		if slices.Contains(existing, otherID) {
			return false, nil
		}

		if err := gallery.AddPerformer(ctx, rw, s, otherID); err != nil {
			return false, err
		}

		return true, nil
	})
}

// GalleryStudios tags the provided gallery with the first studio whose name matches the gallery's path.
//
// Gallerys will not be tagged if studio is already set.
func GalleryStudios(ctx context.Context, s *models.Gallery, rw GalleryFinderUpdater, studioReader models.StudioAutoTagQueryer, cache *match.Cache) error {
	if s.StudioID != nil {
		// don't modify
		return nil
	}

	t := getGalleryFileTagger(s, cache)

	return t.tagStudios(ctx, studioReader, func(subjectID, otherID int) (bool, error) {
		return addGalleryStudio(ctx, rw, s, otherID)
	})
}

// GalleryTags tags the provided gallery with tags whose name matches the gallery's path.
func GalleryTags(ctx context.Context, s *models.Gallery, rw GalleryTagUpdater, tagReader models.TagAutoTagQueryer, cache *match.Cache) error {
	t := getGalleryFileTagger(s, cache)

	return t.tagTags(ctx, tagReader, func(subjectID, otherID int) (bool, error) {
		if err := s.LoadTagIDs(ctx, rw); err != nil {
			return false, err
		}
		existing := s.TagIDs.List()

		if slices.Contains(existing, otherID) {
			return false, nil
		}

		if err := gallery.AddTag(ctx, rw, s, otherID); err != nil {
			return false, err
		}

		return true, nil
	})
}

// GalleryDate extracts a date from the gallery path and sets it to gallery.Date if not already set.
func GalleryDate(ctx context.Context, g *models.Gallery, rw models.GalleryUpdater) error {
	t := getGalleryFileTagger(g, nil)
	return t.tagDates(ctx, func(dateStr string) (bool, error) {
		if g.Date != nil {
			return false, nil
		}
		tm, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return false, err
		}
		dateModel := models.Date{Time: tm}
		partial := models.NewGalleryPartial()
		partial.Date = models.NewOptionalDate(dateModel)
		_, err = rw.UpdatePartial(ctx, g.ID, partial)
		return err == nil, err
	})
}
