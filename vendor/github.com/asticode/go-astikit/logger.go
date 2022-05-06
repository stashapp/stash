package astikit

import (
	"context"
)

// CompleteLogger represents a complete logger
type CompleteLogger interface {
	StdLogger
	SeverityLogger
	SeverityCtxLogger
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

type completeLogger struct {
	print, debug, error, fatal, info, warn       func(v ...interface{})
	printf, debugf, errorf, fatalf, infof, warnf func(format string, v ...interface{})
	debugC, errorC, fatalC, infoC, warnC         func(ctx context.Context, v ...interface{})
	debugCf, errorCf, fatalCf, infoCf, warnCf    func(ctx context.Context, format string, v ...interface{})
}

func newCompleteLogger() *completeLogger {
	return &completeLogger{
		debug:   func(v ...interface{}) {},
		debugf:  func(format string, v ...interface{}) {},
		debugC:  func(ctx context.Context, v ...interface{}) {},
		debugCf: func(ctx context.Context, format string, v ...interface{}) {},
		error:   func(v ...interface{}) {},
		errorf:  func(format string, v ...interface{}) {},
		errorC:  func(ctx context.Context, v ...interface{}) {},
		errorCf: func(ctx context.Context, format string, v ...interface{}) {},
		fatal:   func(v ...interface{}) {},
		fatalf:  func(format string, v ...interface{}) {},
		fatalC:  func(ctx context.Context, v ...interface{}) {},
		fatalCf: func(ctx context.Context, format string, v ...interface{}) {},
		info:    func(v ...interface{}) {},
		infof:   func(format string, v ...interface{}) {},
		infoC:   func(ctx context.Context, v ...interface{}) {},
		infoCf:  func(ctx context.Context, format string, v ...interface{}) {},
		print:   func(v ...interface{}) {},
		printf:  func(format string, v ...interface{}) {},
		warn:    func(v ...interface{}) {},
		warnf:   func(format string, v ...interface{}) {},
		warnC:   func(ctx context.Context, v ...interface{}) {},
		warnCf:  func(ctx context.Context, format string, v ...interface{}) {},
	}
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
	} else {
		l.debug = l.print
		l.debugf = l.printf
		l.error = l.print
		l.errorf = l.printf
		l.info = l.print
		l.infof = l.printf
		l.warn = l.print
		l.warnf = l.printf
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
	} else {
		l.debugC = func(ctx context.Context, v ...interface{}) { l.debug(v...) }
		l.debugCf = func(ctx context.Context, format string, v ...interface{}) { l.debugf(format, v...) }
		l.errorC = func(ctx context.Context, v ...interface{}) { l.error(v...) }
		l.errorCf = func(ctx context.Context, format string, v ...interface{}) { l.errorf(format, v...) }
		l.fatalC = func(ctx context.Context, v ...interface{}) { l.fatal(v...) }
		l.fatalCf = func(ctx context.Context, format string, v ...interface{}) { l.fatalf(format, v...) }
		l.infoC = func(ctx context.Context, v ...interface{}) { l.info(v...) }
		l.infoCf = func(ctx context.Context, format string, v ...interface{}) { l.infof(format, v...) }
		l.warnC = func(ctx context.Context, v ...interface{}) { l.warn(v...) }
		l.warnCf = func(ctx context.Context, format string, v ...interface{}) { l.warnf(format, v...) }
	}
	return l
}
