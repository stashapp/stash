package transcoder

import (
	"errors"
	"fmt"

	"github.com/stashapp/stash/pkg/ffmpeg"
)

var ErrUnsupportedFormat = errors.New("unsupported image format")

type ImageThumbnailOptions struct {
	InputFormat   ffmpeg.ImageFormat
	OutputFormat  ffmpeg.ImageFormat
	OutputPath    string
	MaxDimensions int
	Quality       int
	// FlattenAlpha composites the image over a white background using
	// filter_complex, so that transparent areas become white in the JPEG
	// output rather than black.
	FlattenAlpha bool
}

func ImageThumbnail(input string, options ImageThumbnailOptions) ffmpeg.Args {
	var args ffmpeg.Args
	args = append(args, "-hide_banner")
	args = args.LogLevel(ffmpeg.LogLevelError)

	args = args.Overwrite().
		ImageFormat(options.InputFormat).
		Input(input)

	if options.FlattenAlpha {
		// Composite the image over a white background and scale in one pass.
		// scale2ref resizes the generated white source to match the input
		// dimensions before the overlay, so no fixed size is needed.
		fc := fmt.Sprintf(
			"color=white[bg];[0:v]format=rgba[fg];[bg][fg]scale2ref[bg2][fg2];[bg2][fg2]overlay=format=auto,scale=%[1]d:%[1]d:force_original_aspect_ratio=decrease[out]",
			options.MaxDimensions,
		)
		args = args.FilterComplex(fc).Map("[out]")
	} else {
		var videoFilter ffmpeg.VideoFilter
		videoFilter = videoFilter.ScaleMaxSize(options.MaxDimensions)
		args = args.VideoFilter(videoFilter)
	}

	args = args.VideoCodec(ffmpeg.VideoCodecMJpeg)

	args = append(args, "-frames:v", "1")

	if options.Quality > 0 {
		args = args.FixedQualityScaleVideo(options.Quality)
	}

	args = args.ImageFormat(ffmpeg.ImageFormatImage2Pipe).
		Output(options.OutputPath).
		ImageFormat(options.OutputFormat)

	return args
}
