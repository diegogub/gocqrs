package gocqrs

import (
	"encoding/json"
	"log"
	"net/url"
)

type APPDocs struct {
	Name      string              `json:"name"`
	Version   string              `json:"version"`
	Entities  map[string][]string `json:"entities"`
	Endpoints []Endpoint          `json:"endpoints"`
}

type Endpoint struct {
	Path        string `json:"path"`
	Method      string `json:"method"`
	Description string `json:"description,omitempty"`
	ExampleBody string `json:"exampleBody,omitempty"`
}

func NewEndpoint(path, method string, exampleBody interface{}) Endpoint {
	var e Endpoint

	_, err := url.Parse(path)
	if err != nil {
		log.Fatal("Invalid path for doc endpoint:", path)
	}

	e.Path = path
	e.Method = method

	if exampleBody != nil {
		b, _ := json.Marshal(exampleBody)
		e.ExampleBody = string(b)
	}

	return e
}

func (ad APPDocs) GetEvents(e string) []string {
	for entity, events := range ad.Entities {
		if entity == e {
			return events
		}
	}
	return []string{}
}

func GenerateDocs(app *App) APPDocs {
	return app.GenDocs()
}

func (app *App) AddEndpoint(e Endpoint) {
	app.Endpoints = append(app.Endpoints, e)
}

func (app *App) GenDocs() APPDocs {
	var docs APPDocs
	docs.Entities = make(map[string][]string)
	docs.Name = app.Name
	docs.Version = app.Version

	for e, c := range app.Entities {
		docs.Entities[e] = []string{}
		for event, _ := range c.EventHandlers {
			docs.Entities[e] = append(docs.Entities[e], event)
		}
	}

	docs.Endpoints = app.Endpoints

	return docs
}
