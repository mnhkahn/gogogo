package app

import (
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/mnhkahn/gogogo/logger"
	"github.com/mnhkahn/gogogo/panicer"
)

// Handler ...
type Handler struct {
	TimeOut     time.Duration
	RecoverFunc func(w http.ResponseWriter, r *http.Request) `json:"-"`
	Stats       map[string]*Stat
	statLock    sync.RWMutex
}

// NewHandler ...
func NewHandler() *Handler {
	return &Handler{
		RecoverFunc: panicer.RecoverHandler,
		Stats:       make(map[string]*Stat),
	}
}

func (h *Handler) Cost(u string, st time.Time) {
	h.statLock.Lock()
	defer h.statLock.Unlock()

	d := time.Now().Sub(st).Nanoseconds()
	if stat, e := h.Stats[u]; e {
		stat.Cnt++
		stat.SumTime += time.Duration(d)
		stat.AvgTime = stat.SumTime / time.Duration(stat.Cnt)
	} else {
		h.Stats[u] = &Stat{1, time.Duration(d), time.Duration(d)}
	}
}

// DefaultHandler ...
var DefaultHandler = NewHandler()

// SetTimeout ...
func SetTimeout(d time.Duration) {
	DefaultHandler.TimeOut = d
}

// SetRecoverFunc ...
func SetRecoverFunc(r func(w http.ResponseWriter, r *http.Request)) {
	DefaultHandler.RecoverFunc = r
}

// Got ...
type Got struct {
	H func(c *Context) error
}

// ServeHTTP ...
func (h Got) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	st := time.Now()

	if DefaultHandler.RecoverFunc != nil {
		defer DefaultHandler.RecoverFunc(w, r)
	}

	c := AllocContext(w, r)
	defer FreeContext(c)

	if String("seo") == "true" {
		if seo(c) {
			return
		}
	}

	var err error
	c.Params, err = url.ParseQuery(c.Request.URL.RawQuery)
	if err != nil {
		c.Error(err.Error())
		go DefaultHandler.Cost(c.Request.URL.Path, st)
		return
	}

	err = h.H(c)
	logger.Infof("%s %s", c.Request.Method, c.Request.URL.String())
	if err != nil {
		c.Error(err.Error())
		go DefaultHandler.Cost(c.Request.URL.Path, st)
		return
	}

	go DefaultHandler.Cost(c.Request.URL.Path, st)
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
