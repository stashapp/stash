package image

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

var vipsPath string
var once sync.Once

type ThumbnailEncoder struct {
	FFMPEGPath string
	VipsPath   string
}

func GetVipsPath() string {
	once.Do(func() {
		vipsPath, _ = exec.LookPath("vips")
	})
	return vipsPath
}

func NewThumbnailEncoder(ffmpegPath string) ThumbnailEncoder {
	return ThumbnailEncoder{
		FFMPEGPath: ffmpegPath,
		VipsPath:   GetVipsPath(),
	}
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
	if e.VipsPath != "" && runtime.GOOS != "windows" {
		return e.getVipsThumbnail(buf, maxSize)
	} else {
		return e.getFFMPEGThumbnail(buf, format, maxSize, img.Path)
	}
}

func (e *ThumbnailEncoder) getVipsThumbnail(image *bytes.Buffer, maxSize int) ([]byte, error) {
	args := []string{
		"thumbnail_source",
		"[descriptor=0]",
		".jpg[Q=70,strip]",
		fmt.Sprint(maxSize),
		"--size", "down",
	}
	data, err := e.run(e.VipsPath, args, image)

	return []byte(data), err
}

func (e *ThumbnailEncoder) getFFMPEGThumbnail(image *bytes.Buffer, format *string, maxDimensions int, path string) ([]byte, error) {
	// ffmpeg spends a long sniffing image format when data is piped through stdio, so we pass the format explicitly instead
	ffmpegformat := ""
	if format != nil && *format == "jpeg" {
		ffmpegformat = "mjpeg"
	} else if format != nil && *format == "png" {
		ffmpegformat = "png_pipe"
	} else if format != nil && *format == "webp" {
		ffmpegformat = "webp_pipe"
	} else {
		return nil, errors.New("unsupported image format")
	}

	args := []string{
		"-f", ffmpegformat,
		"-i", "-",
		"-vf", fmt.Sprintf("scale=%v:%v:force_original_aspect_ratio=decrease", maxDimensions, maxDimensions),
		"-c:v", "mjpeg",
		"-q:v", "5",
		"-f", "image2pipe",
		"-",
	}
	data, err := e.run(e.FFMPEGPath, args, image)

	return []byte(data), err
}

func (e *ThumbnailEncoder) run(path string, args []string, stdin *bytes.Buffer) (string, error) {
	cmd := exec.Command(path, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = stdin

	if err := cmd.Start(); err != nil {
		return "", err
	}

	err := cmd.Wait()

	if err != nil {
		// error message should be in the stderr stream
		logger.Errorf("image encoder error when running command <%s>: %s", strings.Join(cmd.Args, " "), stderr.String())
		return stdout.String(), err
	}

	return stdout.String(), nil
}
