package ffmpeg

import (
	"fmt"
	"runtime"
)

// Arger is an interface that can be used to append arguments to an Args slice.
type Arger interface {
	Args() []string
}

// Args represents a slice of arguments to be passed to ffmpeg.
type Args []string

// LogLevel sets the LogLevel to l and returns the result.
func (a Args) LogLevel(l LogLevel) Args {
	if l == "" {
		return a
	}

	return append(a, l.Args()...)
}

// XError adds the -xerror flag and returns the result.
func (a Args) XError() Args {
	return append(a, "-xerror")
}

// Overwrite adds the overwrite flag (-y) and returns the result.
func (a Args) Overwrite() Args {
	return append(a, "-y")
}

// Seek adds a seek (-ss) to the given seconds and returns the result.
func (a Args) Seek(seconds float64) Args {
	return append(a, "-ss", fmt.Sprint(seconds))
}

// Duration sets the duration (-t) to the given seconds and returns the result.
func (a Args) Duration(seconds float64) Args {
	return append(a, "-t", fmt.Sprint(seconds))
}

// Input adds the input (-i) and returns the result.
func (a Args) Input(i string) Args {
	return append(a, "-i", i)
}

// Output adds the output o and returns the result.
func (a Args) Output(o string) Args {
	return append(a, o)
}

// NullOutput adds a null output and returns the result.
// On Windows, this outputs to NUL, on everything else, /dev/null.
func (a Args) NullOutput() Args {
	var output string
	if runtime.GOOS == "windows" {
		output = "nul" // https://stackoverflow.com/questions/313111/is-there-a-dev-null-on-windows
	} else {
		output = "/dev/null"
	}

	return a.Output(output)
}

// VideoFrames adds the -frames:v with f and returns the result.
func (a Args) VideoFrames(f int) Args {
	return append(a, "-frames:v", fmt.Sprint(f))
}

// FixedQualityScaleVideo adds the -q:v argument with q and returns the result.
func (a Args) FixedQualityScaleVideo(q int) Args {
	return append(a, "-q:v", fmt.Sprint(q))
}

// VideoFilter adds the vf video filter and returns the result.
func (a Args) VideoFilter(vf VideoFilter) Args {
	return append(a, vf.Args()...)
}

// VSync adds the VsyncMethod and returns the result.
func (a Args) VSync(m VSyncMethod) Args {
	return append(a, m.Args()...)
}

// AudioBitrate adds the -b:a argument with b and returns the result.
func (a Args) AudioBitrate(b string) Args {
	return append(a, "-b:a", b)
}

// MaxMuxingQueueSize adds the -max_muxing_queue_size argument with s and returns the result.
func (a Args) MaxMuxingQueueSize(s int) Args {
	// https://trac.ffmpeg.org/ticket/6375
	return append(a, "-max_muxing_queue_size", fmt.Sprint(s))
}

// SkipAudio adds the skip audio flag (-an) and returns the result.
func (a Args) SkipAudio() Args {
	return append(a, "-an")
}

// VideoCodec adds the given video codec and returns the result.
func (a Args) VideoCodec(c VideoCodec) Args {
	return append(a, c.Args()...)
}

// AudioCodec adds the given audio codec and returns the result.
func (a Args) AudioCodec(c AudioCodec) Args {
	return append(a, c.Args()...)
}

// Format adds the format flag with f and returns the result.
func (a Args) Format(f Format) Args {
	return append(a, f.Args()...)
}

// ImageFormat adds the image format (using -f) and returns the result.
func (a Args) ImageFormat(f ImageFormat) Args {
	return append(a, f.Args()...)
}

// AppendArgs appends the given Arger to the Args and returns the result.
func (a Args) AppendArgs(o Arger) Args {
	return append(a, o.Args()...)
}

// Args returns a string slice of the arguments.
func (a Args) Args() []string {
	return []string(a)
}

// LogLevel represents the log level of ffmpeg.
type LogLevel string

// Args returns the arguments to set the log level in ffmpeg.
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

// VSyncMethod represents the vsync method of ffmpeg.
type VSyncMethod string

// Args returns the arguments to set the vsync method in ffmpeg.
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
