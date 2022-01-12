package redisync

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

var rc redis.Conn

func init() {
	redisUrl, err := url.Parse(os.Getenv("REDIS_URL"))
	if err != nil {
		panic("Must set: $REDIS_URL")
	}
	rc, err = redis.Dial("tcp", redisUrl.Host)
	if err != nil {
		panic("Unable to connect to: $REDIS_URL")
	}
}

func ExampleLock() {
	ttl := time.Second
	m := NewMutex("my-lock", ttl)
	m.Lock(rc)
	defer m.Unlock(rc)

	done := make(chan bool)
	expired := make(chan bool)

	go func(e chan bool) {
		time.Sleep(ttl)
		e <- true
	}(expired)

	go func(d chan bool) {
		fmt.Printf("at=critical-section\n")
		d <- true
	}(done)

	select {
	case <-done:
		fmt.Printf("Finished.\n")
	case <-expired:
		fmt.Printf("Expired.\n")
	}
	// Output:
	// at=critical-section
	// Finished.
}
