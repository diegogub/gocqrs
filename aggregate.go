package gocqrs

type Aggregator interface {
	Aggregate(id string, entity interface{})
}

type EventHandler interface {
	EventName() []string
	Handle(id, userid, role string, e Eventer, entity *Entity, replay bool) (StoreOptions, error)
	CheckBase(e Eventer) bool
}
