package manager

import (
	"bytes"
	"image"
	"image/jpeg"
	"os"

	"github.com/disintegration/imaging"

	// needed to decode other image formats
	_ "image/gif"
	_ "image/png"
)

func writeImage(path string, imageData []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(imageData)
	return err
}

func writeThumbnail(path string, thumbnail image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return jpeg.Encode(f, thumbnail, nil)
}

func SetSceneScreenshot(checksum string, imageData []byte) error {
	thumbPath := instance.Paths.Scene.GetThumbnailScreenshotPath(checksum)
	normalPath := instance.Paths.Scene.GetScreenshotPath(checksum)

	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return err
	}

	// resize to 320 width maintaining aspect ratio, for the thumbnail
	const width = 320
	origWidth := img.Bounds().Max.X
	origHeight := img.Bounds().Max.Y
	height := width / origWidth * origHeight

	thumbnail := imaging.Resize(img, width, height, imaging.Lanczos)
	err = writeThumbnail(thumbPath, thumbnail)
	if err != nil {
		return err
	}

	err = writeImage(normalPath, imageData)

	return err
}
