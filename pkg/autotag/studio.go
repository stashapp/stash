package autotag

import (
	"database/sql"

	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
)

func addSceneStudio(sceneWriter models.SceneReaderWriter, sceneID, studioID int) (bool, error) {
	// don't set if already set
	scene, err := sceneWriter.Find(sceneID)
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

	if _, err := sceneWriter.Update(scenePartial); err != nil {
		return false, err
	}
	return true, nil
}

func addImageStudio(imageWriter models.ImageReaderWriter, imageID, studioID int) (bool, error) {
	// don't set if already set
	image, err := imageWriter.Find(imageID)
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

	if _, err := imageWriter.Update(imagePartial); err != nil {
		return false, err
	}
	return true, nil
}

func addGalleryStudio(galleryWriter models.GalleryReaderWriter, galleryID, studioID int) (bool, error) {
	// don't set if already set
	gallery, err := galleryWriter.Find(galleryID)
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

	if _, err := galleryWriter.UpdatePartial(galleryPartial); err != nil {
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

// StudioScenes searches for scenes whose path matches the provided studio name and tags the scene with the studio, if studio is not already set on the scene.
func StudioScenes(p *models.Studio, paths []string, aliases []string, rw models.SceneReaderWriter, cache *match.Cache) error {
	t := getStudioTagger(p, aliases, cache)

	for _, tt := range t {
		if err := tt.tagScenes(paths, rw, func(subjectID, otherID int) (bool, error) {
			return addSceneStudio(rw, otherID, subjectID)
		}); err != nil {
			return err
		}
	}

	return nil
}

// StudioImages searches for images whose path matches the provided studio name and tags the image with the studio, if studio is not already set on the image.
func StudioImages(p *models.Studio, paths []string, aliases []string, rw models.ImageReaderWriter, cache *match.Cache) error {
	t := getStudioTagger(p, aliases, cache)

	for _, tt := range t {
		if err := tt.tagImages(paths, rw, func(subjectID, otherID int) (bool, error) {
			return addImageStudio(rw, otherID, subjectID)
		}); err != nil {
			return err
		}
	}

	return nil
}

// StudioGalleries searches for galleries whose path matches the provided studio name and tags the gallery with the studio, if studio is not already set on the gallery.
func StudioGalleries(p *models.Studio, paths []string, aliases []string, rw models.GalleryReaderWriter, cache *match.Cache) error {
	t := getStudioTagger(p, aliases, cache)

	for _, tt := range t {
		if err := tt.tagGalleries(paths, rw, func(subjectID, otherID int) (bool, error) {
			return addGalleryStudio(rw, otherID, subjectID)
		}); err != nil {
			return err
		}
	}

	return nil
}
