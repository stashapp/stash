package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
)

type LogItem struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

var logger = logrus.New()
var progressLogger = logrus.New()

var LogCache []LogItem
var mutex = &sync.Mutex{}

func addLogItem(l *LogItem) {
	mutex.Lock()
	LogCache = append([]LogItem{*l}, LogCache...)
	if len(LogCache) > 30 {
		LogCache = LogCache[:len(LogCache)-1]
	}
	mutex.Unlock()
}

func init() {
	progressLogger.SetFormatter(new(ProgressFormatter))
}

func Progressf(format string, args ...interface{}) {
	progressLogger.Infof(format, args...)
	l := &LogItem{
		Type: "progress",
		Message: fmt.Sprintf(format, args...),
	}
	addLogItem(l)

}

func Trace(args ...interface{}) {
	logger.Trace(args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
	l := &LogItem{
		Type: "debug",
		Message: fmt.Sprint(args),
	}
	addLogItem(l)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
	l := &LogItem{
		Type: "debug",
		Message: fmt.Sprintf(format, args...),
	}
	addLogItem(l)
}

func Info(args ...interface{}) {
	logger.Info(args...)
	l := &LogItem{
		Type: "info",
		Message: fmt.Sprint(args),
	}
	addLogItem(l)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
	l := &LogItem{
		Type: "info",
		Message: fmt.Sprintf(format, args...),
	}
	addLogItem(l)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
	l := &LogItem{
		Type: "warn",
		Message: fmt.Sprint(args),
	}
	addLogItem(l)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
	l := &LogItem{
		Type: "warn",
		Message: fmt.Sprintf(format, args...),
	}
	addLogItem(l)
}

func Error(args ...interface{}) {
	logger.Error(args...)
	l := &LogItem{
		Type: "error",
		Message: fmt.Sprint(args),
	}
	addLogItem(l)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
	l := &LogItem{
		Type: "error",
		Message: fmt.Sprintf(format, args...),
	}
	addLogItem(l)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

//func WithRequest(req *http.Request) *logrus.Entry {
//	return logger.WithFields(RequestFields(req))
//}