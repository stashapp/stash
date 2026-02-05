// Package log provides an implementation of [logger.LoggerImpl], using logrus.
package log

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type LogItem struct {
	Time    time.Time `json:"time"`
	Type    string    `json:"type"`
	Message string    `json:"message"`
}

type Logger struct {
	logger         *logrus.Logger
	progressLogger *logrus.Logger
	mutex          sync.Mutex
	logCache       []LogItem
	logSubs        []chan []LogItem
	waiting        bool
	lastBroadcast  time.Time
	logBuffer      []LogItem
}

func NewLogger() *Logger {
	ret := &Logger{
		logger:         logrus.New(),
		progressLogger: logrus.New(),
		lastBroadcast:  time.Now(),
	}

	ret.progressLogger.SetFormatter(new(ProgressFormatter))

	return ret
}

// Init initialises the logger based on a logging configuration
func (log *Logger) Init(logFile string, logOut bool, logLevel string, logFileMaxSize int) {
	var logger io.WriteCloser
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.ForceColors = true
	customFormatter.FullTimestamp = true
	log.logger.SetOutput(os.Stderr)
	log.logger.SetFormatter(customFormatter)

	// #1837 - trigger the console to use color-mode since it won't be
	// otherwise triggered until the first log entry
	// this is covers the situation where the logger is only logging to file
	// and therefore does not trigger the console color-mode - resulting in
	// the access log colouring not being applied
	_, _ = customFormatter.Format(logrus.NewEntry(log.logger))

	// if size is 0, disable rotation
	if logFile != "" {
		if logFileMaxSize == 0 {
			var err error
			logger, err = os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "unable to open log file %s: %v\n", logFile, err)
			}
		} else {
			logger = &lumberjack.Logger{
				Filename: logFile,
				MaxSize:  logFileMaxSize, // Megabytes
				Compress: true,
			}
		}
	}

	if logger != nil {
		if logOut {
			// log to file separately disabling colours
			fileFormatter := new(logrus.TextFormatter)
			fileFormatter.TimestampFormat = customFormatter.TimestampFormat
			fileFormatter.FullTimestamp = customFormatter.FullTimestamp
			log.logger.AddHook(&fileLogHook{
				Writer:    logger,
				Formatter: fileFormatter,
			})
		} else {
			// logging to file only
			// turn off the colouring for the file
			customFormatter.ForceColors = false
			log.logger.Out = logger
		}
	}

	// otherwise, output to StdErr

	log.SetLogLevel(logLevel)
}

func (log *Logger) SetLogLevel(level string) {
	log.logger.Level = logLevelFromString(level)
}

func logLevelFromString(level string) logrus.Level {
	ret := logrus.InfoLevel

	switch strings.ToLower(level) {
	case "debug":
		ret = logrus.DebugLevel
	case "warning":
		ret = logrus.WarnLevel
	case "error":
		ret = logrus.ErrorLevel
	case "trace":
		ret = logrus.TraceLevel
	}

	return ret
}

func (log *Logger) addToCache(l *LogItem) {
	// assumes mutex held
	// only add to cache if meets minimum log level
	level := logLevelFromString(l.Type)
	if level <= log.logger.Level {
		log.logCache = append([]LogItem{*l}, log.logCache...)
		if len(log.logCache) > 30 {
			log.logCache = log.logCache[:len(log.logCache)-1]
		}
	}
}

func (log *Logger) addLogItem(l *LogItem) {
	log.mutex.Lock()
	l.Time = time.Now()
	log.addToCache(l)
	log.mutex.Unlock()
	go log.broadcastLogItem(l)
}

func (log *Logger) GetLogCache() []LogItem {
	log.mutex.Lock()

	ret := make([]LogItem, len(log.logCache))
	copy(ret, log.logCache)

	log.mutex.Unlock()

	return ret
}

func (log *Logger) SubscribeToLog(stop chan int) <-chan []LogItem {
	ret := make(chan []LogItem, 100)

	go func() {
		<-stop
		log.unsubscribeFromLog(ret)
	}()

	log.mutex.Lock()
	log.logSubs = append(log.logSubs, ret)
	log.mutex.Unlock()

	return ret
}

func (log *Logger) unsubscribeFromLog(toRemove chan []LogItem) {
	log.mutex.Lock()
	for i, c := range log.logSubs {
		if c == toRemove {
			log.logSubs = append(log.logSubs[:i], log.logSubs[i+1:]...)
		}
	}
	close(toRemove)
	log.mutex.Unlock()
}

func (log *Logger) doBroadcastLogItems() {
	// assumes mutex held

	for _, c := range log.logSubs {
		// don't block waiting to broadcast
		select {
		case c <- log.logBuffer:
		default:
		}
	}

	log.logBuffer = nil
	log.waiting = false
	log.lastBroadcast = time.Now()
}

