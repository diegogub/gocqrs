package dom

import (
	"encoding/json"
	"errors"
	"github.com/diegogub/lib"
	"github.com/rs/xid"
	"log"
)

type EventStatus int

var (
	InvalidEventPkg       = errors.New("Invalid event pkg")
	WrongStreamVersionPkg = errors.New("Invalid version in stream compared with events")
	NoCorrelativeEventPkg = errors.New("Invalid correlativity in event pkg")
)

const (
	Stored EventStatus = iota
	StreamGroupNotFound
	SqlError
	VersionError
)

type EventRepository interface {
	Push(pkgs []*EventPkg, sync bool) error
	Read(streamid string, version uint64) (*Event, error)
	Scan(streamid string, fromVersion, toVersion uint64) chan Event
	Exist(id string) (string, uint64, bool)
	Lock() error
	Unlock() error
	//ReadSnapshot(streamid string) (Event, error)
}

type Event struct {
	Id string `json:"id"`

	StreamId    string   `json:"streamid,omitempty"`
	LinkStreams []string `json:"streams,omitempty"`
	NewStream   bool     `json:"-"`
	Version     uint64   `json:"version"`
	Type        string   `json:"type"`
	Link        bool     `json:"link,omitempty"`
	LinkStream  string   `json:"lstream,omitempty"`
	LinkVersion uint64   `json:"vlink"`

	Data map[string]interface{} `json:"data"`
}

type SEvent struct {
	Id       string                 `json:"i"`
	StreamId string                 `json:"s,omitempty"`
	Version  uint64                 `json:"v"`
	Type     string                 `json:"t,omitempty"`
	Data     map[string]interface{} `json:"d"`
}

type SLinkEvent struct {
	Id          string                 `json:"i"`
	StreamId    string                 `json:"s,omitempty"`
	Version     uint64                 `json:"v"`
	Type        string                 `json:"t,omitempty"`
	Link        uint                   `json:"l,omitempty"`
	LinkStream  string                 `json:"ls,omitempty"`
	LinkVersion uint64                 `json:"lv,omitempty"`
	Data        map[string]interface{} `json:"d"`
}

func (e SLinkEvent) ToEvent() Event {
	event := Event{
		Id:          e.Id,
		StreamId:    e.StreamId,
		Version:     e.Version,
		Type:        e.Type,
		LinkStream:  e.LinkStream,
		LinkVersion: e.LinkVersion,
		Data:        e.Data,
	}
	if e.Link == 1 {
		event.Link = true
	}
	return event
}

func (e *Event) ToSmall() interface{} {
	if e.Link {
		return SLinkEvent{
			Id:          e.Id,
			StreamId:    e.StreamId,
			Version:     e.Version,
			Type:        e.Type,
			Link:        1,
			LinkStream:  e.LinkStream,
			LinkVersion: e.LinkVersion,
			Data:        e.Data,
		}
	} else {
		return SEvent{
			Id:       e.Id,
			StreamId: e.StreamId,
			Version:  e.Version,
			Type:     e.Type,
			Data:     e.Data,
		}
	}
}

func NewEvenLink(stream *Stream, event *Event) *Event {
	var e Event
	e.StreamId = stream.Id
	e.Version = stream.Version

	e.Link = true
	e.LinkStream = event.StreamId
	e.LinkVersion = event.Version
	e.Type = event.Type
	return &e
}

func NewEvent(stream, eventType string, data map[string]interface{}) *Event {
	var e Event
	e.StreamId = stream
	e.Data = data
	e.Type = eventType
	return &e
}

func (e *Event) Validate() error {
	if e.Id == "" {
		e.Id = xid.New().String()
	}

	if e.Type == "" {
		return errors.New("Invalid event type")
	}

	if e.Data == nil {
		e.Data = make(map[string]interface{})
	}

	return nil
}

func (e Event) Decode(i interface{}) error {
	b, err := e.Marshal()
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, i)
	if err != nil {
		return err
	}

	return nil
}

func (e Event) Marshal() ([]byte, error) {
	b, err := json.Marshal(e.Data)
	if err != nil {

		return []byte("{}"), err
	}

	return b, nil
}

type EventPkg struct {
	Stream *Stream
	Events []*Event
}

func NewEventPkg() *EventPkg {
	var e EventPkg
	e.Events = make([]*Event, 0)
	return &e
}

func (ep EventPkg) Validate() error {
	//var curVersion uint64
	var lastVersion uint64
	var err error

	if len(ep.Events) == 0 {
		return InvalidEventPkg
	}

	//Last version check
	lastVersion = ep.Events[len(ep.Events)-1].Version
	if ep.Stream.Version != lastVersion {
		log.Println(ep.Stream.Version, lastVersion)
		return WrongStreamVersionPkg
	}

	// check correlativity
	for i, e := range ep.Events {
		if i != 0 {
			if e.Version != ep.Events[i-1].Version+1 {
				return NoCorrelativeEventPkg
			}
		}
	}

	return err
}

func (e Event) String() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func (e Event) GetId() string {
	return ""
}

func (e Event) GetLinks() []string {
	return []string{}
}

func (e Event) GetType() string {
	return e.Type
}

func (e Event) GetStream() string {
	return e.StreamId
}

func (e Event) GetData() map[string]interface{} {
	return lib.Map(e.Data)
}

func (e Event) Create() bool {
	return e.NewStream
}
