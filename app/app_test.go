// Package app
package app

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStat(t *testing.T) {
	for i := 0; i < 10; i++ {
		DefaultHandler.Cost("a", http.StatusOK, 50*time.Millisecond)
	}
	assert.EqualValues(t, 10, DefaultHandler.Stats["a"].Cnt)
	//assert.EqualValues(t, time.Duration(50*time.Millisecond).String(), time.Duration(DefaultHandler.Stats["a"].AvgTime).String())
}
