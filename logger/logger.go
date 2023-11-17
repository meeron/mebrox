package logger

import (
	"fmt"
	"log"
)

var l = log.Default()

func Error(err error) {
	l.Printf("[Error] %v", err)
}

func Debug(format string, v ...any) {
	l.Printf("[Debug] %s", fmt.Sprintf(format, v...))
}

func Info(format string, v ...any) {
	l.Printf("[Info] %s", fmt.Sprintf(format, v...))
}

func Fatal(err error) {
	l.Fatal(err)
}
