package gocqrs

type Aggregator interface {
	Aggregate(id string, entity interface{})
}

type EventHandler interface {
	EventName() []string
	Handle(e Eventer, entity *Entity) (StoreOptions, error)
}
