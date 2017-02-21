package api

import (
	"bitbucket.org/dgub/evento/dom"
	//	use "bitbucket.org/dgub/evento/usecases"
	"github.com/nsqio/go-nsq"
	"log"
	"regexp"
	"sync"
)

const (
	BOLT_BUCKET = "evento_keys"
)

type ReaderClient struct {
	eventoCli *Client
	Id        string
	wg        sync.WaitGroup
	lock      sync.RWMutex

	running     bool
	dev         bool
	clientType  string
	consumer    *nsq.Consumer
	stringRegex []string
	regexs      []*regexp.Regexp

	applier EventApplier

	chEvents chan *dom.Event
}

func NewReaderClient(copt ConnectOpt, id string, applier EventApplier) *ReaderClient {
	ecli := NewClient(copt)
	var cli ReaderClient
	cli.Id = id
	cli.eventoCli = ecli
	cli.chEvents = make(chan *dom.Event, 1000)
	cli.stringRegex = make([]string, 0)
	cli.regexs = make([]*regexp.Regexp, 0)
	cli.applier = applier
	return &cli
}

type ReaderOptions struct {
	// Replay all streams, will purge read model
	FullReplay bool `json:"replay"`
	// Avoid to catch-up to last version.
	NoCatchUp bool `json:"noCatchUp"`

	Nsqd string `json:"nsqd"`
}

func (cli *ReaderClient) Start(ropt ReaderOptions) {
	cli.lock.Lock()
	defer cli.lock.Unlock()
	if !ropt.NoCatchUp {
		cli.CatchUp()
	}

	go cli.applyEvents()

	err := cli.connectNSQ(ropt.Nsqd)
	if err != nil {
		panic(err)
	}

}

func (cli *ReaderClient) AddHandler(eventer EventApplier) {
	cli.lock.Lock()
	defer cli.lock.Unlock()
	cli.applier = eventer
}

func (cli *ReaderClient) applyEvents() {
	for {
		select {
		case event, ok := <-cli.chEvents:
			if !ok {
				log.Println("Closed event channel")
				return
			}

			currentVersion, sync, err := cli.applier.Apply(event)
			if err != nil {
				log.Println(err)
			}

			if !sync {
				err := cli.syncStream(event.StreamId, currentVersion+1, event.Version)
				if err != nil {
					log.Println(err)
				}
			}

			// done with event
			cli.wg.Done()
		}
	}
}

// Listen to streams matching regex
func (cli *ReaderClient) ListenStreams(regex string) error {
	cli.lock.Lock()
	defer cli.lock.Unlock()
	var err error

	r, err := regexp.Compile(regex)
	if err == nil {
		cli.stringRegex = append(cli.stringRegex, regex)
		cli.regexs = append(cli.regexs, r)
	}

	return err
}

func (cli *ReaderClient) StreamMatch(streamid string) bool {
	var match bool
	if len(cli.regexs) == 0 {
		return true
	}

	for _, r := range cli.regexs {
		match = r.Match([]byte(streamid))
	}
	return match
}

func (cli *ReaderClient) SetDev(dev bool) {
	cli.eventoCli.SetDevMode(dev)
}

func (cli *ReaderClient) CatchUp() {
	/*
		log.Println("Catching up events...")
		v, _ := cli.eventoCli.Version("_log")
		chStreams := cli.eventoCli.RangeStream("_log", 1, v)

		for e := range chStreams {
			var streamCreated use.StreamCreated

			switch e.Type {
			case use.EventoStreamCreated:
				e.Decode(&streamCreated)
				if !cli.StreamMatch(streamCreated.Name) {
					if DevMode {
						log.Println("no match..", streamCreated)
					}
					continue
				}

				serverVersion, _ := cli.eventoCli.Version(streamCreated.Name)
				modelVersion := cli.applier.GetVersion(streamCreated.Name)

				if DevMode {
					log.Println("Server version:", serverVersion, "Local version:", modelVersion)
				}

				if modelVersion != 0 {
					modelVersion++
				}

				if modelVersion < serverVersion {
					cli.syncStream(streamCreated.Name, modelVersion, serverVersion)
				}
			}
		}
	*/
}

func (cli *ReaderClient) syncStream(streamid string, from, to uint64) error {
	ch := cli.eventoCli.RangeStream(streamid, from, to)
	for e := range ch {
		_, _, err := cli.applier.Apply(&e)
		if err != nil {
			return err
		}
	}
	return nil
}
