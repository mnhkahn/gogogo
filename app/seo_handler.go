package app

import (
	"net/http"
)

func seo(c *Context) bool {
	if c.Request.Host == "localhost" {
		return false
	}
	if c.Request.Header.Get("X-Forwarded-Proto") == "http" {
		http.Redirect(c.ResponseWriter, c.Request, "https://"+c.Request.Host+c.Request.RequestURI, http.StatusMovedPermanently)
		return true
	}
	host := String("host")
	if c.Request.Host == String("redirect_host") && host != "" {
		http.Redirect(c.ResponseWriter, c.Request, "https://"+host+c.Request.RequestURI, http.StatusMovedPermanently)
		return true
	}
	return false
}
