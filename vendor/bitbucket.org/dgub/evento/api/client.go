package api

import (
	"bitbucket.org/dgub/evento/ap"
	"bitbucket.org/dgub/evento/dom"
	"bitbucket.org/dgub/evento/proxy/use"
	"errors"
	nap "github.com/jmcvetta/napping"
	"strconv"
	"strings"
	"sync"
)

var (
	DevMode bool
)

const (
	RPC_CLIENT      = "rpc"
	EVENTS_PER_PAGE = 10
)

type Client struct {
	lock sync.RWMutex

	eventoServer string
	eventoProxy  string
	client       *nap.Session
	dev          bool
}

type ConnectOpt struct {
	EventoServer string
	EventoProxy  string
}

func NewClient(copt ConnectOpt) *Client {
	var c Client
	var s nap.Session
	c.eventoProxy = copt.EventoProxy
	c.eventoServer = copt.EventoServer
	c.client = &s
	return &c
}

type StoreOpt struct {
	Sync            bool
	Create          bool
	Lock            bool
	ExpectedVersion uint64
}

func (c *Client) StoreEvent(event Eventer, opt *StoreOpt) (uint64, error) {
	host, err := c.getServer(event.GetStream())
	if err != nil {
		return 0, err
	}

	if opt == nil {
		opt = &StoreOpt{
			Sync: false,
		}
	}
	request := buildRequest(event, opt.Create, opt.Sync, opt.Lock, opt.ExpectedVersion)

	path := host + "/stream/" + event.GetStream()
	pathOpt := make([]string, 0)

	if opt.Create {
		pathOpt = append(pathOpt, "create=true")
	}

	if opt.Lock {
		pathOpt = append(pathOpt, "lock="+strconv.FormatUint(opt.ExpectedVersion, 10))
	}

	if len(pathOpt) > 0 {
		path += "?" + strings.Join(pathOpt, "&")
	}

	response := make(map[string]interface{})

	res, err := c.client.Post(path, request.Event, &response, &response)
	if err != nil {
		return 0, err
	}

	if res.Status() != 200 {
		return 0, errors.New("Failed to store event")
	}

	return uint64(response["version"].(float64)), err
}

type Version struct {
	Version uint64 `json:"version"`
	Error   string `json:"error"`
}

func (c *Client) Version(streamid string) (uint64, error) {
	host, err := c.getServer(streamid)
	if err != nil {
		return 0, err
	}

	v := Version{}
	path := host + "/stream/" + streamid

	c.client.Get(path, nil, &v, &v)
	if v.Error != "" {
		err = errors.New(v.Error)
	}
	return v.Version, err
}

type rangeResponse struct {
	Events []dom.Event `json:"events"`
}

func (c *Client) DeleteStream(streamid string) error {
	host, err := c.getServer(streamid)
	if err != nil {
		return err
	}

	path := host + "/stream/" + streamid
	res, err := c.client.Delete(path, nil, nil, nil)
	if err != nil {
		return err
	}

	switch res.Status() {
	case 200:
		return nil
	default:
		return errors.New("Failed to delete stream, already deleted")
	}
}

func (c *Client) PurgeStream(streamid string) error {
	host, err := c.getServer(streamid)
	if err != nil {
		return err
	}
	path := host + "/stream/" + streamid + "/purge"
	res, err := c.client.Put(path, nil, nil, nil)
	if err != nil {
		return err
	}

	switch res.Status() {
	case 200:
		return nil
	default:
		return errors.New("Failed to purge stream")
	}
}

func (c *Client) RangeStream(streamid string, from, to uint64) chan dom.Event {
	chEvent := make(chan dom.Event, 100000)
	host, err := c.getServer(streamid)
	if err != nil {
		close(chEvent)
		return chEvent
	}

	if from > to {
		close(chEvent)
		return chEvent
	}

	go func() {
		for {
			var stop bool
			// fetch 100
			dTo := from + 100
			if dTo > to {
				dTo = to
				stop = true
			}

			path := host + "/range/" + streamid + "/" + strconv.FormatUint(from, 10) + "/" + strconv.FormatUint(dTo, 10)

			var res rangeResponse
			c.client.Get(path, nil, &res, &res)
			for _, e := range res.Events {
				chEvent <- e
			}

			from = dTo
			if stop {
				close(chEvent)
				break
			}

		}
	}()

	return chEvent
}

type streamResponse struct {
	Streams []dom.Stream
}

func (c *Client) StreamMatch(regex string) chan dom.Stream {
	// TODO parse all servers
	chStream := make(chan dom.Stream, 100000)
	go func() {
		var re streamResponse
		query := ap.StreamQuery{
			Regex: regex,
		}

		path := c.eventoServer + "/stream"
		c.client.Put(path, query, &re, &re)

		for _, s := range re.Streams {
			chStream <- s
		}
		close(chStream)
	}()
	return chStream
}

func (c *Client) SetDevMode(on bool) {
	c.dev = on
	c.client.Log = true
}

func (c *Client) Ping() error {
	// TODO ping all servers in proxy
	res, err := c.client.Get(c.eventoServer+"/ping", nil, nil, nil)
	if err != nil {
		return err
	}

	if res.Status() != 200 {
		return errors.New("Failed to ping, error:")
	}
	return nil
}

func (c *Client) getServer(stream string) (string, error) {
	var eventoHost use.EventoHost
	if c.eventoProxy != "" {
		_, err := c.client.Get(c.eventoProxy+"/reg/"+stream, nil, &eventoHost, &eventoHost)
		if err != nil {
			return "", err
		}
		return eventoHost.URL, nil
	} else {
		return c.eventoServer, nil
	}
}
