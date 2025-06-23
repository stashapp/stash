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

type SceneQueryTagUpdater interface {
	models.SceneQueryer
	models.TagIDLoader
	models.SceneUpdater
}

type ImageQueryTagUpdater interface {
	models.ImageQueryer
	models.TagIDLoader
	models.ImageUpdater
}

type GalleryQueryTagUpdater interface {
	models.GalleryQueryer
	models.TagIDLoader
	models.GalleryUpdater
}

func getTagTaggers(p *models.Tag, aliases []string, cache *match.Cache) []tagger {
	ret := []tagger{{
		ID:    p.ID,
		Type:  "tag",
		Name:  p.Name,
		cache: cache,
	}}

	for _, a := range aliases {
		ret = append(ret, tagger{
			ID:    p.ID,
			Type:  "tag",
			Name:  a,
			cache: cache,
		})
	}

	return ret
}

// TagScenes searches for scenes whose path matches the provided tag name and tags the scene with the tag.
func (tagger *Tagger) TagScenes(ctx context.Context, p *models.Tag, paths []string, aliases []string, rw SceneQueryTagUpdater) error {
	t := getTagTaggers(p, aliases, tagger.Cache)

	for _, tt := range t {
		if err := tt.tagScenes(ctx, paths, rw, func(o *models.Scene) (bool, error) {
			if err := o.LoadTagIDs(ctx, rw); err != nil {
				return false, err
			}
			existing := o.TagIDs.List()

			if slices.Contains(existing, p.ID) {
				return false, nil
			}

			if err := txn.WithTxn(ctx, tagger.TxnManager, func(ctx context.Context) error {
				return scene.AddTag(ctx, rw, o, p.ID)
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

// TagImages searches for images whose path matches the provided tag name and tags the image with the tag.
func (tagger *Tagger) TagImages(ctx context.Context, p *models.Tag, paths []string, aliases []string, rw ImageQueryTagUpdater) error {
	t := getTagTaggers(p, aliases, tagger.Cache)

	for _, tt := range t {
		if err := tt.tagImages(ctx, paths, rw, func(o *models.Image) (bool, error) {
			if err := o.LoadTagIDs(ctx, rw); err != nil {
				return false, err
			}
			existing := o.TagIDs.List()

			if slices.Contains(existing, p.ID) {
				return false, nil
			}

			if err := txn.WithTxn(ctx, tagger.TxnManager, func(ctx context.Context) error {
				return image.AddTag(ctx, rw, o, p.ID)
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

// TagGalleries searches for galleries whose path matches the provided tag name and tags the gallery with the tag.
func (tagger *Tagger) TagGalleries(ctx context.Context, p *models.Tag, paths []string, aliases []string, rw GalleryQueryTagUpdater) error {
	t := getTagTaggers(p, aliases, tagger.Cache)

	for _, tt := range t {
		if err := tt.tagGalleries(ctx, paths, rw, func(o *models.Gallery) (bool, error) {
			if err := o.LoadTagIDs(ctx, rw); err != nil {
				return false, err
			}
			existing := o.TagIDs.List()

			if slices.Contains(existing, p.ID) {
				return false, nil
			}

			if err := txn.WithTxn(ctx, tagger.TxnManager, func(ctx context.Context) error {
				return gallery.AddTag(ctx, rw, o, p.ID)
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
