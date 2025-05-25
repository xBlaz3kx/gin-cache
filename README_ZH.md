# gin-cache 
[![Release](https://img.shields.io/github/release/chenyahui/gin-cache.svg?style=flat-square)](https://github.com/xBlaz3kx/gin-cache/releases)
[![doc](https://img.shields.io/badge/go.dev-doc-007d9c?style=flat-square&logo=read-the-docs)](https://pkg.go.dev/github.com/xBlaz3kx/gin-cache)
[![goreportcard for gin-cache](https://goreportcard.com/badge/github.com/xBlaz3kx/gin-cache)](https://goreportcard.com/report/github.com/xBlaz3kx/gin-cache)
![](https://img.shields.io/badge/license-MIT-green)
[![codecov](https://codecov.io/gh/chenyahui/gin-cache/branch/main/graph/badge.svg?token=MX8Z4D5RZS)](https://codecov.io/gh/chenyahui/gin-cache)

[English](README_ZH.md) | 🇨🇳中文

一个用于缓存http接口内容的gin高性能中间件。相比于官方的gin-contrib/cache，gin-cache有巨大的性能提升。

# 特性
* 相比于gin-contrib/cache，性能提升巨大。
* 同时支持本机内存和redis作为缓存后端。
* 支持用户根据请求来指定cache策略。
* 使用singleflight解决了缓存击穿问题。
* 仅缓存http状态码为2xx的回包

# 用法

## 安装

```
go get -u github.com/xBlaz3kx/gin-cache
```

## 例子
## 使用本地缓存
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

### 使用redis作为缓存
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

# 压测
```
wrk -c 500 -d 1m -t 5 http://127.0.0.1:8080/hello
```

## MemoryStore

![MemoryStore QPS](https://www.cyhone.com/img/gin-cache/memory_cache_qps.png)

## RedisStore

![RedisStore QPS](https://www.cyhone.com/img/gin-cache/redis_cache_qps.png)
