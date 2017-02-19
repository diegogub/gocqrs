package gocqrs

type Validator interface {
	Name() string
	Validate(e Entity) error
}
