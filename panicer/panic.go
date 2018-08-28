// Package panicer
package panicer

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"
)

func RecoverHandler(w http.ResponseWriter, r *http.Request) {
	if err := recover(); err != nil {
		fmt.Fprintf(os.Stderr, "Recover %s %v %s\n", time.Now().Format(time.RFC3339), err, r.URL.String())
		printStack()

		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

// RecoverDebug 支持打印一个调试信息
func RecoverDebug(v interface{}) {
	if err := recover(); err != nil {
		fmt.Fprintf(os.Stderr, "Recover: %s %v %+v\n", time.Now().Format(time.RFC3339), err, v)
		printStack()
	}
}

// Recover ...
func Recover() {
	if err := recover(); err != nil {
		fmt.Fprintf(os.Stderr, "Recover: %s %v\n", time.Now().Format(time.RFC3339), err)
		printStack()
	}
}

func printStack() {
	for i := 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fmt.Fprintf(os.Stderr, "%d %s:%d\n", pc, file, line)
	}
}
