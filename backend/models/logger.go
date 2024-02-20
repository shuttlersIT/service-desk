package models

import (
	"fmt"
)

type Logger interface {
	Error(args ...interface{})
	Info(args ...interface{})
	Printf(format string, v ...interface{})
}

// PrintLogger is an implementation of the Logger interface using fmt.Printf
type PrintLogger struct{}

// NewLogger creates a new instance of Logger
func NewLogger() *PrintLogger {
	// Perform any setup or initialization
	return &PrintLogger{}
}

// Printf implements the Printf method of the Logger interface
func (pl *PrintLogger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

// Error logs an error message
func (l *PrintLogger) Error(args ...interface{}) {
	fmt.Print("ERROR: ")
	fmt.Println(args...)
}

// Info logs an informational message
func (l *PrintLogger) Info(args ...interface{}) {
	fmt.Print("INFO: ")
	fmt.Println(args...)
}

// Info logs an informational message
func (l *PrintLogger) Fatal(args ...interface{}) {
	fmt.Print("FATAL: ")
	fmt.Println(args...)
}

// Assuming an EventPublisher interface is defined elsewhere in your application
type EventPublisher interface {
	Publish(event interface{}) error
}

// EventPublisher is an implementation of the EventPublisher interface using fmt.Printf
type EventPublisherImpl struct {
	EventLog *EventLog
}

// NewLogger creates a new instance of Logger
func NewEventPublisher() *EventPublisherImpl {
	// Perform any setup or initialization
	return &EventPublisherImpl{}
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
