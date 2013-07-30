package redisync

import (
	"time"
	"fmt"
)

func ExampleLock() {
	ttl := time.Second
	m, err := NewMutex("my-lock", ttl, "")
	if err != nil {
		return
	}
	m.Lock()

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
