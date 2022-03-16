package image

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"sync"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/models"
)

var vipsPath string
var once sync.Once

var (
	// ErrNotSupportedForThumbnail is returned if the image format is not supported for thumbnail generation
	ErrNotSupportedForThumbnail = errors.New("unsupported image format for thumbnail")
)

type ThumbnailEncoder struct {
	ffmpeg ffmpeg.Encoder
	vips   *vipsEncoder
}

func GetVipsPath() string {
	once.Do(func() {
		vipsPath, _ = exec.LookPath("vips")
	})
	return vipsPath
}

func NewThumbnailEncoder(ffmpegEncoder ffmpeg.Encoder) ThumbnailEncoder {
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

	if format != nil && !supportThumbnail(*format) {
		return nil, fmt.Errorf("%w: %s", ErrNotSupportedForThumbnail, *format)
	}

	// vips has issues loading files from stdin on Windows
	if e.vips != nil && runtime.GOOS != "windows" {
		return e.vips.ImageThumbnail(buf, maxSize)
	} else {
		return e.ffmpeg.ImageThumbnail(buf, format, maxSize, img.Path)
	}
}

// supportThumbnail returns false if thumbnails should not be generated for the given format.
// Returns false for animated formats: gif and webp
func supportThumbnail(format string) bool {
	return format != "gif" && format != "webp"
}
