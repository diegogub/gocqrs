package api

import (
	"bitbucket.org/dgub/evento/dom"
)

func (c *Client) ReadStream(id string) (chan dom.Event, error) {
	ch := make(chan dom.Event)

	version, err := c.Version(id)
	if err != nil {
		close(ch)
		return ch, err
	}

	go func(ch *chan dom.Event, version uint64) {

		var n uint64 = 100
		var current uint64 = 0
		for {
			events := c.RangeStream(id, current, current+n)
			for e := range events {
				*ch <- e
			}

			current += n
			if current >= version {
				close(*ch)
				break
			}
		}
	}(&ch, version)

	return ch, nil
}

func (c *Client) Exist(stream string) bool {
	_, err := c.Version(stream)
	if err != nil {
		return false
	}
	return true
}

func (c *Client) NewStream(eventer Eventer) error {
	opt := &StoreOpt{
		Create: true,
	}

	_, err := c.StoreEvent(eventer, opt)
	return err
}

func (c *Client) SaveEvent(eventer Eventer) error {
	_, err := c.StoreEvent(eventer, nil)
	return err
}
