package autotag

import (
	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
)

func getMatchingTags(path string, tagReader models.TagReader) ([]*models.Tag, error) {
	words := getPathWords(path)
	tags, err := tagReader.QueryForAutoTag(words)

	if err != nil {
		return nil, err
	}

	var ret []*models.Tag
	for _, p := range tags {
		if nameMatchesPath(p.Name, path) {
			ret = append(ret, p)
		}
	}

	return ret, nil
}

func getTagTagger(p *models.Tag) tagger {
	return tagger{
		ID:   p.ID,
		Type: "tag",
		Name: p.Name,
	}
}

// TagScenes searches for scenes whose path matches the provided tag name and tags the scene with the tag.
func TagScenes(p *models.Tag, paths []string, rw models.SceneReaderWriter) error {
	t := getTagTagger(p)

	return t.tagScenes(paths, rw, func(subjectID, otherID int) (bool, error) {
		return scene.AddTag(rw, otherID, subjectID)
	})
}

// TagImages searches for images whose path matches the provided tag name and tags the image with the tag.
func TagImages(p *models.Tag, paths []string, rw models.ImageReaderWriter) error {
	t := getTagTagger(p)

	return t.tagImages(paths, rw, func(subjectID, otherID int) (bool, error) {
		return image.AddTag(rw, otherID, subjectID)
	})
}

// TagGalleries searches for galleries whose path matches the provided tag name and tags the gallery with the tag.
func TagGalleries(p *models.Tag, paths []string, rw models.GalleryReaderWriter) error {
	t := getTagTagger(p)

	return t.tagGalleries(paths, rw, func(subjectID, otherID int) (bool, error) {
		return gallery.AddTag(rw, otherID, subjectID)
	})
}
