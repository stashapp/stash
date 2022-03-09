package autotag

import (
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
)

func getImageFileTagger(s *models.Image, cache *match.Cache) tagger {
	return tagger{
		ID:    s.ID,
		Type:  "image",
		Name:  s.GetTitle(),
		Path:  s.Path,
		cache: cache,
	}
}

// ImagePerformers tags the provided image with performers whose name matches the image's path.
func ImagePerformers(s *models.Image, rw models.ImageReaderWriter, performerReader models.PerformerReader, cache *match.Cache) error {
	t := getImageFileTagger(s, cache)

	return t.tagPerformers(performerReader, func(subjectID, otherID int) (bool, error) {
		return image.AddPerformer(rw, subjectID, otherID)
	})
}

// ImageStudios tags the provided image with the first studio whose name matches the image's path.
//
// Images will not be tagged if studio is already set.
func ImageStudios(s *models.Image, rw models.ImageReaderWriter, studioReader models.StudioReader, cache *match.Cache) error {
	if s.StudioID.Valid {
		// don't modify
		return nil
	}

	t := getImageFileTagger(s, cache)

	return t.tagStudios(studioReader, func(subjectID, otherID int) (bool, error) {
		return addImageStudio(rw, subjectID, otherID)
	})
}

// ImageTags tags the provided image with tags whose name matches the image's path.
func ImageTags(s *models.Image, rw models.ImageReaderWriter, tagReader models.TagReader, cache *match.Cache) error {
	t := getImageFileTagger(s, cache)

	return t.tagTags(tagReader, func(subjectID, otherID int) (bool, error) {
		return image.AddTag(rw, subjectID, otherID)
	})
}
