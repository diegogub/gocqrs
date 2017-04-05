package gocqrs

import (
	"fmt"
	"log"
	"time"
)

const (
	RebuiltOpt = "rebuild"
	PurgeOpt   = "purge"
)

type View struct {
	Name string `json:"view"`

	All     bool `json:"all"`
	CatchUp bool `json:"catchup"`

	V       Viewer
	Store   EventStore
	wakeUP  chan bool
	running bool

	every time.Duration
}

type Viewer interface {
	Init(params map[string]string, dev bool)
	Purge() error
	Status() (uint64, error)
	Stream() string
	Rebuild() error
	Apply(event Event) error
}

func NewView(name string, i Viewer) *View {
	var v View
	v.Name = name
	v.wakeUP = make(chan bool, 2)
	v.V = i
	// set default wakeup time every 400ms
	v.every = time.Duration(time.Millisecond * 400)
	return &v
}

func (v *View) awake() {
	for {
		time.Sleep(v.every)
		log.Println("Slept for ", v.every, " ,waking up")
		v.wakeUP <- true
	}
}

func (v *View) Run(action string, params map[string]string, dev bool) error {
	var err error

	switch action {
	case RebuiltOpt:
		err := v.V.Purge()
		if err != nil {
			panic("Failed to purge view:" + err.Error())
		}
		// Rebuild view
		// init view
		v.V.Init(params, dev)
	case PurgeOpt:
		v.V.Init(params, dev)
		log.Println("Purging view...")
		err := v.V.Purge()
		if err != nil {
			panic("Failed to purge view:" + err.Error())
		}
		return nil
	default:
		// init view
		v.V.Init(params, dev)
	}

	go v.awake()
	mainStream := v.V.Stream()
	for {
		v.running = true
		fmt.Println("-------------------")
		curVersion, err := v.V.Status()
		if err != nil {
			log.Println("View not started, no status available")
		}
		esVersion, err := v.Store.Version(mainStream)
		if err != nil {
			log.Println("Failed to get current stream version, sleeping for few seconds")
			time.Sleep(time.Second * 5)
			continue
		}

		log.Println("current Version:", curVersion, " evento version:", esVersion)
		log.Println(mainStream)
		if esVersion > curVersion || curVersion == 0 {
			log.Println("Catching up from version ", curVersion, " to version ", esVersion)
			var events chan Event
			if curVersion == 0 {
				events = v.Store.Scan(mainStream, curVersion, esVersion)
			} else {
				events = v.Store.Scan(mainStream, curVersion+1, esVersion)
			}

			for e := range events {
				err := v.V.Apply(e)
				if err != nil {
					log.Println(err)
					log.Fatal("Failed to apply event :", e)
				}
			}
		}
		v.running = false
		fmt.Println("-------------------")

		select {
		case <-v.wakeUP:
		}
	}

	return err
}
