package log

import (
	"fmt"
	"github.com/mgutz/ansi"
)

// Logging-related utility functions.

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
)
var logLevel int

func init() {
	logLevel = INFO
}

// SetLevel sets the level of verbosity from the logging system.
//
// Set it to one of the constants like `INFO`.
func SetLevel(level int) {
	logLevel = level
}

// formatLogLine is a helper function to format a colorful log line.
//
// This prints the log line using a bright version of the color for the log
// level and a darker version for its text.
//
// Parameters:
// - color:   the color string to use, e.g. `cyan` `blue` `yellow`
// - level:   the log level name, e.g. `DEBUG`, `INFO`, `WARN`
// - message: the actual log message
// - a:       parameters for Printf to substitute into the message.
func formatLogLine(color, level, message string, a ...interface{}) {
	bright := fmt.Sprintf("%s+h", color)
	fmt.Printf(
		ansi.Color(fmt.Sprintf("[%s] ", level), bright) +
		ansi.Color(message, color) + "\n",
		a...
	)
}

// Debug emits a debug log message.
func Debug(message string, a ...interface{}) {
	if logLevel <= DEBUG {
		formatLogLine("blue", "DEBUG", message, a...)
	}
}

// Info emits an informational log message.
func Info(message string, a ...interface{}) {
	if logLevel <= INFO {
		formatLogLine("cyan", "INFO", message, a...)
	}
}

// Warn emits a warning log message.
func Warn(message string, a ...interface{}) {
	if logLevel <= WARN {
		formatLogLine("yellow", "WARN", message, a...)
	}
}

// Error emits an error log message.
func Error(message string, a ...interface{}) {
	if logLevel <= ERROR {
		formatLogLine("red", "ERROR", message, a...)
	}
}
