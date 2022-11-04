package astikit

import (
	"context"
)

// LoggerLevel represents a logger level
type LoggerLevel int

// Logger levels
const (
	LoggerLevelDebug LoggerLevel = iota
	LoggerLevelInfo
	LoggerLevelWarn
	LoggerLevelError
	LoggerLevelFatal
)

// LoggerLevelFromString creates a logger level from string
func LoggerLevelFromString(s string) LoggerLevel {
	switch s {
	case "debug":
		return LoggerLevelDebug
	case "error":
		return LoggerLevelError
	case "fatal":
		return LoggerLevelFatal
	case "warn":
		return LoggerLevelWarn
	default:
		return LoggerLevelInfo
	}
}

func (l LoggerLevel) String() string {
	switch l {
	case LoggerLevelDebug:
		return "debug"
	case LoggerLevelError:
		return "error"
	case LoggerLevelFatal:
		return "fatal"
	case LoggerLevelWarn:
		return "warn"
	default:
		return "info"
	}
}

func (l *LoggerLevel) UnmarshalText(b []byte) error {
	*l = LoggerLevelFromString(string(b))
	return nil
}

func (l LoggerLevel) MarshalText() ([]byte, error) {
	b := []byte(l.String())
	return b, nil
}

// CompleteLogger represents a complete logger
type CompleteLogger interface {
	SeverityCtxLogger
	SeverityLogger
	SeverityWriteLogger
	SeverityWriteCtxLogger
	StdLogger
}

// StdLogger represents a standard logger
type StdLogger interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Print(v ...interface{})
	Printf(format string, v ...interface{})
}

// SeverityLogger represents a severity logger
type SeverityLogger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
}

type TestLogger interface {
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Log(v ...interface{})
	Logf(format string, v ...interface{})
}

// SeverityCtxLogger represents a severity with context logger
type SeverityCtxLogger interface {
	DebugC(ctx context.Context, v ...interface{})
	DebugCf(ctx context.Context, format string, v ...interface{})
	ErrorC(ctx context.Context, v ...interface{})
	ErrorCf(ctx context.Context, format string, v ...interface{})
	FatalC(ctx context.Context, v ...interface{})
	FatalCf(ctx context.Context, format string, v ...interface{})
	InfoC(ctx context.Context, v ...interface{})
	InfoCf(ctx context.Context, format string, v ...interface{})
	WarnC(ctx context.Context, v ...interface{})
	WarnCf(ctx context.Context, format string, v ...interface{})
}

type SeverityWriteLogger interface {
	Write(l LoggerLevel, v ...interface{})
	Writef(l LoggerLevel, format string, v ...interface{})
}

type SeverityWriteCtxLogger interface {
	WriteC(ctx context.Context, l LoggerLevel, v ...interface{})
	WriteCf(ctx context.Context, l LoggerLevel, format string, v ...interface{})
}

type completeLogger struct {
	print, debug, error, fatal, info, warn       func(v ...interface{})
	printf, debugf, errorf, fatalf, infof, warnf func(format string, v ...interface{})
	debugC, errorC, fatalC, infoC, warnC         func(ctx context.Context, v ...interface{})
	debugCf, errorCf, fatalCf, infoCf, warnCf    func(ctx context.Context, format string, v ...interface{})
	write                                        func(l LoggerLevel, v ...interface{})
	writeC                                       func(ctx context.Context, l LoggerLevel, v ...interface{})
	writeCf                                      func(ctx context.Context, l LoggerLevel, format string, v ...interface{})
	writef                                       func(l LoggerLevel, format string, v ...interface{})
}

