package cache

import (
	"encoding/gob"
	"net/http"
)

func init() {
	gob.Register(&ResponseCache{})
}

// ResponseCache record the http response cache
type ResponseCache struct {
	Status int
	Header http.Header
	Data   []byte
}

func (c *ResponseCache) fillWithCacheWriter(cacheWriter *responseCacheWriter, cfg *Config) {
	c.Status = cacheWriter.Status()
	c.Data = cacheWriter.body.Bytes()
	if !cfg.withoutHeader {
		c.Header = cacheWriter.Header().Clone()

		for _, headerKey := range cfg.discardHeaders {
			c.Header.Del(headerKey)
		}
	}
}
