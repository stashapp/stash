package autotag

import (
	"context"

	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
)

type SceneQueryPerformerUpdater interface {
	scene.Queryer
	scene.PerformerUpdater
}

type ImageQueryPerformerUpdater interface {
	image.Queryer
	image.PerformerUpdater
}

type GalleryQueryPerformerUpdater interface {
	gallery.Queryer
	gallery.PerformerUpdater
}

func getPerformerTagger(p *models.Performer, cache *match.Cache) tagger {
	return tagger{
		ID:    p.ID,
		Type:  "performer",
		Name:  p.Name.String,
		cache: cache,
	}
}

// PerformerScenes searches for scenes whose path matches the provided performer name and tags the scene with the performer.
func PerformerScenes(ctx context.Context, p *models.Performer, paths []string, rw SceneQueryPerformerUpdater, cache *match.Cache) error {
	t := getPerformerTagger(p, cache)

	return t.tagScenes(ctx, paths, rw, func(subjectID, otherID int) (bool, error) {
		return scene.AddPerformer(ctx, rw, otherID, subjectID)
	})
}

// PerformerImages searches for images whose path matches the provided performer name and tags the image with the performer.
func PerformerImages(ctx context.Context, p *models.Performer, paths []string, rw ImageQueryPerformerUpdater, cache *match.Cache) error {
	t := getPerformerTagger(p, cache)

	return t.tagImages(ctx, paths, rw, func(subjectID, otherID int) (bool, error) {
		return image.AddPerformer(ctx, rw, otherID, subjectID)
	})
}

// PerformerGalleries searches for galleries whose path matches the provided performer name and tags the gallery with the performer.
func PerformerGalleries(ctx context.Context, p *models.Performer, paths []string, rw GalleryQueryPerformerUpdater, cache *match.Cache) error {
	t := getPerformerTagger(p, cache)

	return t.tagGalleries(ctx, paths, rw, func(subjectID, otherID int) (bool, error) {
		return gallery.AddPerformer(ctx, rw, otherID, subjectID)
	})
}
