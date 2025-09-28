// Package app
package app

import (
	"time"
)

type Stat struct {
	Url        string
	StatusCode int
	Cnt        int64
	SumTime    time.Duration
	AvgTime    time.Duration
}
