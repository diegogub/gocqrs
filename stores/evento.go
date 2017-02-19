package stores

import (
	es "bitbucket.org/dgub/evento/api"
	"github.com/diegogub/gocqrs"
	"log"
)

type EventoStore struct {
	URL    string `json:"url"`
	Proxy  bool   `json:"proxy"`
	client *es.Client
}

func NewEventoStore(url string, proxy bool) *EventoStore {
	var e EventoStore
	var cli *es.Client

	if proxy {
		copt := es.ConnectOpt{
			EventoProxy: url,
		}
		cli = es.NewClient(copt)
	} else {
		copt := es.ConnectOpt{
			EventoServer: url,
		}
		cli = es.NewClient(copt)
		err := cli.Ping()
		if err != nil {
			log.Fatal(err)
		}
	}

	e.client = cli
	return &e
}

func (estore EventoStore) Store(e gocqrs.Eventer, opt gocqrs.StoreOptions) (uint64, error) {
	var v uint64
	var err error
	// TODO retry
	if opt.LockVersion > 0 {
		return estore.client.StoreEvent(e, &es.StoreOpt{Create: opt.Create, Lock: true, ExpectedVersion: opt.LockVersion})
	} else {
		return estore.client.StoreEvent(e, &es.StoreOpt{Create: opt.Create})
	}
	return v, err
}

func (es EventoStore) Range(streamid string) chan gocqrs.Eventer {
	ch := make(chan gocqrs.Eventer, 20)
	lastVersion, _ := es.Version(streamid)
	events := es.client.RangeStream(streamid, 0, lastVersion)
	go func() {
		for e := range events {
			ev := gocqrs.NewEvent(e.GetId(), e.GetType(), e.GetData())
			ch <- ev
		}
		close(ch)
	}()
	return ch
}

func (es EventoStore) Version(streamid string) (uint64, error) {
	return es.client.Version(streamid)
}
