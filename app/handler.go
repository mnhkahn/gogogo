package app

import (
	"net/http"
	"net/url"
	"time"

	"github.com/mnhkahn/gogogo/logger"
	"github.com/mnhkahn/gogogo/panicer"
)

type Handler struct {
	TimeOut     time.Duration
	RecoverFunc func(w http.ResponseWriter, r *http.Request) `json:"-"`
}

var DefaultHandler = &Handler{RecoverFunc: panicer.RecoverHandler}

func SetTimeout(d time.Duration) {
	DefaultHandler.TimeOut = d
}

func SetRecoverFunc(r func(w http.ResponseWriter, r *http.Request)) {
	DefaultHandler.RecoverFunc = r
}

type Got struct {
	H func(c *Context) error
}

func (h Got) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if DefaultHandler.RecoverFunc != nil {
		defer DefaultHandler.RecoverFunc(w, r)
	}

	c := AllocContext(w, r)
	defer FreeContext(c)

	var err error
	c.Params, err = url.ParseQuery(c.Request.URL.RawQuery)
	if err != nil {
		c.Error(err.Error())
		return
	}

	err = h.H(c)
	logger.Infof("%s %s", c.Request.Method, c.Request.URL.String())
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

func (h basicAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u, p, ok := r.BasicAuth()
	if !ok || !(u == h.user && p == h.pwd) {
		w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	h.handler.ServeHTTP(w, r)
}

func BacisAuthHandler(handler http.Handler, user, pwd string) http.Handler {
	return basicAuthHandler{handler: handler, user: user, pwd: pwd}
}
