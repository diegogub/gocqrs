package usecases

import (
	"bitbucket.org/dgub/evento/dom"
	"errors"
)

var (
	EventNotFound = errors.New("Event not found")
)

type EventRequest struct {
	StreamId string    `json:"streamid"`
	Event    dom.Event `json:"event"`
	Streams  []string  `json:"streams,omitempty"`

	OptimisticLock   bool   `json:"lock"`
	ExpectedVersion  uint64 `json:"expectedVersion"`
	MustCreateStream bool   `json:"mustCreate"`
}

func ReadEvent(stream string, version uint64) (*dom.Event, error) {
	//TODO cache
	event, err := man.EventRepo.Read(stream, version)
	if err != nil {
		return nil, EventNotFound
	}
	return event, nil
}

func RangeQuery(stream string, from, to uint64) ([]dom.Event, error) {
	list := make([]dom.Event, 0)
	if from > to {
		return list, errors.New("invalid range")
	}

	if !StreamExist(stream) {
		return list, errors.New("stream not found")
	}

	// include last
	to++

	chEvent := man.EventRepo.Scan(stream, from, to)
	for e := range chEvent {
		list = append(list, e)
	}

	return list, nil
}
