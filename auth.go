package gocqrs

import (
	"github.com/diegogub/lib"
	"time"
)

type Sessioner interface {
	Save(s *Session) error
	Valid(id string) (*Session, error)
}

type Session struct {
	ID         string                 `json:"id"`
	Username   string                 `json:"un"`
	Role       string                 `json:"role"`
	ValidUntil time.Time              `json:"ttl"`
	Data       map[string]interface{} `json:"session"`
}

func NewSession(u *User, validity string) *Session {
	var s Session
	s.ID = lib.NewLongId("S")
	s.Username = u.Username
	s.Role = u.Role

	d, err := time.ParseDuration(validity)
	if err != nil {
		d = time.Duration(time.Minute * 10)
	}

	s.ValidUntil = time.Now().UTC().Add(d)

	return &s
}
