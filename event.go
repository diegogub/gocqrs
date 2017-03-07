package gocqrs

import (
	"github.com/diegogub/lib"
	"time"
)

type Eventer interface {
	GetId() string
	GetStream() string
	GetVersion() uint64
	GetType() string
	GetData() map[string]interface{}
	GetLinks() []string
}

type BaseEvent struct {
	EventID        string    `json:"eid"`
	EventTimestamp time.Time `json:"ets"`
	EventType      string    `json:"ety"`
	EventVersion   uint64    `json:"eve"`
}

type Event struct {
	BaseEvent
	Entity            string                 `json:"entity"`
	CorrelationStream string                 `json:"cid"`
	EntityID          string                 `json:"eid"`
	StreamPrefix      string                 `json:"sprefix"`
	EventData         map[string]interface{} `json:"ed"`
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

func (e *Event) GetLinks() []string {
	return []string{e.CorrelationStream}
}

func (e *Event) GetStream() string {
	return e.Entity + "-" + e.EntityID
}

func (e *Event) GetType() string {
	return e.EventType
}

func (e *Event) GetVersion() uint64 {
	return e.EventVersion
}
