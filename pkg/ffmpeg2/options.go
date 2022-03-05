package ffmpeg2

import "fmt"

type arger interface {
	Args() []string
}

type Args []string

func (a Args) LogLevel(l LogLevel) Args {
	if l == "" {
		return a
	}

	return append(a, l.Args()...)
}

func (a Args) XError() Args {
	return append(a, "-xerror")
}

func (a Args) Overwrite() Args {
	return append(a, "-y")
}

func (a Args) Seek(seconds float64) Args {
	return append(a, "-ss", fmt.Sprint(seconds))
}

func (a Args) Duration(seconds float64) Args {
	return append(a, "-t", fmt.Sprint(seconds))
}

func (a Args) Input(i string) Args {
	return append(a, "-i", i)
}

func (a Args) Output(o string) Args {
	return append(a, o)
}

func (a Args) VideoFrames(f int) Args {
	return append(a, "-frames:v", fmt.Sprint(f))
}

func (a Args) FixedQualityScaleVideo(q int) Args {
	return append(a, "-q:v", fmt.Sprint(q))
}

func (a Args) VideoFilter(vf VideoFilter) Args {
	return append(a, vf.Args()...)
}

func (a Args) VSync(m VSyncMethod) Args {
	return append(a, m.Args()...)
}

func (a Args) AudioBitrate(b string) Args {
	return append(a, "-b:a", b)
}

func (a Args) MaxMuxingQueueSize(s int) Args {
	// https://trac.ffmpeg.org/ticket/6375
	return append(a, "-max_muxing_queue_size", fmt.Sprint(s))
}

func (a Args) SkipAudio() Args {
	return append(a, "-an")
}

func (a Args) VideoCodec(c VideoCodec) Args {
	return append(a, c.Args()...)
}

func (a Args) AudioCodec(c AudioCodec) Args {
	return append(a, c.Args()...)
}

func (a Args) Format(f Format) Args {
	return append(a, f.Args()...)
}

func (a Args) AppendArgs(o arger) Args {
	return append(a, o.Args()...)
}

func (a Args) Args() []string {
	return []string(a)
}

type LogLevel string

func (l LogLevel) Args() []string {
	if l == "" {
		return nil
	}

	return []string{"-v", string(l)}
}

// LogLevels for ffmpeg. See -v entry under https://ffmpeg.org/ffmpeg.html#Generic-options
var (
	LogLevelQuiet   LogLevel = "quiet"
	LogLevelPanic   LogLevel = "panic"
	LogLevelFatal   LogLevel = "fatal"
	LogLevelError   LogLevel = "error"
	LogLevelWarning LogLevel = "warning"
	LogLevelInfo    LogLevel = "info"
	LogLevelVerbose LogLevel = "verbose"
	LogLevelDebug   LogLevel = "debug"
	LogLevelTrace   LogLevel = "trace"
)

type VSyncMethod string

func (m VSyncMethod) Args() []string {
	if m == "" {
		return nil
	}

	return []string{"-vsync", string(m)}
}

// Video sync methods for ffmpeg. See -vsync entry under https://ffmpeg.org/ffmpeg.html#Advanced-options
var (
	VSyncMethodPassthrough VSyncMethod = "0"
	VSyncMethodCFR         VSyncMethod = "1"
	VSyncMethodVFR         VSyncMethod = "2"
	VSyncMethodDrop        VSyncMethod = "drop"
	VSyncMethodAuto        VSyncMethod = "-1"
)
