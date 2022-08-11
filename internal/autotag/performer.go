package autotag

import (
	"context"

	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
)

type SceneQueryPerformerUpdater interface {
	scene.Queryer
	models.PerformerIDLoader
	scene.PartialUpdater
}

type ImageQueryPerformerUpdater interface {
	image.Queryer
	models.PerformerIDLoader
	image.PartialUpdater
}

type GalleryQueryPerformerUpdater interface {
	gallery.Queryer
	models.PerformerIDLoader
	gallery.PartialUpdater
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

	return t.tagScenes(ctx, paths, rw, func(o *models.Scene) (bool, error) {
		if err := o.LoadPerformerIDs(ctx, rw); err != nil {
			return false, err
		}
		existing := o.PerformerIDs.List()

		if intslice.IntInclude(existing, p.ID) {
			return false, nil
		}

		if err := scene.AddPerformer(ctx, rw, o, p.ID); err != nil {
			return false, err
		}

		return true, nil
	})
}

// PerformerImages searches for images whose path matches the provided performer name and tags the image with the performer.
func PerformerImages(ctx context.Context, p *models.Performer, paths []string, rw ImageQueryPerformerUpdater, cache *match.Cache) error {
	t := getPerformerTagger(p, cache)

	return t.tagImages(ctx, paths, rw, func(o *models.Image) (bool, error) {
		if err := o.LoadPerformerIDs(ctx, rw); err != nil {
			return false, err
		}
		existing := o.PerformerIDs.List()

		if intslice.IntInclude(existing, p.ID) {
			return false, nil
		}

		if err := image.AddPerformer(ctx, rw, o, p.ID); err != nil {
			return false, err
		}

		return true, nil
	})
}

// PerformerGalleries searches for galleries whose path matches the provided performer name and tags the gallery with the performer.
func PerformerGalleries(ctx context.Context, p *models.Performer, paths []string, rw GalleryQueryPerformerUpdater, cache *match.Cache) error {
	t := getPerformerTagger(p, cache)

	return t.tagGalleries(ctx, paths, rw, func(o *models.Gallery) (bool, error) {
		if err := o.LoadPerformerIDs(ctx, rw); err != nil {
			return false, err
		}
		existing := o.PerformerIDs.List()

		if intslice.IntInclude(existing, p.ID) {
			return false, nil
		}

		if err := gallery.AddPerformer(ctx, rw, o, p.ID); err != nil {
			return false, err
		}

		return true, nil
	})
}
