package image

import (
	"image"
	"io"
	"path/filepath"
	"strings"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
	_ "golang.org/x/image/webp"
)

func DecodeSourceImage(r io.Reader) (*image.Config, *string, error) {
	config, format, err := image.DecodeConfig(r)

	return &config, &format, err
}

func IsCover(img *models.Image) bool {
	return strings.HasSuffix(strings.ToLower(img.Path), "cover.jpg")
}

func GetTitle(s *models.Image) string {
	if s.Title.String != "" {
		return s.Title.String
	}

	return filepath.Base(s.Path)
}

// GetFilename gets the base name of the image file
// If stripExt is set the file extension is omitted from the name
func GetFilename(s *models.Image, stripExt bool) string {
	return utils.GetNameFromPath(s.Path, stripExt)
}
