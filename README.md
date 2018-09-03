# gogogo

### app

A simple web framework compatible with `net/http`.

It also has a predefined 'standard' app engine server called GoEngine accessible through helper functions Handle, Server and so on.

```
func A(c *Context) error {
	c.JSON([]byte("a"))
	return nil
}

func main() {
    app.Handle("/", &app.Got{A})
    
    l, err := net.Listen("tcp", ":"+port)
    if err != nil {
        logger.Errorf("Listen: %v", err)
        return
    }
    app.Serve(l)
}
```

A example site: [Cyeam](http://cyeam.com).

### logger

A logger package that compatible with sdk log package. It supports level logging and file auto rotating.

It also has a predefined 'standard' Logger called StdLogger accessible through helper functions Debug[f], Info[f], Warn[f], Error[f], LogLevel and SetJack.

It supports 4 level:

```
LevelError = iota
LevelWarning
LevelInformational
LevelDebug
```

You can use LogLevel to handle the log level.

File rotating based on package gopkg.in/natefinch/lumberjack.v2, you can control file settings by using SetJack.

Quick start

```
import "github.com/mnhkahn/gogogo/logger"

logger.Info("hello, world.")
```

Defined our own logger:

```
l := logger.NewWriterLogger(w, flag, 3)
l.Info("hello, world")
```

For more information, goto [godoc](https://godoc.org/github.com/mnhkahn/gogogo/logger).

Chinese details, goto [link](http://blog.cyeam.com//golang/2017/07/14/go-log?utm_source=github).

### panicer

A help package that catch panic error and print the panic stack.

Quick start

	import "github.com/mnhkahn/gogogo/panicer"

	func main() {
		defer Recover()
	}

For more information, goto godoc https://godoc.org/github.com/mnhkahn/gogogo/panicer


### util

A util package collection.

For more information, goto godoc https://godoc.org/github.com/mnhkahn/gogogo/util

Chinese details, goto http://blog.cyeam.com//golang/2018/08/27/retry?utm_source=github
