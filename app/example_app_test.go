// Package app_test

package app_test

import (
	"github.com/mnhkahn/gogogo/app"
	"github.com/mnhkahn/gogogo/logger"
	"net"
)

func Index(c *app.Context) error {
	c.WriteString("hello, world")
	return nil
}

func Example() {
	app.Handle("/", &app.Got{Index})

	l, err := net.Listen("tcp", ":1031")
	if err != nil {
		logger.Errorf("Listen: %v", err)
		return
	}
	app.Serve(l)
}
