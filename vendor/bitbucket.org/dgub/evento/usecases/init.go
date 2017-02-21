package usecases

import (
	"bitbucket.org/dgub/evento/dom"
	"log"
	"os"
	"os/signal"
	"time"
)

type Manager struct {
	EventRepo  dom.EventRepository
	StreamRepo dom.StreamRepository
	PubRepo    dom.PubRepository
	SrvRepo    dom.ServerRepository
	//StreamerRepo StreamerRepo
}

var man *Manager
var dev bool
var PurgeEvery time.Duration

func Init(first, update bool, id, version string, eventRepo dom.EventRepository, streamRepo dom.StreamRepository, pubRepo dom.PubRepository, serverRepo dom.ServerRepository, d bool) {
	var m Manager

	m.EventRepo = eventRepo
	m.StreamRepo = streamRepo
	m.PubRepo = pubRepo
	m.SrvRepo = serverRepo
	man = &m

	dev = d
	pm := NewPubManager()
	pubMan = pm

	go pm.Publish()

	// Start database
	EventoDb = &EventoDB{}
	EventoDb.Start(first, update, id, version)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for sig := range c {
		log.Println("Got ", sig)
		EventoDb.Stop()
		break
	}
	log.Println("test")
}
