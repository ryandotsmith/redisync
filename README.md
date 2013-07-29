# Redisync

A Go package which implements synchronization functions on top of Redis. The heavy lifting is done with Lua scripts to make the elimination of race conditions easier.

## Install
```bash
$ go install github.com/ryandotsmith/redisync
```

## Documentation:
[GoDoc](http://godoc.org/github.com/ryandotsmith/redisync)

## Hacking on Redisync

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
