package events

import (
	"fmt"
	"os"
	"testing"
	interfaces "github.com/josemiguelmelo/gocacheable/events/interfaces"
	"github.com/stretchr/testify/assert"
)

const (
	identifier = "testing_manager"
	moduleName = "testing_module"
)

var invokeCalled = 0
var invokeRes = make(chan int)

type EventProvider struct {
}

func (ep *EventProvider) Invoke() error {
	fmt.Println("invoked")
	return nil
}

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
	_, err := eventsManager.SubscribeEvent(
		"new_module", 
		"new_event", 
		func(a interfaces.CacheEvent) {},
	)
	assert.Nil(t, err)
	events, err := eventsManager.SubscribedEvents("new_module")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(events))

	// must return error
	_, err = eventsManager.SubscribeEvent(
		"new_module", 
		"not_existing_event", 
		func(a  interfaces.CacheEvent) {},
	)
	assert.NotNil(t, err)

	// must return error
	_, err = eventsManager.SubscribeEvent(
		"not_existing_module", 
		"new_event",
		 func(a  interfaces.CacheEvent) {},
	)
	assert.NotNil(t, err)

	_, err = eventsManager.SubscribedEvents("not_existing_module")
	assert.NotNil(t, err)
}


func TestEmitEvent(t *testing.T) {
	// Register modules
	eventsManager.RegisterModule("event_module")
	eventsManager.RegisterModule("event_module_2")
	// Register events
	eventsManager.RegisterEvent("new_event")
	eventsManager.RegisterEvent("new_event_2")
	eventsManager.RegisterEvent("new_event_3")

	// Module Subscribe event
	eventsManager.SubscribeEvent(
		"event_module", 
		"new_event", 
		func(a interfaces.CacheEvent) {
			invokeCalled = invokeCalled + 1
			invokeRes <- invokeCalled
		},
	)
	eventsManager.SubscribeEvent(
		"event_module_2",
		"new_event_3", 
		func(a  interfaces.CacheEvent) {
			invokeCalled = invokeCalled + 1
			invokeRes <- invokeCalled
		},
	)
	eventsManager.SubscribeEvent(
		"event_module_2", 
		"new_event", 
		func(a  interfaces.CacheEvent) {
			invokeCalled = invokeCalled + 1
			invokeRes <- invokeCalled
		},
	)

	// Invoke called variable must be 0 before emit event
	assert.Equal(t, 0, invokeCalled)
	eventsManager.EmitEvent("new_event_2", &EventProvider{})
	// Invoke called variable must be 0 after emit event not subscribed
	assert.Equal(t, 0, invokeCalled)

	// emit event subscribed by 1 module, must increment invokeCalled in 1
	eventsManager.EmitEvent("new_event_3", &EventProvider{})
	// Invoke called variable must be 1 after emit event
	invokeCalled = <-invokeRes
	assert.Equal(t, 1, invokeCalled)

	// emit event subscribed by 1 module, must increment invokeCalled by 2
	eventsManager.EmitEvent("new_event", &EventProvider{})
	// Invoke called variable must be 3 after emit event
	invokeCalled = <-invokeRes
	invokeCalled = <-invokeRes
	assert.Equal(t, 3, invokeCalled)
}
