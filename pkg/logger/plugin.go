package logger

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type PluginLogLevel struct {
	char byte
	Name string
}

// Valid Level values.
var (
	TraceLevel = PluginLogLevel{
		char: 't',
		Name: "trace",
	}
	DebugLevel = PluginLogLevel{
		char: 'd',
		Name: "debug",
	}
	InfoLevel = PluginLogLevel{
		char: 'i',
		Name: "info",
	}
	WarningLevel = PluginLogLevel{
		char: 'w',
		Name: "warning",
	}
	ErrorLevel = PluginLogLevel{
		char: 'e',
		Name: "error",
	}
	ProgressLevel = PluginLogLevel{
		char: 'p',
		Name: "progress",
	}
	NoneLevel = PluginLogLevel{
		Name: "none",
	}
)

var validLevels = []PluginLogLevel{
	TraceLevel,
	DebugLevel,
	InfoLevel,
	WarningLevel,
	ErrorLevel,
	ProgressLevel,
	NoneLevel,
}

const startLevelChar byte = 1
const endLevelChar byte = 2

func (l PluginLogLevel) prefix() string {
	return string([]byte{
		startLevelChar,
		byte(l.char),
		endLevelChar,
	})
}

func (l PluginLogLevel) Log(args ...interface{}) {
	if l.char == 0 {
		return
	}

	argsToUse := []interface{}{
		l.prefix(),
	}
	argsToUse = append(argsToUse, args...)
	fmt.Fprintln(os.Stderr, argsToUse...)
}

func (l PluginLogLevel) Logf(format string, args ...interface{}) {
	if l.char == 0 {
		return
	}

	formatToUse := string(l.prefix()) + format + "\n"
	fmt.Fprintf(os.Stderr, formatToUse, args...)
}

// PluginLogLevelFromName returns the Level that matches the provided name or nil if
// the name does not match a valid value.
func PluginLogLevelFromName(name string) *PluginLogLevel {
	for _, l := range validLevels {
		if l.Name == name {
			return &l
		}
	}

	return nil
}

// DetectLogLevel returns the Level and the logging string for a provided line
// of plugin output. It parses the string for logging control characters and
// determines the log level, if present. If not present, the plugin output
// is returned unchanged with a nil Level.
func DetectLogLevel(line string) (*PluginLogLevel, string) {
	if len(line) < 4 || line[0] != startLevelChar || line[2] != endLevelChar {
		return nil, line
	}

	char := line[1]
	var level *PluginLogLevel
	for _, l := range validLevels {
		if l.char == char {
			level = &l
			break
		}
	}

	if level == nil {
		return nil, line
	}

	line = strings.TrimSpace(line[3:])

	return level, line
}

type PluginLogger struct {
	Prefix          string
	DefaultLogLevel *PluginLogLevel
	ProgressChan    chan float64
}

func (log *PluginLogger) HandleStderrLine(line string) {
	level, ll := DetectLogLevel(line)

	// if no log level, just output to info
	if level == nil {
		if log.DefaultLogLevel != nil {
			level = log.DefaultLogLevel
		} else {
			level = &InfoLevel
		}
	}

	switch *level {
	case TraceLevel:
		Trace(log.Prefix, ll)
	case DebugLevel:
		Debug(log.Prefix, ll)
	case InfoLevel:
		Info(log.Prefix, ll)
	case WarningLevel:
		Warn(log.Prefix, ll)
	case ErrorLevel:
		Error(log.Prefix, ll)
	case ProgressLevel:
		p, err := strconv.ParseFloat(ll, 64)
		if err != nil {
			Errorf("Error parsing progress value '%s': %s", ll, err.Error())
		} else {
			// only pass progress through if channel present
			if log.ProgressChan != nil {
				// don't block on this
				select {
				case log.ProgressChan <- p:
				default:
				}
			}
		}
	}
}

func (log *PluginLogger) HandlePluginStdErr(pluginStdErr io.ReadCloser) {
	// pipe plugin stderr to our logging
	scanner := bufio.NewScanner(pluginStdErr)
	for scanner.Scan() {
		str := scanner.Text()
		if str != "" {
			log.HandleStderrLine(str)
		}
	}

	str := scanner.Text()
	if str != "" {
		log.HandleStderrLine(str)
	}

	pluginStdErr.Close()
}
