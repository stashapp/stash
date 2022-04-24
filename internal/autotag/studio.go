package autotag

import (
	"context"

	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
)

func addSceneStudio(ctx context.Context, sceneWriter scene.PartialUpdater, o *models.Scene, studioID int) (bool, error) {
	// don't set if already set
	if o.StudioID != nil {
		return false, nil
	}

	// set the studio id
	scenePartial := models.ScenePartial{
		StudioID: models.NewOptionalInt(studioID),
	}

	if _, err := sceneWriter.UpdatePartial(ctx, o.ID, scenePartial); err != nil {
		return false, err
	}
	return true, nil
}

func addImageStudio(ctx context.Context, imageWriter image.PartialUpdater, i *models.Image, studioID int) (bool, error) {
	// don't set if already set
	if i.StudioID != nil {
		return false, nil
	}

	// set the studio id
	imagePartial := models.ImagePartial{
		StudioID: models.NewOptionalInt(studioID),
	}

	if _, err := imageWriter.UpdatePartial(ctx, i.ID, imagePartial); err != nil {
		return false, err
	}
	return true, nil
}

func addGalleryStudio(ctx context.Context, galleryWriter GalleryFinderUpdater, o *models.Gallery, studioID int) (bool, error) {
	// don't set if already set
	if o.StudioID != nil {
		return false, nil
	}

	// set the studio id
	galleryPartial := models.GalleryPartial{
		StudioID: models.NewOptionalInt(studioID),
	}

	if _, err := galleryWriter.UpdatePartial(ctx, o.ID, galleryPartial); err != nil {
		return false, err
	}
	return true, nil
}

func getStudioTagger(p *models.Studio, aliases []string, cache *match.Cache) []tagger {
	ret := []tagger{{
		ID:    p.ID,
		Type:  "studio",
		Name:  p.Name.String,
		cache: cache,
	}}

	for _, a := range aliases {
		ret = append(ret, tagger{
			ID:   p.ID,
			Type: "studio",
			Name: a,
		})
	}

	return ret
}

type SceneFinderUpdater interface {
	scene.Queryer
	scene.PartialUpdater
}

// StudioScenes searches for scenes whose path matches the provided studio name and tags the scene with the studio, if studio is not already set on the scene.
func StudioScenes(ctx context.Context, p *models.Studio, paths []string, aliases []string, rw SceneFinderUpdater, cache *match.Cache) error {
	t := getStudioTagger(p, aliases, cache)

	for _, tt := range t {
		if err := tt.tagScenes(ctx, paths, rw, func(o *models.Scene) (bool, error) {
			return addSceneStudio(ctx, rw, o, p.ID)
		}); err != nil {
			return err
		}
	}

	return nil
}

type ImageFinderUpdater interface {
	image.Queryer
	Find(ctx context.Context, id int) (*models.Image, error)
	UpdatePartial(ctx context.Context, id int, partial models.ImagePartial) (*models.Image, error)
}

// StudioImages searches for images whose path matches the provided studio name and tags the image with the studio, if studio is not already set on the image.
func StudioImages(ctx context.Context, p *models.Studio, paths []string, aliases []string, rw ImageFinderUpdater, cache *match.Cache) error {
	t := getStudioTagger(p, aliases, cache)

	for _, tt := range t {
		if err := tt.tagImages(ctx, paths, rw, func(i *models.Image) (bool, error) {
			return addImageStudio(ctx, rw, i, p.ID)
		}); err != nil {
			return err
		}
	}

	return nil
}

type GalleryFinderUpdater interface {
	gallery.Queryer
	gallery.PartialUpdater
	Find(ctx context.Context, id int) (*models.Gallery, error)
}

// StudioGalleries searches for galleries whose path matches the provided studio name and tags the gallery with the studio, if studio is not already set on the gallery.
func StudioGalleries(ctx context.Context, p *models.Studio, paths []string, aliases []string, rw GalleryFinderUpdater, cache *match.Cache) error {
	t := getStudioTagger(p, aliases, cache)

	for _, tt := range t {
		if err := tt.tagGalleries(ctx, paths, rw, func(o *models.Gallery) (bool, error) {
			return addGalleryStudio(ctx, rw, o, p.ID)
		}); err != nil {
			return err
		}
	}

	return nil
}
