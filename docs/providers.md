## Providers

A provider is the interface with the caching technology (Redis, in memory, ...). 

Currently, the providers available are:

1) BigCache (https://github.com/allegro/bigcache)


### Create a new provider

Creating a new provider is quite simple. Providers must implement an interface **CacheProviderInterface**, which contains some methods required on a caching system. To create a new provider, it is only required to implement this interface.

```
// CacheProviderInterface interface to implement new cache provider
type CacheProviderInterface interface {
	Init(cacheLife time.Duration) error
	Set(key string, value []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
	Reset() error
}
```

