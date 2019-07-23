package events

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	identifier = "testing_manager"
	moduleName = "testing_module"
)

var eventsManager CacheEventsManager

func CreateEventsManager() CacheEventsManager {
	return NewCacheEventsManager()
}

func setup() {
	eventsManager = CreateEventsManager()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func TestManagerCreatedSuccessfully(t *testing.T) {
	// Check manager does not contain any module
	assert.Equal(t, 0, eventsManager.EventsCount())
}

func TestRegisterEvent(t *testing.T) {
	err := eventsManager.RegisterEvent("new_event")
	assert.Nil(t, err)
	assert.Equal(t, 1, eventsManager.EventsCount())
	assert.Equal(t, true, eventsManager.ContainsEvent("new_event"))

	err = eventsManager.RegisterEvent("new_event")
	assert.NotNil(t, err)
	assert.Equal(t, 1, eventsManager.EventsCount())
	assert.Equal(t, true, eventsManager.ContainsEvent("new_event"))
}

func TestRegisterModule(t *testing.T) {
	err := eventsManager.RegisterModule("new_module")
	assert.Nil(t, err)
	assert.Equal(t, 1, eventsManager.ModuleCount())
	assert.Equal(t, true, eventsManager.ContainsModule("new_module"))

	err = eventsManager.RegisterModule("new_module")
	assert.NotNil(t, err)
	assert.Equal(t, 1, eventsManager.ModuleCount())
}

func TestSubscribeEvent(t *testing.T) {
	// must be successful
	err := eventsManager.SubscribeEvent("new_module", "new_event")
	assert.Nil(t, err)
	events, err := eventsManager.SubscribedEvents("new_module")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(events))

	// must return error
	err = eventsManager.SubscribeEvent("new_module", "not_existing_event")
	assert.NotNil(t, err)

	// must return error
	err = eventsManager.SubscribeEvent("not_existing_module", "new_event")
	assert.NotNil(t, err)

	_, err = eventsManager.SubscribedEvents("not_existing_module")
	assert.NotNil(t, err)
}
