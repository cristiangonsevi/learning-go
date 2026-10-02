// Package logger provides simple leveled logging helpers (Info, Error, Debug)
// for printing formatted messages to standard output.
package logger

import (
	"fmt"
	"strings"
)

func log(format string, a ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format = format + "\n"
	}
	fmt.Printf(format, a...)
}

func Info(format string, args ...interface{}) {
	log("[INFO] "+format, args...)
}

func Error(format string, args ...interface{}) {
	log("[ERROR] "+format, args...)
}

func Debug(format string, args ...interface{}) {
	// if !verbose {
	// 	return
	// }
	log("[DEBUG] "+format, args...)
}
