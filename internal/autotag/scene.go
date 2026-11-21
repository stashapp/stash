package autotag

import (
	"context"
	"slices"

	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/txn"
)

type SceneFinderUpdater interface {
	models.SceneQueryer
	models.SceneUpdater
}

type ScenePerformerUpdater interface {
	models.PerformerIDLoader
	models.SceneUpdater
}

type SceneTagUpdater interface {
	models.TagIDLoader
	models.SceneUpdater
}

func getSceneFileTagger(s *models.Scene, cache *match.Cache) tagger {
	return tagger{
		ID:    s.ID,
		Type:  "scene",
		Name:  s.DisplayName(),
		Path:  s.Path,
		cache: cache,
	}
}

// ScenePerformersAtPath tags the provided scene with performers whose name
// matches the scene's path. The match phase runs using the current context
// (no outer write txn needed); a fresh write txn is opened only when a match
// is applied.
func (tagger *Tagger) ScenePerformersAtPath(ctx context.Context, s *models.Scene, rw ScenePerformerUpdater, performerReader models.PerformerAutoTagQueryer) error {
	t := getSceneFileTagger(s, tagger.Cache)

	return t.tagPerformers(ctx, performerReader, func(subjectID, otherID int) (bool, error) {
		if err := s.LoadPerformerIDs(ctx, rw); err != nil {
			return false, err
		}
		existing := s.PerformerIDs.List()

		if slices.Contains(existing, otherID) {
			return false, nil
		}

		if err := txn.WithTxn(ctx, tagger.TxnManager, func(ctx context.Context) error {
			return scene.AddPerformer(ctx, rw, s, otherID)
		}); err != nil {
			return false, err
		}

		return true, nil
	})
}

// SceneStudiosAtPath tags the provided scene with the first studio whose name
// matches the scene's path.
//
// Scenes will not be tagged if studio is already set.
func (tagger *Tagger) SceneStudiosAtPath(ctx context.Context, s *models.Scene, rw SceneFinderUpdater, studioReader models.StudioAutoTagQueryer) error {
	if s.StudioID != nil {
		// don't modify
		return nil
	}

	t := getSceneFileTagger(s, tagger.Cache)

	return t.tagStudios(ctx, studioReader, func(subjectID, otherID int) (bool, error) {
		var added bool
		if err := txn.WithTxn(ctx, tagger.TxnManager, func(ctx context.Context) error {
			var err error
			added, err = addSceneStudio(ctx, rw, s, otherID)
			return err
		}); err != nil {
			return false, err
		}
		return added, nil
	})
}

// SceneTagsAtPath tags the provided scene with tags whose name matches the
// scene's path.
func (tagger *Tagger) SceneTagsAtPath(ctx context.Context, s *models.Scene, rw SceneTagUpdater, tagReader models.TagAutoTagQueryer) error {
	t := getSceneFileTagger(s, tagger.Cache)

	return t.tagTags(ctx, tagReader, func(subjectID, otherID int) (bool, error) {
		if err := s.LoadTagIDs(ctx, rw); err != nil {
			return false, err
		}
		existing := s.TagIDs.List()

		if slices.Contains(existing, otherID) {
			return false, nil
		}

		if err := txn.WithTxn(ctx, tagger.TxnManager, func(ctx context.Context) error {
			return scene.AddTag(ctx, rw, s, otherID)
		}); err != nil {
			return false, err
		}

		return true, nil
	})
}
