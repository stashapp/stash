package log

import (
	"fmt"
)

// Logger handles logging in a specific context. It includes a bunch of helpers and compatibility
// over the loggerCore.
type Logger struct {
	loggerCore
}

func (l Logger) WithText(f func(Msg) string) Logger {
	l.msgMaps = append(l.msgMaps, func(msg Msg) Msg {
		return msg.WithText(f)
	})
	return l
}

// Helper for compatibility with "log".Logger.
func (l Logger) Printf(format string, a ...interface{}) {
	l.LazyLog(l.defaultLevel, func() Msg {
		return Fmsg(format, a...).Skip(1)
	})
}

func (l Logger) Log(m Msg) {
	l.LogLevel(l.defaultLevel, m.Skip(1))
}

func (l Logger) LogLevel(level Level, m Msg) {
	l.LazyLog(level, func() Msg {
		return m.Skip(1)
	})
}

// Helper for compatibility with "log".Logger.
func (l Logger) Print(v ...interface{}) {
	l.LazyLog(l.defaultLevel, func() Msg {
		return Str(fmt.Sprint(v...)).Skip(1)
	})
}

func (l Logger) WithContextValue(v interface{}) Logger {
	return l.WithText(func(m Msg) string {
		return fmt.Sprintf("%v: %v", v, m)
	})
}

func (l Logger) WithContextText(s string) Logger {
	return l.WithText(func(m Msg) string {
		return s + ": " + m.Text()
	})
}

func (l Logger) SkipCallers(skip int) Logger {
	return l.WithMap(func(m Msg) Msg {
		return m.Skip(skip)
	})
}

func (l Logger) Levelf(level Level, format string, a ...interface{}) {
	l.LazyLog(level, func() Msg {
		return Fmsg(format, a...).Skip(1)
	})
}

func (l Logger) Println(a ...interface{}) {
	l.LazyLogDefaultLevel(func() Msg {
		return Str(fmt.Sprintln(a...)).Skip(1)
	})
}
