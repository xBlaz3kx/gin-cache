package cache

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const ControlHeader = "Cache-Control"

const (
	NoCache              = "no-cache"
	NoStore              = "no-store"
	MaxAge               = "max-age"
	SMaxAge              = "s-maxage"
	Immutable            = "immutable"
	StaleWhileRevalidate = "stale-while-revalidate"
	StaleIfError         = "stale-if-error"
	Public               = "public"
	Private              = "private"
	MustRevalidate       = "must-revalidate"
	NoTransform          = "no-transform"
	ProxyRevalidate      = "proxy-revalidate"
)

// ControlConfig defines a cache-control configuration.
//
// References:
// https://datatracker.ietf.org/doc/html/rfc7234#section-5.2.2
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control
type ControlConfig struct {
	MustRevalidate       bool
	NoCache              bool
	NoStore              bool
	NoTransform          bool
	Public               bool
	Private              bool
	ProxyRevalidate      bool
	maxAge               *time.Duration
	sMaxAge              *time.Duration
	Immutable            bool
	staleWhileRevalidate *time.Duration
	staleIfError         *time.Duration
}

func (c *ControlConfig) build() string {
	values := map[string]struct{}{}

	if c.NoCache {
		values[NoCache] = struct{}{}
	}

	if c.NoStore {
		values[NoStore] = struct{}{}
	}

	if c.MustRevalidate {
		values[MustRevalidate] = struct{}{}
	}

	if c.NoTransform {
		values[NoTransform] = struct{}{}
	}

	if c.Public {
		values[Public] = struct{}{}

		// Cannot be both public and private.
		c.Private = false
	}

	if c.Private {
		values[Private] = struct{}{}
	}

	if c.ProxyRevalidate {
		values[ProxyRevalidate] = struct{}{}
	}

	if c.maxAge != nil {
		maxAge := fmt.Sprintf("%s=%.f", MaxAge, c.maxAge.Seconds())
		values[maxAge] = struct{}{}
	}

	if c.sMaxAge != nil {
		maxAge := fmt.Sprintf("%s=%.f", SMaxAge, c.sMaxAge.Seconds())
		values[maxAge] = struct{}{}
	}

	if c.Immutable {
		values[Immutable] = struct{}{}
	}

	if c.staleWhileRevalidate != nil {
		stale := fmt.Sprintf("%s=%.f", StaleWhileRevalidate, c.staleWhileRevalidate.Seconds())
		values[stale] = struct{}{}
	}

	if c.staleIfError != nil {
		stale := fmt.Sprintf("%s=%.f", StaleIfError, c.staleIfError.Seconds())
		values[stale] = struct{}{}
	}

	return strings.Join(slices.Collect(maps.Keys(values)), ", ")
}

func (c *ControlConfig) apply(ginCtx *gin.Context, value string) {
	ginCtx.Writer.Header().Set(ControlHeader, value)
}

// NewControlHeadersMiddleware creates a new Gin middleware which generates a cache-control header.
// Existing cache-control headers are removed.
// Other caching-related headers, such as `Expires` and `Pragma`, remain unchanged.
func NewControlHeadersMiddleware(config ControlConfig, opts ...ControlOption) gin.HandlerFunc {
	for _, opt := range opts {
		opt(&config)
	}

	value := config.build()

	return func(ginCtx *gin.Context) {
		config.apply(ginCtx, value)
	}
}

// NoCachePreset is a cache-control configuration preset which advices the HTTP client not to cache at all.
var NoCachePreset = ControlConfig{
	MustRevalidate: true,
	NoCache:        true,
	NoStore:        true,
}
