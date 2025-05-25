package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/xBlaz3kx/gin-cache"
	"github.com/xBlaz3kx/gin-cache/persist"
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
