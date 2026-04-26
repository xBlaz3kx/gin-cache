package cache

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xBlaz3kx/gin-cache/persist"
)

// Strategy the cache strategy
type Strategy struct {
	CacheKey string

	// CacheStore if nil, use default cache store instead
	CacheStore persist.CacheStore

	// CacheDuration
	CacheDuration time.Duration
}

// GetCacheStrategyByRequest User can this function to design custom cache strategy by request.
// The first return value bool means whether this request should be cached.
// The second return value Strategy determine the special strategy by this request.
type GetCacheStrategyByRequest func(c *gin.Context) (bool, Strategy)
