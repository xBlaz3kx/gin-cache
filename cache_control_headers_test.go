package cache

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type cacheControlHeadersTestSuite struct {
	suite.Suite
}

func (s *cacheControlHeadersTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (s *cacheControlHeadersTestSuite) Test_build() {
	dur := func(d time.Duration) *time.Duration { return &d }

	tests := []struct {
		name        string
		config      ControlConfig
		mustContain []string
	}{
		{
			name:        "Default empty config",
			config:      ControlConfig{},
			mustContain: []string{""},
		},
		{
			name:        "NoCache",
			config:      ControlConfig{NoCache: true},
			mustContain: []string{NoCache},
		},
		{
			name:        "NoStore",
			config:      ControlConfig{NoStore: true},
			mustContain: []string{NoStore},
		},
		{
			name:        "MustRevalidate",
			config:      ControlConfig{MustRevalidate: true},
			mustContain: []string{MustRevalidate},
		},
		{
			name:        "NoTransform",
			config:      ControlConfig{NoTransform: true},
			mustContain: []string{NoTransform},
		},
		{
			name:        "Public",
			config:      ControlConfig{Public: true},
			mustContain: []string{Public},
		},
		{
			name:        "Private",
			config:      ControlConfig{Private: true},
			mustContain: []string{Private},
		},
		{
			name:        "Public overrides Private",
			config:      ControlConfig{Public: true, Private: true},
			mustContain: []string{Public},
		},
		{
			name:        "ProxyRevalidate",
			config:      ControlConfig{ProxyRevalidate: true},
			mustContain: []string{ProxyRevalidate},
		},
		{
			name:        "Immutable",
			config:      ControlConfig{Immutable: true},
			mustContain: []string{Immutable},
		},
		{
			name:        "MaxAge",
			config:      ControlConfig{maxAge: dur(60 * time.Second)},
			mustContain: []string{"max-age=60"},
		},
		{
			name:        "SMaxAge",
			config:      ControlConfig{sMaxAge: dur(120 * time.Second)},
			mustContain: []string{"s-maxage=120"},
		},
		{
			name:        "StaleWhileRevalidate",
			config:      ControlConfig{staleWhileRevalidate: dur(30 * time.Second)},
			mustContain: []string{"stale-while-revalidate=30"},
		},
		{
			name:        "StaleIfError",
			config:      ControlConfig{staleIfError: dur(45 * time.Second)},
			mustContain: []string{"stale-if-error=45"},
		},
		{
			name: "No cache, must revalidate, no store",
			config: ControlConfig{
				MustRevalidate: true,
				NoCache:        true,
				NoStore:        true,
			},
			mustContain: []string{NoCache, NoStore, MustRevalidate},
		},
		{
			name: "All boolean directives",
			config: ControlConfig{
				MustRevalidate:  true,
				NoCache:         true,
				NoStore:         true,
				NoTransform:     true,
				Public:          true,
				ProxyRevalidate: true,
				Immutable:       true,
			},
			mustContain: []string{MustRevalidate, NoCache, NoStore, NoTransform, Public, ProxyRevalidate, Immutable},
		},
		{
			name: "All time-based directives",
			config: ControlConfig{
				maxAge:               dur(60 * time.Second),
				sMaxAge:              dur(120 * time.Second),
				staleWhileRevalidate: dur(30 * time.Second),
				staleIfError:         dur(45 * time.Second),
			},
			mustContain: []string{"max-age=60", "s-maxage=120", "stale-while-revalidate=30", "stale-if-error=45"},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			actual := tt.config.build()

			for _, want := range tt.mustContain {
				if want != "" {
					s.Contains(actual, want)
				}
			}
		})
	}
}

