package autotag

import (
	"context"

	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
)

type SceneQueryTagUpdater interface {
	scene.Queryer
	scene.PartialUpdater
}

type ImageQueryTagUpdater interface {
	image.Queryer
	image.PartialUpdater
}

type GalleryQueryTagUpdater interface {
	gallery.Queryer
	gallery.PartialUpdater
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
func TagScenes(ctx context.Context, p *models.Tag, paths []string, aliases []string, rw SceneQueryTagUpdater, cache *match.Cache) error {
	t := getTagTaggers(p, aliases, cache)

	for _, tt := range t {
		if err := tt.tagScenes(ctx, paths, rw, func(o *models.Scene) (bool, error) {
			return scene.AddTag(ctx, rw, o, p.ID)
		}); err != nil {
			return err
		}
	}
	return nil
}

// TagImages searches for images whose path matches the provided tag name and tags the image with the tag.
func TagImages(ctx context.Context, p *models.Tag, paths []string, aliases []string, rw ImageQueryTagUpdater, cache *match.Cache) error {
	t := getTagTaggers(p, aliases, cache)

	for _, tt := range t {
		if err := tt.tagImages(ctx, paths, rw, func(i *models.Image) (bool, error) {
			return image.AddTag(ctx, rw, i, p.ID)
		}); err != nil {
			return err
		}
	}
	return nil
}

// TagGalleries searches for galleries whose path matches the provided tag name and tags the gallery with the tag.
func TagGalleries(ctx context.Context, p *models.Tag, paths []string, aliases []string, rw GalleryQueryTagUpdater, cache *match.Cache) error {
	t := getTagTaggers(p, aliases, cache)

	for _, tt := range t {
		if err := tt.tagGalleries(ctx, paths, rw, func(o *models.Gallery) (bool, error) {
			return gallery.AddTag(ctx, rw, o, p.ID)
		}); err != nil {
			return err
		}
	}
	return nil
}
