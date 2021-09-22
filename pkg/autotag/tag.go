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
	for _, t := range tags {
		matches := false
		if nameMatchesPath(t.Name, path) {
			matches = true
		}

		if !matches {
			aliases, err := tagReader.GetAliases(t.ID)
			if err != nil {
				return nil, err
			}
			for _, alias := range aliases {
				if nameMatchesPath(alias, path) {
					matches = true
					break
				}
			}
		}

		if matches {
			ret = append(ret, t)
		}
	}

	return ret, nil
}

func getTagTaggers(p *models.Tag, aliases []string) []tagger {
	ret := []tagger{{
		ID:   p.ID,
		Type: "tag",
		Name: p.Name,
	}}

	for _, a := range aliases {
		ret = append(ret, tagger{
			ID:   p.ID,
			Type: "tag",
			Name: a,
		})
	}

	return ret
}

// TagScenes searches for scenes whose path matches the provided tag name and tags the scene with the tag.
func TagScenes(p *models.Tag, paths []string, aliases []string, rw models.SceneReaderWriter) error {
	t := getTagTaggers(p, aliases)

	for _, tt := range t {
		if err := tt.tagScenes(paths, rw, func(subjectID, otherID int) (bool, error) {
			return scene.AddTag(rw, otherID, subjectID)
		}); err != nil {
			return err
		}
	}
	return nil
}

// TagImages searches for images whose path matches the provided tag name and tags the image with the tag.
func TagImages(p *models.Tag, paths []string, aliases []string, rw models.ImageReaderWriter) error {
	t := getTagTaggers(p, aliases)

	for _, tt := range t {
		if err := tt.tagImages(paths, rw, func(subjectID, otherID int) (bool, error) {
			return image.AddTag(rw, otherID, subjectID)
		}); err != nil {
			return err
		}
	}
	return nil
}

// TagGalleries searches for galleries whose path matches the provided tag name and tags the gallery with the tag.
func TagGalleries(p *models.Tag, paths []string, aliases []string, rw models.GalleryReaderWriter) error {
	t := getTagTaggers(p, aliases)

	for _, tt := range t {
		if err := tt.tagGalleries(paths, rw, func(subjectID, otherID int) (bool, error) {
			return gallery.AddTag(rw, otherID, subjectID)
		}); err != nil {
			return err
		}
	}
	return nil
}
