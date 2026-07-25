package logger

import "fmt"

func Info(args ...interface{}) {
	fmt.Println(args...)
}
