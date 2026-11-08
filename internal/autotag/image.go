package autotag

import (
	"context"
	"slices"

	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
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

func getImageFileTagger(s *models.Image, cache *match.Cache) tagger {
	return tagger{
		ID:    s.ID,
		Type:  "image",
		Name:  s.DisplayName(),
		Path:  s.Path,
		cache: cache,
	}
}

// ImagePerformersAtPath tags the provided image with performers whose name
// matches the image's path. A fresh write txn is opened only when a match is
// applied.
func (tagger *Tagger) ImagePerformersAtPath(ctx context.Context, s *models.Image, rw ImagePerformerUpdater, performerReader models.PerformerAutoTagQueryer) error {
	t := getImageFileTagger(s, tagger.Cache)

	return t.tagPerformers(ctx, performerReader, func(subjectID, otherID int) (bool, error) {
		if err := s.LoadPerformerIDs(ctx, rw); err != nil {
			return false, err
		}
		existing := s.PerformerIDs.List()

		if slices.Contains(existing, otherID) {
			return false, nil
		}

		if err := txn.WithTxn(ctx, tagger.TxnManager, func(ctx context.Context) error {
			return image.AddPerformer(ctx, rw, s, otherID)
		}); err != nil {
			return false, err
		}

		return true, nil
	})
}

// ImageStudiosAtPath tags the provided image with the first studio whose
// name matches the image's path.
//
// Images will not be tagged if studio is already set.
func (tagger *Tagger) ImageStudiosAtPath(ctx context.Context, s *models.Image, rw ImageFinderUpdater, studioReader models.StudioAutoTagQueryer) error {
	if s.StudioID != nil {
		// don't modify
		return nil
	}

	t := getImageFileTagger(s, tagger.Cache)

	return t.tagStudios(ctx, studioReader, func(subjectID, otherID int) (bool, error) {
		var added bool
		if err := txn.WithTxn(ctx, tagger.TxnManager, func(ctx context.Context) error {
			var err error
			added, err = addImageStudio(ctx, rw, s, otherID)
			return err
		}); err != nil {
			return false, err
		}
		return added, nil
	})
}

// ImageTagsAtPath tags the provided image with tags whose name matches the
// image's path.
func (tagger *Tagger) ImageTagsAtPath(ctx context.Context, s *models.Image, rw ImageTagUpdater, tagReader models.TagAutoTagQueryer) error {
	t := getImageFileTagger(s, tagger.Cache)

	return t.tagTags(ctx, tagReader, func(subjectID, otherID int) (bool, error) {
		if err := s.LoadTagIDs(ctx, rw); err != nil {
			return false, err
		}
		existing := s.TagIDs.List()

		if slices.Contains(existing, otherID) {
			return false, nil
		}

		if err := txn.WithTxn(ctx, tagger.TxnManager, func(ctx context.Context) error {
			return image.AddTag(ctx, rw, s, otherID)
		}); err != nil {
			return false, err
		}

		return true, nil
	})
}
