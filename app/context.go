package app

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/Masterminds/sprig"
	"github.com/mnhkahn/gogogo/logger"
)

// Context ...
type Context struct {
	Params         url.Values
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

// NewContext ...
func NewContext(rw http.ResponseWriter, req *http.Request) *Context {
	c := new(Context)
	c.ResponseWriter = rw
	c.Request = req

	return c
}

var ctxPool = sync.Pool{
	New: func() interface{} {
		return new(Context)
	},
}

// AllocContext ...
func AllocContext(rw http.ResponseWriter, req *http.Request) *Context {
	c := ctxPool.Get().(*Context)
	c.ResponseWriter = rw
	c.Request = req

	return c
}

// FreeContext ...
func FreeContext(c *Context) {
	ctxPool.Put(c)
}

// URL ...
func (c *Context) URL() string {
	return c.Request.URL.String()
}

// Query ...
func (c *Context) Query() url.Values {
	return c.Params
}

func (c *Context) get(key string) string {
	return c.Query().Get(key)
}

func (c *Context) gets(key string) []string {
	return c.Query()[key]
}

// GetBool ...
func (c *Context) GetBool(key string) (bool, error) {
	return strconv.ParseBool(c.get(key))
}

// GetInt ...
func (c *Context) GetInt(key string, def ...int) (int, error) {
	strv := c.get(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	return strconv.Atoi(strv)
}

// GetUInt8 ...
func (c *Context) GetUInt8(key string, def ...uint8) (uint8, error) {
	strv := c.get(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	ui, err := strconv.ParseUint(strv, 10, 0)
	return uint8(ui), err
}

// GetUInt ...
func (c *Context) GetUInt(key string, def ...uint) (uint, error) {
	strv := c.get(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	ui, err := strconv.ParseUint(strv, 10, 0)
	return uint(ui), err
}

// GetUInt32 ...
func (c *Context) GetUInt32(key string, def ...uint32) (uint32, error) {
	strv := c.get(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	ui, err := strconv.ParseUint(strv, 10, 0)
	return uint32(ui), err
}

// GetUInt64 ...
func (c *Context) GetUInt64(key string, def ...uint64) (uint64, error) {
	strv := c.get(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	ui, err := strconv.ParseUint(strv, 10, 0)
	return ui, err
}

// GetInt64 ...
func (c *Context) GetInt64(key string, def ...int64) (int64, error) {
	strv := c.get(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	return strconv.ParseInt(strv, 10, 64)
}

// GetString ...
func (c *Context) GetString(key string, def ...string) string {
	if v := c.get(key); v != "" {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

// GetStrings ...
func (c *Context) GetStrings(key string) []string {
	if v := c.gets(key); len(v) > 0 {
		return v
	}
	return nil
}

// Serve ...
func (c *Context) Serve(v interface{}) {
	json, err := json.Marshal(v)
	if err != nil {
		c.Error(err.Error())
	}
	c.WriteBytes(json)
}

// WriteBytes ...
func (c *Context) WriteBytes(raw []byte) {
	c.ResponseWriter.Write(raw)
}

// WriteString ...
func (c *Context) WriteString(str string) {
	c.ResponseWriter.Write([]byte(str))
}

// JSON ...
func (c *Context) JSON(v interface{}) {
	c.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")

	content, err := json.Marshal(v)
	if err != nil {
		content = []byte(err.Error())
	}

	callback := c.GetString("callback")
	if callback == "" {
		c.ResponseWriter.Write(content)
	} else {
		callback_content := bytes.NewBufferString(" " + template.JSEscapeString(callback))
		callback_content.WriteString("(")
		callback_content.Write(content)
		callback_content.WriteString(");\r\n")
		c.ResponseWriter.Write(callback_content.Bytes())
	}
}

// HTML ...
func (c *Context) HTML(filenames []string, data interface{}) {
	tmpl := template.Must(template.ParseFiles(filenames...))

	err := tmpl.Execute(c.ResponseWriter, data)
	if err != nil {
		DefaultHandler.ErrorMsgFunc(c, err.Error())
		c.Error(err.Error())
	}
}

// HTMLFunc ...
func (c *Context) HTMLFunc(filenames []string, data interface{}, funcs template.FuncMap) {
	filename := filepath.Base(filenames[0])
	tmpl := template.Must(template.New(filename).Funcs(sprig.FuncMap()).Funcs(funcs).ParseFiles(filenames...))
	err := tmpl.ExecuteTemplate(c.ResponseWriter, filename, data)
	if err != nil {
		DefaultHandler.ErrorMsgFunc(c, err.Error())
	}
}

// Error ...
func (c *Context) Error(msg string) {
	logger.Warn("got error", msg)

	c.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.ResponseWriter.Header().Set("X-Content-Type-Options", "nosniff")
	if c.Debug() {
		c.WriteString(msg)
	}
}

// Debug ...
func (c *Context) Debug() bool {
	if c.get("debug") == "yes" {
		return true
	}
	return false
}