func (s *cacheControlHeadersTestSuite) Test_apply() {
	tests := []struct {
		name            string
		config          ControlConfig
		expectedHeaders []string
	}{
		{
			name:            "Default empty config",
			config:          ControlConfig{},
			expectedHeaders: nil,
		},
		{
			name: "No cache, must revalidate, no store",
			config: ControlConfig{
				MustRevalidate: true,
				NoCache:        true,
				NoStore:        true,
			},
			expectedHeaders: []string{MustRevalidate, NoCache, NoStore},
		},
		{
			name: "No cache",
			config: ControlConfig{
				NoCache: true,
			},
			expectedHeaders: []string{NoCache},
		},
		{
			name: "Public with max-age",
			config: ControlConfig{
				Public: true,
				maxAge: func() *time.Duration { d := 300 * time.Second; return &d }(),
			},
			expectedHeaders: []string{Public, "max-age=300"},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

			handler := NewControlHeadersMiddleware(tt.config)
			handler(ctx)

			actual := ctx.Writer.Header().Get(ControlHeader)
			for _, want := range tt.expectedHeaders {
				s.Contains(actual, want)
			}
		})
	}
}

func (s *cacheControlHeadersTestSuite) TestNewControlHeadersMiddleware_withOpts() {
	tests := []struct {
		name        string
		config      ControlConfig
		opts        []ControlOption
		mustContain []string
	}{
		{
			name:        "WithMaxAge option",
			config:      ControlConfig{Public: true},
			opts:        []ControlOption{WithMaxAge(60 * time.Second)},
			mustContain: []string{Public, "max-age=60"},
		},
		{
			name:        "WithSMaxAge option",
			config:      ControlConfig{Public: true},
			opts:        []ControlOption{WithSMaxAge(120 * time.Second)},
			mustContain: []string{Public, "s-maxage=120"},
		},
		{
			name:        "WithStaleWhileRevalidate option",
			config:      ControlConfig{Public: true},
			opts:        []ControlOption{WithStaleWhileRevalidate(30 * time.Second)},
			mustContain: []string{Public, "stale-while-revalidate=30"},
		},
		{
			name:        "WithStaleIfError option",
			config:      ControlConfig{Public: true},
			opts:        []ControlOption{WithStaleIfError(45 * time.Second)},
			mustContain: []string{Public, "stale-if-error=45"},
		},
		{
			name:   "Multiple opts combined",
			config: ControlConfig{Public: true},
			opts: []ControlOption{
				WithMaxAge(60 * time.Second),
				WithSMaxAge(120 * time.Second),
				WithStaleWhileRevalidate(30 * time.Second),
				WithStaleIfError(45 * time.Second),
			},
			mustContain: []string{Public, "max-age=60", "s-maxage=120", "stale-while-revalidate=30", "stale-if-error=45"},
		},
		{
			name:        "Opt overrides config field",
			config:      ControlConfig{maxAge: func() *time.Duration { d := 10 * time.Second; return &d }()},
			opts:        []ControlOption{WithMaxAge(999 * time.Second)},
			mustContain: []string{"max-age=999"},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

			handler := NewControlHeadersMiddleware(tt.config, tt.opts...)
			handler(ctx)

			actual := ctx.Writer.Header().Get(ControlHeader)
			for _, want := range tt.mustContain {
				s.Contains(actual, want)
			}
		})
	}
}

func (s *cacheControlHeadersTestSuite) TestNewControlHeadersMiddleware_setsHeader() {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Writer.Header().Set(ControlHeader, "old-value")

	handler := NewControlHeadersMiddleware(ControlConfig{NoCache: true})
	handler(ctx)

	actual := ctx.Writer.Header().Get(ControlHeader)
	s.Equal(NoCache, actual)
	s.NotContains(actual, "old-value")
}

func (s *cacheControlHeadersTestSuite) TestNoCachePreset() {
	expected := []string{MustRevalidate, NoCache, NoStore}
	actual := NoCachePreset.build()

	for _, want := range expected {
		s.Contains(actual, want)
	}
}

func (s *cacheControlHeadersTestSuite) TestNoCachePreset_asMiddleware() {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	handler := NewControlHeadersMiddleware(NoCachePreset)
	handler(ctx)

	actual := ctx.Writer.Header().Get(ControlHeader)
	for _, want := range strings.Split("must-revalidate, no-cache, no-store", ", ") {
		s.Contains(actual, want)
	}
}

func TestCacheControlHeaders(t *testing.T) {
	suite.Run(t, new(cacheControlHeadersTestSuite))
}
