package gocqrs

type Aggregator interface {
	Aggregate(id string, entity interface{})
}

type EventHandler interface {
	EventName() []string
	Handle(id string, e Eventer, entity *Entity, replay bool) (StoreOptions, error)
}
