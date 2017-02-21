package usecases

import (
	"bitbucket.org/dgub/evento/dom"
	"log"
	"sync"
)

var pubMan *PubManager

type PubManager struct {
	lock     sync.RWMutex
	pubs     map[string]dom.Pub
	chEvents chan *dom.Event
}

func NewPubManager() *PubManager {
	var pMan PubManager
	pMan.pubs = make(map[string]dom.Pub, 1000)
	pMan.chEvents = make(chan *dom.Event, 10000)

	ps := ListPub()
	for _, pub := range ps {
		pMan.Add(pub)
	}

	return &pMan
}

func (pMan *PubManager) Publish() {
	for {
		select {
		case event := <-pMan.chEvents:
			for _, p := range pMan.pubs {
				log.Println(p.Publish(event.StreamId, *event))
				if dev {
				}
			}
		}
	}
}

func (pMan *PubManager) Add(pub dom.Pub) {
	pMan.lock.Lock()
	defer pMan.lock.Unlock()
	pMan.pubs[pub.Id] = pub
}

func (pMan *PubManager) Delete(id string) {
	pMan.lock.Lock()
	defer pMan.lock.Unlock()
	delete(pMan.pubs, id)
}

func CreatePub(id, regex, url string) error {
	pub := dom.Pub{
		Id:    id,
		Regex: regex,
		Url:   url,
	}

	err := pub.Validate()
	if err != nil {
		return err
	}

	err = man.PubRepo.Store(pub)
	if err == nil {
		go func() {
			pubMan.Add(pub)
			if dev {
				log.Println("added to pub cache:", pubMan)
			}
		}()
	}
	return err
}

func ListPub() []dom.Pub {
	return man.PubRepo.List()
}

func DeletePub(id string) error {
	err := man.PubRepo.Delete(id)
	if err == nil {

	}
	return err
}
