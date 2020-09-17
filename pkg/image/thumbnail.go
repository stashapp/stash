package image

import (
	"bytes"
	"image"
	"image/jpeg"

	"github.com/disintegration/imaging"
)

// GetThumbnail returns the thumbnail image of the provided image resized to
// the provided width. It returns nil and an error if an error occurs reading,
// decoding or encoding the image.
func GetThumbnail(srcImage image.Image, width int) ([]byte, error) {
	resizedImage := imaging.Resize(srcImage, width, 0, imaging.Box)
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, resizedImage, nil)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
