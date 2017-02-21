package dom

import (
	"errors"
)

type ServerRepository interface {
	Init(update bool, id string, version string) error
	Get() (*Server, error)
	Update(version string) error
}

type Server struct {
	Id      string `json:"id"`
	Version string `json:"version"`
}

func (s *Server) Validate() error {
	if s.Id == "" {
		return errors.New("Invalid server id")
	}

	return nil
}
