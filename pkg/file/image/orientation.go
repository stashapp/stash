package image

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

func adjustForOrientation(fs models.FS, path string, f *models.ImageFile) {
	isFlipped, err := areDimensionsFlipped(fs, path)
	if err != nil {
		logger.Warnf("Error determining image orientation for %s: %v", path, err)
		// isFlipped is false by default
	}

	if isFlipped {
		f.Width, f.Height = f.Height, f.Width
	}
}

// areDimensionsFlipped returns true if the image dimensions are flipped.
// This is determined by the EXIF orientation tag.
func areDimensionsFlipped(fs models.FS, path string) (bool, error) {
	r, err := fs.Open(path)
	if err != nil {
		return false, fmt.Errorf("reading image file %q: %w", path, err)
	}
	defer r.Close()

	x, err := exif.Decode(r)
	if err != nil {
		if errors.Is(err, io.EOF) || strings.Contains(err.Error(), "failed to find exif") {
			// no exif data
			return false, nil
		}

		return false, fmt.Errorf("decoding exif data: %w", err)
	}

	o, err := x.Get(exif.Orientation)
	if err != nil {
		// assume not present
		return false, nil
	}

	oo, err := o.Int(0)
	if err != nil {
		return false, fmt.Errorf("decoding orientation: %w", err)
	}

	return isOrientationDimensionsFlipped(oo), nil
}

// isOrientationDimensionsFlipped returns true if the image orientation is flipped based on the input orientation EXIF value.
// From https://sirv.com/help/articles/rotate-photos-to-be-upright/
// 1 = 0 degrees: the correct orientation, no adjustment is required.
// 2 = 0 degrees, mirrored: image has been flipped back-to-front.
// 3 = 180 degrees: image is upside down.
// 4 = 180 degrees, mirrored: image has been flipped back-to-front and is upside down.
// 5 = 90 degrees: image has been flipped back-to-front and is on its side.
// 6 = 90 degrees, mirrored: image is on its side.
// 7 = 270 degrees: image has been flipped back-to-front and is on its far side.
// 8 = 270 degrees, mirrored: image is on its far side.
func isOrientationDimensionsFlipped(o int) bool {
	switch o {
	case 5, 6, 7, 8:
		return true
	default:
		return false
	}
}
