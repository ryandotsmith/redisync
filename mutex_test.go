package redisync

import (
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	ttl := time.Second * 2
	m, err := NewMutex("redisync.test.1", ttl, "")
	if err != nil {
		t.Error(err)
	}
	m.Lock()
	time.Sleep(ttl)
	ok := m.TryLock()
	if !ok {
		t.Error("Expected mutex to be lockable.")
		t.FailNow()
	}
}
