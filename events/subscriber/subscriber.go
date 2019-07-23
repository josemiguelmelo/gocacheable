package queue

import (
	"fmt"

	gei "github.com/josemiguelmelo/gocacheable/events/interfaces"
)

// CacheEventSubscriber to handle events that will affect cache
type CacheEventSubscriber struct {
	channel *chan gei.CacheEvent
}

// Subscribe subscribes to a channel
func (subs *CacheEventSubscriber) Subscribe(channel *chan gei.CacheEvent) {
	subs.channel = channel
	go subs.listenForChannel()
}

func (subs *CacheEventSubscriber) listenForChannel() {
	fmt.Println("LISTEN")
	for {
		event, ok := <-(*subs.channel)
		// if channel was closed
		if ok == false {
			fmt.Print(subs.channel)
			fmt.Println(" channel was closed")
			return
		}

		event.Invoke()
	}
}
