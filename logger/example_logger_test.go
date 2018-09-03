// Package logger_test

package logger_test

import (
	"github.com/mnhkahn/gogogo/logger"
	"log"
	"os"
)

func Example() {
	logger.Info("hello, world.")
}

func ExampleNewLogger() {
	w := os.Stdout
	flag := log.Llongfile
	l := logger.NewWriterLogger(w, flag, 3)
	l.Info("hello, world")
}
