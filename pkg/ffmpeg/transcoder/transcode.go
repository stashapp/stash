package transcoder

import "github.com/stashapp/stash/pkg/ffmpeg"

type TranscodeOptions struct {
	OutputPath string
	Format     ffmpeg.Format

	VideoCodec ffmpeg.VideoCodec
	VideoArgs  ffmpeg.Args

	AudioCodec ffmpeg.AudioCodec
	AudioArgs  ffmpeg.Args

	// if XError is true, then ffmpeg will fail on warnings
	XError bool

	StartTime float64
	SlowSeek  bool
	Duration  float64

	// Verbosity is the logging verbosity. Defaults to LogLevelError if not set.
	Verbosity ffmpeg.LogLevel

	// arguments added before the input argument
	ExtraInputArgs []string
	// arguments added before the output argument
	ExtraOutputArgs []string
}

func (o *TranscodeOptions) setDefaults() {
	if o.Verbosity == "" {
		o.Verbosity = ffmpeg.LogLevelError
	}
}

func Transcode(input string, options TranscodeOptions) ffmpeg.Args {
	options.setDefaults()

	// TODO - this should probably be generalised and applied to all operations. Need to verify impact on phash algorithm.
	const fallbackMinSlowSeek = 20.0

	var fastSeek float64
	var slowSeek float64

	if !options.SlowSeek {
		fastSeek = options.StartTime
		slowSeek = 0
	} else {
		// In slowseek mode, try a combination of fast/slow seek instead of just fastseek
		// Commonly with avi/wmv ffmpeg doesn't seem to always predict the right start point to begin decoding when
		// using fast seek. If you force ffmpeg to decode more, it avoids the "blocky green artifact" issue.
		if options.StartTime > fallbackMinSlowSeek {
			// Handle seeks longer than fallbackMinSlowSeek with fast/slow seeks
			// Allow for at least fallbackMinSlowSeek seconds of slow seek
			fastSeek = options.StartTime - fallbackMinSlowSeek
			slowSeek = fallbackMinSlowSeek
		} else {
			// Handle seeks shorter than fallbackMinSlowSeek with only slow seeks.
			slowSeek = options.StartTime
			fastSeek = 0
		}
	}

	var args ffmpeg.Args
	args = args.LogLevel(options.Verbosity).Overwrite()
	args = append(args, options.ExtraInputArgs...)

	if options.XError {
		args = args.XError()
	}

	if fastSeek > 0 {
		args = args.Seek(fastSeek)
	}

	args = args.Input(input)

	if slowSeek > 0 {
		args = args.Seek(slowSeek)
	}

	if options.Duration > 0 {
		args = args.Duration(options.Duration)
	}

	// https://trac.ffmpeg.org/ticket/6375
	args = args.MaxMuxingQueueSize(1024)

	args = args.VideoCodec(options.VideoCodec)
	args = args.AppendArgs(options.VideoArgs)

	// if audio codec is not provided, then skip it
	if options.AudioCodec == "" {
		args = args.SkipAudio()
	} else {
		args = args.AudioCodec(options.AudioCodec)
	}
	args = args.AppendArgs(options.AudioArgs)

	args = append(args, options.ExtraOutputArgs...)

	args = args.Format(options.Format)
	args = args.Output(options.OutputPath)

	return args
}
