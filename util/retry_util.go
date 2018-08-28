// Package util
package util

import (
	"github.com/mnhkahn/gogogo/logger"
	"time"
)

// Retry function. fn is the retry function.
func Retry(attempts int, sleep time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		if s, ok := err.(stop); ok {
			return s.error
		}

		if attempts--; attempts > 0 {
			logger.Warnf("retry func error: %s. attemps #%d after %s.", err.Error(), attempts, sleep)
			time.Sleep(sleep)
			return Retry(attempts, 2*sleep, fn)
		}
		return err
	}
	return nil
}

type stop struct {
	error
}

func NoRetryError(err error) stop {
	return stop{err}
}
