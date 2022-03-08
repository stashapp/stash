package ffmpeg2

type SpliceOptions struct {
	OutputPath string
	Format     Format

	VideoCodec VideoCodec
	VideoArgs  Args

	AudioCodec AudioCodec
	AudioArgs  Args

	// Verbosity is the logging verbosity. Defaults to LogLevelError if not set.
	Verbosity LogLevel
}

func (o *SpliceOptions) setDefaults() {
	if o.Verbosity == "" {
		o.Verbosity = LogLevelError
	}
}

func Splice(concatFile string, options SpliceOptions) Args {
	options.setDefaults()

	var args Args
	args = args.LogLevel(options.Verbosity)
	args = args.Format(FormatConcat)
	args = args.Input(concatFile)
	args = args.Overwrite()

	// if video codec is not provided, then use copy
	if options.VideoCodec == "" {
		options.VideoCodec = VideoCodecCopy
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
