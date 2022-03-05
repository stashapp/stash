package video

import (
	"github.com/stashapp/stash/pkg/ffmpeg2"
)

type SpliceOptions struct {
	OutputPath string
	Format     ffmpeg2.Format

	VideoCodec ffmpeg2.VideoCodec
	VideoArgs  ffmpeg2.Args

	AudioCodec ffmpeg2.AudioCodec
	AudioArgs  ffmpeg2.Args

	// Verbosity is the logging verbosity. Defaults to LogLevelError if not set.
	Verbosity ffmpeg2.LogLevel
}

func (o *SpliceOptions) setDefaults() {
	if o.Verbosity == "" {
		o.Verbosity = ffmpeg2.LogLevelError
	}
}

func Splice(concatFile string, options SpliceOptions) ffmpeg2.Args {
	options.setDefaults()

	var args ffmpeg2.Args
	args = args.LogLevel(options.Verbosity)
	args = args.Format(ffmpeg2.FormatConcat)
	args = args.Input(concatFile)
	args = args.Overwrite()

	// if video codec is not provided, then use copy
	if options.VideoCodec == "" {
		options.VideoCodec = ffmpeg2.VideoCodecCopy
	}

	args = args.VideoCodec(options.VideoCodec)
	args = args.AppendArgs(options.VideoArgs)

	// if audio codec is not provided, then skip it
	if options.AudioCodec == "" {
		args = args.SkipAudio()
	} else {
		args = args.AudioCodec(options.AudioCodec)
	}
	args = args.AppendArgs(options.AudioArgs)

	args = args.Format(options.Format)
	args = args.Output(options.OutputPath)

	return args
}
