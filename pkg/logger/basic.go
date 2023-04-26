package logger

import (
	"fmt"
	"os"
)

// BasicLogger logs all messages to stdout
type BasicLogger struct{}

var _ LoggerImpl = &BasicLogger{}

func (log *BasicLogger) print(level string, args ...interface{}) {
	fmt.Print(level + ": ")
	fmt.Println(args...)
}

func (log *BasicLogger) printf(level string, format string, args ...interface{}) {
	fmt.Printf(level+": "+format+"\n", args...)
}

func (log *BasicLogger) Progressf(format string, args ...interface{}) {
	log.printf("Progress", format, args...)
}

func (log *BasicLogger) Trace(args ...interface{}) {
	log.print("Trace", args...)
}

func (log *BasicLogger) Tracef(format string, args ...interface{}) {
	log.printf("Trace", format, args...)
}

func (log *BasicLogger) TraceFunc(fn func() (string, []interface{})) {
	format, args := fn()
	log.printf("Trace", format, args...)
}

func (log *BasicLogger) Debug(args ...interface{}) {
	log.print("Debug", args...)
}

func (log *BasicLogger) Debugf(format string, args ...interface{}) {
	log.printf("Debug", format, args...)
}

func (log *BasicLogger) DebugFunc(fn func() (string, []interface{})) {
	format, args := fn()
	log.printf("Debug", format, args...)
}

func (log *BasicLogger) Info(args ...interface{}) {
	log.print("Info", args...)
}

func (log *BasicLogger) Infof(format string, args ...interface{}) {
	log.printf("Info", format, args...)
}

func (log *BasicLogger) InfoFunc(fn func() (string, []interface{})) {
	format, args := fn()
	log.printf("Info", format, args...)
}

func (log *BasicLogger) Warn(args ...interface{}) {
	log.print("Warn", args...)
}

func (log *BasicLogger) Warnf(format string, args ...interface{}) {
	log.printf("Warn", format, args...)
}

func (log *BasicLogger) WarnFunc(fn func() (string, []interface{})) {
	format, args := fn()
	log.printf("Warn", format, args...)
}

func (log *BasicLogger) Error(args ...interface{}) {
	log.print("Error", args...)
}

func (log *BasicLogger) Errorf(format string, args ...interface{}) {
	log.printf("Error", format, args...)
}

func (log *BasicLogger) ErrorFunc(fn func() (string, []interface{})) {
	format, args := fn()
	log.printf("Error", format, args...)
}

func (log *BasicLogger) Fatal(args ...interface{}) {
	log.print("Fatal", args...)
	os.Exit(1)
}

func (log *BasicLogger) Fatalf(format string, args ...interface{}) {
	log.printf("Fatal", format, args...)
	os.Exit(1)
}
