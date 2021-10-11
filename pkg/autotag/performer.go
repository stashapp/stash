package autotag

import (
	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
)

func getPerformerTagger(p *models.Performer) tagger {
	return tagger{
		ID:   p.ID,
		Type: "performer",
		Name: p.Name.String,
	}
}

// PerformerScenes searches for scenes whose path matches the provided performer name and tags the scene with the performer.
func PerformerScenes(p *models.Performer, paths []string, rw models.SceneReaderWriter) error {
	t := getPerformerTagger(p)

	return t.tagScenes(paths, rw, func(subjectID, otherID int) (bool, error) {
		return scene.AddPerformer(rw, otherID, subjectID)
	})
}

// PerformerImages searches for images whose path matches the provided performer name and tags the image with the performer.
func PerformerImages(p *models.Performer, paths []string, rw models.ImageReaderWriter) error {
	t := getPerformerTagger(p)

	return t.tagImages(paths, rw, func(subjectID, otherID int) (bool, error) {
		return image.AddPerformer(rw, otherID, subjectID)
	})
}

// PerformerGalleries searches for galleries whose path matches the provided performer name and tags the gallery with the performer.
func PerformerGalleries(p *models.Performer, paths []string, rw models.GalleryReaderWriter) error {
	t := getPerformerTagger(p)

	return t.tagGalleries(paths, rw, func(subjectID, otherID int) (bool, error) {
		return gallery.AddPerformer(rw, otherID, subjectID)
	})
}
