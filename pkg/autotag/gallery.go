package autotag

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/models"
)

func galleryPathsFilter(paths []string) *models.GalleryFilterType {
	if paths == nil {
		return nil
	}

	sep := string(filepath.Separator)

	var ret *models.GalleryFilterType
	var or *models.GalleryFilterType
	for _, p := range paths {
		newOr := &models.GalleryFilterType{}
		if or != nil {
			or.Or = newOr
		} else {
			ret = newOr
		}

		or = newOr

		if !strings.HasSuffix(p, sep) {
			p = p + sep
		}

		or.Path = &models.StringCriterionInput{
			Modifier: models.CriterionModifierEquals,
			Value:    p + "%",
		}
	}

	return ret
}

func getMatchingGalleries(name string, paths []string, galleryReader models.GalleryReader) ([]*models.Gallery, error) {
	regex := getPathQueryRegex(name)
	organized := false
	filter := models.GalleryFilterType{
		Path: &models.StringCriterionInput{
			Value:    "(?i)" + regex,
			Modifier: models.CriterionModifierMatchesRegex,
		},
		Organized: &organized,
	}

	filter.And = galleryPathsFilter(paths)

	pp := models.PerPageAll
	gallerys, _, err := galleryReader.Query(&filter, &models.FindFilterType{
		PerPage: &pp,
	})

	if err != nil {
		return nil, fmt.Errorf("error querying gallerys with regex '%s': %s", regex, err.Error())
	}

	var ret []*models.Gallery
	for _, p := range gallerys {
		if nameMatchesPath(name, p.Path.String) {
			ret = append(ret, p)
		}
	}

	return ret, nil
}

func getGalleryFileTagger(s *models.Gallery) tagger {
	return tagger{
		ID:   s.ID,
		Type: "gallery",
		Name: s.GetTitle(),
		Path: s.Path.String,
	}
}

// GalleryPerformers tags the provided gallery with performers whose name matches the gallery's path.
func GalleryPerformers(s *models.Gallery, rw models.GalleryReaderWriter, performerReader models.PerformerReader) error {
	t := getGalleryFileTagger(s)

	return t.tagPerformers(performerReader, func(subjectID, otherID int) (bool, error) {
		return gallery.AddPerformer(rw, subjectID, otherID)
	})
}

// GalleryStudios tags the provided gallery with the first studio whose name matches the gallery's path.
//
// Gallerys will not be tagged if studio is already set.
func GalleryStudios(s *models.Gallery, rw models.GalleryReaderWriter, studioReader models.StudioReader) error {
	if s.StudioID.Valid {
		// don't modify
		return nil
	}

	t := getGalleryFileTagger(s)

	return t.tagStudios(studioReader, func(subjectID, otherID int) (bool, error) {
		return addGalleryStudio(rw, subjectID, otherID)
	})
}

// GalleryTags tags the provided gallery with tags whose name matches the gallery's path.
func GalleryTags(s *models.Gallery, rw models.GalleryReaderWriter, tagReader models.TagReader) error {
	t := getGalleryFileTagger(s)

	return t.tagTags(tagReader, func(subjectID, otherID int) (bool, error) {
		return gallery.AddTag(rw, subjectID, otherID)
	})
}
