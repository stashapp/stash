package transcoder

import (
	"runtime"
	"strings"

	"github.com/stashapp/stash/pkg/ffmpeg"
)

type SpliceOptions struct {
	OutputPath string
	Format     ffmpeg.Format

	VideoCodec *ffmpeg.VideoCodec
	VideoArgs  ffmpeg.Args

	AudioCodec ffmpeg.AudioCodec
	AudioArgs  ffmpeg.Args

	// Verbosity is the logging verbosity. Defaults to LogLevelError if not set.
	Verbosity ffmpeg.LogLevel
}

func (o *SpliceOptions) setDefaults() {
	if o.Verbosity == "" {
		o.Verbosity = ffmpeg.LogLevelError
	}
}

// fixWindowsPath replaces \ with / in the given path because the \ isn't recognized as valid on windows ffmpeg
func fixWindowsPath(str string) string {
	if runtime.GOOS == "windows" {
		return strings.ReplaceAll(str, `\`, "/")
	}
	return str
}

func Splice(concatFile string, options SpliceOptions) ffmpeg.Args {
	options.setDefaults()

	var args ffmpeg.Args
	args = args.LogLevel(options.Verbosity)
	args = args.Format(ffmpeg.FormatConcat)
	args = args.Input(fixWindowsPath(concatFile))
	args = args.Overwrite()

	// if video codec is not provided, then use copy
	if options.VideoCodec == nil {
		options.VideoCodec = &ffmpeg.VideoCodecCopy
	}

	args = args.VideoCodec(*options.VideoCodec)
	args = args.AppendArgs(options.VideoArgs)

	// if audio codec is not provided, then use copy
	if options.AudioCodec == "" {
		options.AudioCodec = ffmpeg.AudioCodecCopy
	}

	args = args.AudioCodec(options.AudioCodec)
	args = args.AppendArgs(options.AudioArgs)

	args = args.Format(options.Format)
	args = args.Output(options.OutputPath)

	return args
}
