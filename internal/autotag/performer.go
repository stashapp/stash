package autotag

import (
	"context"
	"slices"

	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/txn"
)

type SceneQueryPerformerUpdater interface {
	models.SceneQueryer
	models.PerformerIDLoader
	models.SceneUpdater
}

type ImageQueryPerformerUpdater interface {
	models.ImageQueryer
	models.PerformerIDLoader
	models.ImageUpdater
}

type GalleryQueryPerformerUpdater interface {
	models.GalleryQueryer
	models.PerformerIDLoader
	models.GalleryUpdater
}

func getPerformerTaggers(p *models.Performer, cache *match.Cache) []tagger {
	ret := []tagger{{
		ID:    p.ID,
		Type:  "performer",
		Name:  p.Name,
		cache: cache,
	}}

	// TODO - disabled until we can have finer control over alias matching
	// for _, a := range p.Aliases.List() {
	// 	ret = append(ret, tagger{
	// 		ID:    p.ID,
	// 		Type:  "performer",
	// 		Name:  a,
	// 		cache: cache,
	// 	})
	// }

	return ret
}

// PerformerScenes searches for scenes whose path matches the provided performer name and tags the scene with the performer.
// Performer aliases must be loaded.
func (tagger *Tagger) PerformerScenes(ctx context.Context, p *models.Performer, paths []string, rw SceneQueryPerformerUpdater) error {
	t := getPerformerTaggers(p, tagger.Cache)

	for _, tt := range t {
		if err := tt.tagScenes(ctx, paths, rw, func(o *models.Scene) (bool, error) {
			if err := o.LoadPerformerIDs(ctx, rw); err != nil {
				return false, err
			}
			existing := o.PerformerIDs.List()

			if slices.Contains(existing, p.ID) {
				return false, nil
			}

			if err := txn.WithTxn(ctx, tagger.TxnManager, func(ctx context.Context) error {
				return scene.AddPerformer(ctx, rw, o, p.ID)
			}); err != nil {
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
func (tagger *Tagger) PerformerImages(ctx context.Context, p *models.Performer, paths []string, rw ImageQueryPerformerUpdater) error {
	t := getPerformerTaggers(p, tagger.Cache)

	for _, tt := range t {
		if err := tt.tagImages(ctx, paths, rw, func(o *models.Image) (bool, error) {
			if err := o.LoadPerformerIDs(ctx, rw); err != nil {
				return false, err
			}
			existing := o.PerformerIDs.List()

			if slices.Contains(existing, p.ID) {
				return false, nil
			}

			if err := txn.WithTxn(ctx, tagger.TxnManager, func(ctx context.Context) error {
				return image.AddPerformer(ctx, rw, o, p.ID)
			}); err != nil {
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
func (tagger *Tagger) PerformerGalleries(ctx context.Context, p *models.Performer, paths []string, rw GalleryQueryPerformerUpdater) error {
	t := getPerformerTaggers(p, tagger.Cache)

	for _, tt := range t {
		if err := tt.tagGalleries(ctx, paths, rw, func(o *models.Gallery) (bool, error) {
			if err := o.LoadPerformerIDs(ctx, rw); err != nil {
				return false, err
			}
			existing := o.PerformerIDs.List()

			if slices.Contains(existing, p.ID) {
				return false, nil
			}

			if err := txn.WithTxn(ctx, tagger.TxnManager, func(ctx context.Context) error {
				return gallery.AddPerformer(ctx, rw, o, p.ID)
			}); err != nil {
				return false, err
			}

			return true, nil
		}); err != nil {
			return err
		}
	}
	return nil
}
