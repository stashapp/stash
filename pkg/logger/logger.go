// Package logger provides methods and interfaces used by other stash packages for logging purposes.
package logger

import (
	"os"
)

// LoggerImpl is the interface that groups logging methods.
//
// Progressf logs using a specific progress format.
// Trace, Debug, Info, Warn and Error log to the applicable log level. Arguments are handled in the manner of fmt.Print.
// Tracef, Debugf, Infof, Warnf, Errorf log to the applicable log level. Arguments are handled in the manner of fmt.Printf.
// Fatal and Fatalf log to the applicable log level, then call os.Exit(1).
type LoggerImpl interface {
	Progressf(format string, args ...interface{})

	Trace(args ...interface{})
	Tracef(format string, args ...interface{})
	TraceFunc(fn func() (string, []interface{}))

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	DebugFunc(fn func() (string, []interface{}))

	Info(args ...interface{})
	Infof(format string, args ...interface{})
	InfoFunc(fn func() (string, []interface{}))

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	WarnFunc(fn func() (string, []interface{}))

	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	ErrorFunc(fn func() (string, []interface{}))

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

// Logger is the LoggerImpl used when calling the global Logger functions.
// It is suggested to use the LoggerImpl interface directly, rather than calling global log functions.
var Logger LoggerImpl

// Progressf calls Progressf with the Logger registered using RegisterLogger.
// If no logger has been registered, then this function is a no-op.
func Progressf(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Progressf(format, args...)
	}
}

// Trace calls Trace with the Logger registered using RegisterLogger.
// If no logger has been registered, then this function is a no-op.
func Trace(args ...interface{}) {
	if Logger != nil {
		Logger.Trace(args...)
	}
}

// Tracef calls Tracef with the Logger registered using RegisterLogger.
// If no logger has been registered, then this function is a no-op.
func Tracef(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Tracef(format, args...)
	}
}

// TraceFunc calls TraceFunc with the Logger registered using RegisterLogger.
// If no logger has been registered, then this function is a no-op.
func TraceFunc(fn func() (string, []interface{})) {
	if Logger != nil {
		Logger.TraceFunc(fn)
	}
}

// Debug calls Debug with the Logger registered using RegisterLogger.
// If no logger has been registered, then this function is a no-op.
func Debug(args ...interface{}) {
	if Logger != nil {
		Logger.Debug(args...)
	}
}

// Debugf calls Debugf with the Logger registered using RegisterLogger.
// If no logger has been registered, then this function is a no-op.
func Debugf(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Debugf(format, args...)
	}
}

// DebugFunc calls DebugFunc with the Logger registered using RegisterLogger.
// If no logger has been registered, then this function is a no-op.
func DebugFunc(fn func() (string, []interface{})) {
	if Logger != nil {
		Logger.DebugFunc(fn)
	}
}

// Info calls Info with the Logger registered using RegisterLogger.
// If no logger has been registered, then this function is a no-op.
func Info(args ...interface{}) {
	if Logger != nil {
		Logger.Info(args...)
	}
}

// Infof calls Infof with the Logger registered using RegisterLogger.
// If no logger has been registered, then this function is a no-op.
func Infof(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Infof(format, args...)
	}
}

// InfoFunc calls InfoFunc with the Logger registered using RegisterLogger.
// If no logger has been registered, then this function is a no-op.
func InfoFunc(fn func() (string, []interface{})) {
	if Logger != nil {
		Logger.InfoFunc(fn)
	}
}

// Warn calls Warn with the Logger registered using RegisterLogger.
// If no logger has been registered, then this function is a no-op.
func Warn(args ...interface{}) {
	if Logger != nil {
		Logger.Warn(args...)
	}
}

// Warnf calls Warnf with the Logger registered using RegisterLogger.
// If no logger has been registered, then this function is a no-op.
func Warnf(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Warnf(format, args...)
	}
}

// WarnFunc calls WarnFunc with the Logger registered using RegisterLogger.
// If no logger has been registered, then this function is a no-op.
func WarnFunc(fn func() (string, []interface{})) {
	if Logger != nil {
		Logger.WarnFunc(fn)
	}
}

// Error calls Error with the Logger registered using RegisterLogger.
// If no logger has been registered, then this function is a no-op.
func Error(args ...interface{}) {
	if Logger != nil {
		Logger.Error(args...)
	}
}

// Errorf calls Errorf with the Logger registered using RegisterLogger.
// If no logger has been registered, then this function is a no-op.
func Errorf(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Errorf(format, args...)
	}
}

// ErrorFunc calls ErrorFunc with the Logger registered using RegisterLogger.
// If no logger has been registered, then this function is a no-op.
func ErrorFunc(fn func() (string, []interface{})) {
	if Logger != nil {
		Logger.ErrorFunc(fn)
	}
}

// Fatal calls Fatal with the Logger registered using RegisterLogger.
// If no logger has been registered, then this function is a no-op.
func Fatal(args ...interface{}) {
	if Logger != nil {
		Logger.Fatal(args...)
	} else {
		os.Exit(1)
	}
}

// Fatalf calls Fatalf with the Logger registered using RegisterLogger.
// If no logger has been registered, then this function is a no-op.
func Fatalf(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Fatalf(format, args...)
	} else {
		os.Exit(1)
	}
}
