package ffmpeg2

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
