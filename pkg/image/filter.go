package image

import (
	"path/filepath"
	"strings"

	"github.com/stashapp/stash/pkg/models"
)

func PathsFilter(paths []string) *models.ImageFilterType {
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
			p += sep
		}

		or.Path = &models.StringCriterionInput{
			Modifier: models.CriterionModifierEquals,
			Value:    p + "%",
		}
	}

	return ret
}
