// Package log provides centralized logging infrastructure
package log

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	// DefaultLogger is the default application logger
	DefaultLogger *log.Logger
)

// Init initializes the logging system
func Init(level, format string, production bool) {
	DefaultLogger = log.New()

	// Set log level
	lvl, err := log.ParseLevel(level)
	if err != nil {
		lvl = log.InfoLevel
	}
	DefaultLogger.SetLevel(lvl)

	// Set format
	if production {
		DefaultLogger.SetFormatter(&log.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	} else {
		DefaultLogger.SetFormatter(&log.TextFormatter{
			TimestampFormat: time.RFC3339,
			FullTimestamp:   true,
		})
	}

	// Set output
	DefaultLogger.SetOutput(os.Stdout)

	// Add common fields
	DefaultLogger.SetReportCaller(false)
}

// SetLevel sets the log level
func SetLevel(level string) error {
	lvl, err := log.ParseLevel(level)
	if err != nil {
		return err
	}
	DefaultLogger.SetLevel(lvl)
	return nil
}

// WithField returns a logger with a single field
func WithField(key string, value interface{}) *log.Entry {
	return DefaultLogger.WithField(key, value)
}

// WithFields returns a logger with multiple fields
func WithFields(fields log.Fields) *log.Entry {
	return DefaultLogger.WithFields(fields)
}

// Debug logs a debug message
func Debug(msg string) {
	DefaultLogger.Debug(msg)
}

// Info logs an info message
func Info(msg string) {
	DefaultLogger.Info(msg)
}

// Warn logs a warning message
func Warn(msg string) {
	DefaultLogger.Warn(msg)
}

// Error logs an error message
func Error(msg string) {
	DefaultLogger.Error(msg)
}

// Fatal logs a fatal message and exits
func Fatal(msg string) {
	DefaultLogger.Fatal(msg)
}

// Debugf logs a formatted debug message
func Debugf(format string, args ...interface{}) {
	DefaultLogger.Debugf(format, args...)
}

// Infof logs a formatted info message
func Infof(format string, args ...interface{}) {
	DefaultLogger.Infof(format, args...)
}

// Warnf logs a formatted warning message
func Warnf(format string, args ...interface{}) {
	DefaultLogger.Warnf(format, args...)
}

// Errorf logs a formatted error message
func Errorf(format string, args ...interface{}) {
	DefaultLogger.Errorf(format, args...)
}

// Fatalf logs a formatted fatal message and exits
func Fatalf(format string, args ...interface{}) {
	DefaultLogger.Fatalf(format, args...)
}
