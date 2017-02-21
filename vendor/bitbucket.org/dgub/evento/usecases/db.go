package usecases

import (
	"bitbucket.org/dgub/evento/dom"
	"bitbucket.org/dgub/evento/proxy/use"
	"gopkg.in/h2non/gentleman.v1"
	"gopkg.in/h2non/gentleman.v1/plugins/body"
	"os"

	"errors"
	"log"
	"sync"
	"time"
)

var (
	StreamNotCreated     = errors.New("stream was not created, already exist")
	StreamNotFound       = errors.New("stream do not exist")
	InvalidStreamVersion = errors.New("invalid stream version")
	IdEventExist         = errors.New("Event exist")
)

const MainStream = "#all"

var (
	Proxy     bool
	ProxyHost string
	EventoURL string
)

type EventoDB struct {
	lock    sync.RWMutex
	running bool `json:"running"`
	stop    chan bool

	ID   string `json:"name"`
	Host string `json:"host"`

	Streams uint64 `json:"streams"`

	// Main stream data
	Main *dom.Stream

	//Cache all used streams
	StreamCache map[string]*dom.Stream

	Pubs       *PubManager
	FlushCache time.Duration
}

var EventoDb *EventoDB

func NewEventoDB(id, host string, flush time.Duration) *EventoDB {
	var edb EventoDB
	return &edb
}

func DbRunning() bool {
	return EventoDb.running
}

func (edb *EventoDB) Start(first, update bool, id, host string) {
	edb.lock.Lock()
	defer edb.lock.Unlock()
	var s *dom.Server
	var err error
	var proxyHost use.EventoHost

	edb.ID = id
	edb.Host = host
	edb.stop = make(chan bool, 1)
	edb.StreamCache = make(map[string]*dom.Stream)
	edb.FlushCache = time.Second

	// Create main stream
	edb.findStream(MainStream, true)
	if first {
		err := man.SrvRepo.Init(update, id, "")
		if err != nil {
			log.Fatal(err)
		}
		s, err = man.SrvRepo.Get()
		if err != nil {
			log.Fatal("Failed to init, no id set: try to init with --first --id --ip")
		}
		log.Println(s)
	} else {
		s, err = man.SrvRepo.Get()
		if err != nil {
			log.Fatal("Failed to init, no id set: try to init with --first --id --ip")
		}
		log.Println(s)
	}

	if Proxy {
		if s.Id == "" {
			proxyHost.ID = id
		} else {
			proxyHost.ID = s.Id
		}

		if EventoURL == "" {
			n, _ := os.Hostname()
			EventoURL = "http://" + n + ":6060"
		}
		log.Println("EventoHost:", EventoURL)

		proxyHost.URL = EventoURL
		c := gentleman.New()
		c.URL(ProxyHost)
		go func(h use.EventoHost) {
			for {
				req := c.Request()
				req.Method("POST")
				req.Path("/reg/host")
				req.Use(body.JSON(proxyHost))
				log.Println("Ping proxy..")
				res, err := req.Send()
				if err != nil {
					log.Println("Failed to ping: ", err.Error())
				} else {
					log.Println("Response: ", res)
				}
				time.Sleep(5 * time.Second)
			}
		}(proxyHost)
	}

	// Starting purging every...
	go func() {
		for {
			log.Println("purging..")
			pstatus := GetPurges()
			for _, ps := range pstatus {
				log.Println(ps)
				man.StreamRepo.DoPurge(ps)
			}
			log.Println("sleeping...")
			time.Sleep(PurgeEvery)
		}
	}()

	log.Println("Starting evento db.....")
	edb.Pubs = NewPubManager()
	edb.running = true
	go edb.Pubs.Publish()
}

func (edb EventoDB) Stop() {
	edb.lock.Lock()
	defer edb.lock.Unlock()
	edb.running = false
}

// StoreEvent stores events into eventostore
func (edb *EventoDB) StoreEvent(event dom.Event, streams []string, lock bool, expectedVersion uint64, create bool) (id string, version uint64, err error) {
	edb.lock.Lock()
	defer edb.lock.Unlock()

	err = event.Validate()
	if err != nil {
		return id, version, err
	}

	// check if ID exist
	_, v, exist := man.EventRepo.Exist(event.Id)
	if exist {
		id = event.Id
		version = v
		err = IdEventExist
		return id, version, err
	}

	// check locks and create
	if lock {
		err = edb.checkLocks(event.StreamId, expectedVersion, create)
		if err != nil {
			return id, version, err
		}
	}

	// get stream from db or cache
	stream, created, err := edb.findStream(event.StreamId, create)
	if err != nil {
		return id, version, err
	}

	if stream.Deleted == dom.StreamDeleted {
		return id, version, errors.New("stream deleted")
	}

	// check if it was created
	if create {
		if !created {
			err = StreamNotCreated
			return id, version, err
		}
	}

	if err == nil {
		// increment version
		if created {
			event.Version = 0
		} else {
			event.Version = stream.Version + 1
			stream.Version = stream.Version + 1
		}

		pkgs := make([]*dom.EventPkg, 0)

		pkg := dom.NewEventPkg()
		pkg.Stream = stream
		pkg.Events = append(pkg.Events, &event)

		pkgs = append(pkgs, pkg)

		// link all streams
		for _, s := range streams {
			s, created, _ := edb.findStream(s, create)
			if s.Deleted == dom.StreamDeleted {
				continue
			}

			if !created {
				s.Version = s.Version + 1
			}

			el := dom.NewEvenLink(s, &event)
			el.Validate()

			p := dom.NewEventPkg()
			p.Stream = s
			p.Events = append(p.Events, el)
			if err != nil {
				return id, version, err
			}
			pkgs = append(pkgs, p)
		}

		//push all events into store
		err = man.EventRepo.Push(pkgs, true)
	}

	if err == nil {
		id = event.Id
		version = event.Version

		edb.StreamCache[stream.Id] = stream

		edb.Pubs.chEvents <- &event
	}

	return id, version, err
}

func (edb EventoDB) checkLocks(id string, expVersion uint64, create bool) error {
	v, ok := edb.StreamCache[id]
	if ok {
		if create {
			return StreamNotCreated
		}

		if v.Version != expVersion {
			return InvalidStreamVersion
		}
	} else {
		if create {
			return nil
		} else {
			return StreamNotFound
		}
	}

	return nil
}

func (edb *EventoDB) findStream(streamid string, create bool) (*dom.Stream, bool, error) {
	// try to find stream version and create id
	var created bool
	//if create {
	s, _ := dom.NewStream(streamid, true, []string{}, false)
	err := man.StreamRepo.Create(*s)
	if err != nil {
		created = false
	} else {
		created = true
	}
	//}

	stream, err := man.StreamRepo.Get(streamid)
	if err == nil {
		edb.StreamCache[streamid] = stream
	} else {
		return stream, created, errors.New("stream not found")
	}
	return stream, created, nil
}
