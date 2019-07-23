package events

import (
	"errors"

	"github.com/josemiguelmelo/gocacheable/events/channel"
	gei "github.com/josemiguelmelo/gocacheable/events/interfaces"
)

// CacheEventsManager manages cache events
type CacheEventsManager struct {
	modules map[string]*channel.CacheEventChannel
	events  []string
}

// NewCacheEventsManager create new channel event manager
func NewCacheEventsManager() CacheEventsManager {
	return CacheEventsManager{
		modules: map[string]*channel.CacheEventChannel{},
		events:  []string{},
	}
}

// RegisterModule register a new module to listen for events
func (cm *CacheEventsManager) RegisterModule(moduleName string) error {
	if cm.ContainsModule(moduleName) {
		return errors.New("Module already exists")
	}
	cm.modules[moduleName] = &channel.CacheEventChannel{
		Identifier:       moduleName,
		Channel:          make(chan gei.CacheEvent),
		SubscribedEvents: []string{},
	}
	return nil
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
func (cm *CacheEventsManager) SubscribeEvent(moduleName string, event string) error {
	if !cm.ContainsModule(moduleName) {
		return errors.New("Module not found")
	}

	if !cm.ContainsEvent(event) {
		return errors.New("Event not found")
	}

	if cm.modules[moduleName].IsSubscribedTo(event) {
		return errors.New("Module already subscribed to this event")
	}

	moduleEvents := cm.modules[moduleName].SubscribedEvents
	cm.modules[moduleName].SubscribedEvents = append(moduleEvents, event)
	return nil
}

// SubscribedEvents returns list of subscribed event by module name
func (cm *CacheEventsManager) SubscribedEvents(moduleName string) ([]string, error) {
	if !cm.ContainsModule(moduleName) {
		return nil, errors.New("Module not found")
	}

	return cm.modules[moduleName].SubscribedEvents, nil
}

// RegisterEvent register a new event type if it does not exist
func (cm *CacheEventsManager) RegisterEvent(event string) error {
	if cm.ContainsEvent(event) {
		return errors.New("Event already exists")
	}

	cm.events = append(cm.events, event)
	return nil
}

// ContainsEvent returns true if event already exists
func (cm *CacheEventsManager) ContainsEvent(event string) bool {
	for _, e := range cm.events {
		if e == event {
			return true
		}
	}
	return false
}

// EventsCount returns the number of event types available
func (cm *CacheEventsManager) EventsCount() int {
	return len(cm.events)
}
