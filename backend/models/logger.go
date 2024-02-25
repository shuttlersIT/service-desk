package models

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

// LogLevel represents different logging levels
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	TraceLevel
)

// LogFormat represents the log format type
type LogFormat int

const (
	TextFormat LogFormat = iota
	JSONFormat
)

// LoggerConfig holds the configuration for the logger
type LoggerConfig struct {
	MaxSizeMB    int
	MaxBackups   int
	RotatePeriod time.Duration
	LogLevel     LogLevel
	LogOutput    io.Writer
	LogFormat    LogFormat
}

// Logger interface defines the contract for logging methods
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Trace(args ...interface{})
	Printf(level LogLevel, format string, v ...interface{})
	SetOutput(output io.Writer)
	SetFlags(flags int)
	SetPrefix(prefix string)
	LogWithFields(level LogLevel, fields map[string]interface{}, format string, v ...interface{})
	LogWithStackTrace(level LogLevel, err error, format string, v ...interface{})
}

// PrintLogger is an implementation of the Logger interface using log.Printf
type PrintLogger struct {
	*log.Logger
	mu sync.Mutex
}

// NewLogger creates a new instance of Logger
func NewLogger(output io.Writer, prefix string, flags int) Logger {
	return &PrintLogger{
		Logger: log.New(output, prefix, flags),
	}
}

// Debug logs a debug message
func (pl *PrintLogger) Debug(args ...interface{}) {
	pl.Printf(DebugLevel, "DEBUG: %v", args...)
}

// Info logs an informational message
func (pl *PrintLogger) Info(args ...interface{}) {
	pl.Printf(InfoLevel, "INFO: %v", args...)
}

// Warn logs a warning message
func (pl *PrintLogger) Warn(args ...interface{}) {
	pl.Printf(WarnLevel, "WARN: %v", args...)
}

// Error logs an error message
func (pl *PrintLogger) Error(args ...interface{}) {
	pl.Printf(ErrorLevel, "ERROR: %v", args...)
}

// Fatal logs a fatal message and exits the application
func (pl *PrintLogger) Fatal(args ...interface{}) {
	pl.Printf(FatalLevel, "FATAL: %v", args...)
	os.Exit(1)
}

// Trace logs a trace message
func (pl *PrintLogger) Trace(args ...interface{}) {
	pl.Printf(TraceLevel, "TRACE: %v", args...)
}

// Printf implements the Printf method of the Logger interface
func (pl *PrintLogger) Printf(level LogLevel, format string, v ...interface{}) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	if level >= InfoLevel {
		pl.Logger.Printf(format, v...)
	}
}

// SetOutput sets the logger output destination
func (pl *PrintLogger) SetOutput(output io.Writer) {
	pl.Logger.SetOutput(output)
}

// SetFlags sets the logger output flags
func (pl *PrintLogger) SetFlags(flags int) {
	pl.Logger.SetFlags(flags)
}

// SetPrefix sets the logger output prefix
func (pl *PrintLogger) SetPrefix(prefix string) {
	pl.Logger.SetPrefix(prefix)
}

// LogWithFields logs a message with additional fields
func (pl *PrintLogger) LogWithFields(level LogLevel, fields map[string]interface{}, format string, v ...interface{}) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	if level >= InfoLevel {
		// Add custom fields to the log entry
		for key, value := range fields {
			format += " %s=%v"
			v = append(v, key, value)
		}

		msg := fmt.Sprintf(format, v...)
		pl.Logger.Printf(msg)
	}
}

// LogWithStackTrace logs an error or fatal message with error stack trace
func (pl *PrintLogger) LogWithStackTrace(level LogLevel, err error, format string, v ...interface{}) {
	if level == ErrorLevel || level == FatalLevel {
		pl.Printf(level, format, v...)
		pl.Printf(level, "Error Stack Trace:\n%s", getStackTrace(err))
	}
}

// getStackTrace returns the stack trace of an error
func getStackTrace(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("%+v", err)
}

// NewAsyncLogger creates a new asynchronous logger
func NewAsyncLogger(output io.Writer, prefix string, flags int, bufferSize int) Logger {
	logger := NewLogger(output, prefix, flags)
	asyncLogger := &AsyncLogger{
		Logger: logger,
		ch:     make(chan logEntry, bufferSize),
	}

	go asyncLogger.processEntries()
	return asyncLogger
}