func newCompleteLogger() *completeLogger {
	l := &completeLogger{}
	l.debug = func(v ...interface{}) { l.print(v...) }
	l.debugf = func(format string, v ...interface{}) { l.printf(format, v...) }
	l.debugC = func(ctx context.Context, v ...interface{}) { l.debug(v...) }
	l.debugCf = func(ctx context.Context, format string, v ...interface{}) { l.debugf(format, v...) }
	l.error = func(v ...interface{}) { l.print(v...) }
	l.errorf = func(format string, v ...interface{}) { l.printf(format, v...) }
	l.errorC = func(ctx context.Context, v ...interface{}) { l.error(v...) }
	l.errorCf = func(ctx context.Context, format string, v ...interface{}) { l.errorf(format, v...) }
	l.fatal = func(v ...interface{}) { l.print(v...) }
	l.fatalf = func(format string, v ...interface{}) { l.printf(format, v...) }
	l.fatalC = func(ctx context.Context, v ...interface{}) { l.fatal(v...) }
	l.fatalCf = func(ctx context.Context, format string, v ...interface{}) { l.fatalf(format, v...) }
	l.info = func(v ...interface{}) { l.print(v...) }
	l.infof = func(format string, v ...interface{}) { l.printf(format, v...) }
	l.infoC = func(ctx context.Context, v ...interface{}) { l.info(v...) }
	l.infoCf = func(ctx context.Context, format string, v ...interface{}) { l.infof(format, v...) }
	l.print = func(v ...interface{}) {}
	l.printf = func(format string, v ...interface{}) {}
	l.warn = func(v ...interface{}) { l.print(v...) }
	l.warnf = func(format string, v ...interface{}) { l.printf(format, v...) }
	l.warnC = func(ctx context.Context, v ...interface{}) { l.warn(v...) }
	l.warnCf = func(ctx context.Context, format string, v ...interface{}) { l.warnf(format, v...) }
	l.write = func(lv LoggerLevel, v ...interface{}) {
		switch lv {
		case LoggerLevelDebug:
			l.debug(v...)
		case LoggerLevelError:
			l.error(v...)
		case LoggerLevelFatal:
			l.fatal(v...)
		case LoggerLevelWarn:
			l.warn(v...)
		default:
			l.info(v...)
		}
	}
	l.writeC = func(ctx context.Context, lv LoggerLevel, v ...interface{}) {
		switch lv {
		case LoggerLevelDebug:
			l.debugC(ctx, v...)
		case LoggerLevelError:
			l.errorC(ctx, v...)
		case LoggerLevelFatal:
			l.fatalC(ctx, v...)
		case LoggerLevelWarn:
			l.warnC(ctx, v...)
		default:
			l.infoC(ctx, v...)
		}
	}
	l.writeCf = func(ctx context.Context, lv LoggerLevel, format string, v ...interface{}) {
		switch lv {
		case LoggerLevelDebug:
			l.debugCf(ctx, format, v...)
		case LoggerLevelError:
			l.errorCf(ctx, format, v...)
		case LoggerLevelFatal:
			l.fatalCf(ctx, format, v...)
		case LoggerLevelWarn:
			l.warnCf(ctx, format, v...)
		default:
			l.infoCf(ctx, format, v...)
		}
	}
	l.writef = func(lv LoggerLevel, format string, v ...interface{}) {
		switch lv {
		case LoggerLevelDebug:
			l.debugf(format, v...)
		case LoggerLevelError:
			l.errorf(format, v...)
		case LoggerLevelFatal:
			l.fatalf(format, v...)
		case LoggerLevelWarn:
			l.warnf(format, v...)
		default:
			l.infof(format, v...)
		}
	}
	return l
}

