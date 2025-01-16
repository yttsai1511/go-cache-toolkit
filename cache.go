package cache

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	Second = time.Second
	Minute = 60 * Second
	Hour   = 60 * Minute
	Day    = 24 * Hour
)

type CacheItem struct {
	value   interface{}
	expires time.Time
}

type AtomicTime struct {
	timestamp int64
}

func (at *AtomicTime) Get() time.Time {
	return time.Unix(atomic.LoadInt64(&at.timestamp), 0)
}

func (at *AtomicTime) Set(t time.Time) {
	atomic.StoreInt64(&at.timestamp, t.Unix())
}

var caches sync.Map
var refreshTime AtomicTime
var defaultTTL time.Duration = Minute

// Set adds a new item to the cache or updates an existing one with a default TTL.
func Set(value interface{}, keys ...interface{}) {
	SetWithTTL(value, defaultTTL, keys...)
}

// SetWithTTL adds a new item to the cache or updates an existing one with a custom TTL.
func SetWithTTL(value interface{}, ttl time.Duration, keys ...interface{}) {
	key := GenerateKey(keys...)
	caches.Store(key, &CacheItem{
		value:   value,
		expires: time.Now().Add(ttl),
	})
}

// SetDefaultTTL sets the default TTL for cached items.
func SetDefaultTTL(ttl time.Duration) {
	defaultTTL = ttl
}

// SetRefreshTime sets the next scheduled refresh time for the cache.
func SetRefreshTime(targetTime time.Time) {
	refreshTime.Set(targetTime)
}

// Get gets an item from the cache.
func Get[T any](keys ...interface{}) (*T, error) {
	key := GenerateKey(keys...)
	value, ok := caches.Load(key)
	if !ok {
		return nil, fmt.Errorf("Item not found for key: %v", key)
	}

	item, ok := value.(*CacheItem)
	if !ok {
		caches.Delete(key)
		return nil, fmt.Errorf("Failed to assert cache item type for key: %v", key)
	}

	result, ok := item.value.(T)
	if !ok {
		caches.Delete(key)
		return nil, fmt.Errorf("Expected type %T but got %T for key: %v", result, item.value, key)
	}

	now := time.Now()
	if now.After(item.expires) {
		caches.Delete(key)
		return nil, fmt.Errorf("Item for key %v has expired", key)
	}

	go CheckRefresh(now)

	return &result, nil
}

// Delete removes an item from the cache.
func Delete(keys ...interface{}) {
	key := GenerateKey(keys...)
	caches.Delete(key)
}

// CheckRefresh checks if a refresh of the cache is needed based on the given targetTime.
func CheckRefresh(targetTime time.Time) {
	if targetTime.Before(refreshTime.Get()) {
		return
	}

	SetRefreshTime(targetTime.Add(Hour))
	CleanExpired(targetTime)
}

// CleanExpired iterates over all cache items and deletes those that have expired based on the given targetTime.
func CleanExpired(targetTime time.Time) {
	caches.Range(func(k, v interface{}) bool {
		item, ok := v.(*CacheItem)
		if !ok {
			caches.Delete(k)
			return true
		}

		if targetTime.After(item.expires) {
			caches.Delete(k)
		}
		return true
	})
}

// GenerateKey creates a unique key by concatenating the provided parts
func GenerateKey(keys ...interface{}) string {
	parts := make([]string, len(keys))
	for k, v := range keys {
		parts[k] = fmt.Sprintf("%v", v)
	}
	return strings.Join(parts, "|")
}
