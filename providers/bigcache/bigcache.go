package bigcache

import (
	"time"

	"github.com/allegro/bigcache"
)

// BigCacheProvider is a storage provider based on bigcache caching system
type BigCacheProvider struct {
	cacheStorage *bigcache.BigCache
	Lifetime     time.Duration
}

// Init initializes bigcache storage
func (bigcacheProvider *BigCacheProvider) Init() error {
	configuration := bigcache.Config{
		Shards:             1024,
		LifeWindow:         bigcacheProvider.Lifetime * time.Minute,
		CleanWindow:        bigcacheProvider.Lifetime * time.Minute,
		MaxEntriesInWindow: 1000 * 10 * 60,
		Verbose:            true,
		HardMaxCacheSize:   8192,
	}

	bigcacheProvider.cacheStorage, _ = bigcache.NewBigCache(configuration)
	return nil
}

// Set adds a new value to cache or updates if it already exists
func (bigcacheProvider *BigCacheProvider) Set(key string, value []byte) error {
	return bigcacheProvider.cacheStorage.Set(key, value)
}

// Get returns a cached value or error if it does not exist
func (bigcacheProvider *BigCacheProvider) Get(key string) ([]byte, error) {
	return bigcacheProvider.cacheStorage.Get(key)
	// Output: Hello, nil
}

// Delete removes a value from the cache
func (bigcacheProvider *BigCacheProvider) Delete(key string) error {
	return bigcacheProvider.cacheStorage.Delete(key)
}

// Reset empties cache storage
func (bigcacheProvider *BigCacheProvider) Reset() error {
	return bigcacheProvider.cacheStorage.Reset()
}

// HasKey checks if the key exists
func (bigcacheProvider *BigCacheProvider) HasKey(key string) bool {
	_, err := bigcacheProvider.cacheStorage.Get(key)
	return err == nil
}
