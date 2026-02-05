package imagephash

import (
	"bytes"
	"fmt"
	"image"

	"github.com/corona10/goimagehash"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/models"
)

// Generate computes a perceptual hash for an image file.
func Generate(imageFile *models.ImageFile) (*uint64, error) {
	img, err := loadImage(imageFile)
	if err != nil {
		return nil, fmt.Errorf("loading image: %w", err)
	}

	hash, err := goimagehash.PerceptionHash(img)
	if err != nil {
		return nil, fmt.Errorf("computing phash from image: %w", err)
	}

	hashValue := hash.GetHash()
	return &hashValue, nil
}

// loadImage loads an image from disk and decodes it.
func loadImage(imageFile *models.ImageFile) (image.Image, error) {
	reader, err := imageFile.Open(&file.OsFS{})
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(reader); err != nil {
		return nil, err
	}

	img, _, err := image.Decode(buf)
	if err != nil {
		return nil, fmt.Errorf("decoding image: %w", err)
	}

	return img, nil
}
