package bigcache

import (
	"os"
	"testing"
	"time"

	"github.com/allegro/bigcache"
	"github.com/stretchr/testify/assert"
)

var bigcacheStorage BigCacheProvider

const (
	initialExpectedValue = "test"
	updatedExpectedValue = "updated_test"
	cacheKey             = "Testing"
)

func NewCacheableStorage() BigCacheProvider {
	config := bigcache.Config{
		Shards:             1024,
		LifeWindow:         2 * time.Minute,
		CleanWindow:        2 * time.Minute,
		MaxEntriesInWindow: 1000 * 10 * 60,
		Verbose:            true,
		HardMaxCacheSize:   8192,
	}

	bigCacheStorage, _ := bigcache.NewBigCache(config)

	return BigCacheProvider{
		cacheStorage: bigCacheStorage,
	}
}

func setup() {
	bigcacheStorage = NewCacheableStorage()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func TestBigCacheInit(t *testing.T) {
	cacheProvider := BigCacheProvider{}
	// Cache storage still not initialized
	assert.Nil(t, cacheProvider.cacheStorage)

	cacheProvider.Init(2)

	// Cache storage already initialized
	assert.NotNil(t, cacheProvider.cacheStorage)
	// cache storage must be empty
	assert.Equal(t, 0, cacheProvider.cacheStorage.Len())
}

func TestBigCacheStorageAddAndGetMethods(t *testing.T) {
	// Add inital value to cache
	err := bigcacheStorage.Set(cacheKey, []byte(initialExpectedValue))
	assert.Nil(t, err)

	// Verify cached value equals to initia value
	value, _ := bigcacheStorage.Get(cacheKey)
	assert.Equal(t, initialExpectedValue, string(value))
}

func TestBigCacheStorageReplaceValue(t *testing.T) {
	// Verify cached value equals to initial value
	value, _ := bigcacheStorage.Get(cacheKey)
	assert.Equal(t, initialExpectedValue, string(value))

	// Update cached value to updatedExpectedValue
	err := bigcacheStorage.Set(cacheKey, []byte(updatedExpectedValue))
	assert.Nil(t, err)

	// Verify cached value equals to updated value
	value, _ = bigcacheStorage.Get(cacheKey)
	assert.Equal(t, updatedExpectedValue, string(value))
}

func TestBigCacheStorageDeleteValue(t *testing.T) {
	// Verify cached value exists
	_, err := bigcacheStorage.Get(cacheKey)
	assert.Nil(t, err)

	// Delete must return null
	err = bigcacheStorage.Delete(cacheKey)
	assert.Nil(t, err)

	// Verify cached value does not exist
	_, err = bigcacheStorage.Get(cacheKey)
	assert.NotNil(t, err)
}

func TestBigCacheStorageReset(t *testing.T) {
	// Add inital value to cache
	err := bigcacheStorage.Set(cacheKey, []byte(initialExpectedValue))
	assert.Nil(t, err)

	// Verify cached value exists
	_, err = bigcacheStorage.Get(cacheKey)
	assert.Nil(t, err)

	// Delete must return null
	err = bigcacheStorage.Reset()
	assert.Nil(t, err)

	// Verify cached value does not exist
	_, err = bigcacheStorage.Get(cacheKey)
	assert.NotNil(t, err)
}
