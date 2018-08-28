// Package app
package app

import (
	"runtime"

	"github.com/mnhkahn/gogogo/logger"
)

var (
	VERSION = "N/A"
	BRANCH  = "N/A"
	BUILD   = "N/A"
)

func AppInfo() {
	logger.Info("Git Version:", VERSION)
	logger.Info("Build Branch:", BRANCH)
	logger.Info("Build Date:", BUILD)
	logger.Info("Go Version:", runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
