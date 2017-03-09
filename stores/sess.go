package stores

import (
	"errors"
	"github.com/diegogub/gocqrs"
	"log"
	"sync"
	"time"
)

type MemSessions struct {
	lock     sync.Mutex
	sessions map[string]gocqrs.Session
}

func NewMemSessions() *MemSessions {
	var m MemSessions
	m.sessions = make(map[string]gocqrs.Session)

	go func(m *MemSessions) {
		for {
			m.Clear()
			log.Println("cleaning sessions..")
			time.Sleep(time.Minute * 60)
		}
	}(&m)

	return &m
}

func (m *MemSessions) Clear() {
	m.lock.Lock()
	defer m.lock.Unlock()

	for id, s := range m.sessions {
		if s.ValidUntil.After(time.Now().UTC()) {
			delete(m.sessions, id)
		}
	}
}

func (m *MemSessions) Save(s *gocqrs.Session) error {
	m.lock.Lock()
	m.sessions[s.ID] = *s
	m.lock.Unlock()
	return nil
}

func (m *MemSessions) Valid(id string) (*gocqrs.Session, error) {
	m.lock.Lock()
	s := m.sessions[id]
	if s.ValidUntil.After(time.Now().UTC()) {
		delete(m.sessions, id)
		return nil, errors.New("Invalid session")
	}
	m.lock.Unlock()
	return &s, nil
}
