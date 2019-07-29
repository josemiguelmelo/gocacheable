package queue

import (
	"log"

	gei "github.com/josemiguelmelo/gocacheable/events/interfaces"
)

// CacheEventSubscriber to handle events that will affect cache
type CacheEventSubscriber struct {
	channel *chan gei.CacheEvent
	callback func(gei.CacheEvent)
}

// Subscribe subscribes to a channel with a callback
func (subs *CacheEventSubscriber) Subscribe(channel *chan gei.CacheEvent, callback func(gei.CacheEvent) ) {
	subs.channel = channel
	subs.callback = callback
	go subs.listenForChannel()
}

func (subs *CacheEventSubscriber) listenForChannel() {
	for {
		event, ok := <-(*subs.channel)
		// if channel was closed
		if !ok {
			log.Print(subs.channel)
			log.Println(" channel was closed")
			return
		}

		subs.callback(event)
	}
}
