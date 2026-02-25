package main

type Logger interface {
	LLog(message string)
}

type ConsoleLogger struct{}

func (cl ConsoleLogger) LLog(message string) {
	println("Console log:", message)
}

func Process(l Logger, msg string) {
	l.LLog("procesando " + msg)
}

func main() {
	logger := ConsoleLogger{}
	Process(logger, "mensaje de prueba")
}
