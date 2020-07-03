package log

import (
	"fmt"
	"os"
	"strings"
)

const startLevelChar byte = 1
const endLevelChar byte = 2

type LogLevel byte

const (
	TraceLevel   LogLevel = 't'
	DebugLevel   LogLevel = 'd'
	InfoLevel    LogLevel = 'i'
	WarningLevel LogLevel = 'w'
	ErrorLevel   LogLevel = 'e'
)

func (l LogLevel) valid() bool {
	return l == TraceLevel || l == DebugLevel || l == InfoLevel || l == WarningLevel || l == ErrorLevel
}

func (l LogLevel) prefix() string {
	return string([]byte{
		startLevelChar,
		byte(l),
		endLevelChar,
	})
}

func (l LogLevel) log(args ...interface{}) {
	argsToUse := []interface{}{
		l.prefix(),
	}
	argsToUse = append(argsToUse, args...)
	fmt.Fprintln(os.Stderr, argsToUse...)
}

func (l LogLevel) logf(format string, args ...interface{}) {
	formatToUse := string(l.prefix()) + format + "\n"
	fmt.Fprintf(os.Stderr, formatToUse, args...)
}

func Trace(args ...interface{}) {
	TraceLevel.log(args...)
}

func Tracef(format string, args ...interface{}) {
	TraceLevel.logf(format, args...)
}

func Debug(args ...interface{}) {
	DebugLevel.log(args...)
}

func Debugf(format string, args ...interface{}) {
	DebugLevel.logf(format, args...)
}

func Info(args ...interface{}) {
	InfoLevel.log(args...)
}

func Infof(format string, args ...interface{}) {
	InfoLevel.logf(format, args...)
}

func Warn(args ...interface{}) {
	WarningLevel.log(args...)
}

func Warnf(format string, args ...interface{}) {
	WarningLevel.logf(format, args...)
}

func Error(args ...interface{}) {
	ErrorLevel.log(args...)
}

func Errorf(format string, args ...interface{}) {
	ErrorLevel.logf(format, args...)
}

func DetectLogLevel(line string) (*LogLevel, string) {
	if len(line) < 4 || line[0] != startLevelChar || line[2] != endLevelChar {
		return nil, line
	}

	level := LogLevel(line[1])
	if !level.valid() {
		return nil, line
	}

	line = strings.TrimSpace(line[3:])

	return &level, line
}
