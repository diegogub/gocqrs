package api

import (
	"bitbucket.org/dgub/evento/dom"
	use "bitbucket.org/dgub/evento/usecases"
)

type EventApplier interface {
	GetVersion(streamid string) uint64
	// returns error if something fail
	// true if stream is in sync or false to resync
	Apply(e *dom.Event) (uint64, bool, error)
	PurgeReadModel() error
}

// All events should implement eventer interface to be stored into evento
type Eventer interface {
	GetId() string
	GetStream() string
	GetType() string
	GetData() map[string]interface{}
	GetLinks() []string
}

func buildRequest(ev Eventer, create, sync, lock bool, expversion uint64) *use.EventRequest {
	event := buildEvent(ev)
	er := use.EventRequest{
		StreamId:         ev.GetStream(),
		MustCreateStream: create,
		OptimisticLock:   lock,
		Event:            *event,
		Streams:          ev.GetLinks(),
	}
	return &er
}

func buildEvent(ev Eventer) *dom.Event {
	var event dom.Event
	event.Id = ev.GetId()
	event.Type = ev.GetType()
	event.LinkStreams = ev.GetLinks()

	data := ev.GetData()
	if data == nil {
		data = make(map[string]interface{})
	}
	event.Data = data
	return &event
}
