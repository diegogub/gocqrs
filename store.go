package gocqrs

import (
	"errors"
)

var (
	FailStoreError      = errors.New("Failed to store event, db issue")
	LockVersionError    = errors.New("Invalid lock version")
	StreamNotFoundError = errors.New("Stream don't exist")
)

// Eventstore interface
type EventStore interface {
	Store(e Eventer, opt StoreOptions) (uint64, error)
	Range(streamid string) (chan Eventer, uint64)
	Version(streamid string) (uint64, error)
	Scan(streamid string, from, to uint64) chan Event
}

// Storing event options
type StoreOptions struct {
	LockVersion uint64 `json:"lockversion"`
	Retry       bool   `json:"retry"`
	Create      bool   `json:"create"`
}
