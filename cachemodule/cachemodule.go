package cachemodule

import (
	"encoding/json"
	"strings"

	gcInterfaces "github.com/josemiguelmelo/gocacheable/interfaces"
)

// CacheModule represents an applicational module that contains cache
type CacheModule struct {
	Identifier   string
	Name         string
	cacheStorage gcInterfaces.CacheProviderInterface
}

// New create and returns new CacheModule object
func New(name string, cacheStorage gcInterfaces.CacheProviderInterface) CacheModule {
	return CacheModule{
		Identifier:   generateIdentifier(name),
		Name:         name,
		cacheStorage: cacheStorage,
	}
}

func generateIdentifier(name string) string {
	identifier := strings.ReplaceAll(name, " ", "_")
	return strings.ToLower(identifier)
}

// IsCacheStorageCreated returns true if module cache storage is already created
func (cm CacheModule) IsCacheStorageCreated() bool {
	return cm.cacheStorage == nil
}

// Get returns a cached value
func (cm CacheModule) Get(key string, out interface{}) error {
	var obj interface{}
	err := getFromCache(cm.cacheStorage, key, &obj)
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
func (cm *CacheModule) Set(key string, value interface{}) error {
	return addToCache(cm.cacheStorage, key, value)
}

// Delete removes a key from the cache storage
func (cm *CacheModule) Delete(key string) error {
	return cm.cacheStorage.Delete(key)
}

// Reset empties the cache storage
func (cm *CacheModule) Reset() error {
	return cm.cacheStorage.Reset()
}

func getFromCache(cacheStorage gcInterfaces.CacheProviderInterface, key string, out interface{}) error {
	valueByte, err := cacheStorage.Get(key)
	if err == nil {
		return json.Unmarshal(valueByte, &out)
	}
	return err
}

func addToCache(cacheStorage gcInterfaces.CacheProviderInterface, key string, value interface{}) error {
	valueByte, _ := json.Marshal(value)
	err := cacheStorage.Set(key, valueByte)
	if  err != nil {
		return err
	}
	return nil
}
