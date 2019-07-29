package interfaces

// CacheProviderInterface interface to implement new cache provider
type CacheProviderInterface interface {
	Init() error
	Set(key string, value []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
	Reset() error
}