func (log *Logger) broadcastLogItem(l *LogItem) {
	log.mutex.Lock()

	log.logBuffer = append(log.logBuffer, *l)

	// don't send more than once per second
	if !log.waiting {
		// if last broadcast was under a second ago, wait until a second has
		// passed
		timeSinceBroadcast := time.Since(log.lastBroadcast)
		if timeSinceBroadcast.Seconds() < 1 {
			log.waiting = true
			time.AfterFunc(time.Second-timeSinceBroadcast, func() {
				log.mutex.Lock()
				log.doBroadcastLogItems()
				log.mutex.Unlock()
			})
		} else {
			log.doBroadcastLogItems()
		}
	}

	// if waiting then adding it to the buffer is sufficient
	log.mutex.Unlock()
}

func (log *Logger) Progressf(format string, args ...interface{}) {
	log.progressLogger.Infof(format, args...)
	l := &LogItem{
		Type:    "progress",
		Message: fmt.Sprintf(format, args...),
	}
	log.addLogItem(l)
}

func (log *Logger) Trace(args ...interface{}) {
	log.logger.Trace(args...)
	l := &LogItem{
		Type:    "trace",
		Message: fmt.Sprint(args...),
	}
	log.addLogItem(l)
}

func (log *Logger) Tracef(format string, args ...interface{}) {
	log.logger.Tracef(format, args...)
	l := &LogItem{
		Type:    "trace",
		Message: fmt.Sprintf(format, args...),
	}
	log.addLogItem(l)
}

func (log *Logger) TraceFunc(fn func() (string, []interface{})) {
	if log.logger.Level >= logrus.TraceLevel {
		msg, args := fn()
		log.Tracef(msg, args...)
	}
}

func (log *Logger) Debug(args ...interface{}) {
	log.logger.Debug(args...)
	l := &LogItem{
		Type:    "debug",
		Message: fmt.Sprint(args...),
	}
	log.addLogItem(l)
}

func (log *Logger) Debugf(format string, args ...interface{}) {
	log.logger.Debugf(format, args...)
	l := &LogItem{
		Type:    "debug",
		Message: fmt.Sprintf(format, args...),
	}
	log.addLogItem(l)
}

func (log *Logger) logFunc(level logrus.Level, logFn func(format string, args ...interface{}), fn func() (string, []interface{})) {
	if log.logger.Level >= level {
		msg, args := fn()
		logFn(msg, args...)
	}
}

func (log *Logger) DebugFunc(fn func() (string, []interface{})) {
	log.logFunc(logrus.DebugLevel, log.logger.Debugf, fn)
}

func (log *Logger) Info(args ...interface{}) {
	log.logger.Info(args...)
	l := &LogItem{
		Type:    "info",
		Message: fmt.Sprint(args...),
	}
	log.addLogItem(l)
}

func (log *Logger) Infof(format string, args ...interface{}) {
	log.logger.Infof(format, args...)
	l := &LogItem{
		Type:    "info",
		Message: fmt.Sprintf(format, args...),
	}
	log.addLogItem(l)
}

func (log *Logger) InfoFunc(fn func() (string, []interface{})) {
	log.logFunc(logrus.InfoLevel, log.logger.Infof, fn)
}

func (log *Logger) Warn(args ...interface{}) {
	log.logger.Warn(args...)
	l := &LogItem{
		Type:    "warn",
		Message: fmt.Sprint(args...),
	}
	log.addLogItem(l)
}

func (log *Logger) Warnf(format string, args ...interface{}) {
	log.logger.Warnf(format, args...)
	l := &LogItem{
		Type:    "warn",
		Message: fmt.Sprintf(format, args...),
	}
	log.addLogItem(l)
}

func (log *Logger) WarnFunc(fn func() (string, []interface{})) {
	log.logFunc(logrus.WarnLevel, log.logger.Warnf, fn)
}

func (log *Logger) Error(args ...interface{}) {
	log.logger.Error(args...)
	l := &LogItem{
		Type:    "error",
		Message: fmt.Sprint(args...),
	}
	log.addLogItem(l)
}

func (log *Logger) Errorf(format string, args ...interface{}) {
	log.logger.Errorf(format, args...)
	l := &LogItem{
		Type:    "error",
		Message: fmt.Sprintf(format, args...),
	}
	log.addLogItem(l)
}

func (log *Logger) ErrorFunc(fn func() (string, []interface{})) {
	log.logFunc(logrus.ErrorLevel, log.logger.Errorf, fn)
}

func (log *Logger) Fatal(args ...interface{}) {
	log.logger.Fatal(args...)
}

func (log *Logger) Fatalf(format string, args ...interface{}) {
	log.logger.Fatalf(format, args...)
}
