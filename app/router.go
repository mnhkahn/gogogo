package app

// InitRouter ...
func InitRouter() {
	Handle("/debug/goapp", AuthHandler(&Got{H: GoAppHandler}))
	Handle("/debug/log/level", AuthHandler(Got{LogLevelHandler}))
	Handle("/debug/stat", AuthHandler(Got{StatHandler}))
}
