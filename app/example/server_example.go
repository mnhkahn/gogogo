// Package example
package main

import (
	"net"

	"github.com/mnhkahn/gogogo/app"
	"github.com/mnhkahn/gogogo/logger"
)

func main() {
	l, err := net.Listen("tcp", ":1031")
	if err != nil {
		logger.Errorf("Listen: %v", err)
		return
	}
	app.Serve(l)
}
