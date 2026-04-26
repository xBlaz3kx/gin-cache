package cache

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Config contains all options
type Config struct {
	logger Logger

	getCacheStrategyByRequest GetCacheStrategyByRequest

	hitCacheCallback  OnHitCacheCallback
	missCacheCallback OnMissCacheCallback

	beforeReplyWithCacheCallback BeforeReplyWithCacheCallback

	singleFlightForgetTimeout time.Duration
	shareSingleFlightCallback OnShareSingleFlightCallback

	ignoreQueryOrder bool
	prefixKey        string
	withoutHeader    bool
	discardHeaders   []string
}

func newConfigByOpts(opts ...Option) *Config {
	cfg := &Config{
		logger:                       NoopLogger{},
		hitCacheCallback:             defaultHitCacheCallback,
		missCacheCallback:            defaultMissCacheCallback,
		beforeReplyWithCacheCallback: defaultBeforeReplyWithCacheCallback,
		shareSingleFlightCallback:    defaultShareSingleFlightCallback,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

// Option represents the optional function.
type Option func(c *Config)

// WithLogger set the custom logger
func WithLogger(l Logger) Option {
	return func(c *Config) {
		if l != nil {
			c.logger = l
		}
	}
}

// WithCacheStrategyByRequest set up the custom strategy per request
func WithCacheStrategyByRequest(getGetCacheStrategyByRequest GetCacheStrategyByRequest) Option {
	return func(c *Config) {
		if getGetCacheStrategyByRequest != nil {
			c.getCacheStrategyByRequest = getGetCacheStrategyByRequest
		}
	}
}

// OnHitCacheCallback define the callback when cache was hit
type OnHitCacheCallback func(c *gin.Context)

var defaultHitCacheCallback = func(c *gin.Context) {}

// WithOnHitCache will be called when cache hit.
func WithOnHitCache(cb OnHitCacheCallback) Option {
	return func(c *Config) {
		if cb != nil {
			c.hitCacheCallback = cb
		}
	}
}

// OnMissCacheCallback callback when cache was missed
type OnMissCacheCallback func(c *gin.Context)

// Noop default callback
var defaultMissCacheCallback = func(c *gin.Context) {}

// WithOnMissCache will be called when cache miss.
func WithOnMissCache(cb OnMissCacheCallback) Option {
	return func(c *Config) {
		if cb != nil {
			c.missCacheCallback = cb
		}
	}
}

// BeforeReplyWithCacheCallback define the callback before replying with cache
type BeforeReplyWithCacheCallback func(c *gin.Context, cache *ResponseCache)

var defaultBeforeReplyWithCacheCallback = func(c *gin.Context, cache *ResponseCache) {}

// WithBeforeReplyWithCache will be called before replying with cache.
func WithBeforeReplyWithCache(cb BeforeReplyWithCacheCallback) Option {
	return func(c *Config) {
		if cb != nil {
			c.beforeReplyWithCacheCallback = cb
		}
	}
}

// OnShareSingleFlightCallback define the callback when share the singleflight result
type OnShareSingleFlightCallback func(c *gin.Context)

var defaultShareSingleFlightCallback = func(c *gin.Context) {}

// WithOnShareSingleFlight will be called when share the singleflight result
func WithOnShareSingleFlight(cb OnShareSingleFlightCallback) Option {
	return func(c *Config) {
		if cb != nil {
			c.shareSingleFlightCallback = cb
		}
	}
}

// WithSingleFlightForgetTimeout to reduce the impact of long tail requests.
// singleflight.Forget will be called after the timeout has reached for each backend request when timeout is greater than zero.
func WithSingleFlightForgetTimeout(forgetTimeout time.Duration) Option {
	return func(c *Config) {
		if forgetTimeout > 0 {
			c.singleFlightForgetTimeout = forgetTimeout
		}
	}
}

// IgnoreQueryOrder will ignore the queries order in url when generate cache key.
// This option only takes effect in CacheByRequestURI strategy
func IgnoreQueryOrder() Option {
	return func(c *Config) {
		c.ignoreQueryOrder = true
	}
}

// WithPrefixKey will prefix the key
func WithPrefixKey(prefix string) Option {
	return func(c *Config) {
		c.prefixKey = prefix
	}
}

func WithoutHeader() Option {
	return func(c *Config) {
		c.withoutHeader = true
	}
}

func WithDiscardHeaders(headers []string) Option {
	return func(c *Config) {
		c.discardHeaders = headers
	}
}

func CorsHeaders() []string {
	return []string{
		"Access-Control-Allow-Credentials",
		"Access-Control-Expose-Headers",
		"Access-Control-Allow-Origin",
		"Vary",
	}
}