// AsyncLogger is an asynchronous logger implementation
type AsyncLogger struct {
	Logger
	ch chan logEntry
}

type logEntry struct {
	level  LogLevel
	format string
	v      []interface{}
}

// Printf implements the Printf method of the Logger interface for asynchronous logger
func (al *AsyncLogger) Printf(level LogLevel, format string, v ...interface{}) {
	al.ch <- logEntry{level, format, v}
}

func (al *AsyncLogger) processEntries() {
	for entry := range al.ch {
		al.Logger.Printf(entry.level, entry.format, entry.v...)
	}
}

// Close stops the asynchronous logger and closes the channel
func (al *AsyncLogger) Close() {
	close(al.ch)
}

// JSONLogger is an implementation of Logger with JSON format
type JSONLogger struct {
	Logger
	mu sync.Mutex
}

// NewJSONLogger creates a new instance of JSONLogger
func NewJSONLogger(output io.Writer, prefix string, flags int) Logger {
	logger := NewLogger(output, prefix, flags)
	return &JSONLogger{Logger: logger}
}

// Printf implements the Printf method of the Logger interface with JSON format
func (jl *JSONLogger) Printf(level LogLevel, format string, v ...interface{}) {
	jl.mu.Lock()
	defer jl.mu.Unlock()

	entry := make(map[string]interface{})
	entry["timestamp"] = time.Now().Format(time.RFC3339)
	entry["severity"] = level.String()

	// Add custom fields to the log entry
	for i := 0; i < len(v); i += 2 {
		key, ok := v[i].(string)
		if !ok {
			continue
		}
		entry[key] = v[i+1]
	}

	entry["message"] = fmt.Sprintf(format, v...)
	logEntry, err := json.Marshal(entry)
	if err != nil {
		jl.Logger.Error("Error marshalling log entry:", err)
		return
	}

	// Corrected line - use jl.Logger.Printf with level and formatted log entry
	jl.Logger.Printf(level, string(logEntry))
}

// LogLevelString returns the string representation of a log level
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
	case TraceLevel:
		return "TRACE"
	default:
		return "UNKNOWN"
	}
}

// LogRotate performs log rotation based on log file size
func LogRotate(filePath string, maxSizeMB int, maxBackups int) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("error getting file info: %v", err)
	}

	if fileInfo.Size() < int64(maxSizeMB*1024*1024) {
		return nil
	}

	for i := maxBackups - 1; i > 0; i-- {
		src := fmt.Sprintf("%s.%d", filePath, i)
		dst := fmt.Sprintf("%s.%d", filePath, i+1)
		if err := os.Rename(src, dst); err != nil {
			return fmt.Errorf("error renaming log file: %v", err)
		}
	}

	dst := fmt.Sprintf("%s.1", filePath)
	if err := os.Rename(filePath, dst); err != nil {
		return fmt.Errorf("error renaming log file: %v", err)
	}

	return nil
}

// NewLoggerWithConfig creates a new instance of Logger with configuration
func NewLoggerWithConfig(config LoggerConfig) Logger {
	logger := NewLogger(config.LogOutput, "", log.LstdFlags)
	logger.SetFlags(log.LstdFlags)

	if config.LogFormat == JSONFormat {
		logger = &JSONLogger{Logger: logger}
	}

	if config.RotatePeriod > 0 {
		go func() {
			for range time.Tick(config.RotatePeriod) {
				if err := LogRotate("app.log", config.MaxSizeMB, config.MaxBackups); err != nil {
					logger.Error("Error during log rotation:", err)
				}
			}
		}()
	}

	return logger
}

// EventPublisher is an implementation of the EventPublisher interface using fmt.Printf
type EventPublisherImpl struct {
	EventLog *EventLog
}

// NewLogger creates a new instance of Logger
func NewEventPublisher() EventPublisherImpl {
	// Perform any setup or initialization
	return EventPublisherImpl{}
}

// Info logs an informational message
func (l *EventPublisherImpl) Publish(event ...interface{}) error {
	if event != nil {
		return fmt.Errorf("event publish failed")
	}
	fmt.Print("INFO: ")
	fmt.Println(event...)
	return nil
}
