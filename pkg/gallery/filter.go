package gallery

import (
	"path/filepath"
	"strings"

	"github.com/stashapp/stash/pkg/models"
)

func PathsFilter(paths []string) *models.GalleryFilterType {
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
			p += sep
		}

		or.Path = &models.StringCriterionInput{
			Modifier: models.CriterionModifierEquals,
			Value:    p + "%",
		}
	}

	return ret
}
