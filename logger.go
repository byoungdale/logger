// Package logger provides a simple logging framework with different log levels and customizable output.
package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var requestedLevel = InfoLevel
var displayDateTime = false
var outputDest io.Writer = os.Stderr

// LogLevel type represents the different logging levels.
type LogLevel uint32

const (
	// FatalLevel should be used in fatal situations, and the application will exit.
	FatalLevel LogLevel = iota

	// ErrorLevel should be used when someone should really look at the error.
	ErrorLevel

	// InfoLevel should be used during normal operations.
	InfoLevel

	// WarnLevel should be used for things that should be addressed at some point.
	WarnLevel

	// DebugLevel should be used only during development.
	DebugLevel
)

// ANSI escape codes for console text coloring.
const (
	LogLColorReset = "\033[0m"

	LogLColorRed    = "\033[31m"
	LogLColorGreen  = "\033[32m"
	LogLColorYellow = "\033[33m"
	LogLColorBlue   = "\033[34m"
	LogLColorPurple = "\033[35m"
	LogLColorCyan   = "\033[36m"
	LogLColorWhite  = "\033[37m"

	LogLColorBold      = "\033[1m"
	LogLColorBoldReset = "\033[22m"
	LogLColorUnderline = "\033[4m"
	LogLColorReversed  = "\033[7m"
)

// LogPrefix returns the formatted prefix for log messages based on the log level.
func LogPrefix(level LogLevel) string {
	boldSet := ""
	boldReset := ""

	switch level {
	case DebugLevel:
		return boldSet + "debug" + boldReset
	case InfoLevel:
		return boldSet + "info" + boldReset
	case WarnLevel:
		return boldSet + "warn" + boldReset
	case ErrorLevel:
		return boldSet + "error" + boldReset
	case FatalLevel:
		return boldSet + "fatal" + boldReset
	default:
		return boldSet + "none" + boldReset
	}
}

// LogColorSet returns the ANSI escape code for text color based on the log level.
func LogColorSet(level LogLevel) string {
	switch level {
	case DebugLevel:
		return ""
	case InfoLevel:
		return LogLColorCyan
	case WarnLevel:
		return LogLColorYellow
	case ErrorLevel:
		return LogLColorRed
	case FatalLevel:
		return LogLColorRed
	default:
		return LogLColorGreen
	}
}

// LogColorReset returns the ANSI escape code to reset text color based on the log level.
func LogColorReset(level LogLevel) string {
	if level == DebugLevel {
		return ""
	}
	return LogLColorReset
}

// LogTimestamp returns the current timestamp in a specific format.
func LogTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05.0000")
}

// String returns a string representation of the LogLevel.
func (level LogLevel) String() string {
	switch level {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// SetOutputToFile sets the log output destination to a file with the given path.
func SetOutputToFile(logFilePath *string) {
	var err error
	outputDest, err = os.OpenFile(*logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Fatal("could not open log file %s", *logFilePath)
	}
}

// EnableDateTime enables date time in log messages.
func EnableDateTime() {
	displayDateTime = true
}

// EnableDebug increases logging, more verbose (debug)
func EnableDebug() {
	requestedLevel = DebugLevel
	formatMessage(InfoLevel, "Debug mode enabled")
}

// Debug sends a debug log message.
func Debug(format string, v ...interface{}) {
	if requestedLevel >= DebugLevel {
		formatMessage(DebugLevel, format, v...)
	}
}

// Info sends an info log message.
func Info(format string, v ...interface{}) {
	if requestedLevel >= InfoLevel {
		formatMessage(InfoLevel, format, v...)
	}
}

// Warn sends an info log message.
func Warn(format string, v ...interface{}) {
	if requestedLevel >= InfoLevel {
		formatMessage(InfoLevel, format, v...)
	}
}

// Error sends an error log message.
func Error(format string, v ...interface{}) {
	if requestedLevel >= ErrorLevel {
		formatMessage(ErrorLevel, format, v...)
	}
}

// Fatal sends a fatal log message and stop the execution of the program.
func Fatal(format string, v ...interface{}) {
	if requestedLevel >= FatalLevel {
		formatMessage(FatalLevel, format, v...)
		os.Exit(1)
	}
}

// formatMessage formats and writes the log message to the output destination.
func formatMessage(level LogLevel, format string, v ...interface{}) {
	pc, filename, line, _ := runtime.Caller(2)
	logmsg := fmt.Sprintf(format, v...)
	fmt.Fprintf(outputDest, "%s%s [%s] [%s:%d] %s(): %s%s\n", LogColorSet(level),
		LogTimestamp(),
		LogPrefix(level),
		filepath.Base(filename), line, runtime.FuncForPC(pc).Name(),
		logmsg,
		LogColorReset(level))
}
