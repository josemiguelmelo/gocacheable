package interfaces

import "time"

type CacheInterface interface {
	Init(cacheLife time.Duration) error
	Set(key string, value []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
	Reset() error
}

