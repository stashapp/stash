package logger

import (
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()
var progressLogger = logrus.New()

func init() {
	progressLogger.SetFormatter(new(ProgressFormatter))
}

func Progressf(format string, args ...interface{}) {
	progressLogger.Infof(format, args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
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