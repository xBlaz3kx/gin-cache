# gin-cache

[![Release](https://img.shields.io/github/release/chenyahui/gin-cache.svg?style=flat-square)](https://github.com/xBlaz3kx/gin-cache/releases)
[![doc](https://img.shields.io/badge/go.dev-doc-007d9c?style=flat-square&logo=read-the-docs)](https://pkg.go.dev/github.com/xBlaz3kx/gin-cache)
[![goreportcard for gin-cache](https://goreportcard.com/badge/github.com/xBlaz3kx/gin-cache)](https://goreportcard.com/report/github.com/xBlaz3kx/gin-cache)
![](https://img.shields.io/badge/license-MIT-green)
[![codecov](https://codecov.io/gh/chenyahui/gin-cache/branch/main/graph/badge.svg?token=MX8Z4D5RZS)](https://codecov.io/gh/chenyahui/gin-cache)

A high performance Gin middleware to cache http responses.

# Features

* Has a performance improvement compared to `gin-contrib/cache`.
* Multiple cache store implementations; including `in-memory` and `Redis`.
* Customizable cache key (with prefix)
* Sets appropriate cache control headers by default
* Use singleflight to avoid cache breakdown problem.
* Only caches successful responses (status code 200-299)

# How To Use

## Install

```
go get -u github.com/xBlaz3kx/gin-cache
```

## Example

### Cache In Local Memory

```go
package main

import (
	"time"

	"github.com/xBlaz3kx/gin-cache"
	"github.com/xBlaz3kx/gin-cache/persist"
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.New()

	memoryStore := persist.NewMemoryStore(1 * time.Minute)

	app.GET("/hello",
		cache.CacheByRequestURI(memoryStore, 2*time.Second),
		func(c *gin.Context) {
			c.String(200, "hello world")
		},
	)

	if err := app.Run(":8080"); err != nil {
		panic(err)
	}
}
```

### Cache In Redis

```go
package main

import (
	"time"

	"github.com/xBlaz3kx/gin-cache"
	"github.com/xBlaz3kx/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func main() {
	app := gin.New()

	redisStore := persist.NewRedisStore(redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "127.0.0.1:6379",
	}))

	app.GET("/hello",
		cache.CacheByRequestURI(redisStore, 2*time.Second),
		func(c *gin.Context) {
			c.String(200, "hello world")
		},
	)
	if err := app.Run(":8080"); err != nil {
		panic(err)
	}
}
```

## MemoryStore

![MemoryStore QPS](https://www.cyhone.com/img/gin-cache/memory_cache_qps.png)

## RedisStore

![RedisStore QPS](https://www.cyhone.com/img/gin-cache/redis_cache_qps.png)