func (l *completeLogger) Debug(v ...interface{})                       { l.debug(v...) }
func (l *completeLogger) Debugf(format string, v ...interface{})       { l.debugf(format, v...) }
func (l *completeLogger) DebugC(ctx context.Context, v ...interface{}) { l.debugC(ctx, v...) }
func (l *completeLogger) DebugCf(ctx context.Context, format string, v ...interface{}) {
	l.debugCf(ctx, format, v...)
}
func (l *completeLogger) Error(v ...interface{})                       { l.error(v...) }
func (l *completeLogger) Errorf(format string, v ...interface{})       { l.errorf(format, v...) }
func (l *completeLogger) ErrorC(ctx context.Context, v ...interface{}) { l.errorC(ctx, v...) }
func (l *completeLogger) ErrorCf(ctx context.Context, format string, v ...interface{}) {
	l.errorCf(ctx, format, v...)
}
func (l *completeLogger) Fatal(v ...interface{})                       { l.fatal(v...) }
func (l *completeLogger) Fatalf(format string, v ...interface{})       { l.fatalf(format, v...) }
func (l *completeLogger) FatalC(ctx context.Context, v ...interface{}) { l.fatalC(ctx, v...) }
func (l *completeLogger) FatalCf(ctx context.Context, format string, v ...interface{}) {
	l.fatalCf(ctx, format, v...)
}
func (l *completeLogger) Info(v ...interface{})                       { l.info(v...) }
func (l *completeLogger) Infof(format string, v ...interface{})       { l.infof(format, v...) }
func (l *completeLogger) InfoC(ctx context.Context, v ...interface{}) { l.infoC(ctx, v...) }
func (l *completeLogger) InfoCf(ctx context.Context, format string, v ...interface{}) {
	l.infoCf(ctx, format, v...)
}
func (l *completeLogger) Print(v ...interface{})                      { l.print(v...) }
func (l *completeLogger) Printf(format string, v ...interface{})      { l.printf(format, v...) }
func (l *completeLogger) Warn(v ...interface{})                       { l.warn(v...) }
func (l *completeLogger) Warnf(format string, v ...interface{})       { l.warnf(format, v...) }
func (l *completeLogger) WarnC(ctx context.Context, v ...interface{}) { l.warnC(ctx, v...) }
func (l *completeLogger) WarnCf(ctx context.Context, format string, v ...interface{}) {
	l.warnCf(ctx, format, v...)
}
func (l *completeLogger) Write(lv LoggerLevel, v ...interface{}) { l.write(lv, v...) }
func (l *completeLogger) Writef(lv LoggerLevel, format string, v ...interface{}) {
	l.writef(lv, format, v...)
}
func (l *completeLogger) WriteC(ctx context.Context, lv LoggerLevel, v ...interface{}) {
	l.writeC(ctx, lv, v...)
}
func (l *completeLogger) WriteCf(ctx context.Context, lv LoggerLevel, format string, v ...interface{}) {
	l.writeCf(ctx, lv, format, v...)
}

// AdaptStdLogger transforms an StdLogger into a CompleteLogger if needed
func AdaptStdLogger(i StdLogger) CompleteLogger {
	if v, ok := i.(CompleteLogger); ok {
		return v
	}
	l := newCompleteLogger()
	if i == nil {
		return l
	}
	l.fatal = i.Fatal
	l.fatalf = i.Fatalf
	l.print = i.Print
	l.printf = i.Printf
	if v, ok := i.(SeverityLogger); ok {
		l.debug = v.Debug
		l.debugf = v.Debugf
		l.error = v.Error
		l.errorf = v.Errorf
		l.info = v.Info
		l.infof = v.Infof
		l.warn = v.Warn
		l.warnf = v.Warnf
	}
	if v, ok := i.(SeverityCtxLogger); ok {
		l.debugC = v.DebugC
		l.debugCf = v.DebugCf
		l.errorC = v.ErrorC
		l.errorCf = v.ErrorCf
		l.fatalC = v.FatalC
		l.fatalCf = v.FatalCf
		l.infoC = v.InfoC
		l.infoCf = v.InfoCf
		l.warnC = v.WarnC
		l.warnCf = v.WarnCf
	}
	if v, ok := i.(SeverityWriteLogger); ok {
		l.write = v.Write
		l.writef = v.Writef
	}
	if v, ok := i.(SeverityWriteCtxLogger); ok {
		l.writeC = v.WriteC
		l.writeCf = v.WriteCf
	}
	return l
}

// AdaptTestLogger transforms a TestLogger into a CompleteLogger if needed
func AdaptTestLogger(i TestLogger) CompleteLogger {
	if v, ok := i.(CompleteLogger); ok {
		return v
	}
	l := newCompleteLogger()
	if i == nil {
		return l
	}
	l.error = i.Error
	l.errorf = i.Errorf
	l.fatal = i.Fatal
	l.fatalf = i.Fatalf
	l.print = i.Log
	l.printf = i.Logf
	l.debug = l.print
	l.debugf = l.printf
	l.info = l.print
	l.infof = l.printf
	l.warn = l.print
	l.warnf = l.printf
	return l
}
