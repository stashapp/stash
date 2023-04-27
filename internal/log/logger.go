package log

import (
	"container/list"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const CACHE_SIZE = 1000

type LogItem struct {
	Time    time.Time `json:"time"`
	Type    string    `json:"type"`
	Message string    `json:"message"`
}

type Logger struct {
	logger         *logrus.Logger
	progressLogger *logrus.Logger
	mutex          sync.Mutex
	logCache       logCache
	logSubs        []chan []LogItem
	waiting        bool
	lastBroadcast  time.Time
	logBuffer      []LogItem
}

type logCache struct {
	trace   *list.List
	debug   *list.List
	info    *list.List
	warning *list.List
	error   *list.List
}

func NewLogger() *Logger {
	ret := &Logger{
		logger:         logrus.New(),
		progressLogger: logrus.New(),
		logCache: logCache{
			trace:   list.New(),
			debug:   list.New(),
			info:    list.New(),
			warning: list.New(),
			error:   list.New(),
		},
		lastBroadcast: time.Now(),
	}

	ret.progressLogger.SetFormatter(new(ProgressFormatter))

	return ret
}

// Init initialises the logger based on a logging configuration
func (log *Logger) Init(logFile string, logOut bool, logLevel string) {
	var file *os.File
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

	if logFile != "" {
		var err error
		file, err = os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

		if err != nil {
			fmt.Printf("Could not open '%s' for log output due to error: %s\n", logFile, err.Error())
		}
	}

	if file != nil {
		if logOut {
			// log to file separately disabling colours
			fileFormatter := new(logrus.TextFormatter)
			fileFormatter.TimestampFormat = customFormatter.TimestampFormat
			fileFormatter.FullTimestamp = customFormatter.FullTimestamp
			log.logger.AddHook(&fileLogHook{
				Writer:    file,
				Formatter: fileFormatter,
			})
		} else {
			// logging to file only
			// turn off the colouring for the file
			customFormatter.ForceColors = false
			log.logger.Out = file
		}
	}

	// otherwise, output to StdErr

	log.SetLogLevel(logLevel)
}

func (log *Logger) SetLogLevel(level string) {
	log.logger.Level = logLevelFromString(level, logrus.InfoLevel)
}

func logLevelFromString(level string, defaultLevel logrus.Level) logrus.Level {
	switch strings.ToLower(level) {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warning":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "trace":
		return logrus.TraceLevel
	}

	return defaultLevel
}

func addToCacheList(cache *list.List, l *LogItem) {
	// remove last item if list is full
	if cache.Len() == CACHE_SIZE {
		cache.Remove(cache.Back())
	}
	// prepend new item
	cache.PushFront(l)
}

// assumes mutex held
func (log *Logger) addToCache(l *LogItem) {
	switch logLevelFromString(l.Type, logrus.DebugLevel) {
	case logrus.TraceLevel:
		addToCacheList(log.logCache.trace, l)
	case logrus.DebugLevel:
		addToCacheList(log.logCache.trace, l)
		addToCacheList(log.logCache.debug, l)
	case logrus.InfoLevel:
		addToCacheList(log.logCache.trace, l)
		addToCacheList(log.logCache.debug, l)
		addToCacheList(log.logCache.info, l)
	case logrus.WarnLevel:
		addToCacheList(log.logCache.trace, l)
		addToCacheList(log.logCache.debug, l)
		addToCacheList(log.logCache.info, l)
		addToCacheList(log.logCache.warning, l)
	case logrus.ErrorLevel:
		addToCacheList(log.logCache.trace, l)
		addToCacheList(log.logCache.debug, l)
		addToCacheList(log.logCache.info, l)
		addToCacheList(log.logCache.warning, l)
		addToCacheList(log.logCache.error, l)
	}
}

func (log *Logger) addLogItem(l *LogItem) {
	log.mutex.Lock()
	l.Time = time.Now()
	log.addToCache(l)
	log.logBuffer = append(log.logBuffer, *l)
	log.mutex.Unlock()
	go log.broadcastLogItem(l)
}

// Returns a list of recent log items, with more recent items first.
// minLevel sets the minimum log level of items to return.
// set minLevel to the empty string to return all items.
func (log *Logger) GetLogCache(minLevel string) []*LogItem {
	level := logLevelFromString(minLevel, logrus.TraceLevel)

	log.mutex.Lock()

	var items *list.List

	switch level {
	case logrus.TraceLevel:
		items = log.logCache.trace
	case logrus.DebugLevel:
		items = log.logCache.debug
	case logrus.InfoLevel:
		items = log.logCache.info
	case logrus.WarnLevel:
		items = log.logCache.warning
	case logrus.ErrorLevel:
		items = log.logCache.error
	}

	ret := make([]*LogItem, items.Len())
	i := 0
	for e := items.Front(); e != nil; e = e.Next() {
		ret[i] = e.Value.(*LogItem)
		i++
	}

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

	log.mutex.Unlock()
}

func (log *Logger) Progressf(format string, args ...interface{}) {
	log.progressLogger.Infof(format, args...)
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
	msg, args := fn()
	log.Tracef(msg, args...)
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

func (log *Logger) DebugFunc(fn func() (string, []interface{})) {
	msg, args := fn()
	log.Debugf(msg, args...)
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
	msg, args := fn()
	log.Infof(msg, args...)
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
	msg, args := fn()
	log.Warnf(msg, args...)
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
	msg, args := fn()
	log.Errorf(msg, args...)
}

func (log *Logger) Fatal(args ...interface{}) {
	log.logger.Fatal(args...)
}

func (log *Logger) Fatalf(format string, args ...interface{}) {
	log.logger.Fatalf(format, args...)
}
