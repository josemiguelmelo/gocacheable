package storages

import (
	"github.com/allegro/bigcache"
	"time"
)

type BigCacheCache struct {
	cacheStorage *bigcache.BigCache
}

func (bigcacheCache *BigCacheCache) Init(life time.Duration) error {
	configuration := bigcache.Config{
		Shards:             1024,
		LifeWindow:         life * time.Minute,
		CleanWindow:        life * time.Minute,
		MaxEntriesInWindow: 1000 * 10 * 60,
		Verbose:            true,
		HardMaxCacheSize:   8192,
	}

	bigcacheCache.cacheStorage, _ = bigcache.NewBigCache(configuration)
	return nil
}

func (bigcacheCache *BigCacheCache) Set(key string, value []byte) error {
	return bigcacheCache.cacheStorage.Set(key, value)
}

func (bigcacheCache *BigCacheCache) Get(key string) ([]byte, error) {
	return bigcacheCache.cacheStorage.Get(key)
}

func (bigcacheCache *BigCacheCache) Delete(key string) error {
	return bigcacheCache.cacheStorage.Delete(key)
}

func (bigcacheCache *BigCacheCache) Reset() error {
	return bigcacheCache.cacheStorage.Reset()	
}

