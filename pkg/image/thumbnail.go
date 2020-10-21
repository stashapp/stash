package image

import (
	"bytes"
	"image"
	"image/jpeg"

	"github.com/disintegration/imaging"
)

func ThumbnailNeeded(srcImage image.Image, maxSize int) bool {
	dim := srcImage.Bounds().Max
	w := dim.X
	h := dim.Y

	return w > maxSize || h > maxSize
}

// GetThumbnail returns the thumbnail image of the provided image resized to
// the provided max size. It resizes based on the largest X/Y direction.
// It returns nil and an error if an error occurs reading, decoding or encoding
// the image.
func GetThumbnail(srcImage image.Image, maxSize int) ([]byte, error) {
	var resizedImage image.Image

	// if height is longer then resize by height instead of width
	dim := srcImage.Bounds().Max
	if dim.Y > dim.X {
		resizedImage = imaging.Resize(srcImage, 0, maxSize, imaging.Box)
	} else {
		resizedImage = imaging.Resize(srcImage, maxSize, 0, imaging.Box)
	}

	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, resizedImage, nil)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
