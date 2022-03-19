package logger

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// PluginLogLevel represents a logging level for plugins to send log messages to stash.
type PluginLogLevel struct {
	char byte
	name string
}

// Valid Level values.
var (
	TraceLevel = PluginLogLevel{
		char: 't',
		name: "trace",
	}
	DebugLevel = PluginLogLevel{
		char: 'd',
		name: "debug",
	}
	InfoLevel = PluginLogLevel{
		char: 'i',
		name: "info",
	}
	WarningLevel = PluginLogLevel{
		char: 'w',
		name: "warning",
	}
	ErrorLevel = PluginLogLevel{
		char: 'e',
		name: "error",
	}
	ProgressLevel = PluginLogLevel{
		char: 'p',
		name: "progress",
	}
	NoneLevel = PluginLogLevel{
		name: "none",
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

// Log prints the provided message to os.Stderr in a format that provides the correct LogLevel for stash.
// The message is formatted in the same way as fmt.Println.
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

// Logf prints the provided message to os.Stderr in a format that provides the correct LogLevel for stash.
// The message is formatted in the same way as fmt.Printf.
func (l PluginLogLevel) Logf(format string, args ...interface{}) {
	if l.char == 0 {
		return
	}

	formatToUse := string(l.prefix()) + format + "\n"
	fmt.Fprintf(os.Stderr, formatToUse, args...)
}

// PluginLogLevelFromName returns the PluginLogLevel that matches the provided name or nil if
// the name does not match a valid value.
func PluginLogLevelFromName(name string) *PluginLogLevel {
	for _, l := range validLevels {
		if l.name == name {
			return &l
		}
	}

	return nil
}

// detectLogLevel returns the Level and the logging string for a provided line
// of plugin output. It parses the string for logging control characters and
// determines the log level, if present. If not present, the plugin output
// is returned unchanged with a nil Level.
func detectLogLevel(line string) (*PluginLogLevel, string) {
	if len(line) < 4 || line[0] != startLevelChar || line[2] != endLevelChar {
		return nil, line
	}

	char := line[1]
	var level *PluginLogLevel
	for _, l := range validLevels {
		if l.char == char {
			l := l // Make a copy of the loop variable
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

// PluginLogger interprets incoming log messages from plugins and logs to the appropriate log level.
type PluginLogger struct {
	// Logger is the LoggerImpl to forward log messages to.
	Logger LoggerImpl
	// Prefix is the prefix to prepend to log messages.
	Prefix string
	// DefaultLogLevel is the log level used if a log level prefix is not present in the received log message.
	DefaultLogLevel *PluginLogLevel
	// ProgressChan is a channel that receives float64s indicating the current progress of an operation.
	ProgressChan chan float64
}

func (log *PluginLogger) handleStderrLine(line string) {
	if log.Logger == nil {
		return
	}

	logger := log.Logger

	level, ll := detectLogLevel(line)

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
		logger.Trace(log.Prefix, ll)
	case DebugLevel:
		logger.Debug(log.Prefix, ll)
	case InfoLevel:
		logger.Info(log.Prefix, ll)
	case WarningLevel:
		logger.Warn(log.Prefix, ll)
	case ErrorLevel:
		logger.Error(log.Prefix, ll)
	case ProgressLevel:
		p, err := strconv.ParseFloat(ll, 64)
		if err != nil {
			logger.Errorf("Error parsing progress value '%s': %s", ll, err.Error())
		} else if log.ProgressChan != nil { // only pass progress through if channel present
			// don't block on this
			select {
			case log.ProgressChan <- p:
			default:
			}
		}
	}
}

// ReadLogMessages reads plugin log messages from src, forwarding them to the PluginLoggers Logger.
// ProgressLevel messages are parsed as float64 and forwarded to ProgressChan. If ProgressChan is full,
// then the progress message is not forwarded.
// This method only returns when it reaches the end of src or encounters an error while reading src.
// This method closes src before returning.
func (log *PluginLogger) ReadLogMessages(src io.ReadCloser) {
	// pipe plugin stderr to our logging
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		str := scanner.Text()
		if str != "" {
			log.handleStderrLine(str)
		}
	}

	str := scanner.Text()
	if str != "" {
		log.handleStderrLine(str)
	}

	src.Close()
}
