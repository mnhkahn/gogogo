package panicer

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"time"
)

func RecoverHandlerWithFunc(w http.ResponseWriter, r *http.Request, fn func(string)) {
	if err := recover(); err != nil {
		res := bytes.NewBuffer(nil)
		rec := fmt.Sprintf("Recover %s %v %s\n", time.Now().Format(time.RFC3339), err, r.URL.String())
		_, _ = fmt.Fprintf(os.Stderr, "request: %s", string(dump))
		stack := printStack()

		w.WriteHeader(http.StatusServiceUnavailable)

		res.WriteString(rec)
		res.WriteString(stack)

		dump, _ := httputil.DumpRequest(r, true)
		_, _ = res.Write(dump)

		fn(res.String())
	}
}

// RecoverHandler help to recover in the handler.
func RecoverHandler(w http.ResponseWriter, r *http.Request) {
	if err := recover(); err != nil {
		fmt.Fprintf(os.Stderr, "Recover %s %v %s\n", time.Now().Format(time.RFC3339), err, r.URL.String())
		printStack()

		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

// RecoverDebug help to recover and print a debug info to stderr.
func RecoverDebug(v interface{}) {
	if err := recover(); err != nil {
		fmt.Fprintf(os.Stderr, "Recover: %s %v %+v\n", time.Now().Format(time.RFC3339), err, v)
		printStack()
	}
}

// Recover recover panic and print time info to stderr.
func Recover() {
	if err := recover(); err != nil {
		fmt.Fprintf(os.Stderr, "Recover: %s %v\n", time.Now().Format(time.RFC3339), err)
		printStack()
	}
}

// print panic's stack.
func printStack() string {
	res := bytes.NewBuffer(nil)
	for i := 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		d := fmt.Sprintf("%d %s:%d\n", pc, file, line)
		fmt.Fprintf(os.Stderr, d)
		res.WriteString(d)
	}
	return res.String()
}
