package autotag

import (
	"context"

	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
)

// the following functions aren't used in Tagger because they assume
// use within a transaction

func addSceneStudio(ctx context.Context, sceneWriter models.SceneUpdater, o *models.Scene, studioID int) (bool, error) {
	// don't set if already set
	if o.StudioID != nil {
		return false, nil
	}

	// set the studio id
	scenePartial := models.NewScenePartial()
	scenePartial.StudioID = models.NewOptionalInt(studioID)

	if _, err := sceneWriter.UpdatePartial(ctx, o.ID, scenePartial); err != nil {
		return false, err
	}
	return true, nil
}

func addImageStudio(ctx context.Context, imageWriter models.ImageUpdater, i *models.Image, studioID int) (bool, error) {
	// don't set if already set
	if i.StudioID != nil {
		return false, nil
	}

	// set the studio id
	imagePartial := models.NewImagePartial()
	imagePartial.StudioID = models.NewOptionalInt(studioID)

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
	galleryPartial := models.NewGalleryPartial()
	galleryPartial.StudioID = models.NewOptionalInt(studioID)

	if _, err := galleryWriter.UpdatePartial(ctx, o.ID, galleryPartial); err != nil {
		return false, err
	}
	return true, nil
}

func getStudioTagger(p *models.Studio, aliases []string, cache *match.Cache) []tagger {
	ret := []tagger{{
		ID:    p.ID,
		Type:  "studio",
		Name:  p.Name,
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

// StudioScenes searches for scenes whose path matches the provided studio name and tags the scene with the studio, if studio is not already set on the scene.
func (tagger *Tagger) StudioScenes(ctx context.Context, p *models.Studio, paths []string, aliases []string, rw SceneFinderUpdater) error {
	t := getStudioTagger(p, aliases, tagger.Cache)

	for _, tt := range t {
		if err := tt.tagScenes(ctx, paths, rw, func(o *models.Scene) (bool, error) {
			// don't set if already set
			if o.StudioID != nil {
				return false, nil
			}

			// set the studio id
			scenePartial := models.NewScenePartial()
			scenePartial.StudioID = models.NewOptionalInt(p.ID)

			if err := txn.WithTxn(ctx, tagger.TxnManager, func(ctx context.Context) error {
				_, err := rw.UpdatePartial(ctx, o.ID, scenePartial)
				return err
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

// StudioImages searches for images whose path matches the provided studio name and tags the image with the studio, if studio is not already set on the image.
func (tagger *Tagger) StudioImages(ctx context.Context, p *models.Studio, paths []string, aliases []string, rw ImageFinderUpdater) error {
	t := getStudioTagger(p, aliases, tagger.Cache)

	for _, tt := range t {
		if err := tt.tagImages(ctx, paths, rw, func(i *models.Image) (bool, error) {
			// don't set if already set
			if i.StudioID != nil {
				return false, nil
			}

			// set the studio id
			imagePartial := models.NewImagePartial()
			imagePartial.StudioID = models.NewOptionalInt(p.ID)

			if err := txn.WithTxn(ctx, tagger.TxnManager, func(ctx context.Context) error {
				_, err := rw.UpdatePartial(ctx, i.ID, imagePartial)
				return err
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

// StudioGalleries searches for galleries whose path matches the provided studio name and tags the gallery with the studio, if studio is not already set on the gallery.
func (tagger *Tagger) StudioGalleries(ctx context.Context, p *models.Studio, paths []string, aliases []string, rw GalleryFinderUpdater) error {
	t := getStudioTagger(p, aliases, tagger.Cache)

	for _, tt := range t {
		if err := tt.tagGalleries(ctx, paths, rw, func(o *models.Gallery) (bool, error) {
			// don't set if already set
			if o.StudioID != nil {
				return false, nil
			}

			// set the studio id
			galleryPartial := models.NewGalleryPartial()
			galleryPartial.StudioID = models.NewOptionalInt(p.ID)

			if err := txn.WithTxn(ctx, tagger.TxnManager, func(ctx context.Context) error {
				_, err := rw.UpdatePartial(ctx, o.ID, galleryPartial)
				return err
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
