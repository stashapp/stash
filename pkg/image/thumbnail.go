package image

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
	"runtime"
	"sync"

	"github.com/stashapp/stash/pkg/ffmpeg"
)

var vipsPath string
var once sync.Once

var (
	ErrUnsupportedFormat = errors.New("unsupported image format")
	ErrTooSmall          = errors.New("image too small")
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

// GetThumbnailReader returns the thumbnail image of the provided source resized to
// the provided max size. It resizes based on the largest X/Y direction.
// It returns nil and an error if an error occurs reading, decoding or encoding
// the image. Returns nil and ErrImageTooSmall if the image is smaller than the max size.
func (e *ThumbnailEncoder) GetThumbnail(reader io.Reader, maxSize int) ([]byte, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(reader); err != nil {
		return nil, err
	}

	config, format, err := DecodeSourceImage(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, err
	}

	if config.Width < maxSize && config.Height < maxSize {
		return nil, ErrTooSmall
	}

	if format != nil && *format == "gif" {
		return buf.Bytes(), nil
	}

	// vips has issues loading files from stdin on Windows
	if e.vips != nil && runtime.GOOS != "windows" {
		return e.vips.ImageThumbnail(buf, maxSize)
	} else {
		return e.ffmpeg.ImageThumbnail(buf, format, maxSize, "")
	}
}
