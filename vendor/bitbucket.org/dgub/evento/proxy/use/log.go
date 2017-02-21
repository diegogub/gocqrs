package use

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"gopkg.in/gin-gonic/gin.v1"
	"net/url"
	"sort"
	"sync"
	"time"
)

var plog *ProxyLog

var (
	HostNotRegistered = errors.New("Host not registered")
	NoAvailableHost   = errors.New("No available host")
)

const LogStream = "plog"
const (
	StreamPrefix = "@s#"
	HostPrefix   = "@h#"
)

type ProxyLog struct {
	lock          sync.Mutex `json:"-"`
	db            *bolt.DB
	Status        uint64                `json:"status"`
	Hosts         map[string]EventoHost `json:"hosts"`
	StreamPerHost map[string]uint64     `json:"streamsPerHost"`
	// All stream log mapping
	Streams map[string]string `json:"streams"`
}

type EventoHost struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Closed bool   `json:"closed"`

	Down     bool      `json:"down"`
	LastPing time.Time `json:"lastPing"`
	Streams  uint64    `json:"streams"`
}

func (e *EventoHost) Validate() error {
	if e.ID == "" {
		return errors.New("Invalid evento ID")
	}

	if e.URL == "" {
		return errors.New("Invalid evento host")
	}

	_, err := url.Parse(e.URL)
	if err != nil {
		return err
	}

	return nil
}

type ByStream []EventoHost

func (a ByStream) Len() int           { return len(a) }
func (a ByStream) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByStream) Less(i, j int) bool { return a[i].Streams < a[j].Streams }

func NewProxyLog(db *bolt.DB) *ProxyLog {
	var pl ProxyLog
	pl.db = db
	pl.Hosts = make(map[string]EventoHost)
	pl.Streams = make(map[string]string)
	pl.StreamPerHost = make(map[string]uint64)
	return &pl
}

func (l *ProxyLog) AddHost(h *EventoHost) *EventoHost {
	var resultHost EventoHost
	plog.db.Update(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte("streams"))
		if err != nil {
			return fmt.Errorf("no exist bucket: %s", err)
		}

		hdata := b.Get([]byte(HostPrefix + h.ID))
		if len(hdata) == 0 {
			hostData, err := json.Marshal(h)
			if err != nil {
				return err
			}
			b.Put([]byte(HostPrefix+h.ID), hostData)
			return nil
		} else {
			var host EventoHost
			err = json.Unmarshal(hdata, &host)
			if err != nil {
				return err
			}
			host.URL = h.URL
			host.LastPing = time.Now().UTC()
			host.Closed = h.Closed
			resultHost = host
			hostData, err := json.Marshal(host)
			if err != nil {
				return err
			}
			b.Put([]byte(HostPrefix+h.ID), hostData)

		}
		return nil
	})
	return &resultHost
}

func (l *ProxyLog) AssignStream(s string) (*EventoHost, error) {
	var host EventoHost
	var err error
	plog.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("streams"))
		if err != nil {
			return fmt.Errorf("no exist bucket: %s", err)
		}

		HostID := b.Get([]byte(StreamPrefix + s))
		if len(HostID) != 0 {
			// Get host and return
			hdata := b.Get([]byte(HostPrefix + string(HostID)))
			if len(hdata) == 0 {
				err = HostNotRegistered
				return nil
			} else {
				err = json.Unmarshal(hdata, &host)
			}

			return nil
		} else {
			// TODO get all hosts and select one
			c := b.Cursor()
			openHost := make([]EventoHost, 0)
			prefix := []byte(HostPrefix)
			for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
				fmt.Printf("key=%s, value=%s\n", k, v)
				var h EventoHost
				json.Unmarshal(v, &h)
				if !h.Closed {
					openHost = append(openHost, h)
				}
			}

			sort.Sort(ByStream(openHost))

			if len(openHost) == 0 {
				err = NoAvailableHost
				return nil
			} else {
				openHost[0].Streams++
				host = openHost[0]
				updatedHost, _ := json.Marshal(openHost[0])
				b.Put([]byte(HostPrefix+openHost[0].ID), updatedHost)
				b.Put([]byte(StreamPrefix+s), []byte(openHost[0].ID))
			}
		}

		return nil
	})

	return &host, err
}

func Init(db *bolt.DB) {
	l := NewProxyLog(db)
	plog = l
	// create stream bucket
	plog.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("streams"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	engine := gin.Default()
	engine.POST("/reg/host", RegHost)
	engine.GET("/reg/:stream", QueryStream)
	engine.Run(":6262")
}
