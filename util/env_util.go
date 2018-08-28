// Package util
package util

import (
	"runtime"
	"runtime/pprof"
)

// Goroutine ...
func Goroutine() map[string]interface{} {
	res := map[string]interface{}{}
	res["goroutines"] = runtime.NumGoroutine()
	res["OS threads"] = pprof.Lookup("threadcreate").Count()
	res["GOMAXPROCS"] = runtime.GOMAXPROCS(0)
	res["num CPU"] = runtime.NumCPU()
	return res
}
