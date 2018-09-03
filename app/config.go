package app

import (
	"errors"

	"github.com/sasbury/mini"
)

var (
	appconfig = new(mini.Config)

	ErrConfNotInit = errors.New("appconfig is not init, forgot to call InitAppConf()? ")
)

// InitAppConf ...
func InitAppConf() error {
	err := appconfig.InitializeFromPath("./conf/app.conf")
	if err != nil { // Handle errors reading the config file
		return err
	}

	return nil
}

// Int ...
func Int(key string) int {
	return int(appconfig.Integer(key, 0))
}

// String ...
func String(key string) string {
	return appconfig.String(key, "")
}

// DefaultString ...
func DefaultString(key, def string) string {
	return appconfig.String(key, def)
}

// Strings ...
func Strings(key string) []string {
	return appconfig.Strings(key)
}

// GetConfigAuth ...
func GetConfigAuth() (string, string, error) {
	if appconfig == nil {
		return "", "", ErrConfNotInit
	}

	user := String("auth.user")
	pwd := String("auth.pwd")

	if user == "" {
		return "jia", "jia", nil
	}
	return user, pwd, nil
}
