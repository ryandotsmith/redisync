package redisync

import (
	"github.com/garyburd/redigo/redis"
	"net/url"
	"os"
	"testing"
	"time"
)

func newConn() (redis.Conn, error) {
	redisUrl, err := url.Parse(os.Getenv("REDIS_URL"))
	if err != nil {
		return nil, err
	}
	c, err := redis.Dial("tcp", redisUrl.Host)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func TestLock(t *testing.T) {
	rc, err := newConn()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer rc.Close()
	ttl := time.Second
	m := NewMutex("redisync.test.1", ttl)
	m.Lock(rc)
	time.Sleep(ttl)
	ok := m.TryLock(rc)
	if !ok {
		t.Error("Expected mutex to be lockable.")
		t.FailNow()
	}
	m.Unlock(rc)
}

func TestLockLocked(t *testing.T) {
	rc, err := newConn()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer rc.Close()

	ttl := time.Second
	m1 := NewMutex("redisync.test.1", ttl)
	if ok := m1.TryLock(rc); !ok {
		t.Error("Expected mutex to be lockable.")
		t.FailNow()
	}

	m2 := NewMutex("redisync.test.1", ttl)
	if ok := m2.TryLock(rc); ok {
		t.Error("Expected mutex not to be lockable.")
		t.FailNow()
	}
	m1.Unlock(rc)
}

func TestUnlockOtherLocked(t *testing.T) {
	rc, err := newConn()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer rc.Close()

	ttl := time.Second
	m1 := NewMutex("redisync.test.1", ttl)
	if ok := m1.TryLock(rc); !ok {
		t.Error("Expected mutex to be lockable.")
		t.FailNow()
	}

	m2 := NewMutex("redisync.test.1", ttl)
	if ok, _ := m2.Unlock(rc); ok {
		t.Error("Expected mutex not to be unlockable.")
		t.FailNow()
	}
	m1.Unlock(rc)
}

func TestLockExpired(t *testing.T) {
	rc, err := newConn()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer rc.Close()

	ttl := time.Second
	m1 := NewMutex("redisync.test.1", ttl)
	if ok := m1.TryLock(rc); !ok {
		t.Error("Expected mutex to be lockable.")
		t.FailNow()
	}
	time.Sleep(ttl)

	m2 := NewMutex("redisync.test.1", ttl)
	if ok := m2.TryLock(rc); !ok {
		t.Error("Expected mutex to be lockable.")
		t.FailNow()
	}
	m2.Unlock(rc)
}
