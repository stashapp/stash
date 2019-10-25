package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type LogItem struct {
	Time    time.Time `json:"time"`
	Type    string    `json:"type"`
	Message string    `json:"message"`
}

var logger = logrus.New()
var progressLogger = logrus.New()

var LogCache []LogItem
var mutex = &sync.Mutex{}
var logSubs []chan []LogItem
var waiting = false
var lastBroadcast = time.Now()
var logBuffer []LogItem

// Init initialises the logger based on a logging configuration
func Init(logFile string, logOut bool, logLevel string) {
	var file *os.File

	if logFile != "" {
		var err error
		file, err = os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

		if err != nil {
			fmt.Printf("Could not open '%s' for log output due to error: %s\n", logFile, err.Error())
			logFile = ""
		}
	}

	if file != nil && logOut {
		mw := io.MultiWriter(os.Stderr, file)
		logger.Out = mw
	} else if file != nil {
		logger.Out = file
	}

	// otherwise, output to StdErr

	SetLogLevel(logLevel)
}

func SetLogLevel(level string) {
	logger.Level = logLevelFromString(level)
}

func logLevelFromString(level string) logrus.Level {
	ret := logrus.InfoLevel

	if level == "Debug" {
		ret = logrus.DebugLevel
	} else if level == "Warning" {
		ret = logrus.WarnLevel
	} else if level == "Error" {
		ret = logrus.ErrorLevel
	}

	return ret
}

func addLogItem(l *LogItem) {
	mutex.Lock()
	l.Time = time.Now()
	LogCache = append([]LogItem{*l}, LogCache...)
	if len(LogCache) > 30 {
		LogCache = LogCache[:len(LogCache)-1]
	}
	mutex.Unlock()
	go broadcastLogItem(l)
}

func GetLogCache() []LogItem {
	mutex.Lock()

	ret := make([]LogItem, len(LogCache))
	copy(ret, LogCache)

	mutex.Unlock()

	return ret
}

func SubscribeToLog(stop chan int) <-chan []LogItem {
	ret := make(chan []LogItem, 100)

	go func() {
		<-stop
		unsubscribeFromLog(ret)
	}()

	mutex.Lock()
	logSubs = append(logSubs, ret)
	mutex.Unlock()

	return ret
}

func unsubscribeFromLog(toRemove chan []LogItem) {
	mutex.Lock()
	for i, c := range logSubs {
		if c == toRemove {
			logSubs = append(logSubs[:i], logSubs[i+1:]...)
		}
	}
	close(toRemove)
	mutex.Unlock()
}

func doBroadcastLogItems() {
	// assumes mutex held

	for _, c := range logSubs {
		// don't block waiting to broadcast
		select {
		case c <- logBuffer:
		default:
		}
	}

	logBuffer = nil
	waiting = false
	lastBroadcast = time.Now()
}

func broadcastLogItem(l *LogItem) {
	mutex.Lock()

	logBuffer = append(logBuffer, *l)

	// don't send more than once per second
	if !waiting {
		// if last broadcast was under a second ago, wait until a second has
		// passed
		timeSinceBroadcast := time.Since(lastBroadcast)
		if timeSinceBroadcast.Seconds() < 1 {
			waiting = true
			time.AfterFunc(time.Second-timeSinceBroadcast, func() {
				mutex.Lock()
				doBroadcastLogItems()
				mutex.Unlock()
			})
		} else {
			doBroadcastLogItems()
		}
	}
	// if waiting then adding it to the buffer is sufficient

	mutex.Unlock()
}

func init() {
	progressLogger.SetFormatter(new(ProgressFormatter))
}

func Progressf(format string, args ...interface{}) {
	progressLogger.Infof(format, args...)
	l := &LogItem{
		Type:    "progress",
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
		Type:    "debug",
		Message: fmt.Sprint(args...),
	}
	addLogItem(l)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
	l := &LogItem{
		Type:    "debug",
		Message: fmt.Sprintf(format, args...),
	}
	addLogItem(l)
}

func Info(args ...interface{}) {
	logger.Info(args...)
	l := &LogItem{
		Type:    "info",
		Message: fmt.Sprint(args...),
	}
	addLogItem(l)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
	l := &LogItem{
		Type:    "info",
		Message: fmt.Sprintf(format, args...),
	}
	addLogItem(l)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
	l := &LogItem{
		Type:    "warn",
		Message: fmt.Sprint(args...),
	}
	addLogItem(l)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
	l := &LogItem{
		Type:    "warn",
		Message: fmt.Sprintf(format, args...),
	}
	addLogItem(l)
}

func Error(args ...interface{}) {
	logger.Error(args...)
	l := &LogItem{
		Type:    "error",
		Message: fmt.Sprint(args...),
	}
	addLogItem(l)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
	l := &LogItem{
		Type:    "error",
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
