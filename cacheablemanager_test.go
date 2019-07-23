package gocacheable

import (
	"errors"
	"os"
	"testing"

	gcCacheModule "github.com/josemiguelmelo/gocacheable/cachemodule"
	bcProvider "github.com/josemiguelmelo/gocacheable/providers/bigcache"
	"github.com/stretchr/testify/assert"
)

const (
	identifier = "testing_manager"
	moduleName = "testing_module"
)

var cacheableManager CacheableManager

func CreateCacheableManager() CacheableManager {
	return NewCacheableManager(identifier)
}

func setup() {
	cacheableManager = CreateCacheableManager()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func TestManagerCreatedSuccessfully(t *testing.T) {
	// Check manager does not contain any module
	assert.Equal(t, 0, cacheableManager.ModulesCount())
	assert.Equal(t, identifier, cacheableManager.Identifier)
}

func TestAddModuleToManager(t *testing.T) {
	storageProvider := &bcProvider.BigCacheProvider{}

	err := cacheableManager.AddModule(moduleName, storageProvider, 2)

	assert.Nil(t, err)
	assert.Equal(t, 1, cacheableManager.ModulesCount())
}

func TestAddAlreadyExistingModuleToManager(t *testing.T) {
	storageProvider := &bcProvider.BigCacheProvider{}
	err := cacheableManager.AddModule(moduleName, storageProvider, 2)

	assert.NotNil(t, err)
	assert.Equal(t, "Module already exists", err.Error())
	assert.Equal(t, 1, cacheableManager.ModulesCount())
}

func TestManagerFindModule(t *testing.T) {
	// Find a module that exists
	module, err := cacheableManager.FindModule(moduleName)
	assert.Nil(t, err)
	assert.Equal(t, moduleName, module.Identifier)

	// Check if module exists
	containsModule := cacheableManager.ContainsModule(*module)
	assert.Equal(t, true, containsModule)

	// Find a module that does not exist
	module, err = cacheableManager.FindModule("not_found")
	assert.NotNil(t, err)
	assert.Equal(t, "", module.Identifier)

	// Check if a non existent module is present
	containsModule = cacheableManager.ContainsModule(gcCacheModule.CacheModule{})
	assert.Equal(t, false, containsModule)
}

func methodToCache() int {
	return 9
}

func secondMethodToCache() int {
	panic(errors.New("Should not be here"))
	return 10
}
func TestCacheableMethod(t *testing.T) {
	var outValue int
	err := cacheableManager.Cacheable(moduleName, "test", func() (interface{}, error) {
		return methodToCache(), nil
	}, &outValue)

	assert.Nil(t, err)
	assert.Equal(t, 9, outValue)

	err = cacheableManager.Cacheable(moduleName, "test", func() (interface{}, error) {
		return secondMethodToCache(), nil
	}, &outValue)

	assert.Nil(t, err)
	assert.Equal(t, 9, outValue)

	err = cacheableManager.Cacheable("moduleName", "test", func() (interface{}, error) {
		return secondMethodToCache(), nil
	}, &outValue)

	assert.NotNil(t, err)
	assert.Equal(t, "Module not found", err.Error())
}

type ExampleObj struct {
	St       string
	notShown string
}

func TestCacheableMethodWithObject(t *testing.T) {
	var outValue ExampleObj
	err := cacheableManager.Cacheable(moduleName, "test_object", func() (interface{}, error) {
		return ExampleObj{St: "yes", notShown: "here"}, nil
	}, &outValue)

	assert.Nil(t, err)
	assert.Equal(t, "yes", outValue.St)
	// Is not serializable as it is private, so is not cached.
	assert.Equal(t, "", outValue.notShown)

	err = cacheableManager.Cacheable(moduleName, "test_object", func() (interface{}, error) {
		return ExampleObj{St: "no", notShown: "here"}, nil
	}, &outValue)

	assert.Nil(t, err)
	assert.Equal(t, "yes", outValue.St)
	// Is not serializable as it is private, so is not cached.
	assert.Equal(t, "", outValue.notShown)
}
