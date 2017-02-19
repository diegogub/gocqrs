package gocqrs

import (
	"log"
	"strings"
)

const (
	Created = "Created"
	Updated = "Updated"
	Deleted = "Deleted"
)

var eventsNames = []string{"Created",
	"Updated",
	"Deleted",
}

type CRUDHandler struct {
	EntityName string `json:"entityName"`
}

func NewCRUDHandler(name string) CRUDHandler {
	var ch CRUDHandler
	if name == "" {
		log.Fatal("Invalid entity name to create CRUD handler")
		return ch
	}

	ch.EntityName = name
	return ch
}

// Handler CRUD events
func (ch CRUDHandler) Handle(ev Eventer, en *Entity) (StoreOptions, error) {
	var opt StoreOptions
	var err error

	switch ev.GetType() {
	case ch.CreateEvent():
		en.Version = 0
		en.ID = ev.GetId()
		en.Data = ev.GetData()
		opt.Create = true
	case ch.UpdateEvent():
		data := ev.GetData()
		log.Println(data)
		for k, d := range data {
			en.Data[k] = d
		}
	case ch.DeletedEvent():
	}
	return opt, err
}

func (ch CRUDHandler) EventName() []string {
	events := make([]string, 0)
	for _, p := range eventsNames {
		e := strings.Title(ch.EntityName) + p
		events = append(events, e)
	}

	return events
}

func (ch CRUDHandler) CreateEvent() string {
	return strings.Title(ch.EntityName) + "Created"
}

func (ch CRUDHandler) UpdateEvent() string {
	return strings.Title(ch.EntityName) + "Updated"
}

func (ch CRUDHandler) DeletedEvent() string {
	return strings.Title(ch.EntityName) + "Deleted"
}
