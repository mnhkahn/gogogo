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
	stat := DefaultHandler.Stats["[OK]a"]
	if assert.NotNil(t, stat) {
		assert.EqualValues(t, 10, stat.Cnt)
	}
	//assert.EqualValues(t, time.Duration(50*time.Millisecond).String(), time.Duration(DefaultHandler.Stats["a"].AvgTime).String())
}
