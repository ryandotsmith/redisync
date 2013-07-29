package redisync

import (
	"testing"
)

func TestLock(t *testing.T) {
	m, err := New("redisync.test.1")
	if err != nil {
		t.Error(err)
	}
	m.Lock()
	m.Unlock()
}
