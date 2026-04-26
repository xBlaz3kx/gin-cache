package cache

import "time"

type ControlOption func(*ControlConfig)

// WithMaxAge sets the max-age directive in the Cache-Control cacheControlHeader.
func WithMaxAge(d time.Duration) ControlOption {
	return func(c *ControlConfig) {
		c.maxAge = &d
	}
}

// WithSMaxAge sets the s-maxage directive in the Cache-Control cacheControlHeader.
func WithSMaxAge(d time.Duration) ControlOption {
	return func(c *ControlConfig) {
		c.sMaxAge = &d
	}
}

// WithStaleWhileRevalidate sets the stale-while-revalidate directive in the Cache-Control cacheControlHeader.
func WithStaleWhileRevalidate(d time.Duration) ControlOption {
	return func(c *ControlConfig) {
		c.staleWhileRevalidate = &d
	}
}

// WithStaleIfError sets the stale-if-error directive in the Cache-Control cacheControlHeader.
func WithStaleIfError(d time.Duration) ControlOption {
	return func(c *ControlConfig) {
		c.staleIfError = &d
	}
}
