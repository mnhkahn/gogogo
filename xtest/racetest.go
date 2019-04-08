// Package xtest
package xtest

import (
	"sync"
	"testing"
)

func RaceTestCounter(t *testing.T, count int, f func(i int) interface{}) {
	counterMap := sync.Map{}

	RaceTest(t, count, func(i int) {
		key := f(i)
		a, _ := counterMap.Load(key)
		if a == nil {
			counterMap.Store(key, 1)
		} else {
			counterMap.Store(key, a.(int)+1)
		}
	})

	_count := 0
	counterMap.Range(func(key, value interface{}) bool {
		_count++
		return true
	})
	if _count != count {
		t.Errorf("race test error, want len %d, but got %d", count, _count)
		counterMap.Range(func(key, value interface{}) bool {
			t.Log("key =", key, "value =", value)
			return true
		})
	}
}

func RaceTest(t *testing.T, count int, f func(i int)) {
	if count <= 0 {
		count = 10000
		t.Logf("count use %d instead\n", count)
	}
	var wait sync.WaitGroup

	defer func() {
		if err := recover(); err != nil {
			t.Error(err)
		}
	}()

	for i := 0; i < count; i++ {
		wait.Add(1)
		go func() {
			f(i)
			wait.Done()
		}()
	}

	wait.Wait()
}
