package gocqrs

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
)

type EntityConf struct {
	Name string `json:"name"`

	CRUD bool `json:"crud"`

	Description string `json:"desc"`
	// basic entity prefix
	StreamPrefix      string `json:"stream_prefix"`
	CorrelationStream string `json:"correlation_stream"`

	// Entity ID references
	EntityReferences []EntityReference `json:"references"`

	// Event Handlers/Aggregators
	EventHandlers map[string]EventHandler `json:"handlers"`

	Validators map[string]Validator `json:"validators"`
}

type EntityReference struct {
	Key  string `json:"key"`
	Null bool   `json:"null"`
}

func NewEntityConf(name string) *EntityConf {
	var e EntityConf
	name = strings.ToLower(name)
	e.Name = name
	e.Validators = make(map[string]Validator)
	e.EventHandlers = make(map[string]EventHandler)
	return &e
}

func (e *EntityConf) AddEventHandler(eh ...EventHandler) error {
	var err error
	if e.EventHandlers == nil {
		e.EventHandlers = make(map[string]EventHandler)
	}

	for _, h := range eh {
		for _, event := range h.EventName() {
			_, exist := e.EventHandlers[event]
			if exist {
				err = errors.New("Could not add event handler, already set:" + event)
			} else {
				e.EventHandlers[event] = h
			}
		}
	}

	return err
}

func (e *EntityConf) AddCRUD() *EntityConf {
	e.CRUD = true
	// TODO add basic CRUD
	ch := NewCRUDHandler(e.Name)
	e.AddEventHandler(ch)
	return e
}

func (ec *EntityConf) Aggregate(id string, events chan Eventer) (*Entity, error) {
	var err error
	var entity Entity
	entity.ID = id
	log.Println("---------->", entity)
	for e := range events {
		log.Println("events:", e)
		eventHandler, has := ec.EventHandlers[e.GetType()]
		if !has {
			return &entity, errors.New("Event " + e.GetType() + " not handled")
		}

		if e.GetVersion() == entity.Version+1 || e.GetVersion() == 0 {
			_, err = eventHandler.Handle(e, &entity)
		} else {
			return &entity, errors.New("Failed to aggregate entity " + id + " , unorder events")
		}
	}

	log.Println("--->", entity)

	return &entity, err
}

type Entity struct {
	ID      string                 `json:"id"`
	Version uint64                 `json:"version"`
	Data    map[string]interface{} `json:"data"`
}

func (e *Entity) Decode(i interface{}) error {
	b, err := json.Marshal(e.Data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &i)
	if err != nil {
		return err
	}

	return err
}
