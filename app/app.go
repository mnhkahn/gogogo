package app

import (
	"net"
	"net/http"
	"sort"
	"time"

	"golang.org/x/net/netutil"

	"github.com/mnhkahn/gogogo/logger"
)

// Engine ...
type Engine struct {
	mux      *http.ServeMux
	patterns []string
	server   *http.Server
	l        net.Listener
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
	e.mux.Handle(pattern, h)
	e.patterns = append(e.patterns, pattern)
	sort.Slice(e.patterns, func(i, j int) bool {
		return e.patterns[i] < e.patterns[j]
	})
}

func (e *Engine) Patterns() []string {
	return e.patterns
}

// Default predefined engine
var GoEngine = NewEngine()

// Handle ...
func Handle(pattern string, h http.Handler) {
	GoEngine.Handle(pattern, h)
}

// Serve ...
func Serve(l net.Listener) {
	InitRouter()
	GoEngine.Serve(l)
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
