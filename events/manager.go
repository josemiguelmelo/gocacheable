package events

import (
	"errors"

	"github.com/josemiguelmelo/gocacheable/events/channel"
	gei "github.com/josemiguelmelo/gocacheable/events/interfaces"
	subscriber "github.com/josemiguelmelo/gocacheable/events/subscriber"
)

// CacheEventsManager manages cache events
type CacheEventsManager struct {
	modules    map[string]*channel.CacheEventChannel
	eventTypes []string
}

// NewCacheEventsManager create new channel event manager
func NewCacheEventsManager() CacheEventsManager {
	return CacheEventsManager{
		modules:    map[string]*channel.CacheEventChannel{},
		eventTypes: []string{},
	}
}

// RegisterModule register a new module to listen for events
func (cm *CacheEventsManager) RegisterModule(moduleName string) (*chan gei.CacheEvent, error) {
	if cm.ContainsModule(moduleName) {
		return nil, errors.New("Module already exists")
	}
	cm.modules[moduleName] = &channel.CacheEventChannel{
		Identifier:       moduleName,
		Channel:          make(chan gei.CacheEvent),
		SubscribedEvents: []string{},
	}
	return &cm.modules[moduleName].Channel, nil
}

// ContainsModule returns true if module already exists
func (cm *CacheEventsManager) ContainsModule(moduleName string) bool {
	if _, ok := cm.modules[moduleName]; ok {
		return true
	}
	return false
}

// ModuleCount returns the number of registered modules
func (cm *CacheEventsManager) ModuleCount() int {
	return len(cm.modules)
}

// SubscribeEvent subscribes a module to a event
func (cm *CacheEventsManager) SubscribeEvent(moduleName string, eventType string, callback func(gei.CacheEvent)) (subscriber.CacheEventSubscriber, error) {
	if !cm.ContainsModule(moduleName) {
		return subscriber.CacheEventSubscriber{}, errors.New("Module not found")
	}

	if !cm.ContainsEventType(eventType) {
		return subscriber.CacheEventSubscriber{}, errors.New("Event not found")
	}

	if cm.modules[moduleName].IsSubscribedTo(eventType) {
		return subscriber.CacheEventSubscriber{}, errors.New("Module already subscribed to this event")
	}

	moduleEvents := cm.modules[moduleName].SubscribedEvents
	cm.modules[moduleName].SubscribedEvents = append(moduleEvents, eventType)

	eventSubscriber := subscriber.CacheEventSubscriber{}
	eventSubscriber.Subscribe(&cm.modules[moduleName].Channel, callback)
	return eventSubscriber, nil
}

// SubscribedEvents returns list of subscribed event by module name
func (cm *CacheEventsManager) SubscribedEvents(moduleName string) ([]string, error) {
	if !cm.ContainsModule(moduleName) {
		return nil, errors.New("Module not found")
	}

	return cm.modules[moduleName].SubscribedEvents, nil
}

// RegisterEvent register a new event type if it does not exist
func (cm *CacheEventsManager) RegisterEvent(eventType string) error {
	if cm.ContainsEventType(eventType) {
		return errors.New("Event already exists")
	}

	cm.eventTypes = append(cm.eventTypes, eventType)
	return nil
}

// ContainsEventType returns true if event already exists
func (cm *CacheEventsManager) ContainsEventType(eventType string) bool {
	for _, e := range cm.eventTypes {
		if e == eventType {
			return true
		}
	}
	return false
}

// EventTypesCount returns the number of event types available
func (cm *CacheEventsManager) EventTypesCount() int {
	return len(cm.eventTypes)
}

// EmitEvent emits the event to all subscribers
func (cm *CacheEventsManager) EmitEvent(eventType string, event gei.CacheEvent) {
	for _, module := range cm.modules {
		if module.IsSubscribedTo(eventType) {
			module.Channel <- event
		}
	}
}
