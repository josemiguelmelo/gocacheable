package gocacheable

import (
	"encoding/json"
	"errors"
	"time"

	gcInterfaces "github.com/josemiguelmelo/gocacheable/interfaces"
)

// CacheableManager is responsible to manage cache storage
type CacheableManager struct {
	cacheStorage gcInterfaces.CacheProviderInterface
}

// New Create new CacheableManager object
func New(cacheStorage gcInterfaces.CacheProviderInterface, cacheLife time.Duration) CacheableManager {
	cacheStorage.Init(cacheLife)
	return CacheableManager{
		cacheStorage: cacheStorage,
	}
}

// Get returns a cached value
func (cs *CacheableManager) Get(key string, out interface{}) error {
	var obj interface{}
	err := getFromCache(cs.cacheStorage, key, &obj)
	if err == nil {
		if jsonData, err := json.Marshal(obj); err == nil {
			json.Unmarshal(jsonData, &out)
			return err
		}
		out = obj
		return err
	}
	return err
}

// Set caches a value
func (cs *CacheableManager) Set(key string, value interface{}) {
	addToCache(cs.cacheStorage, key, value)
}

// Delete removes a key from the cache storage
func (cs *CacheableManager) Delete(key string) error {
	return cs.cacheStorage.Delete(key)
}

// Reset empties the cache storage
func (cs *CacheableManager) Reset() error {
	return cs.cacheStorage.Reset()
}

// CacheableManager adds cache to the function passed as parameter
func (cs *CacheableManager) Cacheable(key string, f func() (interface{}, error), out interface{}) error {
	if cs.cacheStorage == nil {
		return errors.New("Cache storage not created")
	}

	var obj interface{}
	// Check on cache and return if found
	err := cs.Get(key, &out)
	if err == nil {
		return err
	}

	obj, err = f()
	if err == nil {
		cs.Set(key, obj)
		jsonData, _ := json.Marshal(obj)
		json.Unmarshal(jsonData, &out)
	}

	return err
}

func getFromCache(cacheStorage gcInterfaces.CacheProviderInterface, key string, out interface{}) error {
	valueByte, err := cacheStorage.Get(key)
	if err == nil {
		return json.Unmarshal(valueByte, &out)
	}
	return err
}

func addToCache(cacheStorage gcInterfaces.CacheProviderInterface, key string, value interface{}) {
	valueByte, _ := json.Marshal(value)
	cacheStorage.Set(key, valueByte)
}
