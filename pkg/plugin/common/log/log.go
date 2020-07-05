// Package log provides a number of logging utility functions for encoding and
// decoding log messages between a stash server and a plugin instance.
//
// Log messages sent from a plugin instance are transmitted via stderr and are
// encoded with a prefix consisting of special character SOH, then the log
// level (one of t, d, i, w, e, or p - corresponding to trace, debug, info,
// warning, error and progress levels respectively), then special character
// STX.
//
// The Trace, Debug, Info, Warning, and Error methods, and their equivalent
// formatted methods are intended for use by plugin instances to transmit log
// messages. The Progress method is also intended for sending progress data.
//
// Conversely, LevelFromName and DetectLogLevel are intended for use by the
// stash server.
package log

import (
	"fmt"
	"math"
	"os"
	"strings"
)

// Level represents a logging level for plugin outputs.
type Level struct {
	char byte
	Name string
}

// Valid Level values.
var (
	TraceLevel = Level{
		char: 't',
		Name: "trace",
	}
	DebugLevel = Level{
		char: 'd',
		Name: "debug",
	}
	InfoLevel = Level{
		char: 'i',
		Name: "info",
	}
	WarningLevel = Level{
		char: 'w',
		Name: "warning",
	}
	ErrorLevel = Level{
		char: 'e',
		Name: "error",
	}
	ProgressLevel = Level{
		char: 'p',
	}
	NoneLevel = Level{
		Name: "none",
	}
)

var validLevels = []Level{
	TraceLevel,
	DebugLevel,
	InfoLevel,
	WarningLevel,
	ErrorLevel,
	ProgressLevel,
	NoneLevel,
}

const startLevelChar byte = 1
const endLevelChar byte = 2

func (l Level) prefix() string {
	return string([]byte{
		startLevelChar,
		byte(l.char),
		endLevelChar,
	})
}

func (l Level) log(args ...interface{}) {
	if l.char == 0 {
		return
	}

	argsToUse := []interface{}{
		l.prefix(),
	}
	argsToUse = append(argsToUse, args...)
	fmt.Fprintln(os.Stderr, argsToUse...)
}

func (l Level) logf(format string, args ...interface{}) {
	if l.char == 0 {
		return
	}

	formatToUse := string(l.prefix()) + format + "\n"
	fmt.Fprintf(os.Stderr, formatToUse, args...)
}

// Trace outputs a trace logging message to os.Stderr. Message is encoded with a
// prefix that signifies to the server that it is a trace message.
func Trace(args ...interface{}) {
	TraceLevel.log(args...)
}

// Tracef is the equivalent of Printf outputting as a trace logging message.
func Tracef(format string, args ...interface{}) {
	TraceLevel.logf(format, args...)
}

// Debug outputs a debug logging message to os.Stderr. Message is encoded with a
// prefix that signifies to the server that it is a debug message.
func Debug(args ...interface{}) {
	DebugLevel.log(args...)
}

// Debugf is the equivalent of Printf outputting as a debug logging message.
func Debugf(format string, args ...interface{}) {
	DebugLevel.logf(format, args...)
}

// Info outputs an info logging message to os.Stderr. Message is encoded with a
// prefix that signifies to the server that it is an info message.
func Info(args ...interface{}) {
	InfoLevel.log(args...)
}

// Infof is the equivalent of Printf outputting as an info logging message.
func Infof(format string, args ...interface{}) {
	InfoLevel.logf(format, args...)
}

// Warn outputs a warning logging message to os.Stderr. Message is encoded with a
// prefix that signifies to the server that it is a warning message.
func Warn(args ...interface{}) {
	WarningLevel.log(args...)
}

// Warnf is the equivalent of Printf outputting as a warning logging message.
func Warnf(format string, args ...interface{}) {
	WarningLevel.logf(format, args...)
}

// Error outputs an error logging message to os.Stderr. Message is encoded with a
// prefix that signifies to the server that it is an error message.
func Error(args ...interface{}) {
	ErrorLevel.log(args...)
}

// Errorf is the equivalent of Printf outputting as an error logging message.
func Errorf(format string, args ...interface{}) {
	ErrorLevel.logf(format, args...)
}

// Progress logs the current progress value. The progress value should be
// between 0 and 1.0 inclusively, with 1 representing that the task is
// complete. Values outside of this range will be clamp to be within it.
func Progress(progress float64) {
	progress = math.Min(math.Max(0, progress), 1)
	ProgressLevel.log(progress)
}

// LevelFromName returns the Level that matches the provided name or nil if
// the name does not match a valid value.
func LevelFromName(name string) *Level {
	for _, l := range validLevels {
		if l.Name == name {
			return &l
		}
	}

	return nil
}

// DetectLogLevel returns the Level and the logging string for a provided line
// of plugin output. It parses the string for logging control characters and
// determines the log level, if present. If not present, the plugin output
// is returned unchanged with a nil Level.
func DetectLogLevel(line string) (*Level, string) {
	if len(line) < 4 || line[0] != startLevelChar || line[2] != endLevelChar {
		return nil, line
	}

	char := line[1]
	var level *Level
	for _, l := range validLevels {
		if l.char == char {
			level = &l
			break
		}
	}

	if level == nil {
		return nil, line
	}

	line = strings.TrimSpace(line[3:])

	return level, line
}
