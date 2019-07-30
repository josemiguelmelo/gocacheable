package redis

import (
	"os"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
)

var redisProvider *RedisProvider
var redisServer *miniredis.Miniredis

const (
	existingKey      = "key"
	notExistingKey   = "not_key"
	existingKeyValue = "key_value"
)

func NewCacheableStorage() *RedisProvider {
	redisProvider, err := NewRedisProvider(redisServer.Addr(), 10, 100)
	if err != nil {
		return nil
	}
	return redisProvider
}

func setup() {
	var err error
	redisServer, err = miniredis.Run()
	if err != nil {
		panic(err)
	}
	redisProvider = NewCacheableStorage()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	redisServer.Close()
	os.Exit(code)
}

func TestInitRedis(t *testing.T) {
	redProv, err := NewRedisProvider(redisServer.Addr(), 10, 100)
	assert.Nil(t, err)
	assert.NotNil(t, redProv)

	err = redisProvider.Init()
	assert.Nil(t, err)

	err = redProv.Ping()
	assert.Nil(t, err)
}

func TestPingRedis(t *testing.T) {
	redProv, err := NewRedisProvider(":80", 10, 100)
	assert.NotNil(t, err)
	assert.Nil(t, redProv)
}

func TestSetAndGetCache(t *testing.T) {
	err := redisProvider.Set(existingKey, []byte("first_value"))
	assert.Nil(t, err)

	val, err := redisProvider.Get(existingKey)
	assert.Nil(t, err)
	assert.Equal(t, "first_value", string(val))

	// Update value
	err = redisProvider.Set(existingKey, []byte(existingKeyValue))
	assert.Nil(t, err)

	val, err = redisProvider.Get(existingKey)
	assert.Nil(t, err)
	assert.Equal(t, existingKeyValue, string(val))
}

func TestGetNotExistingKeyCache(t *testing.T) {
	_, err := redisProvider.Get(notExistingKey)
	assert.NotNil(t, err)
}

func TestDeleteKeyCache(t *testing.T) {
	// Delete not existing key must return error
	err := redisProvider.Delete(notExistingKey)
	assert.Nil(t, err)

	// Test existingKey exists
	val, err := redisProvider.Get(existingKey)
	assert.Nil(t, err)
	assert.Equal(t, existingKeyValue, string(val))

	// Delete existingKey
	err = redisProvider.Delete(existingKey)
	assert.Nil(t, err)

	// Test existingKey not exists anymore
	_, err = redisProvider.Get(existingKey)
	assert.NotNil(t, err)
}

func TestResetKeyCache(t *testing.T) {
	// Add keys
	err := redisProvider.Set(existingKey, []byte(existingKeyValue))
	assert.Nil(t, err)
	err = redisProvider.Set(notExistingKey, []byte(existingKeyValue))
	assert.Nil(t, err)

	// Test keys were added
	_, err = redisProvider.Get(existingKey)
	assert.Nil(t, err)
	_, err = redisProvider.Get(notExistingKey)
	assert.Nil(t, err)

	// Reset cache
	err = redisProvider.Reset()
	assert.Nil(t, err)

	// Test existingKey not exists anymore
	_, err = redisProvider.Get(existingKey)
	assert.NotNil(t, err)
	_, err = redisProvider.Get(notExistingKey)
	assert.NotNil(t, err)
}
