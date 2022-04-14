package autotag

import (
	"context"

	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
)

func getSceneFileTagger(s *models.Scene, cache *match.Cache) tagger {
	return tagger{
		ID:    s.ID,
		Type:  "scene",
		Name:  s.GetTitle(),
		Path:  s.Path,
		cache: cache,
	}
}

// ScenePerformers tags the provided scene with performers whose name matches the scene's path.
func ScenePerformers(ctx context.Context, s *models.Scene, rw scene.PerformerUpdater, performerReader match.PerformerAutoTagQueryer, cache *match.Cache) error {
	t := getSceneFileTagger(s, cache)

	return t.tagPerformers(ctx, performerReader, func(subjectID, otherID int) (bool, error) {
		return scene.AddPerformer(ctx, rw, subjectID, otherID)
	})
}

// SceneStudios tags the provided scene with the first studio whose name matches the scene's path.
//
// Scenes will not be tagged if studio is already set.
func SceneStudios(ctx context.Context, s *models.Scene, rw SceneFinderUpdater, studioReader match.StudioAutoTagQueryer, cache *match.Cache) error {
	if s.StudioID.Valid {
		// don't modify
		return nil
	}

	t := getSceneFileTagger(s, cache)

	return t.tagStudios(ctx, studioReader, func(subjectID, otherID int) (bool, error) {
		return addSceneStudio(ctx, rw, subjectID, otherID)
	})
}

// SceneTags tags the provided scene with tags whose name matches the scene's path.
func SceneTags(ctx context.Context, s *models.Scene, rw scene.TagUpdater, tagReader match.TagAutoTagQueryer, cache *match.Cache) error {
	t := getSceneFileTagger(s, cache)

	return t.tagTags(ctx, tagReader, func(subjectID, otherID int) (bool, error) {
		return scene.AddTag(ctx, rw, subjectID, otherID)
	})
}
