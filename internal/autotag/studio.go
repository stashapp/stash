package autotag

import (
	"context"

	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
)

// the following functions aren't used in Tagger because they assume
// use within a transaction

func addImageStudio(ctx context.Context, imageWriter ImageStudioUpdater, i *models.Image, studioID int) (bool, error) {
	// load existing studios
	if err := i.LoadStudioIDs(ctx, imageWriter); err != nil {
		return false, err
	}

	existing := i.StudioIDs.List()

	// don't add if already present
	for _, id := range existing {
		if id == studioID {
			return false, nil
		}
	}

	// add the studio id
	imagePartial := models.NewImagePartial()
	imagePartial.StudioIDs = &models.UpdateIDs{
		IDs:  []int{studioID},
		Mode: models.RelationshipUpdateModeAdd,
	}

	if _, err := imageWriter.UpdatePartial(ctx, i.ID, imagePartial); err != nil {
		return false, err
	}
	return true, nil
}

func addGalleryStudio(ctx context.Context, galleryWriter GalleryStudioUpdater, o *models.Gallery, studioID int) (bool, error) {
	// load existing studios
	if err := o.LoadStudioIDs(ctx, galleryWriter); err != nil {
		return false, err
	}

	existing := o.StudioIDs.List()

	// don't add if already present
	for _, id := range existing {
		if id == studioID {
			return false, nil
		}
	}

	// add the studio id
	galleryPartial := models.NewGalleryPartial()
	galleryPartial.StudioIDs = &models.UpdateIDs{
		IDs:  []int{studioID},
		Mode: models.RelationshipUpdateModeAdd,
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
			// load existing studios
			if err := o.LoadStudioIDs(ctx, rw); err != nil {
				return false, err
			}

			existing := o.StudioIDs.List()

			// don't add if already present
			for _, id := range existing {
				if id == p.ID {
					return false, nil
				}
			}

			// add the studio id
			scenePartial := models.NewScenePartial()
			scenePartial.StudioIDs = &models.UpdateIDs{
				IDs:  []int{p.ID},
				Mode: models.RelationshipUpdateModeAdd,
			}

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
func (tagger *Tagger) StudioImages(ctx context.Context, p *models.Studio, paths []string, aliases []string, rw ImageStudioUpdater) error {
	t := getStudioTagger(p, aliases, tagger.Cache)

	for _, tt := range t {
		if err := tt.tagImages(ctx, paths, rw, func(i *models.Image) (bool, error) {
			// load existing studios
			if err := i.LoadStudioIDs(ctx, rw); err != nil {
				return false, err
			}

			existing := i.StudioIDs.List()

			// don't add if already present
			for _, id := range existing {
				if id == p.ID {
					return false, nil
				}
			}

			// add the studio id
			imagePartial := models.NewImagePartial()
			imagePartial.StudioIDs = &models.UpdateIDs{
				IDs:  []int{p.ID},
				Mode: models.RelationshipUpdateModeAdd,
			}

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
func (tagger *Tagger) StudioGalleries(ctx context.Context, p *models.Studio, paths []string, aliases []string, rw GalleryStudioUpdater) error {
	t := getStudioTagger(p, aliases, tagger.Cache)

	for _, tt := range t {
		if err := tt.tagGalleries(ctx, paths, rw, func(o *models.Gallery) (bool, error) {
			// load existing studios
			if err := o.LoadStudioIDs(ctx, rw); err != nil {
				return false, err
			}

			existing := o.StudioIDs.List()

			// don't add if already present
			for _, id := range existing {
				if id == p.ID {
					return false, nil
				}
			}

			// add the studio id
			galleryPartial := models.NewGalleryPartial()
			galleryPartial.StudioIDs = &models.UpdateIDs{
				IDs:  []int{p.ID},
				Mode: models.RelationshipUpdateModeAdd,
			}

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
