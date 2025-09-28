package app

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/mnhkahn/gogogo/logger"
	"github.com/mnhkahn/gogogo/panicer"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// Handler ...
type Handler struct {
	TimeOut      time.Duration
	Metrics      *newrelic.Application
	RecoverFunc  func(c *Context)             `json:"-"`
	ErrorMsgFunc func(c *Context, msg string) `json:"-"`
	Stats        map[string]*Stat
	statLock     sync.RWMutex
}

// NewHandler ...
func NewHandler() *Handler {
	return &Handler{
		RecoverFunc: func(c *Context) {
			panicer.RecoverHandler(c.ResponseWriter, c.Request)
		},
		ErrorMsgFunc: func(c *Context, msg string) {
			logger.Warn("ErrorMsgFunc", msg)
		},
		Stats: make(map[string]*Stat),
	}
}

func (h *Handler) Cost(u string, statusCode int, d time.Duration) {
	h.statLock.Lock()
	defer h.statLock.Unlock()

	key := fmt.Sprintf("[%s]%s", http.StatusText(statusCode), u)
	if stat, e := h.Stats[key]; e {
		stat.Cnt++
		stat.StatusCode = statusCode
		stat.SumTime += time.Duration(d)
		stat.AvgTime = stat.SumTime / time.Duration(stat.Cnt)
	} else {
		h.Stats[key] = &Stat{u, statusCode, 1, time.Duration(d), time.Duration(d)}
	}
}

// DefaultHandler ...
var DefaultHandler = NewHandler()

// SetTimeout ...
func SetTimeout(d time.Duration) {
	DefaultHandler.TimeOut = d
}

// SetRecoverFunc ...
func SetRecoverFunc(recoverFunc func(c *Context)) {
	DefaultHandler.RecoverFunc = recoverFunc
}

func SetErrorMsgFunc(fn func(c *Context, msg string)) {
	DefaultHandler.ErrorMsgFunc = fn
}

// Got ...
type Got struct {
	H func(c *Context) error
}

// ServeHTTP ...
func (h Got) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := AllocContext(w, r)
	defer FreeContext(c)

	// https://golang.dbwu.tech/traps/defer_with_recover/#:~:text=错误的原因在于%3A%20defer%20以匿名函数的方式运行，本身就等于包装了一层函数，%20内部的%20myRecover%20函数包装了%20recover%20函数，等于又加了一层包装，变成了两,panic%20就无法被捕获了%E3%80%82%20defer%20直接调用%20myRecover%20函数，这样减去了一层包装，%20panic%20就可以被捕获了%E3%80%82
	if DefaultHandler.RecoverFunc != nil {
		defer DefaultHandler.RecoverFunc(c)
	}

	if String("seo") == "true" {
		if seo(c) {
			return
		}
	}

	var err error
	c.Params, err = url.ParseQuery(c.Request.URL.RawQuery)
	if err != nil {
		c.Error(err.Error())
		return
	}

	err = h.H(c)
	if err != nil {
		c.Error(err.Error())
		return
	}
}

type basicAuthHandler struct {
	handler http.Handler
	user    string
	pwd     string
}

// ServeHTTP ...
func (h basicAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u, p, ok := r.BasicAuth()
	if !ok || !(u == h.user && p == h.pwd) {
		w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	h.handler.ServeHTTP(w, r)
}

// BasicAuthHandler ...
func BasicAuthHandler(handler http.Handler, user, pwd string) http.Handler {
	return basicAuthHandler{handler: handler, user: user, pwd: pwd}
}

var AuthHandler = func(handler http.Handler) http.Handler {
	usr, pwd, _ := GetConfigAuth()
	return BasicAuthHandler(handler, usr, pwd)
}
