package transcoder

import "github.com/stashapp/stash/pkg/ffmpeg"

type ScreenshotOptions struct {
	OutputPath string
	OutputType ScreenshotOutputType

	// Quality is the quality scale. See https://ffmpeg.org/ffmpeg.html#Main-options
	Quality int

	Width int

	// Verbosity is the logging verbosity. Defaults to LogLevelError if not set.
	Verbosity ffmpeg.LogLevel

	UseSelectFilter bool
}

func (o *ScreenshotOptions) setDefaults() {
	if o.Verbosity == "" {
		o.Verbosity = ffmpeg.LogLevelError
	}
}

type ScreenshotOutputType struct {
	codec  *ffmpeg.VideoCodec
	format ffmpeg.Format
}

func (t ScreenshotOutputType) Args() []string {
	var ret []string
	if t.codec != nil {
		ret = append(ret, t.codec.Args()...)
	}
	if t.format != "" {
		ret = append(ret, t.format.Args()...)
	}

	return ret
}

var (
	ScreenshotOutputTypeImage2 = ScreenshotOutputType{
		format: "image2",
	}
	ScreenshotOutputTypeBMP = ScreenshotOutputType{
		codec:  &ffmpeg.VideoCodecBMP,
		format: "rawvideo",
	}
)

func ScreenshotTime(input string, t float64, options ScreenshotOptions) ffmpeg.Args {
	options.setDefaults()

	var args ffmpeg.Args
	args = args.LogLevel(options.Verbosity)
	args = args.Overwrite()
	args = args.Seek(t)

	args = args.Input(input)
	args = args.VideoFrames(1)

	if options.Quality > 0 {
		args = args.FixedQualityScaleVideo(options.Quality)
	}

	var vf ffmpeg.VideoFilter

	if options.Width > 0 {
		vf = vf.ScaleWidth(options.Width)
		args = args.VideoFilter(vf)
	}

	args = args.AppendArgs(options.OutputType)
	args = args.Output(options.OutputPath)

	return args
}

// ScreenshotFrame uses the select filter to get a single frame from the video.
// It is very slow and should only be used for files with very small duration in secs / frame count.
func ScreenshotFrame(input string, frame int, options ScreenshotOptions) ffmpeg.Args {
	options.setDefaults()

	var args ffmpeg.Args
	args = args.LogLevel(options.Verbosity)
	args = args.Overwrite()

	args = args.Input(input)
	args = args.VideoFrames(1)

	args = args.VSync(ffmpeg.VSyncMethodPassthrough)

	var vf ffmpeg.VideoFilter
	// keep only frame number options.Frame)
	vf = vf.Select(frame)

	if options.Width > 0 {
		vf = vf.ScaleWidth(options.Width)
	}

	args = args.VideoFilter(vf)

	args = args.AppendArgs(options.OutputType)
	args = args.Output(options.OutputPath)

	return args
}
