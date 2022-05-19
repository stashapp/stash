package autotag

import (
	"context"
	"database/sql"

	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
)

func addSceneStudio(ctx context.Context, sceneWriter SceneFinderUpdater, sceneID, studioID int) (bool, error) {
	// don't set if already set
	scene, err := sceneWriter.Find(ctx, sceneID)
	if err != nil {
		return false, err
	}

	if scene.StudioID.Valid {
		return false, nil
	}

	// set the studio id
	s := sql.NullInt64{Int64: int64(studioID), Valid: true}
	scenePartial := models.ScenePartial{
		ID:       sceneID,
		StudioID: &s,
	}

	if _, err := sceneWriter.Update(ctx, scenePartial); err != nil {
		return false, err
	}
	return true, nil
}

func addImageStudio(ctx context.Context, imageWriter ImageFinderUpdater, imageID, studioID int) (bool, error) {
	// don't set if already set
	image, err := imageWriter.Find(ctx, imageID)
	if err != nil {
		return false, err
	}

	if image.StudioID.Valid {
		return false, nil
	}

	// set the studio id
	s := sql.NullInt64{Int64: int64(studioID), Valid: true}
	imagePartial := models.ImagePartial{
		ID:       imageID,
		StudioID: &s,
	}

	if _, err := imageWriter.Update(ctx, imagePartial); err != nil {
		return false, err
	}
	return true, nil
}

func addGalleryStudio(ctx context.Context, galleryWriter GalleryFinderUpdater, galleryID, studioID int) (bool, error) {
	// don't set if already set
	gallery, err := galleryWriter.Find(ctx, galleryID)
	if err != nil {
		return false, err
	}

	if gallery.StudioID.Valid {
		return false, nil
	}

	// set the studio id
	s := sql.NullInt64{Int64: int64(studioID), Valid: true}
	galleryPartial := models.GalleryPartial{
		ID:       galleryID,
		StudioID: &s,
	}

	if _, err := galleryWriter.UpdatePartial(ctx, galleryPartial); err != nil {
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
	Find(ctx context.Context, id int) (*models.Scene, error)
	Update(ctx context.Context, updatedScene models.ScenePartial) (*models.Scene, error)
}

// StudioScenes searches for scenes whose path matches the provided studio name and tags the scene with the studio, if studio is not already set on the scene.
func StudioScenes(ctx context.Context, p *models.Studio, paths []string, aliases []string, rw SceneFinderUpdater, cache *match.Cache) error {
	t := getStudioTagger(p, aliases, cache)

	for _, tt := range t {
		if err := tt.tagScenes(ctx, paths, rw, func(subjectID, otherID int) (bool, error) {
			return addSceneStudio(ctx, rw, otherID, subjectID)
		}); err != nil {
			return err
		}
	}

	return nil
}

type ImageFinderUpdater interface {
	image.Queryer
	Find(ctx context.Context, id int) (*models.Image, error)
	Update(ctx context.Context, updatedImage models.ImagePartial) (*models.Image, error)
}

// StudioImages searches for images whose path matches the provided studio name and tags the image with the studio, if studio is not already set on the image.
func StudioImages(ctx context.Context, p *models.Studio, paths []string, aliases []string, rw ImageFinderUpdater, cache *match.Cache) error {
	t := getStudioTagger(p, aliases, cache)

	for _, tt := range t {
		if err := tt.tagImages(ctx, paths, rw, func(subjectID, otherID int) (bool, error) {
			return addImageStudio(ctx, rw, otherID, subjectID)
		}); err != nil {
			return err
		}
	}

	return nil
}

type GalleryFinderUpdater interface {
	gallery.Queryer
	Find(ctx context.Context, id int) (*models.Gallery, error)
	UpdatePartial(ctx context.Context, updatedGallery models.GalleryPartial) (*models.Gallery, error)
}

// StudioGalleries searches for galleries whose path matches the provided studio name and tags the gallery with the studio, if studio is not already set on the gallery.
func StudioGalleries(ctx context.Context, p *models.Studio, paths []string, aliases []string, rw GalleryFinderUpdater, cache *match.Cache) error {
	t := getStudioTagger(p, aliases, cache)

	for _, tt := range t {
		if err := tt.tagGalleries(ctx, paths, rw, func(subjectID, otherID int) (bool, error) {
			return addGalleryStudio(ctx, rw, otherID, subjectID)
		}); err != nil {
			return err
		}
	}

	return nil
}
