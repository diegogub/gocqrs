package gocqrs

import ()

const (
	GroupEntity = "groups"
)

type Group struct {
	BasicEntity
	ID   string `json:"id"`
	Name string `json:"name"`
}
