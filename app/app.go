package app

import (
	"net"
	"net/http"
	"time"

	"github.com/mnhkahn/gogogo/logger"
	"github.com/newrelic/go-agent/v3/newrelic"
	"golang.org/x/net/netutil"
)

// Engine ...
type Engine struct {
	mux    *http.ServeMux
	server *http.Server
	l      net.Listener
}

// NewEngine ...
func NewEngine() *Engine {
	e := new(Engine)
	e.mux = http.NewServeMux()

	if err := InitAppConf(); err != nil {
		// panic(err)
		logger.Warn(err)
	}
	return e
}

// NewEngine2 ...
func NewEngine2() *Engine {
	e := new(Engine)
	e.mux = http.NewServeMux()

	return e
}

// Serve ...
func (e *Engine) Serve(l net.Listener) {
	e.l = l
	e.server = &http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      e.mux,
	}
	logger.Infof("Listening and serving HTTP on %s", e.l.Addr().String())
	logger.Error(e.server.Serve(e.l))
}

// ServeDefault ...
func (e *Engine) ServeDefault(l net.Listener) {
	e.l = l
	e.server = &http.Server{
		Handler: e.mux,
	}
	logger.Infof("Listening and serving HTTP on %s", e.l.Addr().String())
	logger.Error(e.server.Serve(e.l))
}

// ServeMux ...
func (e *Engine) ServeMux(l net.Listener, handler http.Handler) {
	e.l = l
	e.server = &http.Server{
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		Handler:      handler,
	}
	logger.Infof("Listening and serving HTTP on %s", e.l.Addr().String())
	logger.Error(e.server.Serve(e.l))
}

// Handle ...
func (e *Engine) Handle(pattern string, h http.Handler) {
	if DefaultHandler.TimeOut > 0 {
		h = http.TimeoutHandler(h, DefaultHandler.TimeOut, "")
	}
	if DefaultHandler.Metrics != nil {
		p, q := newrelic.WrapHandle(DefaultHandler.Metrics, pattern, h)
		e.mux.Handle(p, q)
	} else {
		e.mux.Handle(pattern, h)
	}
}

// Default predefined engine
var GoEngine = NewEngine()

// Handle ...
func Handle(pattern string, h http.Handler) {
	GoEngine.Handle(pattern, h)
}

func HandleGot(pattern string, h func(c *Context) error) {
	GoEngine.Handle(pattern, &Got{h})
}

// Serve ...
func Serve(l net.Listener) {
	InitRouter()
	GoEngine.Serve(l)
}

// ServeDefault ...
func ServeDefault(l net.Listener) {
	InitRouter()
	GoEngine.ServeDefault(l)
}

// ServeMux ...
func ServeMux(l net.Listener, handler http.Handler) {
	GoEngine.ServeMux(l, handler)
}

// DefaultLimit ...
const DefaultLimit = 20

// LimitServe ...
func LimitServe(limit int) {
	if limit <= 0 {
		limit = DefaultLimit
	}
	logger.Info("start limit server with limit:", limit)

	l, err := net.Listen("tcp", ":"+appconfig.String("port", "1031"))
	if err != nil {
		logger.Errorf("Listen: %v", err)
		return
	}
	// defer l.Close()

	l = netutil.LimitListener(l, limit)

	Serve(l)
}
