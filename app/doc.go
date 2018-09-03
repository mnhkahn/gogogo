/*
This package is a web framework that compatible with sdk http package.
It also has a predefined 'standard' app engine server called GoEngine accessible through helper functions Handle, Server and so on.

Quick start:

	import "github.com/mnhkahn/gogogo/app"
	import "github.com/mnhkahn/gogogo/logger"

	func Index(c *app.Context) error {
		c.WriteString("hello, world")
		return nil
	}

	app.Handle("/", &app.Got{Index})

	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Errorf("Listen: %v", err)
		return
	}
	app.Serve(l)


For more information, goto godoc https://godoc.org/github.com/mnhkahn/gogogo/app.

A example site: [Cyeam](http://cyeam.com).

 */
// Package logger
package app
