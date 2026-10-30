package autotag

import (
	"context"
	"slices"

	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
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

// GalleryPerformersAtPath tags the provided gallery with performers whose
// name matches the gallery's path. A fresh write txn is opened only when a
// match is applied.
func (tagger *Tagger) GalleryPerformersAtPath(ctx context.Context, s *models.Gallery, rw GalleryPerformerUpdater, performerReader models.PerformerAutoTagQueryer) error {
	t := getGalleryFileTagger(s, tagger.Cache)

	return t.tagPerformers(ctx, performerReader, func(subjectID, otherID int) (bool, error) {
		if err := s.LoadPerformerIDs(ctx, rw); err != nil {
			return false, err
		}
		existing := s.PerformerIDs.List()

		if slices.Contains(existing, otherID) {
			return false, nil
		}

		if err := txn.WithTxn(ctx, tagger.TxnManager, func(ctx context.Context) error {
			return gallery.AddPerformer(ctx, rw, s, otherID)
		}); err != nil {
			return false, err
		}

		return true, nil
	})
}

// GalleryStudiosAtPath tags the provided gallery with the first studio whose
// name matches the gallery's path.
//
// Galleries will not be tagged if studio is already set.
func (tagger *Tagger) GalleryStudiosAtPath(ctx context.Context, s *models.Gallery, rw GalleryFinderUpdater, studioReader models.StudioAutoTagQueryer) error {
	if s.StudioID != nil {
		// don't modify
		return nil
	}

	t := getGalleryFileTagger(s, tagger.Cache)

	return t.tagStudios(ctx, studioReader, func(subjectID, otherID int) (bool, error) {
		var added bool
		if err := txn.WithTxn(ctx, tagger.TxnManager, func(ctx context.Context) error {
			var err error
			added, err = addGalleryStudio(ctx, rw, s, otherID)
			return err
		}); err != nil {
			return false, err
		}
		return added, nil
	})
}

// GalleryTagsAtPath tags the provided gallery with tags whose name matches
// the gallery's path.
func (tagger *Tagger) GalleryTagsAtPath(ctx context.Context, s *models.Gallery, rw GalleryTagUpdater, tagReader models.TagAutoTagQueryer) error {
	t := getGalleryFileTagger(s, tagger.Cache)

	return t.tagTags(ctx, tagReader, func(subjectID, otherID int) (bool, error) {
		if err := s.LoadTagIDs(ctx, rw); err != nil {
			return false, err
		}
		existing := s.TagIDs.List()

		if slices.Contains(existing, otherID) {
			return false, nil
		}

		if err := txn.WithTxn(ctx, tagger.TxnManager, func(ctx context.Context) error {
			return gallery.AddTag(ctx, rw, s, otherID)
		}); err != nil {
			return false, err
		}

		return true, nil
	})
}
