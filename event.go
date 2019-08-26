package gocqrs

import (
	"encoding/json"
	"github.com/diegogub/lib"
	"time"
)

type Eventer interface {
	GetId() string
	GetStream() string
	GetVersion() uint64
	GetType() string
	GetData() map[string]interface{}
	SetData(k string, i interface{})
	ClearData()
	GetLinks() []string
	Has(string) bool
}

type BaseEvent struct {
	EventID        string    `json:"eid"`
	EventTimestamp time.Time `json:"timestamp"`
	EventStream    string    `json:"stream,omitempty"`
	EventType      string    `json:"type"`
	EventVersion   uint64    `json:"version"`
}

type Event struct {
	BaseEvent
	Group             string                 `json:"group,omitempty"`
	Entity            string                 `json:"ent,omitepty"`
	CorrelationStream string                 `json:"cid,omitempty"`
	EntityID          string                 `json:"id,omitempty"`
	StreamPrefix      string                 `json:"streamPre,omitempty"`
	EventData         map[string]interface{} `json:"data,omitempty"`
}

func NewEvent(id, t string, data map[string]interface{}) *Event {
	var e Event
	if id == "" {
		e.EventID = lib.NewShortId("")
	} else {
		e.EventID = id
	}

	e.EventData = data
	e.EventType = t
	e.EventTimestamp = time.Now().UTC()
	return &e
}

func (e *Event) GetId() string {
	return e.EventID
}

func (e *Event) GetData() map[string]interface{} {
	return e.EventData
}

func (e *Event) Has(k string) bool {
	if e.EventData == nil {
		return false
	} else {
		_, has := e.EventData[k]
		return has
	}
}

func (e *Event) GetLinks() []string {
	return []string{e.CorrelationStream}
}

func (e *Event) GetStream() string {
	if e.Group == "" {
		return e.Entity + "-" + e.EntityID
	} else {
		return e.Group + "-" + e.Entity + "-" + e.EntityID
	}
}

func (e *Event) GetType() string {
	return e.EventType
}

func (e *Event) GetVersion() uint64 {
	return e.EventVersion
}

func (e *Event) SetStructData(i interface{}) error {
	e.EventData[k] = i
}

func (e *Event) SetData(k string, i interface{}) {
	e.EventData[k] = i
}

func (e *Event) ClearData() {
	e.EventData = make(map[string]interface{})
}

func DecodeEvent(e Eventer, i interface{}) error {
	b, err := json.Marshal(e.GetData())
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &i)
	if err != nil {
		return err
	}

	return err
}

func (event *Event) String() string {
	b, _ := json.Marshal(event)
	return string(b)
}
