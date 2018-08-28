// Package app
package app

import (
	"os"
	"reflect"
	"runtime"
	"sort"

	"github.com/mnhkahn/gogogo/logger"
	"github.com/mnhkahn/gogogo/util"
)

// GoAppHandler is a handler to show debug information.
func GoAppHandler(c *Context) error {
	res := make(map[string]interface{}, 4)
	res["app.VERSION"] = VERSION
	res["app.BUILD"] = BUILD
	res["app.BRANCH"] = BRANCH
	res["go.version"] = runtime.Version()
	res["os.args"] = os.Args
	res["os.env"] = os.Environ()
	res["os.goroutine"] = util.Goroutine()
	res["os.pwd"], _ = os.Getwd()
	res["log.level"] = logger.StdLogger.GetLevel()
	c.JSON(res)

	return nil
}

// DebugRouter is a handler to show all routers.
func DebugRouter(c *Context) error {
	v := reflect.ValueOf(GoEngine.mux)
	m := v.Elem().FieldByName("m")

	keys := m.MapKeys()

	routers := make([]string, 0, len(keys))
	for _, key := range keys {
		routers = append(routers, key.String())
	}

	sort.Strings(routers)
	c.JSON(routers)

	return nil
}

// LogLevelHandler is a handler to set log level for StdLogger.
func LogLevelHandler(c *Context) error {
	l, err := c.GetInt("level")
	if err != nil {
		return err
	}
	res := logger.StdLogger.SetLevel(l)
	c.JSON(res)

	return nil
}
