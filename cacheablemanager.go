package gocacheable

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/josemiguelmelo/gocacheable/events"

	gcCacheModule "github.com/josemiguelmelo/gocacheable/cachemodule"
	gcInterfaces "github.com/josemiguelmelo/gocacheable/interfaces"
)

// CacheableManager is responsible to manage cache storage
type CacheableManager struct {
	Identifier    string
	modules       []gcCacheModule.CacheModule
	EventsManager events.CacheEventsManager
}

// NewCacheableManager Create new CacheableManager object without modules
func NewCacheableManager(identifier string) CacheableManager {
	return CacheableManager{
		Identifier: identifier,
		modules:    []gcCacheModule.CacheModule{},
	}
}

// ModulesCount returns number of modules
func (cs *CacheableManager) ModulesCount() int {
	return len(cs.modules)
}

// ContainsModule verifies if manager contains the module with identifier=:identifier
func (cs *CacheableManager) ContainsModule(module gcCacheModule.CacheModule) bool {
	for _, m := range cs.modules {
		if m.Identifier == module.Identifier {
			return true
		}
	}
	return false
}

// AddModule adds a new module if it still does not exists
func (cs *CacheableManager) AddModule(name string, storageProvider gcInterfaces.CacheProviderInterface, cacheLife time.Duration) error {
	err := storageProvider.Init(cacheLife)
	if err != nil {
		return err
	}

	module := gcCacheModule.New(name, storageProvider)

	if cs.ContainsModule(module) {
		return errors.New("Module already exists")
	}

	cs.modules = append(cs.modules, module)
	return nil
}

// FindModule finds a module by its identifier
func (cs *CacheableManager) FindModule(identifier string) (*gcCacheModule.CacheModule, error) {
	for _, m := range cs.modules {
		if m.Identifier == identifier {
			return &m, nil
		}
	}
	return &gcCacheModule.CacheModule{}, errors.New("Module not found")
}


// DeleteKey removes a key from a module
func (cs *CacheableManager) Get(moduleID string, key string, out interface{}) error {
	module, err := cs.FindModule(moduleID)
	if err != nil {
		return err
	}

	return module.Get(key, &out)
}

// DeleteKey removes a key from a module
func (cs *CacheableManager) DeleteKey(moduleID string, key string) error {
	module, err := cs.FindModule(moduleID)
	if err != nil {
		return err
	}

	return module.Delete(key)
}

// Reset resets a module cache
func (cs *CacheableManager) Reset(moduleID string) error {
	module, err := cs.FindModule(moduleID)
	if err != nil {
		return err
	}

	return module.Reset()
}

// Cacheable adds cache to the function passed as parameter
func (cs *CacheableManager) Cacheable(moduleID string, key string, f func() (interface{}, error), out interface{}) error {
	module, err := cs.FindModule(moduleID)
	if err != nil {
		return err
	}

	if module.IsCacheStorageCreated() {
		return errors.New("Cache storage not created")
	}

	var obj interface{}
	// Check on cache and return if found
	err = module.Get(key, &out)
	if err == nil {
		return err
	}

	obj, err = f()
	if err == nil {
		module.Set(key, obj)
		jsonData, _ := json.Marshal(obj)
		err := json.Unmarshal(jsonData, &out)
		if err != nil {
			return err
		}
	}

	return err
}
