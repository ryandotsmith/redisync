# Redisync

A Go package which implements synchronization functions on top of Redis. The heavy lifting is done with Lua scripts to make the elimination of race conditions easier.

Note: When using a TTL, the user should take care to finsih execution before the TTL.

## Install
```bash
$ go install github.com/ryandotsmith/redisync
```

## Usage
```go
package main

import "github.com/ryandotsmith/redisync"

func main() {
	ttl := time.Minute
	m := redisync.NewMutex("my-lock", ttl, "redis://u:p@localhost:6379")
	m.Lock()
	defer m.Unlock()
	print("at=critical-section\n")
}
```

## Documentation
[GoDoc](http://godoc.org/github.com/ryandotsmith/redisync)

## Hacking on Redisync

[![Build Status](https://drone.io/github.com/ryandotsmith/redisync/status.png)](https://drone.io/github.com/ryandotsmith/redisync/latest)

```bash
$ go version
go version go1.1.1 darwin/amd64
$ ./redis-server --version
Redis server v=2.6.14 sha=00000000:0 malloc=libc bits=64
```

```bash
$ git clone git://github.com/ryandotsmith/redisync.git
$ go get ./...
$ ./redis-server &
$ go test
```
