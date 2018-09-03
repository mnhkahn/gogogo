package app

// InitRouter ...
func InitRouter() {
	usr, pwd, _ := GetConfigAuth()
	Handle("/debug/goapp", BacisAuthHandler(&Got{H: GoAppHandler}, usr, pwd))
	Handle("/debug/router", BacisAuthHandler(&Got{H: DebugRouter}, usr, pwd))
	Handle("/debug/log/level", BacisAuthHandler(Got{LogLevelHandler}, usr, pwd))
}
