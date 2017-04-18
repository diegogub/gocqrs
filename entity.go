package gocqrs

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"strings"
	"time"
)

const (
	AllRole = "all"
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
	BaseStruct interface{}          `json:"base,omitempty"`
	BaseSeted  bool

	Auther    ReadAuther
	ReadRoles []string `json:"roles,omitempty"`
}

type BasicEntity struct {
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type EntityReference struct {
	Entity string `json:"entity"`
	Key    string `json:"key"`
	Null   bool   `json:"null"`
}

func NewEntityConf(name string) *EntityConf {
	var e EntityConf
	name = strings.ToLower(name)
	e.Name = name
	e.Validators = make(map[string]Validator)
	e.EventHandlers = make(map[string]EventHandler)
	return &e
}

func (e *EntityConf) Reference(en, k string, null bool) {
	e.EntityReferences = append(e.EntityReferences, EntityReference{en, k, null})
}

func (e *EntityConf) AddValidator(v ...Validator) error {
	var err error

	for _, validator := range v {
		_, exist := e.Validators[validator.GetName()]
		if exist {
			log.Fatal("Could not add validator, already set:" + validator.GetName())
		} else {
			e.Validators[validator.GetName()] = validator
		}
	}
	return err
}

func (e *EntityConf) SetBaseStruct(i interface{}) {
	t := reflect.TypeOf(i)
	for t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Struct:
	default:
		panic("Invalid type to be base struct")
	}

	e.BaseStruct = i
	e.BaseSeted = true
}

func (e *EntityConf) AddEventHandler(eh ...EventHandler) error {
	var err error
	if e.EventHandlers == nil {
		e.EventHandlers = make(map[string]EventHandler)
	}

	for _, h := range eh {
		for _, event := range h.EventName() {
			// replace handler if needed
			e.EventHandlers[event] = h
		}
	}

	return err
}

func (e *EntityConf) AddCRUD(checkVersion bool) *EntityConf {
	e.CRUD = true
	// TODO add basic CRUD
	ch := NewCRUDHandler(e.Name)
	ch.CheckVersion = checkVersion

	e.AddEventHandler(ch)
	return e
}

func (ec *EntityConf) Aggregate(id string, events chan Eventer) (*Entity, error) {
	var err error
	var entity Entity
	entity.Data = make(map[string]interface{})
	for e := range events {
		eventHandler, has := ec.EventHandlers[e.GetType()]
		if !has {
			return &entity, errors.New("Event " + e.GetType() + " not handled")
		}

		if e.GetVersion() == entity.Version+1 || e.GetVersion() == 0 {
			_, err = eventHandler.Handle(id, e, &entity, true)
		} else {
			return &entity, errors.New("Failed to aggregate entity " + id + " , unorder events")
		}
	}
	entity.ID = id

	return &entity, err
}

func (ev EntityConf) checkBase(data map[string]interface{}) error {
	i := ev.BaseStruct
	t := reflect.TypeOf(i)
	for t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		t = t.Elem()
	}

	n := t.NumField()
	for k, _ := range data {
		has := false
		for f := 0; f < n; f++ {
			v, hasTag := t.Field(f).Tag.Lookup("json")
			if !hasTag {
				v = t.Field(f).Name
			}
			if v == k {
				has = true
				break
			}
		}

		if !has {
			return errors.New("Invalid field, do not exist in base struct: " + k)
		}
	}

	return nil
}

type Entity struct {
	ID      string                 `json:"id"`
	Version uint64                 `json:"version"`
	Deleted bool                   `json:"deleted"`
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

func ToMap(i interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	d, _ := json.Marshal(i)
	json.Unmarshal(d, &m)
	return m
}
