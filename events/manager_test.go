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
	assert.Equal(t, 0, eventsManager.EventTypesCount())
}

func TestRegisterEvent(t *testing.T) {
	err := eventsManager.RegisterEvent("new_event")
	assert.Nil(t, err)
	assert.Equal(t, 1, eventsManager.EventTypesCount())
	assert.Equal(t, true, eventsManager.ContainsEventType("new_event"))

	err = eventsManager.RegisterEvent("new_event")
	assert.NotNil(t, err)
	assert.Equal(t, 1, eventsManager.EventTypesCount())
	assert.Equal(t, true, eventsManager.ContainsEventType("new_event"))
}

func TestRegisterModule(t *testing.T) {
	_, err := eventsManager.RegisterModule("new_module")
	assert.Nil(t, err)
	assert.Equal(t, 1, eventsManager.ModuleCount())
	assert.Equal(t, true, eventsManager.ContainsModule("new_module"))

	_, err = eventsManager.RegisterModule("new_module")
	assert.NotNil(t, err)
	assert.Equal(t, 1, eventsManager.ModuleCount())
}

func TestSubscribeEvent(t *testing.T) {
	// must be successful
	_, err := eventsManager.SubscribeEvent("new_module", "new_event")
	assert.Nil(t, err)
	events, err := eventsManager.SubscribedEvents("new_module")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(events))

	// must return error
	_, err = eventsManager.SubscribeEvent("new_module", "not_existing_event")
	assert.NotNil(t, err)

	// must return error
	_, err = eventsManager.SubscribeEvent("not_existing_module", "new_event")
	assert.NotNil(t, err)

	_, err = eventsManager.SubscribedEvents("not_existing_module")
	assert.NotNil(t, err)
}

type EventProvider struct {
}

var invokeCalled = false

func (ep *EventProvider) Invoke() error {
	invokeCalled = true
	return nil
}

func TestEmitEvent(t *testing.T) {
	eventsManager.RegisterModule("event_module")
	eventsManager.RegisterEvent("new_event")
	// Module Subscribe event
	eventsManager.SubscribeEvent("event_module", "new_event")

	// Invoke called variable must be false before emit event
	assert.False(t, invokeCalled)

	eventsManager.EmitEvent("new_event", &EventProvider{})

	// Invoke called variable must be true after emit event
	assert.True(t, invokeCalled)
}
