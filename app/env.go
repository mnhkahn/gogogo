package app

import (
	"runtime"

	"github.com/mnhkahn/gogogo/logger"
)

// go build version, use build flags to assign value. -ldflags='-X "github.com/mnhkahn/gogogo/app.VERSION=$(VERSION)" -X "github.com/mnhkahn/gogogo/app.BRANCH=$(branch)" -X "github.com/mnhkahn/gogogo/app.BUILD=$(DATE)"'
var (
	VERSION = "N/A"
	BRANCH  = "N/A"
	BUILD   = "N/A"
)

// AppInfo prints app info with git version, git branch, build date and go version.
func AppInfo() {
	logger.Info("Git Version:", VERSION)
	logger.Info("Build Branch:", BRANCH)
	logger.Info("Build Date:", BUILD)
	logger.Info("Go Version:", runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
