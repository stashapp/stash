package image

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"runtime"
	"sync"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/ffmpeg/transcoder"
	"github.com/stashapp/stash/pkg/models"
)

const ffmpegImageQuality = 5

var vipsPath string
var once sync.Once

var ErrUnsupportedImageFormat = errors.New("unsupported image format")

type ThumbnailEncoder struct {
	ffmpeg ffmpeg.FFMpeg
	vips   *vipsEncoder
}

func GetVipsPath() string {
	once.Do(func() {
		vipsPath, _ = exec.LookPath("vips")
	})
	return vipsPath
}

func NewThumbnailEncoder(ffmpegEncoder ffmpeg.FFMpeg) ThumbnailEncoder {
	ret := ThumbnailEncoder{
		ffmpeg: ffmpegEncoder,
	}

	vipsPath := GetVipsPath()
	if vipsPath != "" {
		vipsEncoder := vipsEncoder(vipsPath)
		ret.vips = &vipsEncoder
	}

	return ret
}

// GetThumbnail returns the thumbnail image of the provided image resized to
// the provided max size. It resizes based on the largest X/Y direction.
// It returns nil and an error if an error occurs reading, decoding or encoding
// the image.
func (e *ThumbnailEncoder) GetThumbnail(img *models.Image, maxSize int) ([]byte, error) {
	reader, err := openSourceImage(img.Path)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(reader); err != nil {
		return nil, err
	}

	_, format, err := DecodeSourceImage(img)
	if err != nil {
		return nil, err
	}

	if format != nil && *format == "gif" {
		return buf.Bytes(), nil
	}

	// vips has issues loading files from stdin on Windows
	if e.vips != nil && runtime.GOOS != "windows" {
		return e.vips.ImageThumbnail(buf, maxSize)
	} else {
		return e.ffmpegImageThumbnail(buf, format, maxSize)
	}
}

func (e *ThumbnailEncoder) ffmpegImageThumbnail(image *bytes.Buffer, format *string, maxSize int) ([]byte, error) {
	if format == nil {
		return nil, ErrUnsupportedImageFormat
	}

	var ffmpegFormat ffmpeg.ImageFormat

	switch *format {
	case "jpeg":
		ffmpegFormat = ffmpeg.ImageFormatJpeg
	case "png":
		ffmpegFormat = ffmpeg.ImageFormatPng
	case "webp":
		ffmpegFormat = ffmpeg.ImageFormatWebp
	default:
		return nil, ErrUnsupportedImageFormat
	}

	args := transcoder.ImageThumbnail("-", transcoder.ImageThumbnailOptions{
		InputFormat:   ffmpegFormat,
		OutputPath:    "-",
		MaxDimensions: maxSize,
		Quality:       ffmpegImageQuality,
	})

	return ffmpeg.GenerateOutput(context.TODO(), e.ffmpeg, args)
}
