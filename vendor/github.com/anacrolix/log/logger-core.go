package log

import (
	"fmt"
)

// loggerCore is the essential part of Logger.
type loggerCore struct {
	nonZero      bool
	names        []string
	values       []interface{}
	defaultLevel Level
	filterLevel  Level
	msgMaps      []func(Msg) Msg
	Handlers     []Handler
}

func (l loggerCore) asLogger() Logger {
	return Logger{l}
}

// Returns a logger that adds the given values to logged messages.
func (l loggerCore) WithValues(v ...interface{}) Logger {
	l.values = append(l.values, v...)
	return l.asLogger()
}

// Returns a logger that for a given message propagates the result of `f` instead.
func (l loggerCore) WithMap(f func(m Msg) Msg) Logger {
	l.msgMaps = append(l.msgMaps, f)
	return l.asLogger()
}

func (l loggerCore) WithDefaultLevel(level Level) Logger {
	l.defaultLevel = level
	return l.asLogger()
}

func (l loggerCore) FilterLevel(minLevel Level) Logger {
	if _, ok := levelFromRules(l.names); !ok {
		l.filterLevel = minLevel
	}
	return l.asLogger()
}

func (l loggerCore) IsZero() bool {
	return !l.nonZero
}

func (l loggerCore) IsEnabledFor(level Level) bool {
	return !level.LessThan(l.filterLevel)
}

func (l loggerCore) LazyLog(level Level, f func() Msg) {
	l.lazyLog(level, 1, f)
}

func (l loggerCore) LazyLogDefaultLevel(f func() Msg) {
	l.lazyLog(l.defaultLevel, 1, f)
}

func (l loggerCore) lazyLog(level Level, skip int, f func() Msg) {
	if !l.IsEnabledFor(level) {
		// have a big sook
		//internalLogger.Levelf(Debug, "skipped logging %v for %q", level, l.names)
		return
	}
	r := f().Skip(skip + 1)
	for i := len(l.msgMaps) - 1; i >= 0; i-- {
		r = l.msgMaps[i](r)
	}
	l.handle(level, r)
}

func (l loggerCore) handle(level Level, m Msg) {
	r := Record{
		Msg:   m.Skip(1),
		Level: level,
		Names: l.names,
	}
	if !l.nonZero {
		panic(fmt.Sprintf("Logger uninitialized. names=%q", l.names))
	}
	for _, h := range l.Handlers {
		h.Handle(r)
	}
}

func (l loggerCore) WithNames(names ...string) Logger {
	// Avoid sharing after appending. This might not be enough because some formatters might add
	// more elements concurrently, or names could be empty.
	l.names = append(l.names[:len(l.names):len(l.names)], names...)
	return l.withFilterLevelFromRules()
}

func (l loggerCore) withFilterLevelFromRules() Logger {
	level, ok := levelFromRules(l.names)
	if ok {
		l.filterLevel = level
	}
	return l.asLogger()
}
