package parser

import (
	"fmt"
)

type Logger interface {
	Debug(message string, args ...any)
}

type DefaultLogger struct {
	Logger
}

func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{}
}

func (l *DefaultLogger) Debug(message string, args ...any) {
	fmt.Printf(message+"\n", args...)
}

type NullLogger struct {
	Logger
}

func NewNullLogger() *NullLogger {
	return &NullLogger{}
}

func (l *NullLogger) Debug(message string, args ...any) {
	// Do nothing
}
