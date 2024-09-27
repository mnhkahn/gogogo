package app

// InitRouter ...
func InitRouter() {
	usr, pwd, _ := GetConfigAuth()
	Handle("/debug/goapp", BasicAuthHandler(&Got{H: GoAppHandler}, usr, pwd))
	Handle("/debug/router", BasicAuthHandler(&Got{H: DebugRouter}, usr, pwd))
	Handle("/debug/log/level", BasicAuthHandler(Got{LogLevelHandler}, usr, pwd))
	Handle("/debug/stat", BasicAuthHandler(Got{StatHandler}, usr, pwd))
	// Handle("/sitemap.xml", Got{SiteMapXML})
	Handle("/sitemap", Got{SiteMapRaw})
}
