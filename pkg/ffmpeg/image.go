package ffmpeg

import (
	"bytes"
	"fmt"
)

func (e *Encoder) ImageThumbnail(image *bytes.Buffer, format string, maxDimensions int, path string) ([]byte, error) {
	// ffmpeg spends a long sniffing image format when data is piped through stdio, so we pass the format explicitly instead
	var ffmpegformat string

	switch format {
	case "jpeg":
		ffmpegformat = "mjpeg"
	case "png":
		ffmpegformat = "png_pipe"
	case "webp":
		ffmpegformat = "webp_pipe"
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
