package redisync

import (
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	ttl := time.Second
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
	m.Unlock()
}

func TestLockLocked(t *testing.T) {
	ttl := time.Second
	m1, err := NewMutex("redisync.test.1", ttl, "")
	if err != nil {
		t.Error(err)
	}
	if ok := m1.TryLock(); !ok {
		t.Error("Expected mutex to be lockable.")
		t.FailNow()
	}
	m2, err := NewMutex("redisync.test.1", ttl, "")
	if err != nil {
		t.Error(err)
	}
	if ok := m2.TryLock(); ok {
		t.Error("Expected mutex not to be lockable.")
		t.FailNow()
	}
	m1.Unlock()
}

func TestUnlockOtherLocked(t *testing.T) {
	ttl := time.Second
	m1, err := NewMutex("redisync.test.1", ttl, "")
	if err != nil {
		t.Error(err)
	}
	if ok := m1.TryLock(); !ok {
		t.Error("Expected mutex to be lockable.")
		t.FailNow()
	}

	m2, err := NewMutex("redisync.test.1", ttl, "")
	if err != nil {
		t.Error(err)
	}
	if ok, _ := m2.Unlock(); ok {
		t.Error("Expected mutex not to be unlockable.")
		t.FailNow()
	}
	m1.Unlock()
}

func TestLockExpired(t *testing.T) {
	ttl := time.Second

	m1, err := NewMutex("redisync.test.1", ttl, "")
	if err != nil {
		t.Error(err)
	}
	if ok := m1.TryLock(); !ok {
		t.Error("Expected mutex to be lockable.")
		t.FailNow()
	}
	time.Sleep(ttl)

	m2, err := NewMutex("redisync.test.1", ttl, "")
	if err != nil {
		t.Error(err)
	}
	if ok := m2.TryLock(); !ok {
		t.Error("Expected mutex to be lockable.")
		t.FailNow()
	}
	m2.Unlock()
}
