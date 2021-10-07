package ffmpeg

import (
	"bytes"
	"errors"
	"fmt"
)

func (e *Encoder) ImageThumbnail(image *bytes.Buffer, format *string, maxDimensions int, path string) ([]byte, error) {
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

	data, err := e.run(path, args, image)

	return []byte(data), err
}
