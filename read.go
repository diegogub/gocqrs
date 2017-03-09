package gocqrs

type Reader interface {
	Sync(entity, id string) error
}
