// Synchronization built on top of Redis.
// Depends on github.com/garyburd/redigo/redis
package redisync

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

/*
Locking algorithm:
	if the key exists
		if caller is owner of lock
			update ttl
			return true
		else return false
	if the key does not exist
		set owner
		set ttl
		return true
*/
var luaLock = `
local call_owner = ARGV[1]
local ttl = ARGV[2]
if redis.call('EXISTS', KEYS[1]) == 1 then
	local lock_owner = redis.call('GET', KEYS[1])
	if lock_owner == call_owner then
		redis.call('EXPIRE', KEYS[1], ttl)
		return 1
	end
	return 0
else
	redis.call('SET', KEYS[1], call_owner)
	redis.call('EXPIRE', KEYS[1], ttl)
	return 1
end
`

/*
Unlocking algorithm:
	if the key exists
		if owner of lock
			delete key
			return true
		else return false
	else return false
*/
var luaUnlock = `
local call_owner = ARGV[1]
if redis.call('EXISTS', KEYS[1]) == 1 then
	local lock_owner = redis.call('GET', KEYS[1])
	if lock_owner == call_owner then
		redis.call('DEL', KEYS[1])
		return 1
	end
	return 0
else
	return 0
end
`

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
	// See lock.lua
	lock *redis.Script
	// See unlock.lua
	unlock *redis.Script
}

// Each lock will have a name which corresponds to a key in the Redis server.
// The mutex will also be initialized with a uuid. The mutex uuid
// can be used to extend the TTL for the lock.
func NewMutex(name string, ttl time.Duration) *Mutex {
	m := new(Mutex)
	m.Name = name
	m.Ttl = ttl
	m.Backoff = time.Second
	m.id = uuid()
	m.lock = redis.NewScript(1, luaLock)
	m.unlock = redis.NewScript(1, luaUnlock)
	return m
}

// With similar behaviour to Go's sync pkg,
// this function will sleep until TryLock() returns true.
// The connection will be used once to execute the lock script.
func (m *Mutex) Lock(c redis.Conn) {
	for {
		if m.TryLock(c) {
			return
		}
		time.Sleep(m.Backoff)
	}
}

// Makes a single attempt to acquire the lock.
// Locking a mutex which has already been locked
// using the mutex uuid will result in the TTL of the mutex being extended.
// The connection will be used once to execute the lock script.
func (m *Mutex) TryLock(c redis.Conn) bool {
	m.l.Lock()
	defer m.l.Unlock()

	reply, err := m.lock.Do(c, m.Name, m.id, m.Ttl.Seconds())
	if err != nil {
		return false
	}
	return reply.(int64) == 1
}

// If the local mutex uuid matches the uuid in Redis,
// the lock will be deleted.
// The connection will be used once to execute the unlock script.
func (m *Mutex) Unlock(c redis.Conn) (bool, error) {
	m.l.Lock()
	defer m.l.Unlock()
	reply, err := m.unlock.Do(c, m.Name, m.id)
	if err != nil {
		return false, err
	}
	return reply.(int64) == 1, nil
}

func uuid() string {
	f, _ := os.Open("/dev/urandom")
	b := make([]byte, 16)
	f.Read(b)
	f.Close()
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
