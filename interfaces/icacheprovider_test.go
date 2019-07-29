package interfaces

import (
	"errors"
	"os"
	"testing"
	
	"github.com/stretchr/testify/assert"
)

var errorMsg = "Could not add"

/***
Implement CacheableInterface
**/
type CacheableImplementation struct{}

func (cacheableImplementation *CacheableImplementation) Init() error {
	return errors.New(errorMsg)
}

func (cacheableImplementation *CacheableImplementation) Set(key string, value []byte) error {
	return errors.New(errorMsg)
}

func (cacheableImplementation *CacheableImplementation) Get(key string) ([]byte, error) {
	return []byte(""), nil
}

func (cacheableImplementation *CacheableImplementation) Delete(key string) error {
	return nil
}

func (cacheableImplementation *CacheableImplementation) Reset() error {
	return nil
}

func setup() {}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func exampleCall(interf CacheProviderInterface) error {
	return interf.Set("", []byte(""))
}

func TestCacheableInterfaceImplementation(t *testing.T) {
	// Create object
	implem := &CacheableImplementation{}
	// Test pass object as generic interface param
	err := exampleCall(implem)

	// if it succeeds, err must not be null and error message must equal errorMsg
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), errorMsg)
}
