package ap

import (
	"bitbucket.org/dgub/evento/dom"
	"bitbucket.org/dgub/evento/usecases"
)

type Streamer struct {
}

type EventResponse struct {
	Id      string
	Version uint64
}

func (s *Streamer) StoreEvent(er usecases.EventRequest, res *EventResponse) error {
	_, _, err := usecases.EventoDb.StoreEvent(er.Event, er.Streams, er.OptimisticLock, er.ExpectedVersion, er.MustCreateStream)
	if err != nil {
		return err
	}

	return nil
}

type EventRead struct {
	StreamId string
	Version  uint64
}

func (s *Streamer) ReadEvent(er EventRead, events *[]dom.Event) error {
	ev, err := usecases.ReadEvent(er.StreamId, er.Version)
	if err != nil {
		return err
	}

	*events = append(*events, *ev)

	return nil
}

type RangeEvent struct {
	StreamId string
	From     uint64
	To       uint64
}

func (s *Streamer) Range(er RangeEvent, events *[]dom.Event) error {
	ev, err := usecases.RangeQuery(er.StreamId, er.From, er.To)
	if err != nil {
		return err
	}

	for _, e := range ev {
		*events = append(*events, e)
	}

	return nil
}

func (s *Streamer) Version(sid string, version *uint64) error {
	v, err := usecases.StreamVersion(sid)
	if err != nil {
		return err
	}
	*version = v
	return nil
}

func (s *Streamer) ListStreams(regex string, streams *[]dom.Stream) error {
	list, err := usecases.StreamMatch(regex)
	if err != nil {
		return err
	}

	for _, s := range list {
		*streams = append(*streams, s)
	}

	return nil
}
