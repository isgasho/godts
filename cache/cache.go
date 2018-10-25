package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var local_cache *cache.Cache

func NewCache() {
	local_cache = cache.New(5*time.Minute, 10*time.Minute)
}

func Cache() *cache.Cache {
	if local_cache == nil {
		NewCache()
	}
	return local_cache
}
