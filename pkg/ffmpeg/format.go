package ffmpeg

// Format represents the input/output format for ffmpeg.
type Format string

// Args converts the Format to a slice of arguments to be passed to ffmpeg.
func (f Format) Args() []string {
	if f == "" {
		return nil
	}

	return []string{"-f", string(f)}
}

var (
	FormatConcat   Format = "concat"
	FormatImage2   Format = "image2"
	FormatRawVideo Format = "rawvideo"
	FormatMpegTS   Format = "mpegts"
	FormatMP4      Format = "mp4"
	FormatWebm     Format = "webm"
	FormatMatroska Format = "matroska"
)

// ImageFormat represents the input format for an image for ffmpeg.
type ImageFormat string

// Args converts the ImageFormat to a slice of arguments to be passed to ffmpeg.
func (f ImageFormat) Args() []string {
	if f == "" {
		return nil
	}

	return []string{"-f", string(f)}
}

var (
	ImageFormatJpeg ImageFormat = "mjpeg"
	ImageFormatPng  ImageFormat = "png_pipe"
	ImageFormatWebp ImageFormat = "webp_pipe"

	ImageFormatImage2Pipe ImageFormat = "image2pipe"
)
