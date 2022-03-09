package ffmpeg2

import (
	"errors"
)

var ErrUnsupportedFormat = errors.New("unsupported image format")

type ImageThumbnailOptions struct {
	InputFormat   ImageFormat
	OutputPath    string
	MaxDimensions int
	Quality       int
}

func ImageThumbnail(input string, options ImageThumbnailOptions) Args {
	var videoFilter VideoFilter
	videoFilter = videoFilter.ScaleMaxSize(options.MaxDimensions)

	var videoArgs Args
	videoArgs = videoArgs.VideoFilter(videoFilter)

	var args Args

	args = args.ImageFormat(options.InputFormat).Input(input).
		VideoFilter(videoFilter).
		VideoCodec(VideoCodecMJpeg)

	if options.Quality > 0 {
		args = args.FixedQualityScaleVideo(options.Quality)
	}

	args = args.ImageFormat(ImageFormatImage2Pipe).
		Output(options.OutputPath)

	return args
}
