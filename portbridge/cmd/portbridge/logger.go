package main

import "fmt"

var (
	logger  Logger
	verbose bool
)

type Logger struct{}

func (l *Logger) Info(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	if !verbose {
		return
	}
	fmt.Printf(format, args...)
}
