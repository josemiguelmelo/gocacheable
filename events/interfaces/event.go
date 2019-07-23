package events

// CacheEvent interface that represents an event that disputes an action over the cache storage
type CacheEvent interface {
	Invoke() error
}

// RemoveCacheEvent represent an event that removes a key from the cache storage
type RemoveCacheEvent struct{}

// Invoke receives the key to be removed and removes it from the cache storage
func (removeCacheEvent *RemoveCacheEvent) Invoke(key string) error {
	return nil
}
