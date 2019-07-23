package channel

import (
	gei "github.com/josemiguelmelo/gocacheable/events/interfaces"
)

// CacheEventChannel to handle events that will affect cache
type CacheEventChannel struct {
	Identifier       string
	Channel          chan gei.CacheEvent
	SubscribedEvents []string
}

// IsSubscribedTo verify if is subscribed to event
func (cc *CacheEventChannel) IsSubscribedTo(event string) bool {
	for _, e := range cc.SubscribedEvents {
		if e == event {
			return true
		}
	}
	return false
}
