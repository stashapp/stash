package imagephash

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"path/filepath"
	"strings"

	"github.com/corona10/goimagehash"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/ffmpeg/transcoder"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/models"
)

// Generate computes a perceptual hash for an image file.
func Generate(encoder *ffmpeg.FFMpeg, imageFile *models.ImageFile) (*uint64, error) {
	img, err := loadImage(encoder, imageFile)
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
// For AVIF files, ffmpeg is used to convert to BMP first since Go has no built-in AVIF decoder.
func loadImage(encoder *ffmpeg.FFMpeg, imageFile *models.ImageFile) (image.Image, error) {
	ext := strings.ToLower(filepath.Ext(imageFile.Path))
	if ext == ".avif" {
		// AVIF in zip files is not supported - ffmpeg cannot read files inside zips
		if imageFile.Base().ZipFileID != nil {
			return nil, fmt.Errorf("AVIF images in zip files are not supported for phash generation")
		}
		return loadImageFFmpeg(encoder, imageFile.Path)
	}

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

// loadImageFFmpeg uses ffmpeg to convert an image to BMP and then decodes it.
func loadImageFFmpeg(encoder *ffmpeg.FFMpeg, path string) (image.Image, error) {
	options := transcoder.ScreenshotOptions{
		OutputPath: "-",
		OutputType: transcoder.ScreenshotOutputTypeBMP,
	}

	args := transcoder.ScreenshotTime(path, 0, options)
	data, err := encoder.GenerateOutput(context.Background(), args, nil)
	if err != nil {
		return nil, fmt.Errorf("converting image with ffmpeg: %w", err)
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decoding ffmpeg output: %w", err)
	}

	return img, nil
}
