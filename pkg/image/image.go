package image

import (
	"strings"

	"github.com/stashapp/stash/pkg/models"
	_ "golang.org/x/image/webp"
)

func IsCover(img *models.Image) bool {
	return strings.HasSuffix(img.Path(), "cover.jpg")
}
