// Package app
package app

import "time"

type Stat struct {
	Cnt     int64
	SumTime time.Duration
	AvgTime time.Duration
}
