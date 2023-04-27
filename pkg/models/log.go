package models

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type LogLevel uint32

const (
	LogLevelProgress LogLevel = iota
	LogLevelTrace
	LogLevelDebug
	LogLevelInfo
	LogLevelWarning
	LogLevelError
)

var AllLogLevel = []LogLevel{
	LogLevelProgress,
	LogLevelTrace,
	LogLevelDebug,
	LogLevelInfo,
	LogLevelWarning,
	LogLevelError,
}

func (e LogLevel) IsValid() bool {
	switch e {
	case LogLevelProgress, LogLevelTrace, LogLevelDebug, LogLevelInfo, LogLevelWarning, LogLevelError:
		return true
	}
	return false
}

func (e LogLevel) String() string {
	switch e {
	case LogLevelProgress:
		return "Progress"
	case LogLevelTrace:
		return "Trace"
	case LogLevelDebug:
		return "Debug"
	case LogLevelInfo:
		return "Info"
	case LogLevelWarning:
		return "Warning"
	case LogLevelError:
		return "Error"
	default:
		return "Invalid"
	}
}

func (e *LogLevel) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	switch strings.ToLower(str) {
	case "progress":
		*e = LogLevelProgress
	case "trace":
		*e = LogLevelTrace
	case "debug":
		*e = LogLevelDebug
	case "info":
		*e = LogLevelInfo
	case "warning":
		*e = LogLevelWarning
	case "error":
		*e = LogLevelError
	default:
		return fmt.Errorf("%s is not a valid LogLevel", str)
	}

	return nil
}

func (e LogLevel) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
