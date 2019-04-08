// Package xtest
package xtest

import (
	"strings"
	"testing"

	"github.com/sony/sonyflake"
	"github.com/stretchr/testify/assert"
)

func TestRaceTest(t *testing.T) {
	var b strings.Builder
	RaceTest(t, 10000, func(i int) {
		b.WriteString("1")
	})
	assert.NotEqual(t, 10000, b.Len())
}

func TestRaceTestMap(t *testing.T) {
	m := make(map[string]string, 0)
	go func() {
		for {
			_ = m["a"]
		}
	}()
	RaceTest(t, 10000, func(i int) {
		m["a"] = "a"
	})
}

func TestRaceTestCounter(t *testing.T) {
	sf := sonyflake.NewSonyflake(sonyflake.Settings{
		//StartTime: time.Now(),
		MachineID: func() (u uint16, e error) {
			return 0, nil
		},
	})
	sf2 := sonyflake.NewSonyflake(sonyflake.Settings{
		//StartTime: time.Now(),
		MachineID: func() (u uint16, e error) {
			return 0, nil
		},
	})

	RaceTestCounter(t, 1000, func(i int) interface{} {
		var unique_id uint64
		if i%2 == 0 {
			unique_id, _ = sf.NextID()
			return unique_id
		} else {
			unique_id, _ = sf2.NextID()
		}
		return unique_id
	})
}
