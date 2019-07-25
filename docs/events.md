# Events (Under development)

GoCacheable has a feature that allows modules to subscribe events. This is useful to delete or update cache when something happens.

Suppose you application has two modules A and B. B depends on A and when you update A, it is supposed to update result on B module, which are cached.
With this feature, B will subscribe to an event and A will throw that event when update is done. B will receive the event and update the cached value.

## Event Manager

**CacheEventsManager** is responsible for managing all events inside a **CacheableManager**.

The event manager contains a list of modules which are registered to listen/send events and a list of all event types available. The manager is responsible to handle all event subscriptions and to emit events to subscribers.

## Event Subscriber

**CacheEventSubscriber** is the event subscriber class. It is responsible to listen/read from the channel events are sent to. 


## Add new event

Adding a new event is quite simple, only being necessary to implement the **CacheEvent** interface.

```
// CacheEvent interface that represents an event that disputes an action over the cache storage
type CacheEvent interface {
	Invoke() error
}
```

## How to use

1. Create event manager

```
eventsManager := events.CreateEventsManager()
```

2. Register module

```
eventsManager.RegisterModule("new_module")
```

3. Subscribe module to event

```
err := eventsManager.SubscribeEvent("new_module", "new_event")
```

4. Emit event to all subscribers

```
eventsManager.EmitEvent("new_event", &EventProvider{})
```