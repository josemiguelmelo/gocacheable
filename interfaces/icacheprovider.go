package interfaces

import "time"

// CacheProviderInterface interface to implement new cache provider
type CacheProviderInterface interface {
	Init(cacheLife time.Duration) error
	Set(key string, value []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
	Reset() error
}
