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
	"math"

	"github.com/stashapp/stash/pkg/logger"
)

// Level represents a logging level for plugin outputs.
type Level struct {
	*logger.PluginLogLevel
}

// Valid Level values.
var (
	TraceLevel = Level{
		&logger.TraceLevel,
	}
	DebugLevel = Level{
		&logger.DebugLevel,
	}
	InfoLevel = Level{
		&logger.InfoLevel,
	}
	WarningLevel = Level{
		&logger.WarningLevel,
	}
	ErrorLevel = Level{
		&logger.ErrorLevel,
	}
	ProgressLevel = Level{
		&logger.ProgressLevel,
	}
	NoneLevel = Level{
		&logger.NoneLevel,
	}
)

// Trace outputs a trace logging message to os.Stderr. Message is encoded with a
// prefix that signifies to the server that it is a trace message.
func Trace(args ...interface{}) {
	TraceLevel.Log(args...)
}

// Tracef is the equivalent of Printf outputting as a trace logging message.
func Tracef(format string, args ...interface{}) {
	TraceLevel.Logf(format, args...)
}

// Debug outputs a debug logging message to os.Stderr. Message is encoded with a
// prefix that signifies to the server that it is a debug message.
func Debug(args ...interface{}) {
	DebugLevel.Log(args...)
}

// Debugf is the equivalent of Printf outputting as a debug logging message.
func Debugf(format string, args ...interface{}) {
	DebugLevel.Logf(format, args...)
}

// Info outputs an info logging message to os.Stderr. Message is encoded with a
// prefix that signifies to the server that it is an info message.
func Info(args ...interface{}) {
	InfoLevel.Log(args...)
}

// Infof is the equivalent of Printf outputting as an info logging message.
func Infof(format string, args ...interface{}) {
	InfoLevel.Logf(format, args...)
}

// Warn outputs a warning logging message to os.Stderr. Message is encoded with a
// prefix that signifies to the server that it is a warning message.
func Warn(args ...interface{}) {
	WarningLevel.Log(args...)
}

// Warnf is the equivalent of Printf outputting as a warning logging message.
func Warnf(format string, args ...interface{}) {
	WarningLevel.Logf(format, args...)
}

// Error outputs an error logging message to os.Stderr. Message is encoded with a
// prefix that signifies to the server that it is an error message.
func Error(args ...interface{}) {
	ErrorLevel.Log(args...)
}

// Errorf is the equivalent of Printf outputting as an error logging message.
func Errorf(format string, args ...interface{}) {
	ErrorLevel.Logf(format, args...)
}

// Progress logs the current progress value. The progress value should be
// between 0 and 1.0 inclusively, with 1 representing that the task is
// complete. Values outside of this range will be clamp to be within it.
func Progress(progress float64) {
	progress = math.Min(math.Max(0, progress), 1)
	ProgressLevel.Log(progress)
}

// LevelFromName returns the Level that matches the provided name or nil if
// the name does not match a valid value.
func LevelFromName(name string) *Level {
	l := logger.PluginLogLevelFromName(name)
	if l != nil {
		return &Level{
			l,
		}
	}

	return nil
}
