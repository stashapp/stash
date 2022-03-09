package ffmpeg

type Format string

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

type ImageFormat string

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
