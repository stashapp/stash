package scene

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"os"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/paths"

	"github.com/disintegration/imaging"

	// needed to decode other image formats
	_ "image/gif"
	_ "image/png"
)

type CoverGenerator interface {
	GenerateCover(ctx context.Context, scene *models.Scene, f *file.VideoFile) error
}

type ScreenshotSetter interface {
	SetScreenshot(scene *models.Scene, imageData []byte) error
}

type PathsCoverSetter struct {
	Paths               *paths.Paths
	FileNamingAlgorithm models.HashAlgorithm
}

func (ss *PathsCoverSetter) SetScreenshot(scene *models.Scene, imageData []byte) error {
	// don't set where scene has no file
	if scene.Path == "" {
		return nil
	}
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

func (s *Service) GetCover(ctx context.Context, scene *models.Scene) ([]byte, error) {
	if scene.Path != "" {
		filepath := s.Paths.Scene.GetScreenshotPath(scene.GetHash(s.Config.GetVideoFileNamingAlgorithm()))

		// fall back to the scene image blob if the file isn't present
		screenshotExists, _ := fsutil.FileExists(filepath)
		if screenshotExists {
			return os.ReadFile(filepath)
		}
	}

	return s.Repository.GetCover(ctx, scene.ID)
}
