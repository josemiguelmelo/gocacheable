package redis

import "github.com/gomodule/redigo/redis"

const (
	// ACTION_SET redis set key value action
	ACTION_SET = "SET"
	// ACTION_GET redis get key value action
	ACTION_GET = "GET"
	// ACTION_DELETE redis delete key value action
	ACTION_DELETE = "DEL"
	// ACTION_RESET redis reset cache action
	ACTION_RESET = "FLUSHDB"
	// ACTION_PING redis ping action
	ACTION_PING = "PING"
)

func newRedisPool(addr string, maxIdle int, maxActive int) *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: maxIdle,
		// max number of connections
		MaxActive: maxActive,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr)
		},
	}
}

// NewRedisProvider returns a RedisProvider
func NewRedisProvider(addr string, maxIdle int, maxActive int) (*RedisProvider, error) {
	redisProvider := &RedisProvider{
		Addr:      addr,
		MaxIdle:   maxIdle,
		MaxActive: maxActive,
	}

	err := redisProvider.Init()
	if err != nil {
		return nil, err
	}

	err = redisProvider.Ping()
	if err != nil {
		return nil, err
	}

	return redisProvider, nil
}

// RedisProvider is a storage provider based on redis caching system
type RedisProvider struct {
	redisPool *redis.Pool
	Addr      string
	MaxIdle   int
	MaxActive int
}

// Init initializes redis storage
func (redisProvider *RedisProvider) Init() error {
	redisProvider.redisPool = newRedisPool(redisProvider.Addr, redisProvider.MaxIdle, redisProvider.MaxActive)
	return nil
}

// Set adds a new value to cache or updates if it already exists
func (redisProvider *RedisProvider) Set(key string, value []byte) error {
	conn := redisProvider.redisPool.Get()
	defer conn.Close()

	_, err := conn.Do(ACTION_SET, key, string(value))
	return err
}

// Get returns a cached value or error if it does not exist
func (redisProvider *RedisProvider) Get(key string) ([]byte, error) {
	conn := redisProvider.redisPool.Get()
	defer conn.Close()

	s, err := redis.String(conn.Do(ACTION_GET, key))
	if err != nil {
		return nil, err
	}
	return []byte(s), nil

}

// Delete removes a value from the cache
func (redisProvider *RedisProvider) Delete(key string) error {
	conn := redisProvider.redisPool.Get()
	defer conn.Close()

	_, err := conn.Do(ACTION_DELETE, key)
	return err
}

// Reset empties cache storage
func (redisProvider *RedisProvider) Reset() error {
	conn := redisProvider.redisPool.Get()
	defer conn.Close()

	_, err := conn.Do(ACTION_RESET)
	return err
}

// Ping checks redis connection
func (redisProvider *RedisProvider) Ping() error {
	conn := redisProvider.redisPool.Get()
	defer conn.Close()

	_, err := redis.String(conn.Do(ACTION_PING))
	return err
}
