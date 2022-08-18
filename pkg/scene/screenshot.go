package scene

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"io"
	"os"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/paths"

	"github.com/disintegration/imaging"

	// needed to decode other image formats
	_ "image/gif"
	_ "image/png"
)

type screenshotter interface {
	GenerateScreenshot(ctx context.Context, probeResult *ffmpeg.VideoFile, hash string) error
	GenerateThumbnail(ctx context.Context, probeResult *ffmpeg.VideoFile, hash string) error
}

type ScreenshotSetter interface {
	SetScreenshot(scene *models.Scene, imageData []byte) error
}

type PathsScreenshotSetter struct {
	Paths               *paths.Paths
	FileNamingAlgorithm models.HashAlgorithm
}

func (ss *PathsScreenshotSetter) SetScreenshot(scene *models.Scene, imageData []byte) error {
	checksum := scene.GetHash(ss.FileNamingAlgorithm)
	return SetScreenshot(ss.Paths, checksum, imageData)
}

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

func SetScreenshot(paths *paths.Paths, checksum string, imageData []byte) error {
	thumbPath := paths.Scene.GetThumbnailScreenshotPath(checksum)
	normalPath := paths.Scene.GetScreenshotPath(checksum)

	err := SetThumbnail(thumbPath, bytes.NewReader(imageData))
	if err != nil {
		return err
	}

	err = writeImage(normalPath, imageData)

	return err
}

func SetThumbnail(thumbPath string, reader io.Reader) error {
	screenshot, _, err := image.Decode(reader)
	if err != nil {
		return err
	}

	// resize to 320 width maintaining aspect ratio, for the thumbnail
	const width = 320
	origWidth := screenshot.Bounds().Max.X
	origHeight := screenshot.Bounds().Max.Y
	height := width / origWidth * origHeight

	thumbnail := imaging.Resize(screenshot, width, height, imaging.Lanczos)
	return writeThumbnail(thumbPath, thumbnail)
}
