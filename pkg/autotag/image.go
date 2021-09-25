package autotag

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
)

func imagePathsFilter(paths []string) *models.ImageFilterType {
	if paths == nil {
		return nil
	}

	sep := string(filepath.Separator)

	var ret *models.ImageFilterType
	var or *models.ImageFilterType
	for _, p := range paths {
		newOr := &models.ImageFilterType{}
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

func getMatchingImages(name string, paths []string, imageReader models.ImageReader) ([]*models.Image, error) {
	regex := getPathQueryRegex(name)
	organized := false
	filter := models.ImageFilterType{
		Path: &models.StringCriterionInput{
			Value:    "(?i)" + regex,
			Modifier: models.CriterionModifierMatchesRegex,
		},
		Organized: &organized,
	}

	filter.And = imagePathsFilter(paths)

	pp := models.PerPageAll
	images, _, _, _, err := imageReader.Query(&filter, &models.FindFilterType{
		PerPage: &pp,
	})

	if err != nil {
		return nil, fmt.Errorf("error querying images with regex '%s': %s", regex, err.Error())
	}

	var ret []*models.Image
	for _, p := range images {
		if nameMatchesPath(name, p.Path) {
			ret = append(ret, p)
		}
	}

	return ret, nil
}

func getImageFileTagger(s *models.Image) tagger {
	return tagger{
		ID:   s.ID,
		Type: "image",
		Name: s.GetTitle(),
		Path: s.Path,
	}
}

// ImagePerformers tags the provided image with performers whose name matches the image's path.
func ImagePerformers(s *models.Image, rw models.ImageReaderWriter, performerReader models.PerformerReader) error {
	t := getImageFileTagger(s)

	return t.tagPerformers(performerReader, func(subjectID, otherID int) (bool, error) {
		return image.AddPerformer(rw, subjectID, otherID)
	})
}

// ImageStudios tags the provided image with the first studio whose name matches the image's path.
//
// Images will not be tagged if studio is already set.
func ImageStudios(s *models.Image, rw models.ImageReaderWriter, studioReader models.StudioReader) error {
	if s.StudioID.Valid {
		// don't modify
		return nil
	}

	t := getImageFileTagger(s)

	return t.tagStudios(studioReader, func(subjectID, otherID int) (bool, error) {
		return addImageStudio(rw, subjectID, otherID)
	})
}

// ImageTags tags the provided image with tags whose name matches the image's path.
func ImageTags(s *models.Image, rw models.ImageReaderWriter, tagReader models.TagReader) error {
	t := getImageFileTagger(s)

	return t.tagTags(tagReader, func(subjectID, otherID int) (bool, error) {
		return image.AddTag(rw, subjectID, otherID)
	})
}
