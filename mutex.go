// Synchronization built on top of Redis.
// Depends on github.com/garyburd/redigo/redis
package redisync

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io/ioutil"
	"net/url"
	"os"
	"sync"
	"time"
)

type Mutex struct {
	// The key used in Redis.
	Name string
	// The amount of time before Redis will expire the lock.
	Ttl time.Duration
	// The time to sleep before retrying a lock attempt.
	Backoff time.Duration
	// A uuid representing the local instantiation of the mutex.
	id string
	// Local conrrency controll.
	l sync.Mutex
	// Redis connections are always created in package.
	c redis.Conn
	// See lock.lua
	lock *redis.Script
	// See unlock.lua
	unlock *redis.Script
}

// Each lock will have a name which corrisponds to a key in the Redis server.
// The mutex will also be initialized with a uuid. The mutex uuid
// can be used to extend the TTL for the lock.
// The redis url should take the form: redis://user:pass@host:port.
// If the rurl == "" then a connection will be made on localhsot at port
// 6379 with no authentication.
func NewMutex(name string, ttl time.Duration, rurl string) (*Mutex, error) {
	m := new(Mutex)
	m.Name = name
	m.Ttl = ttl
	m.Backoff = time.Second
	m.id = uuid()
	c, err := establishRedisConn(rurl)
	if err != nil {
		return nil, err
	}
	m.c = c
	m.lock = redis.NewScript(1, readSource("./lock.lua"))
	m.unlock = redis.NewScript(1, readSource("./unlock.lua"))
	return m, nil
}

func establishRedisConn(u string) (redis.Conn, error) {
	if len(u) == 0 {
		c, err := redis.Dial("tcp", ":6379")
		if err != nil {
			return nil, err
		}
		return c, nil
	}

	redisUrl, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	c, err := redis.Dial("tcp", redisUrl.Host)
	if err != nil {
		return nil, err
	}
	pass, usingPass := redisUrl.User.Password()
	if usingPass {
		_, err = c.Do("AUTH", pass)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

// With similar behaviour to Go's sync pkg,
// this function will sleep until TryLock() returns true.
func (m *Mutex) Lock() {
	for {
		if m.TryLock() {
			return
		}
		time.Sleep(m.Backoff)
	}
}

// Makes a single attempt to acquire the lock.
// Locking a mutex which has already been locked
// using the mutex uuid will result in the TTL of the mutex being extended.
func (m *Mutex) TryLock() bool {
	m.l.Lock()
	defer m.l.Unlock()

	reply, err := m.lock.Do(m.c, m.Name, m.id, m.Ttl.Seconds())
	if err != nil {
		return false
	}
	return reply.(int64) == 1
}

// If the local mutex uuid matches the uuid in Redis,
// the lock will be deleted.
func (m *Mutex) Unlock() (bool, error) {
	m.l.Lock()
	defer m.l.Unlock()
	reply, err := m.unlock.Do(m.c, m.Name, m.id)
	if err != nil {
		return false, err
	}
	return reply.(int64) == 1, nil
}

func readSource(name string) string {
	src, err := ioutil.ReadFile(name)
	if err != nil {
		panic("redisync: Unable to read unlock.lua")
	}
	return string(src)
}

func uuid() string {
	f, _ := os.Open("/dev/urandom")
	b := make([]byte, 16)
	f.Read(b)
	f.Close()
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
