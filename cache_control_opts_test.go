package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOpts(t *testing.T) {
	tests := []struct {
		name string
		opts []ControlOption
		want ControlConfig
	}{
		{
			name: "With max-age",
			opts: []ControlOption{
				WithMaxAge(time.Second),
			},
			want: ControlConfig{
				maxAge: toPtr(time.Second),
			},
		},
		{
			name: "With s-maxage",
			opts: []ControlOption{
				WithSMaxAge(time.Second),
			},
			want: ControlConfig{
				sMaxAge: toPtr(time.Second),
			},
		},
		{
			name: "With stale-while-revalidate",
			opts: []ControlOption{
				WithStaleWhileRevalidate(time.Second),
			},
			want: ControlConfig{
				staleWhileRevalidate: toPtr(time.Second),
			},
		},
		{
			name: "With stale-if-error",
			opts: []ControlOption{
				WithStaleIfError(time.Second),
			},
			want: ControlConfig{
				staleIfError: toPtr(time.Second),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := ControlConfig{}

			for _, opt := range tt.opts {
				opt(&config)
			}

			assert.Equal(t, tt.want, config)
		})
	}
}

func toPtr[T any](x T) *T {
	return &x
}
