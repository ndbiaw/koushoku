package cache

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bluele/gcache"
)

// LRU is a wrapper for gcache.Cache.
// It implements type-safe methods that
// only accept string or int64 keys.
type LRU struct {
	instance gcache.Cache
}

// New creates a new LRU cache.
func New(size int, expr time.Duration) *LRU {
	return &LRU{gcache.New(size).LRU().Expiration(expr).Build()}
}

// Get gets a value from the cache by using string key.
func (c *LRU) Get(key string) (any, error) {
	return c.instance.Get(key)
}

// Set sets a value to the cache by using string key.
func (c *LRU) Set(key string, value any, ttl time.Duration) error {
	if ttl > 0 {
		return c.instance.SetWithExpire(key, value, ttl)
	}
	return c.instance.Set(key, value)
}

// Remove removes a value from the cache by using string key.
func (c *LRU) Remove(key string) bool {
	return c.instance.Remove(key)
}

// Has checks if the cache has the key.
func (c *LRU) Has(key string) bool {
	return c.instance.Has(key)
}

// Keys gets all the keys in the cache.
func (c *LRU) Keys() []any {
	return c.instance.Keys(false)
}

// Purge purges all the cache.
func (c *LRU) Purge() {
	c.instance.Purge()
}

// GetWIthInt64 gets a value from the cache by using int64 key.
func (c *LRU) GetWithInt64(cid int64) (any, error) {
	return c.instance.Get(strconv.Itoa(int(cid)))
}

// SetWithInt64 sets a value to the cache by using int64 key.
func (c *LRU) SetWithInt64(cid int64, value any, ttl time.Duration) error {
	if ttl > 0 {
		return c.instance.SetWithExpire(strconv.Itoa(int(cid)), value, ttl)
	}
	return c.instance.Set(strconv.Itoa(int(cid)), value)
}

// RemoveWithInt64 removes a value from the cache by using int64 key.
func (c *LRU) RemoveWithInt64(cid int64) bool {
	return c.instance.Remove(strconv.Itoa(int(cid)))
}

// HasWithInt64 checks if the cache has the key.
func (c *LRU) HasWithInt64(cid int64) bool {
	return c.instance.Has(strconv.Itoa(int(cid)))
}

// GetWithPrefix gets a value from the cache by using prefix and key.
// Type of prefix and key can be string or int64.
func (c *LRU) GetWithPrefix(prefix, key any) (any, error) {
	return c.instance.Get(fmt.Sprintf("/%v/%v", prefix, key))
}

// SetWithPrefix sets a value to the cache by using prefix and key.
// Type of prefix and key can be string or int64.
func (c *LRU) SetWithPrefix(prefix, key, value any, ttl time.Duration) error {
	if ttl > 0 {
		return c.instance.SetWithExpire(fmt.Sprintf("/%v/%v", prefix, key), value, ttl)
	}
	return c.instance.Set(fmt.Sprintf("/%v/%v", prefix, key), value)
}

// RemoveWithPrefix removes a value from the cache by using prefix and key.
// Type of prefix and key can be string or int64.
func (c *LRU) RemoveWithPrefix(prefix, key any) bool {
	return c.instance.Remove(fmt.Sprintf("/%v/%v", prefix, key))
}

// HasWithPrefix checks if the cache has the key.
func (c *LRU) HasWithPrefix(prefix, key any) bool {
	return c.instance.Has(fmt.Sprintf("/%v/%v", prefix, key))
}

// PurgeWithPrefix purges all the cache with the given prefix.
// Type of prefix can be string or int64.
func (c *LRU) PurgeWithPrefix(prefix any) {
	for _, k := range c.instance.Keys(false) {
		if strings.HasPrefix(fmt.Sprintf("%v", k), fmt.Sprintf("/%v/", prefix)) {
			c.instance.Remove(k)
		}
	}
}

// KeysWithPrefix gets all the keys with the given prefix.
// Type of prefix can be string or int64.
func (c *LRU) KeysWithPrefix(prefix any) []string {
	var keys []string
	for _, k := range c.instance.Keys(false) {
		if strings.HasPrefix(fmt.Sprintf("%v", k), fmt.Sprintf("/%v/", prefix)) {
			keys = append(keys, strings.TrimPrefix(fmt.Sprintf("%v", k), fmt.Sprintf("/%v/", prefix)))
		}
	}
	return keys
}

// CacheStats represents the stats of the LRU cache.
type CacheStats struct {
	Size        int     `json:"size"`
	HitCount    uint64  `json:"hitCount"`
	HitRate     float64 `json:"hitRate"`
	MissCount   uint64  `json:"missCount"`
	LookupCount uint64  `json:"lookupCount"`
}

// GetStats returns the stats of the LRU cache.
func (c *LRU) GetStats() *CacheStats {
	return &CacheStats{
		Size:        c.instance.Len(true),
		HitCount:    c.instance.HitCount(),
		HitRate:     c.instance.HitRate(),
		MissCount:   c.instance.MissCount(),
		LookupCount: c.instance.LookupCount(),
	}
}
