package cache

import (
	"net/http/httptest"
	"strings"
	"testing"

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
	tests := []struct {
		name     string
		config   ControlConfig
		expected string
	}{
		{
			name:     "Default",
			config:   ControlConfig{},
			expected: "",
		},
		{
			name: "No cache, must revalidate, no store",
			config: ControlConfig{
				MustRevalidate: true,
				NoCache:        true,
				NoStore:        true,
			},
			expected: "no-cache, no-store, must-revalidate",
		},
		{
			name: "No cache",
			config: ControlConfig{
				NoCache: true,
			},
			expected: "no-cache",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			actual := tt.config.build()

			for _, header := range strings.Split(tt.expected, ", ") {
				s.Contains(actual, header)
			}
		})
	}
}

func (s *cacheControlHeadersTestSuite) Test_apply() {
	tests := []struct {
		name            string
		config          ControlConfig
		expectedHeaders string
	}{
		{
			name:            "Default",
			config:          ControlConfig{},
			expectedHeaders: "",
		},
		{
			name: "No cache, must revalidate, no store",
			config: ControlConfig{
				MustRevalidate: true,
				NoCache:        true,
				NoStore:        true,
			},
			expectedHeaders: "must-revalidate, no-cache, no-store",
		},
		{
			name: "No cache",
			config: ControlConfig{
				NoCache: true,
			},
			expectedHeaders: "no-cache",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

			// Apply headers
			handler := NewControlHeadersMiddleware(tt.config)
			handler(ctx)

			actual := ctx.Writer.Header().Get(ControlHeader)
			for _, header := range strings.Split(tt.expectedHeaders, ", ") {
				s.Contains(actual, header)
			}
		})
	}
}

func (s *cacheControlHeadersTestSuite) TestNoCachePreset() {
	expected := "must-revalidate, no-cache, no-store"
	actual := NoCachePreset.build()

	for _, header := range strings.Split(expected, ", ") {
		s.Contains(actual, header)
	}
}

func TestCacheControlHeaders(t *testing.T) {
	suite.Run(t, new(cacheControlHeadersTestSuite))
}
