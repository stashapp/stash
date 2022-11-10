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

func getPerformerTaggers(p *models.Performer, cache *match.Cache) []tagger {
	ret := []tagger{{
		ID:    p.ID,
		Type:  "performer",
		Name:  p.Name,
		cache: cache,
	}}

	for _, a := range p.Aliases.List() {
		ret = append(ret, tagger{
			ID:    p.ID,
			Type:  "performer",
			Name:  a,
			cache: cache,
		})
	}

	return ret
}

// PerformerScenes searches for scenes whose path matches the provided performer name and tags the scene with the performer.
// Performer aliases must be loaded.
func PerformerScenes(ctx context.Context, p *models.Performer, paths []string, rw SceneQueryPerformerUpdater, cache *match.Cache) error {
	t := getPerformerTaggers(p, cache)

	for _, tt := range t {
		if err := tt.tagScenes(ctx, paths, rw, func(o *models.Scene) (bool, error) {
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
		}); err != nil {
			return err
		}
	}
	return nil
}

// PerformerImages searches for images whose path matches the provided performer name and tags the image with the performer.
// Performer aliases must be loaded.
func PerformerImages(ctx context.Context, p *models.Performer, paths []string, rw ImageQueryPerformerUpdater, cache *match.Cache) error {
	t := getPerformerTaggers(p, cache)

	for _, tt := range t {
		if err := tt.tagImages(ctx, paths, rw, func(o *models.Image) (bool, error) {
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
		}); err != nil {
			return err
		}
	}
	return nil
}

// PerformerGalleries searches for galleries whose path matches the provided performer name and tags the gallery with the performer.
// Performer aliases must be loaded.
func PerformerGalleries(ctx context.Context, p *models.Performer, paths []string, rw GalleryQueryPerformerUpdater, cache *match.Cache) error {
	t := getPerformerTaggers(p, cache)

	for _, tt := range t {
		if err := tt.tagGalleries(ctx, paths, rw, func(o *models.Gallery) (bool, error) {
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
		}); err != nil {
			return err
		}
	}
	return nil
}
