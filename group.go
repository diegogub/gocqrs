package gocqrs

import ()

const (
	GroupEntity = "groups"
)

type Group struct {
	BasicEntity
	ID   string            `json:"id"`
	Data map[string]string `json:"data"`
}
